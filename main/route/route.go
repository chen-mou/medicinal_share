package route

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/controller/file"
	"medicinal_share/main/controller/order"
	"medicinal_share/main/controller/project"
	"medicinal_share/main/controller/tag"
	"medicinal_share/main/controller/treat"
	"medicinal_share/main/controller/user"
	"medicinal_share/main/controller/wares"
	"medicinal_share/main/middleware"
)

type mode int

const (
	Debug mode = iota
	Release
	Deploy
)

func Route(mod int) *gin.Engine {
	var app *gin.Engine
	switch mode(mod) {
	case Release:
		gin.SetMode(gin.ReleaseMode)
	case Deploy:
		gin.SetMode(gin.TestMode)
	}
	app = gin.Default()

	app.Use(middleware.Catch)

	app.Use(middleware.Cross)

	app.OPTIONS("/*path", func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.AbortWithStatus(200)
	})

	user.Route(app)

	file.Route(app)

	wares.Route(app)

	order.Route(app)

	project.Route(app)

	tag.Route(app)

	treat.Route(app)

	return app
}
