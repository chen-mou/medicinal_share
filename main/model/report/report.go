package report

import (
	"bytes"
	"encoding/json"
	"medicinal_share/main/entity/report"
	"medicinal_share/tool/db/elasticsearch"
)

//TODO:报告由mysql储存基本信息，如属于那个索引，elasticsearch储存实体

func Create(projectId int64, userId int64, data map[string]any) {
	define := GetDefineByProjectId(projectId)
	r := &report.Base{}
	r.ProjectId = projectId
	r.UserId = userId
	r.Date = data
	if err := elasticsearch.Save("report-"+define.Indices, r); err != nil {
		panic(err)
	}
}

func GetReport(projectId int64, userId int64) *report.Base {

}

func GetReportByUserId(userId int64, last int64, num int) []*report.Base {
	queryBody := map[string]any{
		"bool": map[string]any{
			"filter": map[string]any{
				"term": map[string]int64{
					"user_id": userId,
				},
			},
		},
	}
	byt, _ := json.Marshal(queryBody)
	res := make([]*report.Base, 0)
	search := elasticsearch.GetClient().Search
	elasticsearch.Get(res, search.WithBody(bytes.NewBuffer(byt)), search.WithIndex("report-*"))
	return res
}

func SaveReport(report report.Base) {
	def := GetDefineByProjectId(report.ProjectId)
	byt, _ := json.Marshal(report)
	if _, err := elasticsearch.GetClient().Create("report"+def.Indices,
		elasticsearch.GetRandomId("report"),
		bytes.NewBuffer(byt)); err != nil {
		panic(err)
	}
}
