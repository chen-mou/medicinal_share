package project

import "github.com/gin-gonic/gin"

func Route(app *gin.Engine) {
	g := app.Group("/hospital")

	g.GET("/getNear", GetNearHospital).
		GET("/:id", GetByHospitalId) //TODO:这个方法逻辑要改成返回制定医院的每个种类前四项项目并添加一个接口用于获取剩下的项目用于展开
}
