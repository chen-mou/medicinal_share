package main

import (
	"flag"
	"medicinal_share/main/resource"
	"medicinal_share/main/route"
	"medicinal_share/tool/socket"
)

func ws() {
	cm := socket.NewConnManager("localhost:15889")

	cm.HeaderHandler("Token", func(conn *socket.Conn, s string) error {
		conn.SetAuth(s)
		return nil
	})
}

func main() {
	var mode int

	flag.IntVar(&mode, "mode", 0, "启动模式 0 debug 1 发行模式 2 测试模式")

	flag.Parse()

	resource.Mode = mode

	app := route.Route(mode)

	app.MaxMultipartMemory = 10 << 20

	ws()

	app.Run(":15888")

}
