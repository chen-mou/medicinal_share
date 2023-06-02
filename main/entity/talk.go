package entity

import (
	"encoding/json"
	"strconv"
	"strings"
)

type RoomStatus uint8

func (r *RoomStatus) UnmarshalJSON(bytes []byte) error {
	i := 0
	json.Unmarshal(bytes, &i)
	*r = RoomStatus(i)
	return nil
}

func (r RoomStatus) MarshalBinary() (data []byte, err error) {
	byt, _ := json.Marshal(int(r))
	return byt, nil
}

type Ints64 []int64

type Tags []*Tag

func (r *Ints64) UnmarshalJSON(bytes []byte) error {
	strs := strings.Split(string(bytes), ",")
	var res = make([]int64, 0)
	for _, v := range strs {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		res = append(res, i)
	}
	*r = res
	return nil
}

func (r Ints64) MarshalBinary() (data []byte, err error) {
	s := make([]string, 0)
	for _, v := range r {
		s = append(s, strconv.FormatInt(v, 10))
	}
	return []byte(strings.Join(s, ",")), nil
}

func (r *Ints64) UnmarshalBinary(byt []byte) error {
	s := string(byt)
	if s == "" {
		return nil
	}
	ss := strings.Split(s, ",")
	res := make([]int64, 0)
	for _, v := range ss {
		num, _ := strconv.ParseInt(v, 10, 64)
		res = append(res, num)
	}
	*r = res
	return nil
}

func (r *Tags) UnmarshalJSON(byt []byte) error {
	strs := strings.Split(string(byt[1:len(byt)-1]), ",")
	val := make([]*Tag, 0)
	for _, s := range strs {
		tag := &Tag{}
		json.Unmarshal([]byte(s), tag)
		val = append(val, tag)
	}
	*r = val
	return nil
}

func (r *Tags) UnmarshalBinary(byt []byte) error {
	r.UnmarshalJSON(byt)
	return nil
}

func (r Tags) MarshalBinary() (data []byte, err error) {
	return r.MarshalJSON()
}

func (r Tags) MarshalJSON() (data []byte, err error) {
	byt := make([]byte, 0)
	for _, tag := range r {
		b, _ := json.Marshal(tag)
		byt = append(byt, b...)
		byt = append(byt, byte(','))
	}
	if len(byt) == 0 {
		return []byte{'[', ']'}, nil
	}
	return append(append([]byte{'['}, byt[:len(byt)-1]...), ']'), nil
}

const (
	Waiting RoomStatus = iota
	Talking
	Closed
)

// Room TODO:用户和医生的聊天室的实体
type Room struct {
	Id       string     `redis:"-" json:"id" gorm:"primaryKey;type:varchar(32)"`
	CustomId int64      `redis:"custom_id" json:"custom_id"`
	DoctorId int64      `redis:"doctor_id" json:"doctor_id"`
	TagsId   Ints64     `redis:"tags_id" json:"tags_id"`
	Tags     Tags       `redis:"-" json:"tags"`
	Result   string     `json:"result" redis:"result"`
	Status   RoomStatus `redis:"status" json:"status"`
	Doctor   DoctorInfo `json:"doctor" redis:"-" gorm:"foreignKey:DoctorId"`
	Custom   UserData   `json:"custom" redis:"-" gorm:"foreignKey:CustomId"`
}

type Message struct {
	Id           int64  `json:"id" gorm:"primaryKey;autoIncrement"'`
	Type         string `json:"type" gorm:"size:16"`
	Main         string `json:"main"`
	RoomId       string `json:"room_id"`
	Sender       int64  `json:"sender"`
	SenderAvatar string `json:"sender_avatar" gorm:"-"`
	Getter       int64  `json:"getter"`
	Time         Time   `json:"time" gorm:"type:datetime;index:status_time_idx"`
	Status       uint8  `json:"status" gorm:"index:status_time_idx"`
}
