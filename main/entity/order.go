package entity

//TODO:完成订单以及交易流水的模型定义

type Order struct {
	Id         int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId     int64   `json:"user_id" gorm:"index"`
	Price      float64 `json:"price"`
	Status     string  `json:"status" gorm:"size:16"`
	Version    string  `json:"version" gorm:"-"`
	ProjectIds []int64 `json:"project_ids" gorm:"-"`
}

type OrderData struct {
	Id        int64   `json:"id" gorm:"primaryKey"`
	OrderId   int64   `json:"order_id"`
	ProjectId int64   `json:"project_id"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectId"`
}

func (Order) TableName() string {
	return "order"
}

type Pay struct {
}
