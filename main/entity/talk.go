package entity

type RoomStatus uint8

const (
	Waiting RoomStatus = iota
	Talking
	Closed
)

// Room TODO:用户和医生的聊天室的实体
type Room struct {
	Custom int64
	Doctor int64
	Status RoomStatus
}

type Message struct {
	Id     int64  `json:"id" gorm:"primaryKey;autoIncrement"'`
	Type   string `json:"type" gorm:"size:16"`
	Main   string `json:"main"`
	Sender int64  `json:"sender"`
	Getter int64  `json:"getter"`
	Time   Time   `json:"time" gorm:"type:datetime;index:status_time_idx"`
	Status uint8  `json:"status" gorm:"index:status_time_idx"`
}
