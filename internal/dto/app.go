package tpl

// AppBody ...
type AppBody struct {
	ID        int64  `json:"-"`          // ID
	AppID     string `json:"app_id"`     // 客户端ID
	AppSecret string `json:"app_secret"` //
}

type AppInfo struct {
	ID                    int64  `json:"-"`                       // ID
	AppID                 string `json:"app_id"`                  // 客户端ID
	AppSecret             string `json:"app_secret"`              // 客户端密钥
	ResourceIDs           string `json:"resource_ids"`            // 资源集合
	Scope                 string `json:"scope"`                   // 授权范围
	AuthorizedGrantTypes  string `json:"authorized_grant_types"`  // 授权类型
	WebServerRedirectURI  string `json:"web_server_redirect_uri"` // 回调地址
	Authorities           string `json:"authorities"`             // 权限
	AccessTokenValidity   int    `json:"access_token_validity"`   // 令牌过期秒数
	RefreshTokenValidity  int    `json:"refresh_token_validity"`  // 刷新令牌过期秒数
	AdditionalInformation string `json:"additional_information"`  // 附件说明
	Autoapprove           string `json:"autoapprove"`             // 自动授权
	CreateUser            string `json:"create_user"`             // 创建人
	UpdateUser            string `json:"update_user"`             // 修改人
	Status                int    `json:"status"`                  // 状态
	IsDeleted             int    `json:"is_deleted"`              // 是否已删除
}

type AppsResponse struct {
	SuccessResponseType
	Data []AppInfo `json:"data"`
}

type AppResponse struct {
	SuccessResponseType
	Data AppInfo `json:"data"`
}
