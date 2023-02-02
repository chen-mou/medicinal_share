package user

import (
	"context"
	"errors"
	redis2 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"medicinal_share/gen/out/dao"
	"medicinal_share/main/entity"
	"medicinal_share/tool"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"medicinal_share/tool/encrypt/md5"
	"strconv"
	"time"
)

const (
	userKeyPrefix     = "USER:"
	userLockKeyPrefix = "USER:LOCK:"
)

func GetByName(username string) *entity.User {
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	res, err := dao.User.WithContext(ctx).GetByUserName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	res.Username = username
	return res
}

func getByIdFromDB(userId int64) *entity.User {
	var user entity.User
	err := mysql.GetConnect().Model(&entity.User{}).
		Joins("UserInfo").Preload("Role").
		Where("user.id = ?", userId).
		Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return &user
}

func getByIdFromCache(userId int64) *entity.User {
	user := &entity.User{}
	err := redis.Get(userKeyPrefix+strconv.FormatInt(userId, 10), user)
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return nil
		}
		panic(err)
	}
	return user
}

func GetById(userId int64) *entity.User {
	id := strconv.FormatInt(userId, 10)
	key := userKeyPrefix + id
	lock := userLockKeyPrefix + id
	c := redis.NewCache(lock, key)
	return c.Get(&entity.User{}, func() any {
		return getByIdFromDB(userId)
	}).(*entity.User)
}

func Create(name, password string, tx *gorm.DB) *entity.User {
	id, err := tool.GetId("user")
	if err != nil {
		panic(err)
	}
	user := entity.User{
		Id:       id,
		Username: name,
		Password: password,
	}
	err = tx.Create(&user).Error
	if err != nil {
		panic(err)
	}
	return &user
}

func CreateData(userId int64, tx *gorm.DB) *entity.UserData {
	data := entity.UserData{
		UserId:   userId,
		Nickname: "新用户" + md5.Hash(time.Now().String())[25:],
		Avatar:   1,
	}
	err := tx.Create(&data).Error
	if err != nil {
		panic(err)
	}
	return &data
}

func GetDataByUserId(userId int64) *entity.UserData {
	data := entity.UserData{}
	err := mysql.GetConnect().Where("user_id = ?", userId).Take(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return &data
}

func UpdateData(data *entity.UserData, tx *gorm.DB) *entity.UserData {
	err := tx.Model(&entity.UserData{}).Where("user_id = ?", data.UserId).Updates(*data).Error
	if err != nil {
		panic(err)
	}
	return data
}

func UpdatePassword(userId int64, password string, tx *gorm.DB) {
	err := tx.Model(&entity.User{}).Where("id = ?", userId).Update("password", password).Error
	if err != nil {
		panic(err)
	}
}
