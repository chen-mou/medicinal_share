package project

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/project"
)

func GetAllProject(ctx *gin.Context) {

}

func GetProjectByType(ctx *gin.Context) {}

func GetNearHospital(ctx *gin.Context) {
	type Param struct {
		Longitude float64 `json:"longitude" binding:"require"`
		Latitude  float64 `json:"latitude" binding:"require"`
		Range     int     `json:"range" binding:"require"`
		Last      int64   `json:"last" binding:"require"`
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
