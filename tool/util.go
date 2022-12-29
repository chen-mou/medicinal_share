package tool

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
)

func GetNowUser(ctx *gin.Context) *entity.User {
	u, _ := ctx.Get("CurrentUser")
	return u.(*entity.User)
}
