package socket

type code int

const (
	Successful code = 1 << iota
	Fail
	Error
)
