package talk

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/talk"
	"medicinal_share/main/model/user"
	"medicinal_share/tool"
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
	var avatar string
	if room.DoctorId == sender {
		doc := user.GetDataByUserId(sender)
		avatar = doc.AvatarFile.File.Uri
		getter = room.CustomId
	} else {
		usr := user.GetDataByUserId(sender)
		avatar = usr.AvatarFile.File.Uri
		getter = room.DoctorId
	}
	message := &entity.Message{
		Sender:       sender,
		Getter:       getter,
		Type:         typ,
		Main:         msg,
		SenderAvatar: avatar,
		Time:         entity.CreateTime(time.Now()),
		Status:       0,
	}
	byt, _ := json.Marshal(message)
	socket.SendTo("talk", string(byt), getter)
	socket.SendTo("talk", string(byt), sender)
	go func() {
		talk.SaveMessage(message, mysql.GetConnect())
	}()
}

func Confirm(doc int64, c *socket.Conn) {
	room := talk.GetUserNowRoom(doc)
	if room.DoctorId != doc {
		socket.SendCustomError(c, errors.New("你没有权限"))
	}
	talk.UpdateRoomStatus(room.Id, entity.Talking)
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		tx.Model(&entity.Tag{}).Where("id in (?)", room.TagsId).Find(&room.Tags)
		tx.Model(&entity.DoctorInfo{}).
			Where("user_id = ?", room.DoctorId).
			Preload("Info").
			Preload("Avatar").
			Preload("Avatar.File").
			Take(&room.Doctor)
		tx.Model(&entity.UserData{}).
			Where("user_id = ?", room.CustomId).
			Preload("RealInfo").
			Preload("AvatarFile").
			Preload("AvatarFile.File").
			Take(&room.Custom)
		return nil
	})
	byt, err := json.Marshal(room)
	if err != nil {
		panic(err)
	}
	socket.SendTo("confirm", "", room.CustomId)
	socket.SendTo("room_info", tool.BytesToString(byt), room.DoctorId)
	socket.SendTo("room_info", tool.BytesToString(byt), room.CustomId)
}

func CreateRoom(userId int64, doctor int64, tags []int64) string {
	roomId := md5.Hash(time.Now().String() + strconv.FormatInt(userId, 10) + strconv.FormatInt(doctor, 10))
	room := &entity.Room{
		CustomId: userId,
		DoctorId: doctor,
		Status:   entity.Waiting,
		TagsId:   tags,
	}
	err := talk.CreateRoom(room, roomId)
	if err != nil {
		panic(err)
	}
	return roomId
}

func GetRoom(roomId string) *entity.Room {
	res := talk.GetRoom(roomId)
	if res == nil {
		return nil
	}
	mysql.GetConnect().Model(&entity.Tag{}).Where("id in (?)", res.TagsId).Find(&res.Tags)
	return res
}

func Treat(userId int64, tags []int64, long float64, latit float64) string {
	usr := user.GetDataByUserId(userId)
	if !usr.IsReal {
		panic(middleware.NewCustomErr(middleware.ERROR, "请先实名"))
	}
	//doctorId := user.GetBestDoctor(tags, long, latit)
	doctorId := user.GetBestDoctor(tags, long, latit)
	room := CreateRoom(userId, doctorId, tags)

	socket.SendTo("start", "", doctorId)
	return room
}
