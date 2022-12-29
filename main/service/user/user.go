package user

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/user"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/encrypt/md5"
)

func Register(username, password string) *entity.User {
	usr := user.GetByName(username)
	if usr != nil {
		panic(middleware.NewCustomErr(middleware.NO_REPETITION, "用户名已存在"))
	}
	password = md5.Hash(password)
	err := mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		usr = user.Create(username, password, tx)
		user.CreateData(usr.Id, tx)
		return nil
	})
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
	}
	return usr
}

func Login(username, password string) *entity.User {
	usr := user.GetByName(username)
	if usr == nil {
		panic(middleware.NewCustomErr(middleware.NOT_FOUND, "用户名不存在"))
	}
	password = md5.Hash(password)
	if password != usr.Password {
		panic(middleware.NewCustomErr(middleware.ERROR, "密码有误"))
	}
	usr.UserInfo = user.GetDataByUserId(usr.Id)
	return usr
}

func GetData(id int64) *entity.UserData {
	return user.GetDataByUserId(id)
}
