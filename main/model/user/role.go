package user

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
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
