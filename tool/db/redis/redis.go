package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"reflect"
	"time"
)

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
	if val == -1 {
		return redis.Nil
	}
	c := DB.HGetAll(context.TODO(), key)
	r, err := c.Result()
	if err != nil {
		return err
	}
	if len(r) == 0 {
		return redis.Nil
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
		DB.Set(context.TODO(), key+":expire", "", expire)
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

func (c Cache) getCache(val any) any {
	cmd := DB.Get(context.TODO(), c.key)
	v, err := cmd.Result()
	if err == redis.Nil {
		return nil
	}
	json.Unmarshal([]byte(v), val)
	return val
}

//Get 安全的从数据库和缓存中拿取对象参数要是指针
func (c Cache) Get(val any, getter func() any) any {
	v := c.getCache(val)
	if v == nil {
		lock := RLock(c.lock)
		if lock.TryLock(3 * time.Second) {
			val = getter()
			if val != nil {
				jsn, _ := json.Marshal(val)
				DB.Set(context.TODO(), c.key, string(jsn), 20*time.Minute)
			} else {
				DB.Set(context.TODO(), c.key, "", 5*time.Minute)
			}
		} else {
			val = c.getCache(val)
			for val == nil {
				time.Sleep(10 * time.Millisecond)
				val = c.getCache(val)
			}
		}
	}
	return val
}

//HGet 安全的拿取对象用HGet参数要是指针 TODO:完成这个方法
func (c Cache) HGet(val any) any {
	panic("方法还没完成")
}

//SafeGet 安全的获取缓存或者数据库中的数据 TODO:重写这个方法
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
