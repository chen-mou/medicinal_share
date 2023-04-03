package project

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/project"
)

func GetAllProject(ctx *gin.Context) {

}

func GetProjectByType(ctx *gin.Context) {}

func GetByHospitalId(ctx *gin.Context) {
	type Param struct {
		HospitalId int64 `form:"hospital_id" binding:"required"`
		Last       int64 `form:"last" binding:"required"`
	}
	param := Param{}
	err := ctx.BindQuery(&param)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
	}
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": project.GetProjectByHospitalId(param.HospitalId, param.Last),
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

func AddProjectReserve(ctx *gin.Context) {

}
