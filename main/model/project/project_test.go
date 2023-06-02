package project

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"medicinal_share/main/entity"
	"medicinal_share/tool/db/mysql"
	"testing"
	"time"
)

func TestGetHospitalByNear(t *testing.T) {
	byt, _ := json.Marshal(GetHospitalByNear(10.5, 200, 1, 500))
	fmt.Println(string(byt))
}

func TestCreateProject(t *testing.T) {
	names := []string{"一般检查", "内科", "尿常规", "生化-肝功", "生化-肾功", "外科", "放射科", "心电图"}
	for _, name := range names {
		value := 5 + rand.Float64()*5
		mysql.GetConnect().Create(&entity.Project{
			HospitalId:  17,
			Name:        name,
			Price:       float64(int(value*100)) / 100,
			Type:        "Normal",
			Description: "检查一下身体",
		})
	}
}

func TestCreateProjectReserve(t *testing.T) {
	now := time.Now()
	for i := 0; i < 10; i++ {
		start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+i, 0, 0, 0, now.Location())
		end := start.Add(time.Minute * 30)
		CreateProjectReserve(&entity.ProjectReserve{
			Start:     entity.CreateTime(start),
			End:       entity.CreateTime(end),
			DoctorId:  2,
			ProjectId: 26,
			Overplus:  20,
		})
	}
}
