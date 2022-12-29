package user

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
)

func Route(app *gin.Engine) {
	g := app.Group("/user")

	g.POST("/login", Login).
		POST("/register", Register)
	g.Use(middleware.Verify).
		GET("/getInfo", GetUserData).
		POST("/uploadAvatar", UploadAvatar).
		POST("/updateInfo", UpdateInfo)
}
