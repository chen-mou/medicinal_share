package entity

// Room TODO:用户和医生的聊天室的实体
type Room struct {
	Id     int64
	Custom int64
	Doctor int64
}

type MessageBase struct {
	Id     int64
	Type   string
	Main   string
	Sender int64
	Getter int64
	Time   Time
}

type MessageGoods struct {
	MessageBase
}

type MessageResult struct {
	MessageBase
}
