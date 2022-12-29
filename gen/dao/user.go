package dao

import "gorm.io/gen"

type User interface {

	// GetByUserName
	//
	// select * from user where username = @username
	GetByUserName(username string) (*gen.T, error)
}
