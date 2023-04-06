package treat

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
	"medicinal_share/tool/socket"
)

func Route(app *gin.Engine) {
	app.Group("/treat").
		Use(middleware.Verify).
		POST("/treat", Treat)
}

func Websocket(manager *socket.ConnManager) {
	manager.Message("/send", Send)
}
