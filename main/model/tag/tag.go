package tag

import (
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
)

func CreateTagRelation(ids []int64, typ string, id int64) []*entity.TagRelation {
	entities := make([]*entity.TagRelation, 0)
	for _, v := range ids {
		entities = append(entities, &entity.TagRelation{
			RelationType: typ,
			RelationId:   id,
			TagId:        v,
		})
	}
	err := mysql.GetConnect().Model(&entity.TagRelation{}).Create(&entities).Error
	if err != nil {
		panic(err)
	}
	return entities
}
