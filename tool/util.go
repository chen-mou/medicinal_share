package tool

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
	"unsafe"
)

func GetNowUser(ctx *gin.Context) *entity.User {
	u, exist := ctx.Get("CurrentUser")
	if !exist {
		return nil
	}
	return u.(*entity.User)
}

// BytesToString 不拷贝数据将字符数组转换为字符串
func BytesToString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

func StringToBytes(s string) []byte {
	return *((*[]byte)(unsafe.Pointer(&s)))
}
