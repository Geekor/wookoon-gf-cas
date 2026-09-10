package wookooncas

// StandardUser 业务系统标准用户模型
type StandardUser struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone,omitempty"`
	Avatar      string   `json:"avatar,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Owner       string   `json:"owner,omitempty"`
}
