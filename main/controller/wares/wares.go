package wares

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
)

func Create(ctx *gin.Context) {
	value := &entity.Wares{}
	err := ctx.BindJSON(value)
	if err != nil {
		panic(err)
	}
	ctx.AbortWithStatusJSON(200, gin.H{})
}
