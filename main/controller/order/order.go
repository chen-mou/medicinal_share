package order

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
	"medicinal_share/main/service/order"
	"medicinal_share/tool"
)

func CreateOrder(ctx *gin.Context) {
	o := &entity.Order{}
	usr := tool.GetNowUser(ctx)
	ctx.BindJSON(o)
	order.Create(o, usr.Id)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
	})
}
