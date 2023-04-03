package treat

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/service/talk"
	"medicinal_share/tool/socket"
)

func Treat(ctx *gin.Context) {
	type param struct {
		Tags   []string
		InfoId int64
	}
}

func Send(conn *socket.Conn, payload string) {
	usr := conn.GetCurrentUser()
	talk.Send(usr.Id, payload, "", conn)
}
