package wares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"medicinal_share/main/service/file"
	"medicinal_share/tool"
)

func UploadPhoto(ctx *gin.Context) {
	user := tool.GetNowUser(ctx)
	f, err := ctx.FormFile("file")
	if err != nil {
		panic(err)
	}
	var id int64
	file.Upload(f, "wares", user.Id, func(i int64, db *gorm.DB) error {
		id = i
		return nil
	})
	ctx.AbortWithStatusJSON(200, gin.H{
		"data": id,
	})
}

//func GetPhotos(ctx *gin.Context) {
//	user := tool.GetNowUser(ctx)
//	typ := ctx.Query("type")
//	ctx.AbortWithStatusJSON(200, gin.H{})
//}
