package paseto

import (
	"github.com/hertz-contrib/paseto"
)

type TokenGenerator struct {
	paseto.GenTokenFunc
}

func NewTokenGenerator() (*TokenGenerator, error) {

	return &TokenGenerator{paseto.DefaultGenTokenFunc()}, nil
}

func (g *TokenGenerator) CreateToken(claims *paseto.StandardClaims) (token string, err error) {
	return g.GenTokenFunc(claims, nil, nil)
}
