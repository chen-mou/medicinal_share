package user

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
)

func Route(app *gin.Engine) {
	g := app.Group("/user")

	g.POST("/login", Login).
		POST("/register", Register).
		GET("/getDockerInfo", GetDoctorInfo)

	g.Use(middleware.Verify).
		GET("/getInfo", GetUserData).
		POST("/uploadAvatar", UploadAvatar).
		POST("/updateInfo", UpdateInfo).
		POST("/realName", CreateInfo).
		POST("/dockerInfo", CreateDoctorInfo).
		POST("/uploadDoctorAvatar", UploadDockerAvatar)
}
