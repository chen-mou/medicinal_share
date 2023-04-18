package entity

//TODO:完成订单以及交易流水的模型定义

type OrderStatus uint8

const (
	Padding OrderStatus = iota
	Expired
	UnUsing
	Complete
)

type Order struct {
	Id      int64       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId  int64       `json:"user_id" gorm:"index"`
	Price   float64     `json:"price"`
	Status  OrderStatus `json:"status" gorm:"size:16"`
	Version string      `json:"version" gorm:"-"`
}

type OrderData struct {
	Id        int64          `json:"id" gorm:"primaryKey"`
	OrderId   int64          `json:"order_id"`
	ReserveId int64          `json:"project_id"`
	Reserve   ProjectReserve `json:"project" gorm:"foreignKey:ProjectId"`
}

func (Order) TableName() string {
	return "order"
}

type Pay struct {
}
