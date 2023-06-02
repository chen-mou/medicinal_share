package entity

type ReserveStatue uint

const (
	Unused = iota
	Using
	Overdue
	Completed //完成
)

// ProjectReserve 用于展示项目可以预约的时间
type ProjectReserve struct {
	Id         int64       `json:"id"`
	Overplus   int         `json:"overplus" gorm:"default:0"`    //人数
	ReserveNum int         `json:"reserve_num" gorm:"default:0"` //预约人数
	Start      Time        `json:"start" bind:"required"`
	End        Time        `json:"end" bind:"required"`
	ProjectId  int64       `json:"project_id" bind:"required"`
	DoctorId   int64       `json:"doctor_id" bind:"required"` //主治医生
	Project    *Project    `json:"project" gorm:"foreignKey:ProjectId"`
	DoctorInfo *DoctorInfo `json:"doctor_info" gorm:"foreignKey:DoctorId"`
}

// Project 项目介绍
type Project struct {
	Id          int64   `json:"id" gorm:"primaryKey"`
	HospitalId  int64   `json:"hospital_id" gorm:"index"`
	Name        string  `json:"name" gorm:"size:64"`
	Price       float64 `json:"price" gorm:"type:decimal"`
	Type        string  `json:"type" gorm:"size:32;index"`
	Description string  `json:"description" gorm:"type:tinytext"`
}

// Reserve 用户预约
type Reserve struct {
	Id               int64           `json:"id" gorm:"primaryKey"`
	ProjectReserveId int64           `json:"project_reserve_id" gorm:"index"`
	Status           ReserveStatue   `json:"status" gorm:"size:16;default:0"`
	UserId           int64           `json:"user_id" gorm:"index"`
	User             *User           `json:"user" gorm:"UserId"`
	ProjectReserve   *ProjectReserve `json:"project_reserve" gorm:"foreignKey:ProjectReserveId"`
	ShareReports     []ShareReport   `json:"share_reports" gorm:"foreignKey:ReserveId"`
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

func (ShareReport) TableName() string {
	return "share_report"
}
