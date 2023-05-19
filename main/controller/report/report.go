package report

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/entity/report"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/file"
	service "medicinal_share/main/service/report"
	"medicinal_share/tool"
	"strings"
)

func UploadReportImage(ctx *gin.Context) {
	f, err := ctx.FormFile("file")
	user := tool.GetNowUser(ctx)
	if err != nil {
		panic(err)
	}
	if f.Size >= 20<<20 {
		panic(middleware.NewCustomErr(middleware.ERROR, "文件过大"))
	}
	suffix := strings.Split(f.Filename, ".")[1]
	if _, ok := file.Suffix[suffix]; !ok {
		panic(middleware.NewCustomErr(middleware.ERROR, "文件类型有误"))
	}
	file.Upload(f, "report_image", user.Id, func(i int64, db *gorm.DB) error {
		res := &entity.FileData{
			Id: i,
		}
		if err := db.Where(res).Joins("File").Take(&res).Error; err != nil {
			panic(err)
		}
		ctx.AbortWithStatusJSON(200, gin.H{
			"code": 0,
			"data": res,
		})
		return nil
	})
}

func CreateReport(ctx *gin.Context) {
	r := &report.Report{}
	err := ctx.ShouldBindJSON(r)
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
	}
	service.UploadReport(r)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
	})
}

func MyReport(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": service.GetUserReport(usr.Id),
	})
}

func DoctorReport(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 200,
		"data": service.GetDoctorReport(usr.Id),
	})
}
