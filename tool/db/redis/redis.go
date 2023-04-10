package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"medicinal_share/main/middleware"
	"reflect"
	"strconv"
	"time"
)

type Empty struct {
}

func (Empty) Error() string {
	return "空"
}

var RedisEmpty = Empty{}

var DB *redis.ClusterClient

func init() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":6380", ":6381", ":6382", ":6480", ":6481", ":6482"},
	})
	DB = rdb
	var err error
	antiShake, err = DB.ScriptLoad(context.TODO(), AntiShakeScript).Result()
	if err != nil {
		panic(err)
	}
}

func Get(key string, v any) error {
	cmd := DB.Get(context.TODO(), key)
	str, err := cmd.Result()
	if err != nil {
		return err
	}
	if str != "" {
		return redis.Nil
	}
	return json.Unmarshal([]byte(str), v)
}

func Set(key string, v any, expire time.Duration) error {
	val, err := json.Marshal(v)
	if err != nil {
		return err
	}
	DB.Set(context.TODO(), key, string(val), expire)
	return nil
}

func HGet(key string, v any) error {
	cmd := DB.TTL(context.TODO(), key+":expire")
	val, err := cmd.Result()
	if err != nil {
		return err
	}
	if val == -1 || val == -2 {
		return redis.Nil
	}
	c := DB.HGetAll(context.TODO(), key)
	r, err := c.Result()
	if err != nil {
		return err
	}
	if len(r) == 0 {
		return RedisEmpty
	}
	return c.Scan(v)
}

func HSet(key string, v any, expire time.Duration) error {
	_, err := DB.Pipelined(context.TODO(), func(pipeliner redis.Pipeliner) error {
		val := reflect.ValueOf(v)
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Pointer {
			val = val.Elem()
			t = t.Elem()
		}
		for i := 0; i < val.NumField(); i++ {
			fieldv := val.Field(i)
			fieldt := t.Field(i)
			name, ok := fieldt.Tag.Lookup("json")
			if !ok {
				name = fieldt.Name
			}
			DB.HSet(context.TODO(), key, name, fieldv.Interface())
		}
		if v == nil {
			DB.Set(context.TODO(), key+":expire", "", expire)
		} else {
			DB.Set(context.TODO(), key+":expire", "value", expire)
		}
		return nil
	})
	return err
}

type Cache struct {
	lock string
	key  string
}

func NewCache(lock string, key string) *Cache {
	return &Cache{lock: lock, key: key}
}

func (c Cache) getCache(val any) (any, error) {
	cmd := DB.Get(context.TODO(), c.key)
	v, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	if v == "" {
		return nil, nil
	}
	json.Unmarshal([]byte(v), val)
	return val, nil
}

func (c Cache) getHCache(val any) (any, error) {
	err := HGet(c.key, val)
	if err == RedisEmpty {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Get 安全的从数据库和缓存中拿取对象参数要是指针
func (c Cache) Get(val any, getter func() any) any {
	v, err := c.getCache(val)
	if err != nil && err == redis.Nil {
		lock := RLock(c.lock)
		if !lock.TryLockWithTime(3*time.Second, 3*time.Second) {
			v, err = c.getCache(val)
			if err != nil && err == redis.Nil {
				v = getter()
				if val != nil {
					jsn, _ := json.Marshal(v)
					DB.Set(context.TODO(), c.key, string(jsn), 20*time.Minute)
				} else {
					DB.Set(context.TODO(), c.key, "", 5*time.Minute)
				}
			}
		} else {
			panic(middleware.NewCustomErr(middleware.ERROR, "服务器繁忙"))
		}
	}
	if err != nil {
		panic(err)
	}
	return v
}

// HGet 安全的拿取对象用HGet参数要是指针 TODO:完成这个方法
func (c Cache) HGet(val any, getter func() any) any {
	val, err := c.getHCache(val)
	if err != nil && errors.Is(err, Empty{}) {
		lock := RLock(c.lock)
		if lock.TryLockWithTime(3*time.Second, 3*time.Second) {
			val, err = c.getHCache(val)
			if err != nil && err == RedisEmpty {
				val = getter()
				if val != nil {
					HSet(c.key, val, 20*time.Minute)
				} else {
					HSet(c.key, val, 5*time.Minute)
				}
			} else {
				panic(middleware.NewCustomErr(middleware.ERROR, "服务器繁忙"))
			}
		}
	}
	if err != nil {
		panic(err)
	}
	return val
}

// LoadInt 加载一个数字到redis里
func (c Cache) LoadInt(getter func() (int, error)) (int, error) {
	var val int
	var s string
	var err error
	s, err = DB.Get(context.TODO(), c.key).Result()
	if err != nil && err == redis.Nil {
		lock := RLock(c.lock)
		if lock.TryLockWithTime(3*time.Second, 3*time.Second) {
			s, err = DB.Get(context.TODO(), c.key).Result()
			if err != nil && err == redis.Nil {
				val, err = getter()
				if err != nil {
					DB.Set(context.TODO(), c.key, strconv.Itoa(val), 20*time.Minute)
				} else {
					DB.Set(context.TODO(), c.key, "", 5*time.Minute)
				}
			}
		} else {
			panic(middleware.NewCustomErr(middleware.ERROR, "服务器繁忙"))
		}
	}
	if err != nil {
		panic(err)
	}
	if s == "" {
		return 0, redis.Nil
	}
	val, _ = strconv.Atoi(s)
	return val, nil
}

// SafeGet 安全的获取缓存或者数据库中的数据 TODO:重写这个方法
func SafeGet(key, lockKey string, cache func() any, getter func() any) any {
	val := cache()
	v := reflect.ValueOf(val)
	if v.IsNil() {
		lock := RLock(lockKey)
		if lock.TryLock(3 * time.Second) {
			val = getter()
			if val != nil {
				jsn, _ := json.Marshal(val)
				DB.Set(context.TODO(), key, string(jsn), 20*time.Minute)
			} else {
				DB.Set(context.TODO(), key, "", 5*time.Minute)
			}
		} else {
			val = cache()
			for val == nil {
				val = cache()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
	return val
}

const AntiShakeScript = `
	if (redis.call('EXISTS', KEYS[1]) == 0) then
		redis.call('SET', KEYS[1], 'ANTI_SHAKE')
		redis.call('EXPIRE', KEYS[1], 60)
		return 1;
	end;
	return 0;
`

var antiShake string

func AntiShake(key string) bool {
	res, _ := DB.EvalSha(context.TODO(), antiShake, []string{key}).Result()
	return res == 1
}
