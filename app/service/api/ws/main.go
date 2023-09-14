package ws

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
)

var upgrader = websocket.HertzUpgrader{
	// TODO: custom option
}

func serve(_ context.Context, c *app.RequestContext) {

}

func main() {

}
