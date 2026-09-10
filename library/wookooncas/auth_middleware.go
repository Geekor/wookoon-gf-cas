package wookooncas

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	CtxUserId       = "user_id"
	CtxUserPassport = "user_passport"
	CtxUserName     = "user_name"
	CtxUserEmail    = "user_email"
	CtxUserRoles    = "user_roles"
)

// JWT 验证中间件
func CasAuthed(jwtCfg *JwtConfig) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "未提供认证信息"))
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			r.SetError(gerror.NewCode(gcode.CodeValidationFailed, "认证格式错误"))
			return
		}

		tokenString := parts[1]

		// 验证 Token
		claims, err := JwtParseToken(tokenString, jwtCfg.Secret)
		fmt.Println(jwtCfg)
		if err != nil {
			message := "认证失败"
			if err == ErrTokenExpired {
				message = "Token 已过期"
			}
			r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, message))
			return
		}

		// 将用户信息注入 Context，供后续 Handler 使用
		r.SetCtxVar(CtxUserId, claims.UserID)
		r.SetCtxVar(CtxUserName, claims.DisplayName)
		r.SetCtxVar(CtxUserPassport, claims.Username)
		r.SetCtxVar(CtxUserEmail, claims.Email)
		r.SetCtxVar(CtxUserRoles, claims.Roles)

		// r.SetCtxVar("tokenString", tokenString)
		// r.SetCtxVar("userId", claims.UserID)
		// r.SetCtxVar("username", claims.Username)
		// r.SetCtxVar("displayName", claims.DisplayName)
		// r.SetCtxVar("email", claims.Email)
		// r.SetCtxVar("roles", claims.Roles)
		// r.SetCtxVar("claims", claims)

		r.Middleware.Next()
	}
}
