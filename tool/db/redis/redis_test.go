package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

//func do(f func()) {
//	f()
//}
//
//func TestLock(t *testing.T) {
//	f := func(val int) {
//		fmt.Println(val)
//	}
//	for i := 0; i < 10; i++ {
//		do(func() {
//			f(i)
//		})
//	}
//}
//
//func TestAntiShake(t *testing.T) {
//	wg := sync.WaitGroup{}
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			if !AntiShake("test") {
//				fmt.Println("防抖成功")
//			}
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//}

func TestFuck(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	cmd := r.Get(context.TODO(), "a")
	res, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println(err.Error())
		}
		panic(err)
	}
	fmt.Println(res)
}
