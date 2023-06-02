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
	r, _ := redis.DB.HGet(context.TODO(), "room-"+res, "tags_id").Result()
	room.TagsId = entity.Ints64{}
	room.TagsId.UnmarshalBinary([]byte(r))
	if err != nil {
		return nil
	}
	room.Id = res
	return room
}

func UpdateRoomStatus(roomId string, status entity.RoomStatus) {
	room := GetRoom(roomId)
	if room == nil {
		return
	}
	room.Status = status
	CreateRoom(room, roomId)
}

func CreateRoom(room *entity.Room, roomId string) error {
	_, err := redis.DB.Pipelined(context.TODO(), func(pipe redis2.Pipeliner) error {
		err := redis.PipeHSet(pipe, "room-"+roomId, room, -1)
		if err != nil {
			return err
		}
		ustr := strconv.FormatInt(room.CustomId, 10)
		dstr := strconv.FormatInt(room.DoctorId, 10)
		_, err = pipe.Do(context.TODO(), "set", "user-now-room-"+ustr, roomId).Result()
		if err != nil {
			return err
		}
		_, err = pipe.Do(context.TODO(), "set", "user-now-room-"+dstr, roomId).Result()
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetRoom(roomId string) *entity.Room {
	res := &entity.Room{}
	err := redis.HGet("room-"+roomId, res)
	if err != nil {
		if err == redis.RedisEmpty {
			return nil
		}
		panic(err)
	}
	return res
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
