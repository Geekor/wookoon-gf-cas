package wookooncas

import (
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

		// 获取 auth 信息
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
		if err == ErrTokenExpired {
			rtoken := r.Header.Get("Wk-Refresh")
			if rtoken != "" {
				_, err2 := JwtParseToken(rtoken, jwtCfg.Secret)
				if err2 == ErrTokenExpired {
					r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "Token 已过期"))
					return
				} else if err2 != nil {
					r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "认证失败"))
					return
				}

				// 重新签发 token
				claims, err = JwtParseClaims(tokenString, jwtCfg.Secret)
				if err == nil {
					ntoken, err3 := JwtGenerateToken(
						&TokenGenParams{
							UserID:      claims.UserID,
							Username:    claims.Username,
							DisplayName: claims.DisplayName,
							Email:       claims.Email,
							Roles:       claims.Roles,
							JWTSecret:   jwtCfg.Secret,
							JWTExpire:   jwtCfg.Expire,
							JWTIssuer:   jwtCfg.Issuer,
						},
					)
					if err3 != nil {
						r.Response.Header().Set("Wk-Token", ntoken)
					}
				}
			}
		}

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

		r.Middleware.Next()
	}
}
