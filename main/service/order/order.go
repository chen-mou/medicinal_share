package order

import (
	"context"
	redis1 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/order"
	"medicinal_share/main/model/project"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"strconv"
	"time"
)

const AntiShake = "OrderAntiShake:"

func Create(version string, reverseId int64, userId int64) {
	if !redis.AntiShake(AntiShake + version) {
		panic(middleware.NewCustomErr(middleware.ERROR, "操作过于频繁"))
	}
	err := project.LoadReserveById(reverseId)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "目标预约不存在"))
	}
	rstr := strconv.FormatInt(reverseId, 10)
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	n, err := redis.DB.Decr(ctx, project.ReserveGet+":"+rstr).Result()
	if err != nil {
		panic(err)
	}
	if n < 0 {
		panic(middleware.NewCustomErr(middleware.ERROR, "预约满了"))
	}
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		p := project.GetProjectReserveById(reverseId, tx)
		ord := &entity.Order{
			UserId: userId,
			Price:  p.Project.Price,
		}
		order.CreateOrder(ord, reverseId, mysql.GetConnect())
		project.UpdateProjectReserveNum(reverseId, tx)
		redis.DB.Set(context.TODO(), "Order-"+strconv.FormatInt(ord.Id, 64), "", 30*time.Minute)
		return nil
	})
}

// Pay 支付
func Pay(orderId int64, userId int64) {
	_, err := redis.DB.Get(context.TODO(), "Order-"+strconv.FormatInt(orderId, 64)).Result()
	if err == redis1.Nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "订单已过期或不存在"))
	}
	if err != nil {
		panic(err)
	}
	err = mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		if !order.ExistOrder(orderId, userId, tx) {
			return middleware.NewCustomErr(middleware.ERROR, "订单不属于你或不存在")
		}
		order.UpdateOrderStatus(orderId, entity.UnUsing, tx)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func GetUserOrder(userId, last int64, status string) []*entity.Order {
	return order.GetUserOrder(userId, last, status)
}
