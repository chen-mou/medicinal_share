package project

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/project"
	"time"
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

// AddProjectReserve 添加可以预约的时间
func AddProjectReserve(ctx *gin.Context) {

}

// Add TODO:添加项目
func Add(ctx *gin.Context) {

}

func GetProjectReserveByDateAndProjectId(ctx *gin.Context) {
	p := &struct {
		Time      time.Time `form:"time"`
		ProjectId int64     `form:"project_id"`
	}{}
	ctx.BindQuery(p)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": project.GetProjectReserveByDataAndProject(p.Time, p.ProjectId),
	})
}
