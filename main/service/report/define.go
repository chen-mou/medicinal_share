package report

import (
	"encoding/json"
	"errors"
	"medicinal_share/main/entity/report"
	"medicinal_share/main/entity/report/column"
	"medicinal_share/main/middleware"
	model "medicinal_share/main/model/report"
	"medicinal_share/tool/encrypt/md5"
	"reflect"
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
