package wookooncas

import (
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

// Client Casdoor 客户端
type CasdoorClient struct {
	sdkClient *casdoorsdk.Client
}

// NewClient 创建 Casdoor 客户端
func NewCasdoorClient(cfg *CasAuthConfig) *CasdoorClient {
	client := casdoorsdk.NewClient(
		cfg.Endpoint,
		cfg.ClientID,
		cfg.ClientSecret,
		cfg.Certificate,
		cfg.Organization,
		cfg.Application,
	)

	return &CasdoorClient{
		sdkClient: client,
	}
}

// GetSigninURL 获取 Casdoor 登录 URL
func (c *CasdoorClient) GetSigninURL(redirectURI string) string {
	return c.sdkClient.GetSigninUrl(redirectURI)
}

// GetSigninURL 获取 Casdoor 登录 URL
func (c *CasdoorClient) GetMyProfileURL() string {
	return c.sdkClient.GetMyProfileUrl("")
}

// ExchangeCodeForUser 使用 code 换取用户信息
func (c *CasdoorClient) ExchangeCodeForUser(code, state string) (*StandardUser, error) {
	// 1. 用 code 换取 Casdoor 的 OAuth Token
	token, err := c.sdkClient.GetOAuthToken(code, state)
	if err != nil {
		return nil, fmt.Errorf("获取 Casdoor Token 失败: %w", err)
	}

	// 2. 解析 Access Token 获取用户 Claims
	claims, err := c.sdkClient.ParseJwtToken(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("解析 Casdoor Token 失败: %w", err)
	}

	// 3. 转换为标准用户模型（防腐层）
	user := adaptFromCasdoor(claims)
	if user == nil {
		return nil, fmt.Errorf("用户信息为空")
	}

	return user, nil
}

// AdaptFromCasdoor 将 Casdoor Claims 转换为标准用户模型
func adaptFromCasdoor(claims *casdoorsdk.Claims) *StandardUser {
	if claims == nil {
		return nil
	}

	// 转换 Roles: []*casdoorsdk.Role -> []string
	roles := make([]string, 0, len(claims.Roles))
	for _, role := range claims.Roles {
		if role != nil {
			roles = append(roles, role.Name)
		}
	}

	return &StandardUser{
		ID:          claims.Id,
		Username:    claims.Name,
		DisplayName: claims.DisplayName,
		Email:       claims.Email,
		Phone:       claims.Phone,
		Avatar:      claims.Avatar,
		Roles:       roles, // 使用转换后的 []string
		Owner:       claims.Owner,
	}
}
