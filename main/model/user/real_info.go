package user

import (
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

func getDoctorInfoByIdFormDB(userId int64) *entity.DoctorInfo {
	info := &entity.DoctorInfo{}
	err := mysql.GetConnect().
		Model(info).
		Where("user_id = ?", userId).
		Joins("Info").
		Preload("Tags").
		Preload("Tags.Tag").
		Preload("Avatar").
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
	lock := DoctorLockKey + id
	c := redis.NewCache(lock, key)
	return c.Get(&entity.DoctorInfo{}, func() any {
		return getDoctorInfoByIdFormDB(userId)
	}).(*entity.DoctorInfo)
}

func GetDoctorById(id int) *entity.DoctorInfo {
	res := &entity.DoctorInfo{}
	db := mysql.GetConnect()
	err := db.Where(&entity.DoctorInfo{Id: id}).
		Joins("Info").
		Preload("Tags", db.Where("tag_ref.relation_type = 'area'")).
		Preload("Tags.Tag").Take(res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
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

func GetDoctors(page int) []*entity.DoctorInfo {
	info := make([]*entity.DoctorInfo, 0)
	mysql.GetConnect().Model(&entity.DoctorInfo{}).
		Joins("Info").Preload("TagsId").
		Limit(20).Offset((page - 1) * 20).Find(&info)
	return info
}
