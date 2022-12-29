package file

import (
	"github.com/gin-gonic/gin"
)

func Route(app *gin.Engine) {
	g := app.Group("/file")

	g.GET("/get/:hash", GetByHash)
}
