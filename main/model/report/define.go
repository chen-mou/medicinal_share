package report

import (
	"errors"
	"gorm.io/gorm"
	"medicinal_share/main/entity/report"
	"medicinal_share/tool/db/mysql"
)

func CreateDefineTx(def *report.Define, tx *gorm.DB) error {
	return tx.Create(def).Error
}

func CreateDefine(def *report.Define) {
	if err := CreateDefineTx(def, mysql.GetConnect()); err != nil {
		panic(err)
	}
}

func GetDefineByProjectId(projectId int64) *report.Define {
	def := &report.Define{}
	if err := mysql.GetConnect().Where("project_id = ?", projectId).Find(&def).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(err)
	}
	return def
}
