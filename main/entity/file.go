package entity

import (
	"fmt"
	"gorm.io/gorm"
)

type File struct {
	Id      int64  `json:"id" gorm:"primaryKey;"`
	Hash    string `json:"hash" gorm:"unique;size:32;not null"`
	Machine string `json:"machine" gorm:"size:64;unique;not null"`
	Status  string `json:"status" gorm:"size:16;not null"`
	Path    string `json:"path" gorm:"size:256;not null"`
	Uri     string `json:"uri" gorm:"size:256;not null"`
}

type FileData struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"size:32;not null"`
	Suffix   string `json:"suffix" gorm:"size:16;not null"`
	FileId   int64  `json:"file_id" gorm:"index;not null"`
	Status   string `json:"status" gorm:"size:16;not null"`
	Type     string `json:"type" gorm:"size:16;index:idx_uploader_type;not null"`
	Uploader int64  `json:"uploader" gorm:"index:idx_uploader_type;not null"`
	File     *File  `json:"file" gorm:"foreignKey:file_id"`
}

func (File) TableName() string {
	return "file"
}

func (FileData) TableName() string {
	return "file_data"
}

func (FileData) AfterFind(tx *gorm.DB) error {
	fmt.Println("fuck you")
	return nil
}
