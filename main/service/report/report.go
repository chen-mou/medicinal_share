package report

import (
	"bytes"
	"encoding/json"
	"medicinal_share/main/entity"
	"medicinal_share/main/entity/report"
	"medicinal_share/main/entity/report/column"
	"medicinal_share/main/middleware"
	model "medicinal_share/main/model/report"
	"medicinal_share/tool"
	"medicinal_share/tool/db/elasticsearch"
	"time"
)

func UploadReport(projectId int64, userId int64, data map[string]any) {
	def := model.GetDefineByProjectId(projectId)
	cols := make([]map[string]any, 0)
	json.Unmarshal(tool.StringToBytes(def.Define), cols)
	//验证数据是否正确

	for _, v := range cols {
		col := column.Factory(v["type"].(string))
		col.Builder(v)
		val := data[col.GetName()]
		if err := col.Verify(val); err != nil {
			panic(middleware.NewCustomErr(middleware.ERROR, err.Error()))
		}
	}

	now := time.Now()
	//TODO:存elasticsearch里
	report := report.Base{
		CreateAt:  entity.CreateTime(now),
		Date:      data,
		ProjectId: projectId,
		UserId:    userId,
	}
	byt, _ := json.Marshal(report)
	elasticsearch.GetClient().Create("report", elasticsearch.GetRandomId("report"), bytes.NewBuffer(byt))
}

func GetReport(id string) *report.Base {
	return nil
}

func GetAllUserReport(userId int64) []*report.Base {
	return nil
}
