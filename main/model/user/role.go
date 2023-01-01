package user

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
)

func CreateRole(userId int64, role string, tx *gorm.DB) *entity.UserRole {
	r := entity.UserRole{
		UserId: userId,
		Name:   role,
	}
	err := tx.Create(&r).Error
	if err != nil {
		panic(err)
	}
	return &r
}

func GetRolesById(userId int64) []*entity.UserRole {
	roles := make([]*entity.UserRole, 0)
	err := mysql.GetConnect().Model(&entity.UserRole{}).Where("user_id = ?", userId).Find(&roles).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
	}
	return roles
}
