package project

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/project"
	"strconv"
)

func GetHospitalById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		panic(middleware.NewCustomErr(middleware.ERROR, "缺少参数ID"))
	}
	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
	}
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": project.GetHospitalById(iid),
	})
}
