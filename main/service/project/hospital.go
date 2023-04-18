package project

import (
	"medicinal_share/main/entity"
	"medicinal_share/main/model/project"
)

func GetHospitalById(id int64) *entity.Hospital {
	return project.GetHospitalById(id)
}
