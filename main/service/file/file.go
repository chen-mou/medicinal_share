package file

import (
	"gorm.io/gorm"
	"io/ioutil"
	"medicinal_share/main/entity"
	"medicinal_share/main/model/file"
	"medicinal_share/main/resource"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/encrypt/md5"
	"mime/multipart"
	"os"
	"strings"
)

const BasePath = "workplace/img/"

var Suffix = map[string]struct{}{
	"jpg":  {},
	"png":  {},
	"bmp":  {},
	"jpeg": {},
}

func init() {
	info, err := os.Stat(BasePath)
	if err != nil || !info.IsDir() {
		err = os.Mkdir("workplace", os.ModeDir)
		os.Mkdir(BasePath, os.ModeDir)
	}
}

func Upload(f *multipart.FileHeader, typ string, uploader int64, callback func(int64, *gorm.DB) error) {
	file1, _ := f.Open()
	byt, err := ioutil.ReadAll(file1)
	if err != nil {
		panic(err)
	}
	v := md5.Hash(string(byt))
	fe := file.GetByHash(v)
	filenames := strings.Split(f.Filename, ".")
	fd := &entity.FileData{
		Name:     filenames[0],
		Suffix:   filenames[1],
		Type:     typ,
		Uploader: uploader,
		Status:   "NORMAL",
	}
	if fe != nil {
		fd.FileId = fe.Id
		mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
			fd = file.CreateData(fd, tx)
			callback(fe.Id, tx)
			return nil
		})
		return
	}
	save(fd, byt, v, callback)
	return
}

func save(fd *entity.FileData, byt []byte, v string, callback func(int64, *gorm.DB) error) {
	file2, _ := os.Create(BasePath + v)
	fe := &entity.File{
		Machine: resource.Machine,
		Status:  "NORMAL",
		Uri:     resource.UriPre + "/" + resource.Machine + "/file/get/" + v,
		Path:    BasePath + v,
		Hash:    v,
	}
	err := mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		fe = file.Create(fe, tx)
		file2.Write(byt)
		return nil
	})
	if err != nil {
		panic(err)
	}
	mysql.GetConnect().Transaction(func(tx *gorm.DB) error {
		fd.FileId = fe.Id
		file.CreateData(fd, tx)
		return callback(fd.Id, tx)
	})
}

func GetByHash(hash string) *entity.File {
	return file.GetByHash(hash)
}

func GetUserFileByType(userId int64, typ string) []*entity.FileData {
	return GetUserFileByType(userId, typ)
}
