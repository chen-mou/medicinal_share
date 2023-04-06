package tag

import (
	"medicinal_share/main/entity"
	"medicinal_share/main/model/tag"
)

func GetByType(typ string, parent int64) []*entity.Tag {
	return tag.GetTagByType(typ, parent)
}
