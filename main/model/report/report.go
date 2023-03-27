package report

import (
	"medicinal_share/main/entity/report"
	"medicinal_share/tool/db/elasticsearch"
	"medicinal_share/tool/db/mysql"
)

//TODO:报告由mysql储存基本信息，如属于那个索引，elasticsearch储存实体

func Create(projectId int64, userId int64, data map[string]any) {
	define := &report.Define{}
	r := &report.Base{}
	if err := mysql.GetConnect().Where("project_id = ?", projectId).Take(define).Error; err != nil {
		panic(err)
	}
	r.ProjectId = projectId
	r.UserId = userId
	r.Date = data
	if err := elasticsearch.Save(define.Indices, r); err != nil {
		panic(err)
	}
}

func GetReport(projectId int64, userId int64) *report.Base {

}

func GetReportByUserId(userId int64) []*report.Base {

}
