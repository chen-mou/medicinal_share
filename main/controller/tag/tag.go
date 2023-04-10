package tag

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/tag"
)

func GetTagByType(ctx *gin.Context) {
	type param struct {
		Type   string `json:"type"`
		Parent int64  `json:"parent"`
	}
	p := &param{}
	ctx.BindQuery(p)
	res := tag.GetTagByType(p.Type, p.Parent)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": res,
	})
}

func GetTagByKey(ctx *gin.Context) {
	key := ctx.Query("key")
	if key == "" {
		panic(middleware.NewCustomErr(middleware.ERROR, "缺少参数key"))
	}
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": tag.SearchByKeyWord(key),
	})
}
