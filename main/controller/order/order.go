package order

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/service/order"
	"medicinal_share/tool"
)

func CreateOrder(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	type param struct {
		Version   string `json:"version"`
		ReverseId int64  `json:"reverse_id"`
	}
	p := &param{}
	ctx.BindJSON(p)
	order.Create(p.Version, p.ReverseId, usr.Id)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
	})
}

// GetOrder 获取用户的订单
func GetOrder(ctx *gin.Context) {
	param := struct {
		Type string `form:"type"`
		Last int64  `form:"last"`
	}{}
	ctx.BindJSON(&param)
	usr := tool.GetNowUser(ctx)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": order.GetUserOrder(usr.Id, param.Last, param.Type),
	})
}

// Pay 支付
func Pay(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	type param struct {
		OrderId int64 `json:"order_id"`
	}
	p := &param{}
	ctx.BindJSON(p)
	order.Pay(p.OrderId, usr.Id)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
	})
}
