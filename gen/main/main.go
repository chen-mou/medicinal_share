package main

import (
	"gorm.io/gen"
	"medicinal_share/gen/dao"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "gen/out/dao",
		Mode:          gen.WithQueryInterface | gen.WithDefaultQuery,
		FieldNullable: true,
		ModelPkgPath:  "dao",
	})

	g.UseDB(mysql.GetConnect())

	g.ApplyInterface(func(filter dao.Filter) {}, entity.User{}, entity.UserData{}, entity.File{}, entity.FileData{})

	g.ApplyInterface(func(user dao.User) {}, entity.User{})

	g.Execute()

}
