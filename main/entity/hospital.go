package entity

type Hospital struct {
	Id             int64     `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"size:64"`
	Address        string    `json:"address" gorm:"size:64"`
	Longitude      float64   `json:"longitude" `
	Latitude       float64   `json:"latitude"`
	Distance       float64   `json:"distance" gorm:"-:migration"`
	Avatar         int64     `json:"avatar" gorm:"not null;default:1"`
	Background     int64     `json:"background" gorm:"not null;default:1"`
	BackgroundFile FileData  `json:"background_file" gorm:"foreignKey:Background"`
	AvatarFile     FileData  `json:"avatar_file" gorm:"foreignKey:Avatar"`
	Description    string    `json:"description" gorm:"type:tinytext"`
	Projects       []Project `json:"projects,omitempty" gorm:"foreignKey:HospitalId"`
}

func (Hospital) TableName() string {
	return "hospital"
}
