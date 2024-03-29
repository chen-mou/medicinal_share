package entity

import (
	"gorm.io/gorm"
)

type Sex Enum

const (
	MAN Sex = iota
	WOMAN
)

type User struct {
	Id         int64       `json:"id" gorm:"primaryKey"`
	Username   string      `json:"username" gorm:"uniqueIndex;not null;size:32"`
	Password   string      `json:"password" gorm:"not null;size:64"`
	UserInfo   *UserData   `json:"user_info" gorm:"foreignKey:UserId"`
	DockerInfo *DoctorInfo `json:"docker_info,omitempty" gorm:"UserId"`
	Role       []*UserRole `json:"role" gorm:"foreignKey:UserId"`
}

type UserData struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Nickname   string    `json:"nickname" gorm:"index;size:64"`
	UserId     int64     `json:"user_id" gorm:"uniqueIndex;not null"`
	InfoId     *int64    `json:"info_id" gorm:"uniqueIndex"`
	Avatar     int64     `json:"avatar" gorm:"default:1"`
	HelpNum    int       `json:"help_num"`
	IsReal     bool      `json:"is_real" gorm:"-"`
	AvatarFile *FileData `json:"avatar_file,omitempty" gorm:"foreignKey:Avatar"`
	RealInfo   *RealInfo `json:"real_info,omitempty" gorm:"foreignKey:InfoId"`
}

// RealInfo TODO: 完成实名认证的模型定义
type RealInfo struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"size:16" binding:"required"`
	Sex      Sex    `json:"sex" gorm:"size:2"`
	IDNumber string `json:"id_number" binding:"required"`
}

type UserRole struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId int64  `json:"user_id" gorm:"uniqueIndex:idx_user_id_name;"`
	Name   string `json:"name" gorm:"uniqueIndex:idx_user_id_name;size:16"`
}

type DoctorInfo struct {
	Id           int            `json:"id" gorm:"primaryKey"`
	DoctorAvatar int64          `json:"doctor_avatar"`
	UserId       int64          `json:"user_id" gorm:"uniqueIndex;not null"`
	HospitalId   int64          `json:"hospital_id" gorm:"index"`             //工作医院
	Position     string         `json:"position" gorm:"size:64"`              //职位
	Status       string         `json:"status" gorm:"size:16;default:normal"` //医生当前状态 normal 空闲 busy 忙碌中 offline 下线
	Description  string         `json:"description" gorm:"type:tinytext"`
	TagsId       []int64        `json:"tags_id" gorm:"-"`
	Tags         []*TagRelation `json:"tags,omitempty" gorm:"foreignKey:RelationId"`
	InfoId       *int64         `json:"infoId,omitempty" gorm:"uniqueIndex;not null"`
	Info         *RealInfo      `json:"info,omitempty" gorm:"foreignKey:InfoId"`
	Avatar       *FileData      `json:"avatar,omitempty" gorm:"foreignKey:DoctorAvatar"`
	Hospital     *Hospital      `json:"hospital" gorm:"foreignKey:HospitalId"`
}

type Worker struct {
	Id         int64    `json:"id" gorm:"primaryKey;"`
	HospitalId int64    `json:"hospital_id"`
	UserId     int64    `json:"user_id" gorm:"uniqueIndex"`
	Hospital   Hospital `json:"hospital" gorm:"HospitalId"`
}

func (Worker) TableName() string {
	return "worker"
}

func (User) TableName() string {
	return "user"
}

func (UserData) TableName() string {
	return "user_data"
}

func (u *UserData) AfterFind(tx *gorm.DB) error {
	u.IsReal = !(u.InfoId == nil)
	return nil
}

func (RealInfo) TableName() string {
	return "real_info"
}

func (DoctorInfo) TableName() string {
	return "doctor_info"
}

func (UserRole) TableName() string {
	return "user_role"
}

func (r *UserRole) MarshalJSON() ([]byte, error) {
	return []byte("\"" + r.Name + "\""), nil
}

func (r *UserRole) UnmarshalJSON(b []byte) error {
	r.Name = string(b[1 : len(b)-1])
	return nil
}
