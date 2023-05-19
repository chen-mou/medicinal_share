package route

import (
	"medicinal_share/main/controller/treat"
	"medicinal_share/tool/socket"
)

func Websocket() {
	cm := socket.NewConnManager("localhost:15777")

	//cm.HeaderHandler("Sec-WebSocket-Protocol", func(conn *socket.Conn, s string) error {
	//	conn.SetAuth(s)
	//	return nil
	//})

	treat.Websocket(cm)

	cm.Run()
}
