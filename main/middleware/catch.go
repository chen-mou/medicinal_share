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

type NoRepetition struct {
	msg    string
	status int
}

func (NoRepetition) GetStatus() int {
	return NO_REPETITION
}

func (n NoRepetition) Error() string {
	return n.msg
}

func NewCustomErr(typ int, msg string) CustomErr {
	switch typ {
	case NO_REPETITION:
		return NoRepetition{
			msg: msg,
		}
	case NOT_FOUND:
	case ERROR:
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
		default:
			ctx.AbortWithStatusJSON(500, err)
		}
	}()

	ctx.Next()
}
