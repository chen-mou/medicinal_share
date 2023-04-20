package project

import "github.com/gin-gonic/gin"

func Route(app *gin.Engine) {

	project := app.Group("project")

	project.GET("/getByHospitalId", GetByHospitalId)

	hospital := app.Group("/hospital")

	hospital.GET("/getNear", GetNearHospital).
		GET("/:id", GetHospitalById).
		GET("/getProjectReserveByDateAndProjectId", GetProjectReserveByDateAndProjectId)
}
