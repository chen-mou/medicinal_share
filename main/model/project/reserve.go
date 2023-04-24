package project

import (
	"errors"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
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
		err := mysql.GetConnect().Select("num").Where("id = ?", reserveId).Take(res).Error
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
	err := tx.Joins("Project").Joins("DoctorInfo").Where("id = ?", id).Take(res).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return res
}

func GetProjectReserveByDateAndProjectId(start, end entity.Time, projectId int64) []*entity.ProjectReserve {
	res := make([]*entity.ProjectReserve, 0)
	err := mysql.GetConnect().
		Joins("Project").
		Joins("DoctorInfo").
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
	tx.Model(&entity.ProjectReserve{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"num":         gorm.Expr("num - 1"),
			"reserve_num": gorm.Expr("reserve_num + 1"),
		})
}
