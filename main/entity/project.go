package entity

import "time"

type Hospital struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:64"`
	Address    string    `json:"address" gorm:"size:64"`
	Longitude  float64   `json:"longitude" `
	Latitude   float64   `json:"latitude"`
	Distance   float64   `json:"distance" gorm:"-:migration"`
	Avatar     int64     `json:"avatar" gorm:"not null;default:1"`
	AvatarFile FileData  `json:"avatar_file" gorm:"foreignKey:avatar"`
	Projects   []Project `json:"projects,omitempty" gorm:"foreignKey:HospitalId"`
}

type Project struct {
	Id          int64   `json:"id" gorm:"primaryKey"`
	HospitalId  int64   `json:"hospital_id" gorm:"index"`
	Name        string  `json:"name" gorm:"size:64"`
	Price       float64 `json:"price"`
	Type        string  `json:"type" gorm:"size:32;index"`
	Description string  `json:"description" gorm:"type:tinytext"`
}

type Reserve struct {
	Id        int64      `json:"id" gorm:"primaryKey"`
	Time      *time.Time `json:"time" gorm:"type:datetime"`
	ProjectId int64      `json:"project_id" gorm:"index"`
	Status    string     `json:"status" gorm:"size:16;default:pending"`
	UserId    int64      `json:"user_id" gorm:"index"`
	Project   Project    `json:"project" gorm:"foreignKey:ProjectId"`
}

func (Hospital) TableName() string {
	return "hospital"
}

func (Project) TableName() string {
	return "project"
}

func (Reserve) TableName() string {
	return "reserve"
}
