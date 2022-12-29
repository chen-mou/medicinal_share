package entity

type User struct {
	Id       int64     `json:"id" gorm:"primaryKey"`
	Username string    `json:"username" gorm:"uniqueIndex;not null;size:32"`
	Password string    `json:"password" gorm:"not null;size:64"`
	UserInfo *UserData `json:"user_info,omitempty" gorm:"-"`
}

type UserData struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Nickname   string    `json:"nickname" gorm:"index;size:64"`
	UserId     int64     `json:"user_id" gorm:"uniqueIndex;not null"`
	InfoId     *int64    `json:"info_id" gorm:"uniqueIndex"`
	Avatar     int64     `json:"avatar"`
	HelpNum    int       `json:"help_num"`
	SellNum    int       `json:"sell_num"`
	IsReal     bool      `json:"is_real" gorm:"-"`
	AvatarFile *File     `json:"avatar_file,omitempty" gorm:"-"`
	RealInfo   *RealInfo `json:"real_info,omitempty" gorm:"-"`
}

// RealInfo TODO: 完成实名认证的模型定义
type RealInfo struct {
}

func (User) TableName() string {
	return "user"
}

func (UserData) TableName() string {
	return "user_data"
}

func (RealInfo) TableName() string {
	return "real_info"
}
