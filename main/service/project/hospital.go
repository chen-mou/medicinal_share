package project

import (
	"medicinal_share/main/entity"
	"medicinal_share/main/model/project"
)

func GetHospitalById(id int64) *entity.Hospital {
	return project.GetHospitalById(id)
}

func GetNearHospital(g1 float64, g2 float64, last int64, rge int) []*entity.Hospital {
	return project.GetHospitalByNear(g1, g2, last, rge)
}
