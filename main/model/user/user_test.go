package user

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/encrypt/md5"
	"testing"
	"time"
)

func TestCreateWorker(t *testing.T) {
	usrs := make([]*entity.User, 0)
	mysql.GetConnect().Model(&entity.User{}).Order("id desc").Limit(5).Find(&usrs)
	for i, user := range usrs {
		mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
			CreateWorker(tx, int64(i), user.Id)
			tx.Create(&entity.UserRole{
				UserId: user.Id,
				Name:   "Worker",
			})
			return nil
		})
	}
}

func TestCreate(t *testing.T) {
	for i := 0; i < 10; i++ {
		Create(md5.Hash(time.Now().String()), md5.Hash(time.Now().String()+"password"), mysql.GetConnect())
	}
}

func TestCreate2(t *testing.T) {
	userId := int64(48)
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		CreateWorker(tx, int64(1), userId)
		tx.Create(&entity.UserRole{
			UserId: userId,
			Name:   "Worker",
		})
		return nil
	})
}
