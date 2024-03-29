package model

import "gorm.io/gorm"

func GetErrorHandler(err error, val any) any {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return any(nil)
		}
		panic(err)
	}
	return val
}
