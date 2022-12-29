package main

import (
	"flag"
	"medicinal_share/main/resource"
	"medicinal_share/main/route"
)

func main() {
	var mode int

	flag.IntVar(&mode, "mode", 0, "启动模式 0 debug 1 发行模式 2 测试模式")

	flag.Parse()

	resource.Mode = mode

	app := route.Route(mode)

	app.MaxMultipartMemory = 10 << 20

	app.Run(":15888")

}
