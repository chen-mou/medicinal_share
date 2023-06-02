package report

import (
	"context"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/model"
	"medicinal_share/tool/db/mysql"
	"time"
)

//func Create(projectId int64, userId int64, data map[string]any) {
//	define := GetDefineByProjectId(projectId)
//	r := &report.Base{}
//	r.ProjectId = projectId
//	r.UserId = userId
//	r.Date = data
//	if err := elasticsearch.Save("report-"+define.Indices, r); err != nil {
//		panic(err)
//	}
//}
//
//func GetReport(projectId, userId int64) *report.Base {
//	Search := elasticsearch.GetClient().Search
//	queryBody := map[string]any{
//		"size": 1,
//		"query": map[string]any{
//			"bool": map[string]any{
//				"term": map[string]int64{
//					"project_id": projectId,
//					"user_id":    userId,
//				},
//			},
//		},
//	}
//	byt, _ := json.Marshal(queryBody)
//	res := &report.Base{}
//	elasticsearch.Get(res,
//		Search.WithBody(bytes.NewBuffer(byt)),
//		Search.WithIndex("report-*"))
//	return res
//}
//
//func GetReportByUserId(userId int64, last entity.Time, num int) []*report.Base {
//	queryBody := map[string]any{
//		"size": num,
//		"query": map[string]any{
//			"bool": map[string]any{
//				"filter": map[string]any{
//					"term": map[string]int64{
//						"user_id": userId,
//					},
//					"range": map[string]any{
//						"le": last,
//					},
//				},
//			},
//		},
//		"sort": map[string]any{
//			"create_at": map[string]any{
//				"order": "desc",
//			},
//		},
//	}
//	byt, _ := json.Marshal(queryBody)
//	res := make([]*report.Base, 0)
//	search := elasticsearch.GetClient().Search
//	elasticsearch.Get(res, search.WithBody(bytes.NewBuffer(byt)), search.WithIndex("report-*"))
//	return res
//}
//
//func SaveReport(report report.Base) {
//	def := GetDefineByProjectId(report.ProjectId)
//	report.CreateAt = entity.CreateTime(time.Now())
//	byt, _ := json.Marshal(report)
//	if _, err := elasticsearch.GetClient().Create("report-"+def.Indices,
//		elasticsearch.GetRandomId("report"),
//		bytes.NewBuffer(byt)); err != nil {
//		panic(err)
//	}
//}

func Create(report *entity.Report, tx *gorm.DB) {
	err := tx.Create(report).Error
	if err != nil {
		panic(err)
	}
}

func GetByUserId(userId, reserveId int64) []*entity.Report {
	res := make([]*entity.Report, 0)
	db := mysql.GetConnect()
	err := db.Model(&entity.Report{}).
		Joins("Reserve").
		Where("Reserve.user_id = ? and reports.id not in(?)", userId,
			db.Model(&entity.ShareReport{}).
				Select("report_id").
				Where("reserve_id = ?", reserveId)).Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Report)
}

func GetById(id int64) *entity.Report {
	res := &entity.Report{}
	err := mysql.GetConnect().Where("id = ?", id).Take(res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
}

func GetByDoctorAndDate(doctorId int, statue entity.ReserveStatue, start, end time.Time) []*entity.Reserve {
	res := make([]*entity.Reserve, 0)
	s := entity.CreateTime(start)
	e := entity.CreateTime(end)
	db := mysql.GetConnect().WithContext(context.TODO())
	err := db.Model(&entity.Reserve{}).
		Joins("ProjectReserve").
		Joins("User").
		Preload("ProjectReserve.Project").
		Preload("User.UserInfo").
		Preload("User.UserInfo.AvatarFile").
		Preload("User.UserInfo.AvatarFile.File").
		Preload("User.UserInfo.RealInfo").
		Where("ProjectReserve.doctor_id = ? and"+
			" Reserve.status = ? and"+
			" ProjectReserve.start between ? and ?", doctorId, statue, s, e).
		Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Reserve)
}

func GetByDoctor(doctorId int, statue entity.ReserveStatue) []*entity.Reserve {
	res := make([]*entity.Reserve, 0)
	db := mysql.GetConnect().WithContext(context.TODO())
	err := db.Model(&entity.Reserve{}).
		Joins("ProjectReserve").
		Joins("User").
		Preload("ProjectReserve.Project").
		Preload("User.UserInfo").
		Preload("User.UserInfo.AvatarFile").
		Preload("User.UserInfo.AvatarFile.File").
		Preload("User.UserInfo.RealInfo").
		Where("ProjectReserve.doctor_id = ? and"+
			" Reserve.status = ?", doctorId, statue).
		Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Reserve)
}

func GetByReserveId(reserveId int64) *entity.Report {
	res := &entity.Report{}
	err := mysql.GetConnect().Where("reserve_id = ?", reserveId).
		Joins("Image").Preload("Image.File").Take(res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
}

func IsUserReports(reports []int64, userId int64) bool {
	res := make([]*entity.Report, 0)
	mysql.GetConnect().
		Model(&entity.Report{}).
		Joins("Reserve").
		Where("Reserve.user_id = ? and reports.id in(?)", userId, reports).Find(&res)
	return len(res) == len(reports)
}

func HaveReportPermission(report int64, doc int) bool {
	db := mysql.GetConnect()
	sp := &entity.ShareReport{}
	err := db.Where("report_id = ? and reserve_id in (?)", report,
		db.Model(&entity.Reserve{}).
			Joins("ProjectReserve").
			Where("ProjectReserve.doctor_id = ?", doc)).Take(sp).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return false
	}
	return true
}

func GetShareReportByReserveId(reserveId int64) []*entity.ShareReport {
	res := make([]*entity.ShareReport, 0)
	err := mysql.GetConnect().
		Where("share_report.reserve_id = ?", reserveId).
		Joins("Report").
		Preload("Report.Image").
		Preload("Report.Image.File").Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.ShareReport)
}
