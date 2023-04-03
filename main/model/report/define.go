package report

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"medicinal_share/main/entity/report"
	"medicinal_share/tool/db/elasticsearch"
)

func CreateDefine(def *report.Define) {
	client := elasticsearch.GetClient()
	byt, _ := json.Marshal(def)
	if _, err := client.Create("report_define",
		elasticsearch.GetRandomId("report_define"),
		bytes.NewBuffer(byt)); err != nil {
		panic(err)
	}
}

func GetDefineByProjectId(projectId int64) *report.Define {
	def := &report.Define{}
	//if err := mysql.GetConnect().Where("project_id = ?", projectId).Find(&def).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil
	//	}
	//	panic(err)
	//}
	query := map[string]any{
		"query": map[string]any{
			"filter": map[string]any{
				"project_id": projectId,
			},
		},
	}
	client := elasticsearch.GetClient()
	byt, _ := json.Marshal(query)
	res, err := client.Search(client.Search.WithIndex("report_define"),
		client.Search.WithBody(bytes.NewBuffer(byt)),
		client.Search.WithSize(1),
		client.Search.WithSource("id", "project_id", "define", "name", "index"))
	if err != nil {
		panic(err)
	}
	byt, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal(byt, def)
	return def
}

func GetAllDefineIndex() []string {
	client := elasticsearch.GetClient()
	body := map[string]any{
		"_source": []string{"indices"},
	}
	byt, _ := json.Marshal(body)
	res, err := client.Search(client.Search.WithRequestCache(true), client.Search.WithBody(bytes.NewBuffer(byt)), client.Search.WithIndex("report_define"))
	if err != nil {
	}
}
