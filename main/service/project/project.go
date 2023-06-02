package project

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/project"
	"medicinal_share/main/model/report"
	"medicinal_share/main/service/user"
	"medicinal_share/tool/db/mysql"
	"strconv"
	"time"
)

func GetProjectByHospitalId(id int64, last int64) []*entity.Project {
	return project.GetProjectByHospitalId(id, last)
}

func GetProjectReserveByDataAndProject(t time.Time, projectId int64) []*entity.ProjectReserve {
	today := strconv.Itoa(t.Day())
	if t.Day() < 10 {
		today = "0" + today
	}
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	return project.GetProjectReserveByDateAndProjectId(entity.CreateTime(start),
		entity.CreateTime(end),
		projectId)
}

func CreateProjectReserve(reserve *entity.ProjectReserve, createBy int64) {
	p := project.GetById(reserve.ProjectId)
	if p == nil {
		panic(middleware.NewCustomErr(middleware.ERROR, "项目不存在"))
	}
	if !user.IsHospitalWorker(createBy, p.HospitalId) {
		panic(middleware.NewCustomErr(middleware.ERROR, "没有添加的权限"))
	}
	project.CreateProjectReserve(reserve)
}

func CreateReserve(projectReserveId []int64, userId int64, tx *gorm.DB) {
	reserve := make([]*entity.Reserve, 0)
	for _, id := range projectReserveId {
		projectReserve := project.GetProjectReserveById(id, tx)
		if projectReserve == nil {
			panic("预约不存在")
		}
		reserve = append(reserve, &entity.Reserve{
			ProjectReserveId: projectReserve.Id,
			UserId:           userId,
		})
	}
	project.CreateReserve(reserve, tx)
}

func GetUserReserve(userId int64) []*entity.Reserve {
	return project.GetReserveByUserId(userId)
}

func GetDoctorReserve(userId int64, date time.Time) []*entity.Reserve {
	doc := user.GetDoctorInfoByUserId(userId)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	end := time.Date(now.Year(), now.Month(), now.Day(), 24, 59, 59, 0, loc)
	return report.GetByDoctorAndDate(doc.Id, entity.Unused, start, end)
}

func Come(userId, reserveId int64) {
	doc := user.GetDoctorInfoByUserId(userId)
	reserve := project.GetReserveById(reserveId)
	now := time.Now()
	if reserve.ProjectReserve.Start.Time().After(now) {
		panic(middleware.NewCustomErr(middleware.ERROR, "不在预约的时间段中"))
	}
	if !project.IsDoctorReserve(doc.Id, reserveId) {
		panic(middleware.NewCustomErr(middleware.FORBID, "没有相关权限"))
	}
	project.UpdateReserveStatus(reserveId, entity.Using, mysql.GetConnect())
}
