package middleware

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go-ssip/app/common/errno"
)

func Recovery() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(
		func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"code":    errno.BadRequest,
				"message": fmt.Sprintf("[Recovery] Panic"),
			})
		},
	))
}
