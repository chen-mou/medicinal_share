package report

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/entity/report"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/project"
	model "medicinal_share/main/model/report"
	"medicinal_share/main/service/user"
	"medicinal_share/tool/db/mysql"
)

//没时间搞这一套了
//func UploadReport(projectId int64, userId int64, data map[string]any) {
//	def := model.GetDefineByProjectId(projectId)
//	cols := make([]map[string]any, 0)
//	json.Unmarshal(tool.StringToBytes(def.Define), cols)
//	//验证数据是否正确
//
//	for _, v := range cols {
//		col := column.Factory(v["type"].(string))
//		col.Builder(v)
//		val := data[col.GetName()]
//		if err := col.Verify(val); err != nil {
//			panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
//		}
//	}
//
//	now := time.Now()
//	//TODO:存elasticsearch里
//	report := report.Base{
//		CreateAt:  entity.CreateTime(now),
//		Date:      data,
//		ProjectId: projectId,
//		UserId:    userId,
//	}
//	model.SaveReport(report)
//}
//
//func GetReport(id string) *report.Base {
//	return nil
//}
//
//func GetAllUserReport(userId int64) []*report.Base {
//	return nil
//}

func UploadReport(report *report.Report) *report.Report {
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		fd := &entity.FileData{}
		err := tx.Model(fd).Where("id = ? and type = 'report_image'", report.ImageId).Find(fd).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				panic(middleware.NewCustomErr(middleware.ERROR, "不存在这个报告图"))
			}
			panic(err)
		}
		model.Create(report, tx)
		project.UpdateReserveStatus(report.ReserveId, entity.Completed, tx)
		return nil
	})
	return report
}

func GetUserReport(userId int64) []*report.Report {
	return model.GetByUserId(userId)
}

func GetDoctorReport(userId int64) []*entity.Reserve {
	info := user.GetDoctorInfoByUserId(userId)
	if info == nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "你不是医生"))
	}
	return model.GetByDoctor(info.Id)
}
