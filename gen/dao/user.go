package dao

import "gorm.io/gen"

type User interface {

	// GetByUserName
	//
	// select password, id from user where username = @username
	GetByUserName(username string) (*gen.T, error)
}
