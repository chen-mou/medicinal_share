package project

import "github.com/gin-gonic/gin"

func Route(app *gin.Engine) {
	g := app.Group("/hospital")

	g.GET("/getNear", GetNearHospital)
}
