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
	client := elasticsearch.GetClient()
	client.Search(client.Search.WithIndex("report-*"), client.Search.WithBody(bytes.NewBuffer(byt)))
}
