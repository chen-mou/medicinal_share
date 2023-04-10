package redis

import (
	"context"
	"math"
	"medicinal_share/tool/encrypt/md5"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type rlock struct {
	state int32
	val   string
	key   string
}

func RLock(key string) *rlock {
	return &rlock{
		key: key,
	}
}

const (
	scriptUnLock = `
		if (redis.call('EXISTS', KEYS[1]) == 0) then
			return -1;
		end;
		if (redis.call('hget', KEYS[1], 'getter') ~= ARGV[1]) then
			return {err = "it not you rlock"};
		end;
		local count = redis.call('hincrby', KEYS[1], 'count', -1);
		if (count == 0) then
			redis.call('del', KEYS[1])
			return -1;
		end;
		return 0;`
	scriptTryLock = `
		if (redis.call('EXISTS', KEYS[1]) == 0) then
			redis.call('hset', KEYS[1], 'count', 1);
			redis.call('hset', KEYS[1], 'getter', ARGV[1]);
			redis.call('expire', KEYS[1], ARGV[2]);
			return 1;
		end;
		if (redis.call('hget', KEYS[1], 'getter') == ARGV[1]) then
			redis.call('hincrby', KEYS[1], 'count', 1);
			return 1;
		end;
		return 0;`
)

var (
	unlock  string
	tryLock string
)

type SoleId struct {
	sync.Mutex
	name string
	val  int64
	step int
}

var soleIdMap sync.Map

func GetSole(name string, step int) *SoleId {
	res := &SoleId{
		name: name,
		val:  0,
		step: step,
	}
	soleIdMap.LoadOrStore(name, res)
	return res
}

func init() {
	for DB == nil {
	}
	var err error
	scripts := []string{scriptTryLock, scriptUnLock}
	addr := []*string{&tryLock, &unlock}
	for i := 0; i < len(scripts); i++ {
		cmd := DB.ScriptLoad(context.TODO(), scripts[i])
		*addr[i], err = cmd.Result()
		if err != nil {
			panic(err)
		}
	}
}

func (sole *SoleId) GetID() string {
	sole.Lock()
	const mod = math.MaxInt64 >> 2
	sole.val = (sole.val + int64(sole.step)) % (mod)
	sole.Unlock()
	str := sole.name + strconv.FormatInt(sole.val, 10) + time.Now().String()
	return md5.Hash(str)
}

//Lock 这个方法有点问题会导致redis卡死
//func (l *rlock) Lock(tim time.Duration) {
//	for atomic.LoadInt32(&l.state) == 1 {
//	}
//	val := GetSole("RedisId", 8).GetID()
//	DB.EvalSha(context.TODO(), lock, []string{l.key}, val, int64(tim.Seconds()))
//	l.state = 1
//	l.val = val
//	go l.watchDog(tim)
//}

// watchDog 超时续
func (l *rlock) watchDog(tim time.Duration) {
	for {
		ctx, _ := context.WithTimeout(context.TODO(), tim-time.Second)
		select {
		case <-ctx.Done():
			if atomic.LoadInt32(&l.state) == 1 {
				DB.Expire(context.TODO(), l.key, tim)
			} else {
				break
			}
		}
	}
}

func (l *rlock) Unlock() {
	res, err := DB.EvalSha(context.TODO(), unlock, []string{l.key}, l.val).Result()
	if err != nil {
		panic(err)
	}
	if res == int64(-1) {
		l.state = 0
	}

}

// TryLockWithTime exp 是锁的过期时间 tim 是过期时间
func (l *rlock) TryLockWithTime(exp, tim time.Duration) bool {
	ctx, _ := context.WithTimeout(context.TODO(), tim)
	if l.TryLock(exp) {
		return true
	}
	for {
		select {
		case <-ctx.Done():
			return false
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (l *rlock) TryLock(exp time.Duration) bool {
	val := GetSole("RedisId", 8).GetID()
	res, err := DB.EvalSha(context.TODO(), tryLock, []string{l.key}, val, int64(exp.Seconds())).Result()
	if err != nil {
		panic(err)
	}
	if res == int64(1) {
		l.val = val
		l.state = 1
		go l.watchDog(exp)
		return true
	}
	return false
}
