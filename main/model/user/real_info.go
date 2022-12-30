package user

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
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

func GetDoctorInfoById(userId int64) *entity.DoctorInfo {
	info := &entity.DoctorInfo{}
	err := mysql.GetConnect().Where("user_id = ?", userId).Take(info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return info
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
