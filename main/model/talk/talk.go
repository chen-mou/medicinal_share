package talk

import (
	"context"
	"errors"
	redis2 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"strconv"
	"time"
)

func GetUserNowRoom(usr int64) *entity.Room {
	room := &entity.Room{}
	ustr := strconv.FormatInt(usr, 10)
	res, err := redis.DB.Get(context.TODO(), "user-now-room-"+ustr).Result()
	if err != nil {
		return nil
	}
	err = redis.HGet("room-"+res, room)
	if err != nil {
		return nil
	}
	return room
}

func CreateRoom(room *entity.Room, roomId string) error {
	redis.HSet("room-"+roomId, room, -1)
	_, err := redis.DB.Pipelined(context.TODO(), func(pipe redis2.Pipeliner) error {
		ustr := strconv.FormatInt(room.Custom, 10)
		dstr := strconv.FormatInt(room.Doctor, 10)
		pipe.Set(context.TODO(), "user-now-room-"+ustr, roomId, -1)
		pipe.Set(context.TODO(), "user-now-room"+dstr, roomId, -1)
		return nil
	})
	return err
}

func SaveMessage(msg *entity.Message, tx *gorm.DB) error {
	return tx.Save(msg).Error
}

func GetLastMessage(user1 int64, user2 int64, tim time.Time) []*entity.Message {
	db := mysql.GetConnect()
	model := &entity.Message{}
	res := make([]*entity.Message, 0)
	err := db.Table("(?) temp", db.Raw("? union ?",
		db.Select("*, 'send' as kind").Model(model).
			Where("sender = ? and getter = ?", user1, user2),
		db.Select("*, 'get' as kind").Model(model).
			Where("getter = ? and sender = ?", user1, user2))).
		Where("time > ?", entity.CreateTime(tim)).
		Order("time desc").Limit(10).Find(&res).Error
	db.Model(model).Update("status", 1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(err)
	}
	return res
}

func CountUnReadMessage(user1 int64, user2 int64) map[int64]int {
	res := map[int64]int{}
	db := mysql.GetConnect()
	model := &entity.Message{}
	err := db.Model(model).
		Select("sender, count(*) as num").Where(map[string]any{
		"status": 0,
		"getter": user1,
	}).
		Group("sender").
		Scan(&res).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}
	return res
}
