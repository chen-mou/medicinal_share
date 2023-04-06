package report

import (
	"bytes"
	"encoding/json"
	"medicinal_share/main/entity/report"
	"medicinal_share/tool/db/elasticsearch"
)

func Create(projectId int64, userId int64, data map[string]any) {
	define := GetDefineByProjectId(projectId)
	r := &report.Base{}
	r.ProjectId = projectId
	r.UserId = userId
	r.Date = data
	if err := elasticsearch.Save(define.Indices, r); err != nil {
		panic(err)
	}
}

func GetUserReports(userId int64) []*report.Base {
	body := map[string]any{
		"query": map[string]any{
			"filter": map[string]any{
				"user_id": userId,
			},
		},
	}
	byt, _ := json.Marshal(body)
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
