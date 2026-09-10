package wookooncas

import (
	"encoding/json"

	"github.com/geekor/wookoon-gf-cas/library/eventbus"
)

// LoginResponse 登录响应
type LoginResponse struct {
	LoginURL string `json:"loginUrl"`
}

// MyProfileResponse 我的账户响应
type MyProfileResponse struct {
	ProfileURL string `json:"profileUrl"`
}

// CallbackRequest 回调请求
type CallbackRequest struct {
	Code  string `json:"code" v:"required#缺少 code 参数"`
	State string `json:"state"`
}

// CallbackResponse 回调响应
type CallbackResponse struct {
	Token string        `json:"token"`
	User  *StandardUser `json:"user"`
}

// AuthService 认证服务
type AuthService struct {
	core *CasServer
}

// NewAuthService 创建认证服务
func (s *CasServer) newAuthService() *AuthService {
	return &AuthService{
		core: s,
	}
}

// GetLoginURL 获取登录 URL
func (s *AuthService) GetLoginURL(redirectURI string) (*LoginResponse, error) {
	loginURL := s.core.casdoorClient.GetSigninURL(redirectURI)

	return &LoginResponse{LoginURL: loginURL}, nil
}

// HandleCallback 处理登录回调
func (s *AuthService) HandleCallback(req *CallbackRequest) (*CallbackResponse, error) {
	// 1. 用 code 换取用户信息
	user, err := s.core.casdoorClient.ExchangeCodeForUser(req.Code, req.State)
	if err != nil {
		return nil, err
	}

	// 2. (可选) 在本地数据库创建或更新用户记录
	// 这里可以调用 UserRepository.Save(user)
	// 实现用户数据的本地化存储
	buf, _ := json.Marshal(user)
	s.core.bus.Publish(eventbus.Event{
		Name: "user.login",
		Data: buf,
	}) // 广播，所有订阅者收到

	// 3. 签发业务系统 JWT
	token, err := JwtGenerateToken(
		&TokenGenParams{
			UserID:      user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Roles:       user.Roles,
			JWTSecret:   s.core.Jwtcfg.Secret,
			JWTExpire:   s.core.Jwtcfg.Expire,
			JWTIssuer:   s.core.Jwtcfg.Issuer,
		},
	)
	if err != nil {
		return nil, err
	}

	return &CallbackResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetLoginURL 获取登录 URL
func (s *AuthService) GetMyProfile() (*MyProfileResponse, error) {
	url := s.core.casdoorClient.GetMyProfileURL()
	return &MyProfileResponse{
		ProfileURL: url,
	}, nil
}

// GetUserFromToken 从 Token 中解析用户信息
func (s *AuthService) GetUserFromToken(tokenString string) (*StandardUser, error) {

	claims, err := JwtParseToken(tokenString, s.core.Jwtcfg.Secret)
	if err != nil {
		return nil, err
	}

	return &StandardUser{
		ID:          claims.UserID,
		Username:    claims.Username,
		DisplayName: claims.DisplayName,
		Email:       claims.Email,
		Roles:       claims.Roles,
	}, nil
}
