package wookooncas

import (
	"fmt"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler 创建认证处理器
func (s *CasServer) newAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GetLoginURL 获取登录 URL
// GET /api/cas/login?redirect_uri=xxx
func (h *AuthHandler) GetLoginURL(r *ghttp.Request) {
	redirectURI := r.Get("redirect_uri").String()

	if redirectURI == "" {
		r.SetError(gerror.NewCode(gcode.CodeMissingParameter))
		return
	}

	resp, err := h.authService.GetLoginURL(redirectURI)
	if err != nil {
		r.SetError(gerror.NewCode(gcode.CodeInternalError, err.Error()))
		return
	}

	r.Response.WriteJson(resp)
}

// RedirectLogin 直接重定向到登录页（可选）
// GET /api/cas/login-auto?redirect_uri=xxx
func (h *AuthHandler) RedirectLogin(r *ghttp.Request) {
	redirectURI := r.Get("redirect_uri").String()

	if redirectURI == "" {
		r.SetError(gerror.NewCode(gcode.CodeMissingParameter))
		return
	}

	resp, err := h.authService.GetLoginURL(redirectURI)
	if err != nil {
		r.SetError(gerror.NewCode(gcode.CodeInternalError, err.Error()))
		return
	}

	// c.Redirect(http.StatusFound, resp.LoginURL)
	fmt.Printf("Re: %s\n", resp.LoginURL)
	r.Response.RedirectTo(resp.LoginURL)
}

// HandleCallback 处理登录回调
// POST /api/cas/callback
func (h *AuthHandler) HandleCallback(r *ghttp.Request) {
	var req CallbackRequest

	if err := r.Parse(&req); err != nil {
		r.SetError(err)
		return
	}

	resp, err := h.authService.HandleCallback(&req)
	if err != nil {
		r.SetError(gerror.NewCode(gcode.CodeInternalError, "登录失败: "+err.Error()))
		return
	}

	r.Response.WriteJson(resp)
}

// RedirectProfile 直接重定向到个人资料页
// GET /api/cas/profile-auto
func (h *AuthHandler) RedirectProfile(r *ghttp.Request) {
	resp, err := h.authService.GetMyProfile()
	if err != nil {
		r.SetError(gerror.NewCode(gcode.CodeInvalidRequest, err.Error()))
		return
	}

	r.Response.RedirectTo(resp.ProfileURL)
}

// GetMe 获取当前用户信息
// GET /api/cas/me
func (h *AuthHandler) GetMe(r *ghttp.Request) {
	// 从中间件注入的 Context 中获取用户信息
	user := &StandardUser{
		ID:          r.GetCtxVar(CtxUserId).String(),
		Username:    r.GetCtxVar(CtxUserPassport).String(),
		DisplayName: r.GetCtxVar(CtxUserName).String(),
		Email:       r.GetCtxVar(CtxUserEmail).String(),
		Roles:       r.GetCtxVar(CtxUserRoles).Strings(),
	}

	// c.JSON(http.StatusOK, user)
	r.Response.WriteJson(user)
}

// Logout 登出
// POST /api/cas/logout
func (h *AuthHandler) Logout(r *ghttp.Request) {
	// 由于 JWT 是无状态的，登出主要由前端清除 Token 实现
	// 如果需要服务端登出，可以维护一个 Token 黑名单（Redis）

	r.Response.WriteJson(g.Map{
		"message": "登出成功",
	})
}
