package entity

type Report struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Result    string    `json:"result" binding:"required"`
	ImageId   int64     `json:"image_id" binding:"required"`
	ReserveId int64     `json:"reserve_id" binding:"required"`
	Image     *FileData `json:"image" gorm:"foreignKey:ImageId""`
	Reserve   *Reserve  `json:"reserve" gorm:"foreignKey:ReserveId"`
}

type ShareReport struct {
	Id        int64  `json:"id" gorm:"primaryKey"`
	ReportId  int64  `json:"report_id"`
	ReserveId int64  `json:"reserve_id"`
	Report    Report `json:"report" gorm:"foreignKey:ReportId"`
}
