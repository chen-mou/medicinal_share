package talk

import (
	"encoding/json"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/talk"
	"medicinal_share/main/model/user"
	"medicinal_share/tool/db/mysql"
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
	if room.DoctorId == sender {
		getter = room.CustomId
	} else {
		getter = room.DoctorId
	}
	message := &entity.Message{
		Sender: sender,
		Getter: getter,
		Type:   typ,
		Main:   msg,
		Time:   entity.CreateTime(time.Now()),
		Status: 0,
	}
	byt, _ := json.Marshal(message)
	socket.SendTo(string(byt), getter)
	socket.SendTo(string(byt), sender)
	go func() {
		talk.SaveMessage(message, mysql.GetConnect())
	}()
}

func CreateRoom(userId int64, doctor int64) string {
	roomId := md5.Hash(time.Now().String() + strconv.FormatInt(userId, 10) + strconv.FormatInt(doctor, 10))
	room := &entity.Room{
		CustomId: userId,
		DoctorId: doctor,
		Status:   entity.Waiting,
	}
	err := talk.CreateRoom(room, roomId)
	if err != nil {
		panic(err)
	}
	return roomId
}

func Treat(userId int64, tags []int64, long float64, latit float64) string {
	usr := user.GetDataByUserId(userId)
	if !usr.IsReal {
		panic(middleware.NewCustomErr(middleware.ERROR, "请先实名"))
	}
	//doctorId := user.GetBestDoctor(tags, long, latit)
	doctorId := user.GetBestDoctorTest()
	room := CreateRoom(userId, doctorId)
	return room
}
