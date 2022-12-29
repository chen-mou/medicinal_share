package dao

import "gorm.io/gen"

type File interface {

	//GetFileByHash
	//
	// select * from file where hash = @hash
	GetFileByHash(hash string) *gen.T
}

type FileData interface {

	//GetFileDataByType
	//
	// select * from file_data where uploader = @uploader and type = @typ
	GetFileDataByType(uploader int64, typ string)
}
