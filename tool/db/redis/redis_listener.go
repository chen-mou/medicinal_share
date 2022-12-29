package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"medicinal_share/tool/pool"
	"regexp"
	"sync"
)

var pb *redis.PubSub

var tasks map[string][]func(payload string)

var listenPool = pool.NewPool(40)

var one = sync.Once{}

func init() {
	for DB == nil {
	}
	pb = DB.PSubscribe(context.TODO(), "*")
	go func() {
		for {
			msg, err := pb.ReceiveMessage(context.TODO())
			if err != nil {
				//TODO: 打印错误
			}
			ch := msg.Channel
			for k, v := range tasks {
				reg, _ := regexp.Compile(k)
				if reg.Match([]byte(ch)) {
					for _, f := range v {
						listenPool.Submit(context.TODO(), func() {
							f(msg.Payload)
						})
					}
				}
			}
		}
	}()
}

func NewTask(pattern string, f func(string)) {
	one.Do(func() {
		_, ok := tasks[pattern]
		if !ok {
			tasks[pattern] = make([]func(string), 0)
		}
		tasks[pattern] = append(tasks[pattern], f)
	})
}
