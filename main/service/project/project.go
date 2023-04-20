package project

import (
	"medicinal_share/main/entity"
	"medicinal_share/main/model/project"
	"strconv"
	"time"
)

func GetProjectByHospitalId(id int64, last int64) []*entity.Project {
	return project.GetProjectByHospitalId(id, last)
}

func GetProjectReserveByDataAndProject(t time.Time, projectId int64) []*entity.ProjectReserve {
	today := strconv.Itoa(t.Day())
	if t.Day() < 10 {
		today = "0" + today
	}
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	return project.GetProjectReserveByDateAndProjectId(start, end, projectId)
}
