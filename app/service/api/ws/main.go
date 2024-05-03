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
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/api/ws/global"
	"go-ssip/app/service/api/ws/initialize"
	"go-ssip/app/service/api/ws/initialize/rpc"
	"net/http"
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
	consumer, topic := initialize.InitMq()
	defer consumer.Close()

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(g.ServerConfig.Name),
		provider.WithExportEndpoint(g.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

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

	hub = newHub(consumer.Messages(), topic)
	defer func() {
		for _, c := range hub.clients {
			hub.unregister <- c
		}
	}()
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
