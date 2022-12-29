package route

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/controller/file"
	"medicinal_share/main/controller/user"
	"medicinal_share/main/controller/wares"
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
	user.Route(app)

	file.Route(app)

	wares.Route(app)

	return app
}
