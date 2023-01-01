package socket

type code int

const (
	Successfull code = 1 << iota
	Fail
	Error
)
