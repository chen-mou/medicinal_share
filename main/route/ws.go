package route

import (
	"medicinal_share/main/controller/treat"
	"medicinal_share/tool/socket"
)

func Websocket() {
	cm := socket.NewConnManager("localhost:15777")

	treat.Websocket(cm)

	cm.Run()
}
