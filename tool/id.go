package tool

import (
	"context"
	redis2 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"time"
)

//TODO: 分布式ID构思MYSQL一个表加步长保证ID唯一

const step = 16

func GetId(name string) (int64, error) {
	key := name + ":ID"
	_, err := redis.DB.Get(context.TODO(), key).Result()
	if err == redis2.Nil {
		lock := redis.RLock(name + ":ID:LOCK")
		if lock.TryLock(3 * time.Second) {
			defer lock.Unlock()
			var id int
			err := mysql.GetConnect().Table(name).Order("id desc").Limit(1).Pluck("id", &id).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return 0, err
				}
			}
			redis.DB.Do(context.TODO(), "set", key, id)

		} else {
			_, err := redis.DB.Get(context.TODO(), key).Result()
			for err == redis2.Nil {
				_, err = redis.DB.Get(context.TODO(), key).Result()
				time.Sleep(10 * time.Millisecond)
			}
		}
	} else if err != nil {
		return 0, err
	}
	res, _ := redis.DB.IncrBy(context.TODO(), key, step).Result()
	return res, nil
}
