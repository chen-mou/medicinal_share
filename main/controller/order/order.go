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

// GetOrder TODO:获取用户的订单
func GetOrder(ctx *gin.Context) {}

// Pay TODO:支付
func Pay(ctx *gin.Context) {}
