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

func GetNearHospital(ctx *gin.Context) {
	type Param struct {
		Longitude float64 `form:"longitude" binding:"required"`
		Latitude  float64 `form:"latitude" binding:"required"`
		Range     int     `form:"range" binding:"required"`
		Last      int64   `form:"last"`
	}
	param := Param{}
	err := ctx.BindQuery(&param)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
	}
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": project.GetNearHospital(param.Longitude,
			param.Latitude,
			param.Last,
			param.Range),
	})
}
