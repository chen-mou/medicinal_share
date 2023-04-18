package tag

import (
	"medicinal_share/main/entity"
	"medicinal_share/main/model/tag"
)

func GetByType(typ string, parent int64) []*entity.Tag {
	return tag.GetTagByType(typ, parent)
}

func Search(key, typ string, parent int64, page *entity.Page) []*entity.Tag {
	return tag.SearchByKeyWord(key, parent, typ, page)
}
