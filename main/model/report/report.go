package report

import (
	"context"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/entity/report"
	"medicinal_share/main/model"
	"medicinal_share/tool/db/mysql"
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

func Create(report *report.Report, tx *gorm.DB) {
	tx.Create(report)
}

func GetByUserId(userId int64) []*report.Report {
	res := make([]*report.Report, 0)
	err := mysql.GetConnect().Model(&report.Report{}).Joins("Order").
		Where("Order.user_id = ?", userId).Find(&res).Error
	return model.GetErrorHandler(err, res).([]*report.Report)
}

func GetByDoctor(doctorId int) []*entity.Reserve {
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
		Where("ProjectReserve.doctor_id = ?", doctorId).
		Find(&res).Error
	return model.GetErrorHandler(err, res).([]*entity.Reserve)
}
