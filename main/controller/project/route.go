package project

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/middleware"
)

func Route(app *gin.Engine) {

	project := app.Group("project")

	project.GET("/getByHospitalId", GetByHospitalId)

	reserve := app.Group("/reserve")

	reserve.GET("/getProjectReserveByDateAndProjectId", GetProjectReserveByDateAndProjectId)

	hospital := app.Group("/hospital")

	hospital.GET("/getNear", GetNearHospital).
		GET("/:id", GetHospitalById).
		Use(middleware.Verify).
		Use(middleware.RoleVerify(middleware.Worker)).
		POST("/updateAvatar", UpdateHospitalAvatar).
		POST("/updateBackground", UpdateHospitalBackground).
		POST("/addProjectReserve", AddProjectReserve)
}
