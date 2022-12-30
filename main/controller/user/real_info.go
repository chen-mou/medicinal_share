package user

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
	"medicinal_share/main/service/user"
	"medicinal_share/tool"
)

func CreateInfo(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	info := entity.RealInfo{}
	ctx.BindJSON(&info)
	user.CreateInfo(usr.Id, &info)
	ctx.AbortWithStatusJSON(200, gin.H{
		"data": 0,
	})
}

func CreateDoctorInfo(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	info := entity.DoctorInfo{}
	ctx.BindJSON(info)
	user.CreateDoctorInfo(usr.Id, &info)
	ctx.AbortWithStatusJSON(200, gin.H{
		"data": 0,
	})
}

func GetDoctorInfo(ctx *gin.Context) {}
