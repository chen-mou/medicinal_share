package dao

import "gorm.io/gen"

type Filter interface {
	// FilterById
	//
	// select * from @@table where id = @id
	FilterById(id int64) (*gen.T, error)
}
