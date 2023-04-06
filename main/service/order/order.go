package order

import (
	"context"
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
	n, err := redis.DB.Incr(ctx, project.ReserveGet+":"+rstr).Result()
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
		return nil
	})
}

// Pay TODO:付款成功创建预约，付款检查redis是否存在目标订单，不存在则支付超时，存在则支付并更新数据库
func Pay(orderId int64) {

}
