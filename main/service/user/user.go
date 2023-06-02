package user

import (
	"gorm.io/gorm"
	user2 "medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/user"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/encrypt/md5"
)

func Register(username, password string) *user2.User {
	usr := user.GetByName(username)
	if usr != nil {
		panic(middleware.NewCustomErr(middleware.NO_REPETITION, "用户名已存在"))
	}
	password = md5.Hash(password)
	err := mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		usr = user.Create(username, password, tx)
		usr.UserInfo = user.CreateData(usr.Id, tx)
		usr.Role = make([]*user2.UserRole, 0)
		usr.Role = append(usr.Role, user.CreateRole(usr.Id, "Custom", tx))
		return nil
	})
	if err != nil {
		panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
	}
	return usr
}

func Login(username, password string) *user2.User {
	usr := user.GetByName(username)
	if usr == nil {
		panic(middleware.NewCustomErr(middleware.NOT_FOUND, "用户名不存在"))
	}
	password = md5.Hash(password)
	if password != usr.Password {
		panic(middleware.NewCustomErr(middleware.ERROR, "密码有误"))
	}
	usr.UserInfo = user.GetDataByUserId(usr.Id)
	usr.Role = user.GetRolesById(usr.Id)
	return usr
}

func GetData(id int64) *user2.UserData {
	return user.GetDataByUserId(id)
}

func IsHospitalWorker(userId, hospitalId int64) bool {
	worker := &user2.Worker{}
	err := mysql.GetConnect().Where("user_id = ? and hospital_id = ?", userId, hospitalId).Take(worker).Error
	if err == nil {
		return true
	} else if err == gorm.ErrRecordNotFound {
		return false
	}
	panic(err.Error())
}
