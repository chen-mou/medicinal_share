package tag

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"strconv"
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

const (
	GetTypeLock = "GET_TAG_TYPE_LOCK"
	GetType     = "GET_TAG_TYPE"
)

func GetTagByType(typ string, parent int64) []*entity.Tag {
	res := make([]*entity.Tag, 0)
	pstr := strconv.FormatInt(parent, 10)
	get := GetType + ":" + typ + ":" + pstr
	lock := GetTypeLock + ":" + typ + ":" + pstr
	c := redis.NewCache(lock, get)
	c.Get(&res, func() any {
		err := mysql.GetConnect().Model(&entity.Tag{}).
			Where("type = ? and parent = ?", typ, parent).Find(&res).Error
		if err != nil {
			panic(err)
		}
		return res
	})
	return res
}

func SearchByKeyWord(key string, parent int64, typ string, page *entity.Page) []*entity.Tag {
	res := make([]*entity.Tag, 0)
	if err := mysql.GetConnect().Model(&entity.Tag{}).
		Where("parent = ? and type = ? and name like ?", parent, typ, key+"%").
		Limit(page.Limit).Offset(page.Offset).
		Find(&res).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return res
}

func AddTag(tag entity.Tag) {
	mysql.GetConnect().Create(&tag)
}
