package entity

//TODO:完成订单以及交易流水的模型定义

type Status uint8

const (
	Padding Status = iota
	Expired
	UnUsing
	Complete
)

type Order struct {
	Id       int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId   int64        `json:"user_id" gorm:"index"`
	Price    float64      `json:"price"`
	Status   Status       `json:"status" gorm:"size:16"`
	Version  string       `json:"version" gorm:"-"`
	Data     []*OrderData `json:"data" gorm:"foreignKey:OrderId"`
	CreateAt Time         `json:"create_at"`
}

type OrderData struct {
	Id               int64          `json:"id" gorm:"primaryKey"`
	OrderId          int64          `json:"order_id"`
	ProjectReserveId int64          `json:"project_reserve_id"`
	ProjectReserve   ProjectReserve `json:"project" gorm:"foreignKey:ProjectReserveId"`
}

func (Order) TableName() string {
	return "order"
}

type Pay struct {
}
