package project

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/model"
	"medicinal_share/tool/db/mysql"
	"time"
)

func GetByType(typ string, page int) {}

func GetAllProject(page int) []*entity.Project {
	res := make([]*entity.Project, 0)
	err := mysql.GetConnect().Model(&entity.Project{}).
		Limit(20).Offset((page - 1) * 20).Find(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
}

func GetProjectByIds(id []int64) []*entity.Project {
	res := make([]*entity.Project, 0)
	err := mysql.GetConnect().Model(&entity.Project{}).Where("id in ?", id).Find(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
}

func GetHospitalByNear(g1 float64, g2 float64, last int64, rg int) []*entity.Hospital {
	res := make([]*entity.Hospital, 0)
	err := mysql.GetConnect().Table("(?) as temp",
		mysql.GetConnect().Table("hospital").
			Select("*, (st_distance(point("+
				"longitude,latitude),"+
				"point(?, ?))*111195/1000"+
				") as distance", g1, g2).
			Where("id > ?", last).
			Order("distance"),
	).Joins("AvatarFile").
		Where("distance < ?", rg).
		Limit(20).
		Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Hospital)
}

func GetHospitalById(id int64) *entity.Hospital {
	res := &entity.Hospital{}
	err := mysql.GetConnect().
		Joins("AvatarFile").
		Joins("Background").
		Preload("Projects").Where("id = ?", id).Take(res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
}

func GetProjectByHospitalId(id int64, last int64) []*entity.Project {
	res := make([]*entity.Project, 0)
	err := mysql.GetConnect().Table("project").
		Where("hospital_id = ? and id > ?", id, last).
		Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Project)

}

func SearchProjectName(key string, hospital, last int64) []*entity.Project {
	res := make([]*entity.Project, 0)
	err := mysql.GetConnect().Table("project").Where("hospital_id= ? and id > ? and name like ?", hospital, last, key+"%").Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Project)
}

func CreateReserve(projectId, userId int64, time time.Time) {}

func UpdateReserve() {}
