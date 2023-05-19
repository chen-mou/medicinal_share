package entity

import "encoding/json"

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
	Tags     []int64    `redis:"tags" json:"tags"`
	Result   string     `json:"result" redis:"result"`
	Status   RoomStatus `redis:"status" json:"status"`
	Doctor   DoctorInfo `redis:"-" gorm:"foreignKey:DoctorId"`
	Custom   UserData   `redis:"-" gorm:"foreignKey:CustomId"`
}

type Message struct {
	Id     int64  `json:"id" gorm:"primaryKey;autoIncrement"'`
	Type   string `json:"type" gorm:"size:16"`
	Main   string `json:"main"`
	RoomId string `json:"room_id"`
	Sender int64  `json:"sender"`
	Getter int64  `json:"getter"`
	Time   Time   `json:"time" gorm:"type:datetime;index:status_time_idx"`
	Status uint8  `json:"status" gorm:"index:status_time_idx"`
}
