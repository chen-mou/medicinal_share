package middleware

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/model/user"
	"medicinal_share/tool/encrypt/jwtutil"
	"strconv"
)

func Verify(ctx *gin.Context) {
	token := ctx.GetHeader("x-token")
	if token == "" {
		panic(NewCustomErr(ERROR, "缺少token"))
	}
	data, err := jwtutil.Parse(token)
	if err != nil {
		panic(NewCustomErr(ERROR, err.Error()))
	}
	id, _ := strconv.ParseInt(data["id"], 10, 64)
	u := user.GetById(id)
	ctx.Set("CurrentUser", u)
	ctx.Next()
}
