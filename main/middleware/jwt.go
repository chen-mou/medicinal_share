package middleware

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/model/user"
	"medicinal_share/tool"
	"medicinal_share/tool/encrypt/jwtutil"
	"strconv"
)

type Role uint8

const (
	User Role = iota
	Doctor
	Worker
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

func RoleVerify(role Role) func(ctx *gin.Context) {
	r := ""
	switch role {
	case User:
		r = "Custom"
	case Doctor:
		r = "Doctor"
	case Worker:
		r = "Worker"
	default:
		panic("权限不存在")
	}
	return func(ctx *gin.Context) {
		u := tool.GetNowUser(ctx)
		if u == nil {
			panic("使用这个中间件前要使用Verify")
		}
		for _, role := range u.Role {
			if role.Name == r {
				ctx.Next()
				return
			}
		}
		panic(NewCustomErr(FORBID, "没有相应权限"))
	}
}
