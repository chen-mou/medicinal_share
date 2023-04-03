package treat

import "medicinal_share/tool/socket"

func Websocket(manager *socket.ConnManager) {
	manager.Message("/send", Send)
}
