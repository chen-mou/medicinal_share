package project

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/project"
	"medicinal_share/tool"
	"time"
)

func GetAllProject(ctx *gin.Context) {

}

func GetProjectByType(ctx *gin.Context) {}

func GetByHospitalId(ctx *gin.Context) {
	type Param struct {
		HospitalId int64 `form:"hospital_id" binding:"required"`
		Last       int64 `form:"last"`
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
	reserve := &entity.ProjectReserve{}
	ctx.BindJSON(reserve)
	u := tool.GetNowUser(ctx)
	project.CreateProjectReserve(reserve, u.Id)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
	})
}

// Add TODO:添加项目
func Add(ctx *gin.Context) {

}

func GetProjectReserveByDateAndProjectId(ctx *gin.Context) {
	p := &struct {
		Time      string `form:"time"`
		ProjectId int64  `form:"project_id"`
	}{}
	ctx.BindQuery(p)
	time, err := time.Parse("2006-01-02 15:04:05", p.Time)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "时间格式错误"))
	}
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": project.GetProjectReserveByDataAndProject(time, p.ProjectId),
	})
}
