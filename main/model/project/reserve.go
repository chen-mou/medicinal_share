package project

import (
	"errors"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/model"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"sort"
	"strconv"
)

const (
	ReserveGetLock = "RESERVE_GET_LOCK"
	ReserveGet     = "RESERVE_GET"
)

func LoadReserveById(reserveId int64) error {
	res := &entity.ProjectReserve{}
	rstr := strconv.FormatInt(reserveId, 10)
	c := redis.NewCache(ReserveGetLock+":"+rstr, ReserveGet+":"+rstr)
	_, err := c.LoadInt(func() (int, error) {
		err := mysql.GetConnect().Select("overplus").Where("id = ?", reserveId).Take(res).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
		if err != nil {
			panic(err)
		}
		return res.Overplus, nil
	})

	return err
}

func CreateProjectReserve(reserve *entity.ProjectReserve) {
	reserve.ReserveNum = 0
	err := mysql.GetConnect().Create(reserve).Error
	if err != nil {
		panic(err)
	}
}

func GetProjectReserveById(id int64, tx *gorm.DB) *entity.ProjectReserve {
	res := &entity.ProjectReserve{}
	err := tx.Joins("Project").Joins("DoctorInfo").Where("project_reserve.id = ?", id).Take(res).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return res
}

func GetProjectReserveByDateAndProjectId(start, end entity.Time, projectId int64) []*entity.ProjectReserve {
	res := make([]*entity.ProjectReserve, 0)
	err := mysql.GetConnect().
		Joins("Project").
		Joins("DoctorInfo").
		Preload("DoctorInfo.Info").
		Where("project_id = ? "+
			"and start between ? and ? "+
			"and overplus > 0", projectId, start, end).Find(&res).Error
	sort.Slice(res, func(i, j int) bool {
		return res[i].End.Time().After(res[j].End.Time())
	})
	if err == nil || err == gorm.ErrRecordNotFound {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return res
	}
	panic(err)
}

func CreateReserve(reserve []*entity.Reserve, tx *gorm.DB) {
	err := tx.Create(&reserve).Error
	if err != nil {
		panic(err)
	}
}

func UpdateProjectReserveNum(id int64, tx *gorm.DB) {
	err := tx.Model(&entity.ProjectReserve{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"overplus":    gorm.Expr("overplus - 1"),
			"reserve_num": gorm.Expr("reserve_num + 1"),
		}).Error
	if err != nil {
		panic(err)
	}
}

func GetReserveByUserId(userId int64) []*entity.Reserve {
	res := make([]*entity.Reserve, 0)
	err := mysql.GetConnect().Model(&entity.Reserve{}).Where(&entity.Reserve{
		UserId: userId,
	}).
		Joins("ProjectReserve").
		Preload("ProjectReserve.Project").
		Preload("ProjectReserve.DoctorInfo").
		Preload("ProjectReserve.DoctorInfo.Avatar").Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Reserve)
}

func GetReserveById(id int64) *entity.Reserve {
	res := &entity.Reserve{}
	err := mysql.GetConnect().
		Where(&entity.Reserve{Id: id}).
		Joins("ProjectReserve").
		Preload("ProjectReserve.Project").Take(res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
	}
	return res
}

func IsUserReserve(reserveId, userId int64) bool {
	err := mysql.GetConnect().Where("id = ? and user_id = ?", reserveId, userId).Take(&entity.Reserve{}).Error
	return model.GetErrorHandler(err, entity.Reserve{}) != nil
}

func IsDoctorReserve(doctorId int, reserveId int64) bool {
	err := mysql.GetConnect().Joins("ProjectReserve").
		Where("reserve.id = ? and ProjectReserve.doctor_id = ?", reserveId, doctorId).
		Take(&entity.Reserve{}).Error
	return model.GetErrorHandler(err, entity.Reserve{}) != nil
}

func UpdateReserveStatus(reserveId int64, status entity.ReserveStatue, tx *gorm.DB) {
	if err := tx.Model(&entity.Reserve{}).
		Where("id = ?", reserveId).
		Update("status", status).Error; err != nil {
		panic(err)
	}
}

func CreateShareReport(reportsId []int64, reserveId int64) {
	sr := make([]*entity.ShareReport, 0)
	for _, reportId := range reportsId {
		sr = append(sr, &entity.ShareReport{
			ReportId:  reportId,
			ReserveId: reserveId,
		})
	}
	if err := mysql.GetConnect().Model(&entity.ShareReport{}).Create(&sr).Error; err != nil {
		panic(err)
	}
}

func GetShareReportByReserveId(reserveId int64) []*entity.ShareReport {
	res := make([]*entity.ShareReport, 0)
	err := mysql.GetConnect().Model(&entity.ShareReport{}).
		Where("reserve_id = ?", reserveId).
		Joins("Report").
		Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.ShareReport)
}
