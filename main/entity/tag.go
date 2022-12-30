package entity

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type" gorm:"index"` //那种类型的tag 如领域，症状，病情，学位
}

type TagRelation struct {
	Id           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	RelationType string `json:"relation_type" gorm:"size:32;index:idx_type_relation"`
	RelationId   int64  `json:"relation_id" gorm:"index:idx_type_relation"`
	TagId        int64  `json:"tag_id" gorm:"index"`
	Tag          Tag    `json:"tag" gorm:"foreignKey:TagId"`
}

func (Tag) TableName() string {
	return "tag_def"
}

func (TagRelation) TableName() string {
	return "tag_ref"
}
