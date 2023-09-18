package main

import (
	gpaseto "aidanwoods.dev/go-paseto"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/paseto"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/websocket"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/api/ws/global"
	"go-ssip/app/service/api/ws/initialize"
	"go-ssip/app/service/api/ws/initialize/rpc"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)

var hub *Hub

var upgrader = websocket.HertzUpgrader{
	// TODO: custom option
}

func main() {
	initialize.InitLogger(consts.WsApiName)
	initialize.InitConfig()
	initialize.InitRdb()
	r, info := initialize.InitRegistry()
	tracer, trcCfg := hertztracing.NewServerTracer()
	rpc.Init()
	conn := initialize.InitMq()
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		g.Logger.Fatal("declare a channel failed", zap.Error(err))
	}
	defer ch.Close()

	// 声明一个队列
	queue, err := ch.QueueDeclare(
		net.JoinHostPort(g.ServerConfig.Host, strconv.Itoa(g.ServerConfig.Port)),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		g.Logger.Fatal("declare a queue failed", zap.Error(err))
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		g.Logger.Fatal("consume msg failed", zap.Error(err))
	}

	sh := func(ctx context.Context, c *app.RequestContext, token *gpaseto.Token) {
		id, err := token.GetString("id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, tools.BuildBaseResp(errno.BadRequest.WithMessage("missing id in token")))
			c.Abort()
		}
		c.Set("ID", id)
	}
	pf, _ := paseto.NewV4PublicParseFunc(paseto.DefaultPublicKey, []byte(paseto.DefaultImplicit))
	ef := func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusUnauthorized, tools.BuildBaseResp(errno.BadRequest.WithMessage("invalid token")))
		c.Abort()
	}
	hub = newHub(msgs)
	go hub.run()

	h := server.New(
		tracer,
		server.WithALPN(true),
		server.WithHostPorts(fmt.Sprintf(":%d", g.ServerConfig.Port)),
		server.WithRegistry(r, info),
	)
	h.NoHijackConnPool = true
	pprof.Register(h)
	h.Use(hertztracing.ServerMiddleware(trcCfg))
	h.GET("/echo", paseto.New(paseto.WithParseFunc(pf), paseto.WithSuccessHandler(sh), paseto.WithErrorFunc(ef)), func(c context.Context, ctx *app.RequestContext) {
		id, _ := ctx.Get("ID")
		ctx.JSON(http.StatusOK, id)
	})
	h.GET("/", paseto.New(paseto.WithParseFunc(pf), paseto.WithSuccessHandler(sh), paseto.WithErrorFunc(ef)), serveWs)
	h.Spin()
}
