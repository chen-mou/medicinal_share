package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"medicinal_share/main/entity"
	"medicinal_share/main/entity/report"
	"medicinal_share/main/entity/report/column"
	"medicinal_share/main/middleware"
	model "medicinal_share/main/model/report"
	"medicinal_share/tool"
	"medicinal_share/tool/db/elasticsearch"
	"medicinal_share/tool/encrypt/md5"
	"reflect"
	"time"
)

func verifyType(val any, kind reflect.Kind) error {
	typ := reflect.TypeOf(val)
	if typ.Kind() != kind {
		return errors.New("类型有误")
	}
	return nil
}

func CreateDefine(define []map[string]any, projectId int64, name string) {
	def := &report.Define{}
	for _, def := range define {
		if typ, ok := def["type"]; ok {
			if err := verifyType(typ, reflect.String); err != nil {
				panic(middleware.NewCustomErr(middleware.ERROR, "字段:type"+err.Error()))
			}
			t := typ.(string)
			col := column.Factory(t)
			col.Builder(def)
		}
	}
	def.Indices = md5.Hash(name)
	def.Name = name
	def.ProjectId = projectId
	var err error
	var byt []byte
	if byt, err = json.Marshal(define); err == nil {
		def.Define = string(byt)
	} else {
		panic(middleware.NewCustomErr(middleware.ERROR, "json序列化失败"))
	}
	model.CreateDefine(def)
}

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

func GetReport(id string) report.Base {}

func GetAllUserReport(userId int64) []report.Base {}
