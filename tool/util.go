package tool

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
	"unsafe"
)

func GetNowUser(ctx *gin.Context) *entity.User {
	u, _ := ctx.Get("CurrentUser")
	return u.(*entity.User)
}

//BytesToString 不拷贝数据将字符数组转换为字符串
func BytesToString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}
