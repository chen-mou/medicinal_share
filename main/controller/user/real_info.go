package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	user2 "medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/file"
	"medicinal_share/main/service/user"
	"medicinal_share/tool"
	"strconv"
)

func CreateInfo(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	info := user2.RealInfo{}
	ctx.BindJSON(&info)
	user.CreateInfo(usr.Id, &info)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
	})
}

func CreateDoctorInfo(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	info := user2.DoctorInfo{}
	ctx.BindJSON(info)
	user.CreateDoctorInfo(usr.Id, &info)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
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
		"code": 0,
		"data": user.GetDoctorInfoByUserId(id),
	})
}

func UploadDockerAvatar(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	f, err := ctx.FormFile("file")
	if err != nil {
		panic(err)
	}
	info := user.GetDoctorInfoByUserId(usr.Id)
	if info == nil {
		panic(middleware.NewCustomErr(middleware.FORBID, "错误操作"))
	}
	var id int64
	file.Upload(f, "doctor_avatar", usr.Id, func(i int64, db *gorm.DB) error {
		id = i
		return nil
	})
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": id,
	})
}

//func GetDoctors(ctx *gin.Context) {
//	p := ctx.Query("page")
//	page, _ := strconv.Atoi(p)
//
//}
