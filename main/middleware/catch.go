package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	NOT_FOUND int = iota
	NO_REPETITION
	FORBID
	ERROR
)

type CustomErr interface {
	error
	GetStatus() int
}

type noRepetition struct {
	msg    string
	status int
}

type errors struct {
	msg string
}

type notFound struct {
	msg string
}

func (notFound) GetStatus() int {
	return 404
}

func (n notFound) Error() string {
	return n.msg
}

func (noRepetition) GetStatus() int {
	return NO_REPETITION
}

func (n noRepetition) Error() string {
	return n.msg
}

func (e errors) Error() string {
	return e.msg
}

func (errors) GetStatus() int {
	return 500
}

func NewCustomErr(typ int, msg string) CustomErr {
	switch typ {
	case NO_REPETITION:
		return noRepetition{
			msg: msg,
		}
	case NOT_FOUND:
		return notFound{
			msg: msg,
		}
	case FORBID:
	case ERROR:
		return errors{
			msg: msg,
		}
	}
	return nil
}

func Catch(ctx *gin.Context) {
	defer func() {
		err := recover()
		switch err.(type) {
		//处理未检查异常
		case error:
			ctx.AbortWithError(500, err.(error))
		case string:
			ctx.AbortWithStatusJSON(500, gin.H{
				"error_msg": err,
			})
			//处理业务异常
		case CustomErr:
			val := err.(CustomErr)
			ctx.AbortWithStatusJSON(val.GetStatus(), val.Error())
		case nil:
			return
		default:
			ctx.AbortWithStatusJSON(500, err)
		}
	}()

	ctx.Next()
}
