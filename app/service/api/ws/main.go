package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/websocket"
	"go-ssip/app/common/consts"
	"go-ssip/app/service/api/ws/initialize"
	"net"
	"strconv"
)

var upgrader = websocket.HertzUpgrader{
	// TODO: custom option
}

func main() {
	initialize.InitLogger(consts.WsApiName)
	initialize.InitConfig()
	IP, Port := initialize.InitFlag()
	initialize.InitRegistry(Port)

	go hub.run()

	h := server.Default(server.WithHostPorts(net.JoinHostPort(IP, strconv.Itoa(Port))))

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		serveWs(ctx)
	})
	h.Spin()
}
