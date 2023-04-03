package main

import (
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
)

func main() {
	db := mysql.GetConnect()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(entity.User{},
		entity.UserData{},
		entity.FileData{},
		entity.File{},
		entity.Reserve{},
		entity.Project{},
		entity.UserRole{},
		entity.DoctorInfo{},
		entity.RealInfo{},
		entity.Tag{},
		entity.TagRelation{},
		entity.Order{},
		entity.OrderData{},
		entity.Hospital{},
		entity.Message{},
		entity.ProjectReserve{})
}
