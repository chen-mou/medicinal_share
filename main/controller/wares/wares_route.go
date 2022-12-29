package wares

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
)

func Route(app *gin.Engine) {
	app.Group("/ware").
		Use(middleware.Verify).
		POST("/uploadPhoto", UploadPhoto)
}
