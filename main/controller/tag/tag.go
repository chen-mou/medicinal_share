package tag

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
	"medicinal_share/main/service/tag"
)

func GetTagByType(ctx *gin.Context) {
	type param struct {
		Type   string `json:"type"`
		Parent int64  `json:"parent"`
	}
	p := &param{}
	ctx.BindQuery(p)
	res := tag.GetByType(p.Type, p.Parent)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": res,
	})
}

func GetTagByKey(ctx *gin.Context) {
	type param struct {
		Type    string `form:"type"`
		Parent  int64  `form:"parent"`
		Key     string `form:"key"`
		PageNum int    `form:"page_num"`
	}
	p := &param{}
	err := ctx.BindQuery(p)
	if err != nil {
		panic(err)
	}
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": tag.Search(p.Key, p.Type, p.Parent, entity.CreatePage(p.PageNum, 20)),
	})
}
