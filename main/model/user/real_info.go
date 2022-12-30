package user

import (
	redis2 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"strconv"
)

const (
	DoctorKey     = "DoctorInfo:UserId:"
	DoctorLockKey = "DoctorInfo:UserId:Lock:"
)

func CreateInfo(info *entity.RealInfo, tx *gorm.DB) *entity.RealInfo {
	err := tx.Create(info).Error
	if err != nil {
		panic(err)
	}
	return info
}

func CreateDoctorInfo(info *entity.DoctorInfo, tx *gorm.DB) *entity.DoctorInfo {
	err := tx.Create(info).Error
	if err != nil {
		panic(err)
	}
	return info
}

func getDoctorByIdFormCache(key string) *entity.DoctorInfo {
	info := &entity.DoctorInfo{}
	err := redis.Get(key, info)
	if err != nil {
		if err == redis2.Nil {
			return nil
		}
		panic(err)
	}
	return info
}

func getDoctorInfoByIdFormDB(userId int64) *entity.DoctorInfo {
	info := &entity.DoctorInfo{}
	err := mysql.GetConnect().
		Model(info).
		Where("user_id = ?", userId).
		Joins("Info").Preload("Tags").
		Take(info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return info
}

func GetDoctorInfoById(userId int64) *entity.DoctorInfo {
	id := strconv.FormatInt(userId, 10)
	key := DoctorKey + id
	return redis.SafeGet(key, DoctorLockKey+id, func() any {
		return getDoctorByIdFormCache(key)
	}, func() any {
		return getByIdFromDB(userId)
	}).(*entity.DoctorInfo)
}

func GetInfoByNameAndIdNumber(name, idNumber string) *entity.RealInfo {
	info := &entity.RealInfo{}
	err := mysql.GetConnect().Where("name = ? and id_number = ?", name, idNumber).Take(info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return info
}
