package wookooncas

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

// CasAuthConfig is the core configuration.
// The first step to use this SDK is to use the InitConfig function to initialize the global authConfig.
type CasAuthConfig struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	Organization string
	Application  string
	Certificate  string
}

type JwtConfig struct {
	Secret string
	Expire time.Duration
	Issuer string
}

func getCasConfigs() *CasAuthConfig {
	ctx := gctx.New()

	certf := g.Cfg().MustGet(ctx, "wookooncas.certificateFile").String()
	certc := gfile.GetContents(certf)

	c := &CasAuthConfig{
		Endpoint:     g.Cfg().MustGet(ctx, "wookooncas.endpoint").String(),
		ClientID:     g.Cfg().MustGet(ctx, "wookooncas.clientId").String(),
		ClientSecret: g.Cfg().MustGet(ctx, "wookooncas.clientSecret").String(),
		Organization: g.Cfg().MustGet(ctx, "wookooncas.organization").String(),
		Application:  g.Cfg().MustGet(ctx, "wookooncas.application").String(),
		Certificate:  certc,
	}

	return c
}

func getJwtConfigs() *JwtConfig {
	ctx := gctx.New()

	return &JwtConfig{
		Secret: g.Cfg().MustGet(ctx, "wookooncas.jwt.secret").String(),
		Expire: g.Cfg().MustGet(ctx, "wookooncas.jwt.expires").Duration() * time.Hour,
		Issuer: g.Cfg().MustGet(ctx, "wookooncas.jwt.issuer").String(),
	}
}
