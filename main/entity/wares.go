package entity

//完成商品的模型定义

type Wares struct {
	Id          int64      `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"size:64"`
	Seller      int64      `json:"seller" gorm:"index"`
	Price       float64    `json:"price"`
	Type        string     `json:"type" gorm:"size:32"`
	Description string     `json:"description" gorm:"type:tinytext"`
	Stock       int        `json:"Stock"`
	PhotoIds	[]int 	`json:"photo_ids" gorm:"-"`
	Photos      []FileData `json:"photos" gorm:"-"`
}

type WaresPhoto struct {
	Id      int   `json:"id" gorm:"primaryKey;autoIncrement"`
	WaresId int64 `json:"wares_id"`
	Cover   int64 `json:"cover"`
}
