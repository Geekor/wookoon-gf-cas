package wookooncas

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func (s *CasServer) RegisterRoutes(g *ghttp.RouterGroup) {
	as := s.newAuthService()
	hd := s.newAuthHandler(as)

	g.Group("/cas", func(gs1 *ghttp.RouterGroup) {
		// 公开接口
		gs1.GET("/login", hd.GetLoginURL)
		gs1.GET("/login-auto", hd.RedirectLogin)
		gs1.POST("/callback", hd.HandleCallback)
		gs1.GET("/me-auto", hd.RedirectProfile)

		// 需要认证的接口
		gs1.Group("/", func(auth *ghttp.RouterGroup) {
			auth.Middleware(CasAuthed(s.Jwtcfg))

			auth.GET("/me", hd.GetMe)
			auth.POST("/logout", hd.Logout)
		})
	})
}
