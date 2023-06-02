package report

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
)

func Route(app *gin.Engine) {
	g := app.Group("/report")

	g.Use(middleware.Verify).
		GET("/myUnShareReport", MyUnShareReport).
		POST("/shareReport", ShareReport).
		GET("/userGetReport", GetReportByReserveId).
		Use(middleware.RoleVerify(middleware.Doctor)).
		GET("/doctorGetReport", GetReportById).
		GET("/doctorReport", DoctorReport).
		GET("/getShareReport", GetShareReport).
		POST("/uploadReportImage", UploadReportImage).
		POST("/create", CreateReport)

}
