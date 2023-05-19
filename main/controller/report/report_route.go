package report

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
)

func Route(app *gin.Engine) {
	g := app.Group("/report")

	g.Use(middleware.Verify).
		GET("/myReport", MyReport).
		Use(middleware.RoleVerify(middleware.Doctor)).
		GET("/doctorReport", DoctorReport).
		POST("/uploadReportImage", UploadReportImage).
		POST("/create", CreateReport)

}
