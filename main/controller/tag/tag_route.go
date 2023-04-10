package tag

import "github.com/gin-gonic/gin"

func Route(app *gin.Engine) {
	app.Group("/tag").
		GET("/getByType", GetTagByType).
		GET("/searchByKey", GetTagByKey)
}
