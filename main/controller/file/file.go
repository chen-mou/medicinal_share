package file

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/service/file"
)

func GetByHash(ctx *gin.Context) {
	hash := ctx.Param("hash")
	fe := file.GetByHash(hash)
	//添加这个就是下载
	//ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", hash))
	ctx.Writer.Header().Add("Content-Type", "image/png")
	ctx.File(fe.Path)
	ctx.Abort()
}
