package order

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool"
)

func CreateOrder(order *entity.Order, tx *gorm.DB) {
	order.Id, _ = tool.GetId("order")
	order.Status = "pending"
	err := tx.Create(order).Error
	if err != nil {
		panic(err)
	}
	datas := make([]*entity.OrderData, len(order.ProjectIds))
	for i, v := range order.ProjectIds {
		datas[i] = &entity.OrderData{
			OrderId:   order.Id,
			ProjectId: v,
		}
	}
	err = tx.Create(datas).Error
	if err != nil {
		panic(err)
	}
}
