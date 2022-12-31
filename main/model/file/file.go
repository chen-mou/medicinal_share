package file

import (
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool"
	"medicinal_share/tool/db/mysql"
)

func GetByHash(hash string) *entity.File {
	file := &entity.File{}
	err := mysql.GetConnect().Where("hash = ?", hash).Take(file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return file
}

func Create(file *entity.File, tx *gorm.DB) *entity.File {
	id, err := tool.GetId("file")
	if err != nil {
		panic(err)
	}
	file.Id = id
	err = tx.Create(file).Error
	if err != nil {
		panic(err)
	}
	return file
}

func CreateData(data *entity.FileData, tx *gorm.DB) *entity.FileData {
	err := tx.Create(data).Error
	if err != nil {
		panic(err)
	}
	return data
}

func UpdateStatus(id int64, status string, tx *gorm.DB) {

}

func GetUserFileByType(userId int64, typ string) []*entity.FileData {
	res := make([]*entity.FileData, 0)
	err := mysql.GetConnect().Model(&entity.FileData{}).
		Where("uploader = ? and type = ?", userId, typ).
		Association("File").Find(&res)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
}

func GetUserFile(userId int64, typ string, fileDataId int64) *entity.FileData {
	res := &entity.FileData{}
	err := mysql.GetConnect().Model(&entity.FileData{}).
		Where("uploader = ? and id = ? and type = ?", userId, fileDataId, typ).
		Joins("Avatar").First(res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return res
}
