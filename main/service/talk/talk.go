package talk

import (
	"context"
	"encoding/json"
	redis2 "github.com/go-redis/redis/v8"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/talk"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"medicinal_share/tool/encrypt/md5"
	"medicinal_share/tool/socket"
	"strconv"
	"time"
)

func Send(sender int64, msg string, typ string, c *socket.Conn) {
	room := talk.GetUserNowRoom(sender)
	if room == nil {
		panic(middleware.NewCustomErr(middleware.NOT_FOUND, "房间不存在"))
	}
	if room.Status != entity.Talking {
		panic(middleware.NewCustomErr(middleware.ERROR, "房间已关闭或医生还未准备好"))
	}
	var getter int64
	if room.Doctor == sender {
		getter = room.Custom
	} else {
		getter = room.Doctor
	}
	now := time.Now()
	message := &entity.Message{
		Sender: sender,
		Getter: getter,
		Type:   typ,
		Main:   msg,
		Time:   &now,
	}
	byt, _ := json.Marshal(message)
	c.SendTo(string(byt), getter)

	go func() {
		talk.SaveMessage(message, mysql.GetConnect())
	}()
}

func CreateRoom(userId int64, doctor int64) string {
	roomId := md5.Hash(time.Now().String())
	room := entity.Room{
		Id:     roomId,
		Custom: userId,
		Doctor: doctor,
		Status: entity.Waiting,
	}
	redis.HSet("room-"+roomId, room, -1)
	redis.DB.Pipelined(context.TODO(), func(pipe redis2.Pipeliner) error {
		ustr := strconv.FormatInt(userId, 10)
		dstr := strconv.FormatInt(doctor, 10)
		pipe.Set(context.TODO(), "user-now-room-"+ustr, roomId, -1)
		pipe.Set(context.TODO(), "user-now-room"+dstr, roomId, -1)
		return nil
	})
	return roomId
}
