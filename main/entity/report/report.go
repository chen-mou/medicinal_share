package report

type Define struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectId int64  `json:"project_id" gorm:"uniqueIndex"`
	Define    string `json:"define"`
	Name      string `json:"name" gorm:"size:64"`
	Indices   string `json:"index" gorm:"size:64"` //对应的elasticsearch的索引
}

//Base 报告的基本信息
type Base struct {
	Id        int   `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectId int64 `json:"project_id" gorm:"index:project_user_idx"`
	UserId    int64 `json:"user_id" gorm:"index:project_user_idx"`
}
