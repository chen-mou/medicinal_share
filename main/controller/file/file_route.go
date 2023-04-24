package file

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/resource"
)

func Route(app *gin.Engine) {
	g := app.Group("/" + resource.Machine + "/file")

	g.GET("/get/:hash", GetByHash)
}
