package report

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
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

func UploadReport(report *entity.Report, userId int64) *entity.Report {
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		doc := user.GetDoctorInfoByUserId(userId)
		if !project.IsDoctorReserve(doc.Id, report.ReserveId) {
			panic(middleware.NewCustomErr(middleware.FORBID, "你没有权限"))
		}
		fd := &entity.FileData{}
		err := tx.Model(fd).Where("id = ? and type = 'report_image'", report.ImageId).Find(fd).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				panic(middleware.NewCustomErr(middleware.ERROR, "不存在这个报告图"))
			}
			panic(err)
		}
		reserve := project.GetReserveById(report.ReserveId)
		report.Name = reserve.ProjectReserve.Start.Time().Format("2006-01-02 15:04:05") + " ~ " +
			reserve.ProjectReserve.End.Time().Format("15:04:05") + " " + reserve.ProjectReserve.Project.Name
		model.Create(report, tx)
		project.UpdateReserveStatus(report.ReserveId, entity.Completed, tx)
		return nil
	})
	return report
}

func GetUserReport(userId, reserveId int64) []*entity.Report {
	res := model.GetByUserId(userId, reserveId)
	for _, v := range res {
		v.Reserve = nil
	}
	return res
}

func GetDoctorReport(userId int64) []*entity.Reserve {
	info := user.GetDoctorInfoByUserId(userId)
	if info == nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "你不是医生"))
	}
	return model.GetByDoctor(info.Id,
		entity.Using)
}

func CreateShareReport(userId, reserveId int64, reportsId []int64) {
	if !model.IsUserReports(reportsId, userId) {
		panic(middleware.NewCustomErr(middleware.ERROR, "不是你的报告"))
	}
	if !project.IsUserReserve(reserveId, userId) {
		panic(middleware.NewCustomErr(middleware.ERROR, "不是你的预约"))
	}
	project.CreateShareReport(reportsId, reserveId)
}

func GetShareReport(userId, reserveId int64) []*entity.ShareReport {
	doc := user.GetDoctorInfoByUserId(userId)
	if !project.IsDoctorReserve(doc.Id, reserveId) {
		panic(middleware.NewCustomErr(middleware.FORBID, "没有查看的权限"))
	}
	return model.GetShareReportByReserveId(reserveId)
}

func GetReportById(reportId, userId int64) *entity.Report {
	doc := user.GetDoctorInfoByUserId(userId)
	if !model.HaveReportPermission(reportId, doc.Id) {
		panic(middleware.NewCustomErr(middleware.FORBID, "你没有权限"))
	}
	return model.GetById(reportId)
}

func GetByReserveId(reserveId, userId int64) *entity.Report {
	if !project.IsUserReserve(reserveId, userId) {
		panic(middleware.NewCustomErr(middleware.FORBID, "你没有权限"))
	}
	return model.GetByReserveId(reserveId)
}
