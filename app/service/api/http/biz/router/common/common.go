package common

import (
	gpaseto "aidanwoods.dev/go-paseto"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/paseto"
	"go-ssip/app/common/errno"
	"go-ssip/app/common/middleware"
	"go-ssip/app/common/tools"
	"net/http"
)

func CommonMW() []app.HandlerFunc {
	return []app.HandlerFunc{
		// use cors mw
		middleware.Cors(),
		// use recovery mw
		//middleware.Recovery(),
		// use gzip mw
		gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".mp4", ".png"})),
	}
}

func PasetoAuth() app.HandlerFunc {
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
	return paseto.New(paseto.WithParseFunc(pf), paseto.WithSuccessHandler(sh), paseto.WithErrorFunc(ef))
}
