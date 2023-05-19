package user

import (
	"gorm.io/gorm"
	user2 "medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/model/file"
	"medicinal_share/main/model/tag"
	"medicinal_share/main/model/user"
	"medicinal_share/tool/db/mysql"
)

func CreateInfo(id int64, real *user2.RealInfo) {
	data := user.GetDataByUserId(id)
	if data.IsReal {
		panic(middleware.NewCustomErr(middleware.ERROR, "该账号已经实名认证了"))
	}
	info := user.GetInfoByNameAndIdNumber(real.Name, real.IDNumber)
	if info != nil {
		panic(middleware.NewCustomErr(middleware.NO_REPETITION, "已被其他账号实名认证"))
	}
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		info := user.CreateInfo(real, tx)
		user.UpdateData(&user2.UserData{
			UserId: id,
			InfoId: &info.Id,
		}, tx)
		return nil
	})
}

func CreateDoctorInfo(userId int64, info *user2.DoctorInfo) {
	data := user.GetDataByUserId(userId)
	if !data.IsReal {
		panic(middleware.NewCustomErr(middleware.FORBID, "未实名认证"))
	}
	in := user.GetDoctorInfoById(userId)
	if in != nil {
		panic(middleware.NewCustomErr(middleware.NO_REPETITION, "已经认证过了"))
	}
	if file.GetUserFile(data.UserId,
		"doctor_avatar",
		info.DoctorAvatar) == nil {
		panic(middleware.NewCustomErr(middleware.NO_REPETITION, "目标图片不存在"))
	}
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		info = user.CreateDoctorInfo(info, tx)
		tag.CreateTagRelation(info.TagsId, "AREA", int64(info.Id))
		user.CreateRole(userId, "Doctor", tx)
		return nil
	})
}

func GetDoctorInfoByUserId(userId int64) *user2.DoctorInfo {
	return user.GetDoctorInfoById(userId)
}

func IsDoctor(userId int64) bool {
	return GetDoctorInfoByUserId(userId) != nil
}
