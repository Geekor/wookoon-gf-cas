package wookooncas

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenGenParams struct {
	UserID      string
	Username    string
	DisplayName string
	Email       string
	Roles       []string
	JWTSecret   string
	JWTExpire   time.Duration
	JWTIssuer   string
}
type RefreshTokenGenParams struct {
	UserID    string
	JWTSecret string
	JWTExpire time.Duration
	JWTIssuer string
}

// CustomClaims 自定义 JWT Claims
type CustomClaims struct {
	UserID      string   `json:"userId"`
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired = errors.New("token 已过期")
	ErrTokenInvalid = errors.New("token 无效")
)

// 生成 刷新TOKEN
func JwtGenerateRefreshToken(param *RefreshTokenGenParams) (string, error) {
	claims := CustomClaims{
		UserID: param.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(param.JWTExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    param.JWTIssuer,
			Subject:   param.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(param.JWTSecret))
}

// 签发业务系统 JWT
func JwtGenerateToken(param *TokenGenParams) (string, error) {

	claims := CustomClaims{
		UserID:      param.UserID,
		Username:    param.Username,
		DisplayName: param.DisplayName,
		Email:       param.Email,
		Roles:       param.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(param.JWTExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    param.JWTIssuer,
			Subject:   param.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(param.JWTSecret))
}

// ParseToken 解析并验证 JWT
func JwtParseToken(tokenString string, secret string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// （忽略过期错误）
func JwtParseClaims(tokenString string, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenInvalid
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
