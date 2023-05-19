package project

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/file"
	"medicinal_share/main/service/project"
	"medicinal_share/main/service/user"
	"medicinal_share/tool"
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

func SearchHospital() {}

func UpdateHospitalAvatar(ctx *gin.Context) {
	updateHospitalPhoto(ctx, "avatar")
}

func UpdateHospitalBackground(ctx *gin.Context) {
	updateHospitalPhoto(ctx, "background")
}

func updateHospitalPhoto(ctx *gin.Context, typ string) {
	f, err := ctx.FormFile("file")
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
	}
	u := tool.GetNowUser(ctx)
	hospitalId, ok := ctx.GetPostForm("id")
	if !ok {
		panic(middleware.NewCustomErr(middleware.ERROR, "参数id不存在"))
	}
	id, err := strconv.ParseInt(hospitalId, 10, 64)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "id不为数字"))
	}
	if !user.IsHospitalWorker(u.Id, id) {
		panic(middleware.NewCustomErr(middleware.ERROR, "你不是这个医院的员工"))
	}
	file.Upload(f, "hospital_"+typ, u.Id, func(i int64, db *gorm.DB) error {
		return db.Model(&entity.Hospital{}).Where("id = ?", hospitalId).Update(typ, i).Error
	})
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
	})
}
