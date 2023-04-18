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

// ProjectReserve 用于展示项目可以预约的时间
type ProjectReserve struct {
	Id         int64      `json:"id"`
	Num        int        `json:"num"`                 //人数
	ReserveNum int        `json:"reserve_num" gorm:""` //预约人数
	Start      Time       `json:"start"`
	End        Time       `json:"end"`
	ProjectId  int64      `json:"project_id"`
	DoctorId   int64      `json:"doctor_id"` //主治医生
	Project    Project    `json:"project" gorm:"foreignKey:ProjectId"`
	DoctorInfo DoctorInfo `json:"doctor_info" gorm:"foreignKey:DoctorId"`
}

// Project 项目介绍
type Project struct {
	Id          int64   `json:"id" gorm:"primaryKey"`
	HospitalId  int64   `json:"hospital_id" gorm:"index"`
	Name        string  `json:"name" gorm:"size:64"`
	Price       float64 `json:"price"`
	Type        string  `json:"type" gorm:"size:32;index"`
	Description string  `json:"description" gorm:"type:tinytext"`
}

// Reserve 用户预约
type Reserve struct {
	Id        int64          `json:"id" gorm:"primaryKey"`
	ProjectId int64          `json:"project_id" gorm:"index"`
	Status    string         `json:"status" gorm:"size:16;default:pending"`
	UserId    int64          `json:"user_id" gorm:"index"`
	Project   ProjectReserve `json:"project" gorm:"foreignKey:ProjectId"`
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

func (ProjectReserve) TableName() string {
	return "project_reserve"
}
