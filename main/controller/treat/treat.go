package treat

import (
	"github.com/gin-gonic/gin"
	"medicinal_share/main/service/talk"
	"medicinal_share/tool"
	"medicinal_share/tool/socket"
)

func Treat(ctx *gin.Context) {
	type param struct {
		Tags      []int64 `json:"tags"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	p := &param{}
	ctx.BindJSON(p)
	usr := tool.GetNowUser(ctx)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": talk.Treat(usr.Id, p.Tags, p.Latitude, p.Longitude),
	})

}

func Send(conn *socket.Conn, payload string) {
	usr := conn.GetCurrentUser()
	talk.Send(usr.Id, payload, "", conn)
}
