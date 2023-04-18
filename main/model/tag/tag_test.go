package tag

import (
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
	"testing"
)

func TestAddTag(t *testing.T) {
	names := []string{"喉咙有痰", "咽喉疼痛", "头痛", "四肢乏力", "食欲不振", "打喷嚏", "咳嗽"}
	for _, name := range names {
		err := mysql.GetConnect().Where("name = ?", name).Find(&entity.Tag{}).Error
		if err == nil {
			continue
		}
		AddTag(entity.Tag{
			Name:   name,
			Type:   "Symptom",
			Parent: 0,
			End:    true,
		})
	}
}
