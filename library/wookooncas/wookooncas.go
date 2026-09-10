package wookooncas

import "github.com/geekor/wookoon-gf-cas/library/eventbus"

type CasServer struct {
	Jwtcfg *JwtConfig
	// 下面是私有变量 ...................
	bus           *eventbus.EventBus
	casdoorClient *CasdoorClient
}

func NewServer(bus *eventbus.EventBus) *CasServer {
	cc := getCasConfigs()
	jc := getJwtConfigs()

	s := &CasServer{
		bus:           bus,
		casdoorClient: NewCasdoorClient(cc),
		Jwtcfg:        jc,
	}

	return s
}
