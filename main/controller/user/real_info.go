package user

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/user"
	"medicinal_share/tool"
	"strconv"
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

func GetDoctorInfo(ctx *gin.Context) {
	idstr, ok := ctx.GetQuery("id")
	if !ok {
		panic(middleware.NewCustomErr(middleware.ERROR, "参数id不存在"))
	}
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "参数类型有误"))
	}
	ctx.AbortWithStatusJSON(200, gin.H{
		"data": user.GetDoctorInfoByUserId(id),
	})
}
