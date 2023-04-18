package order

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/model"
	"medicinal_share/tool"
	"medicinal_share/tool/db/mysql"
)

func CreateOrder(order *entity.Order, reserveId int64, tx *gorm.DB) {
	order.Id, _ = tool.GetId("order")
	order.Status = entity.Padding
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

func UpdateOrderStatus(id int64, status entity.OrderStatus, tx *gorm.DB) {
	err := tx.Model(&entity.Order{}).
		Where("id = ?", id).
		Update("status = ?", status).Error
	if err != nil {
		panic(err)
	}
}

func ExistOrder(orderId, userId int64, tx *gorm.DB) bool {
	err := tx.Model(&entity.Order{}).Select("id").Where("id = ? and user_id = ?", orderId, userId).Take(&entity.Order{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		panic(err)
	}
	return true
}

func GetUserOrder(userId, last int64, typ string) []*entity.Order {
	param := []any{userId, last}
	sql := "user_id = ? and id > last"
	if typ != "" {
		sql += "and status = ?"
		param = append(param, typ)
	}
	orders := make([]*entity.Order, 0)
	err := mysql.GetConnect().Model(&entity.Order{}).Joins("OrderData").Where(sql, param...).Find(&orders).Error
	return model.GetErrorHandler(err, orders).([]*entity.Order)
}
