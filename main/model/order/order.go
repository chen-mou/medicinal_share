package order

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/model"
	"medicinal_share/tool"
	"medicinal_share/tool/db/mysql"
)

func CreateOrder(o *entity.Order, reserveId int64, tx *gorm.DB) {
	o.Id, _ = tool.GetId("order")
	o.Status = entity.Padding
	o.CreateAt = entity.Now()
	err := tx.Create(o).Error
	if err != nil {
		panic(err)
	}
	data := &entity.OrderData{
		OrderId:          o.Id,
		ProjectReserveId: reserveId,
	}
	err = tx.Create(data).Error
	if err != nil {
		panic(err)
	}
}

func UpdateOrderStatus(id int64, status entity.Status, tx *gorm.DB) {
	err := tx.Model(&entity.Order{}).
		Where("id = ?", id).
		Update("status = ?", status).Error
	if err != nil {
		panic(err)
	}
}

func Get(orderId, userId int64, tx *gorm.DB) *entity.Order {
	o := &entity.Order{}
	err := tx.Model(&entity.Order{}).Select("id").
		Where("id = ? and user_id = ?", orderId, userId).Take(o).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return o
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
