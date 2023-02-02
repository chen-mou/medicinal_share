package order

import (
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/order"
	"medicinal_share/main/model/project"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
)

const AntiShake = "OrderAntiShake:"

func Create(ord *entity.Order, userId int64) {
	if !redis.AntiShake(AntiShake + ord.Version) {
		panic(middleware.NewCustomErr(middleware.ERROR, "操作过于频繁"))
	}
	pros := project.GetProjectByIds(ord.ProjectIds)
	if pros == nil || len(ord.ProjectIds) != len(pros) {
		panic(middleware.NewCustomErr(middleware.NOT_FOUND, "所选的项目不存在"))
	}
	price := float64(0)
	for _, v := range pros {
		price += v.Price
	}
	ord.Price = price
	ord.UserId = userId
	order.CreateOrder(ord, mysql.GetConnect())
}
