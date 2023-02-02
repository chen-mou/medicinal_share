package order

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
)

func Route(app *gin.Engine) {
	g := app.Group("order")

	g.Use(middleware.Verify).POST("/order", CreateOrder)

}
