package order

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool"
)

func CreateOrder(order *entity.Order, reserveId int64, tx *gorm.DB) {
	order.Id, _ = tool.GetId("order")
	order.Status = "pending"
	err := tx.Create(order).Error
	if err != nil {
		panic(err)
	}
	data := &entity.OrderData{
		OrderId:   order.Id,
		ReserveId: reserveId,
	}
	err = tx.Create(data).Error
	if err != nil {
		panic(err)
	}
}

func UpdateOrderStatus(id int64, status string, tx *gorm.DB) {
	err := tx.Model(&entity.Order{}).
		Where("id = ?", id).
		Update("status = ?", status).Error
	if err != nil {
		panic(err)
	}
}
