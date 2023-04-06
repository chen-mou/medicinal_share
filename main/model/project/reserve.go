package project

import (
	"errors"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
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
			return 0, nil
		}
		if err != nil {
			panic(err)
		}
		return res.Num, nil
	})

	return err
}

func GetReserveById(id int64, tx *gorm.DB) *entity.ProjectReserve {
	res := &entity.ProjectReserve{}
	err := tx.Where("id = ?", id).Association("Project").Find(res)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return res
}

func UpdateProjectReserveNum(id int64, tx *gorm.DB) {
	tx.Model(&entity.ProjectReserve{}).Where("id = ?", id).Update("num", gorm.Expr("num + 1"))
}
