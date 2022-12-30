package entity

import "gorm.io/gorm"

type User struct {
	Id         int64       `json:"id" gorm:"primaryKey"`
	Username   string      `json:"username" gorm:"uniqueIndex;not null;size:32"`
	Password   string      `json:"password" gorm:"not null;size:64"`
	UserInfo   *UserData   `json:"user_info" gorm:"-"`
	DockerInfo *DoctorInfo `json:"docker_info,omitempty" gorm:"UserId"`
	Role       []*UserRole `json:"role" gorm:"-"`
}

type UserData struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Nickname   string    `json:"nickname" gorm:"index;size:64"`
	UserId     int64     `json:"user_id" gorm:"uniqueIndex;not null"`
	InfoId     *int64    `json:"info_id" gorm:"uniqueIndex"`
	Avatar     int64     `json:"avatar"`
	HelpNum    int       `json:"help_num"`
	IsReal     bool      `json:"is_real" gorm:"-"`
	AvatarFile *File     `json:"avatar_file,omitempty" gorm:"-"`
	RealInfo   *RealInfo `json:"real_info,omitempty" gorm:"foreignKey:InfoId"`
}

// RealInfo TODO: 完成实名认证的模型定义
type RealInfo struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"size:16"`
	IDNumber string `json:"id_number"`
}

type UserRole struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId int64  `json:"user_id" gorm:"uniqueIndex:idx_user_id_name;"`
	Name   string `json:"name" gorm:"uniqueIndex:idx_user_id_name;size:16"`
}

type DoctorInfo struct {
	Id          int            `json:"id" gorm:"primaryKey"`
	UserId      int64          `json:"user_id" gorm:"uniqueIndex;not null"`
	Work        string         `json:"work" gorm:"size:64"`     //工作医院
	Position    string         `json:"position" gorm:"size:64"` //职位
	Description string         `json:"description" gorm:"type:tinytext"`
	TagsId      []int64        `json:"tags_id" gorm:"-"`
	Tags        []*TagRelation `json:"tags" gorm:"foreignKey:RelationId"`
	InfoId      *int64         `json:"infoId" gorm:"uniqueIndex;not null"`
	Info        *RealInfo      `json:"info" gorm:"foreignKey:InfoId"`
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
