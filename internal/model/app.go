package model

import (
	"gorm.io/gorm"
)

// App 应用表
type App struct {
	gorm.Model
	AppID                 int64  `gorm:"column:app_id;type:bigint(64);not null;default:0" json:"app_id"`                                    // 客户端ID
	AppSecret             string `gorm:"column:app_secret;type:varchar(255);not null;default:''" json:"app_secret"`                         // 客户端密钥
	ResourceIDs           string `gorm:"column:resource_ids;type:varchar(255);default:''" json:"resource_ids"`                              // 资源集合
	Scope                 string `gorm:"column:scope;type:varchar(255);not null;default:''" json:"scope"`                                   // 授权范围
	AuthorizedGrantTypes  string `gorm:"column:authorized_grant_types;type:varchar(255);not null;default:''" json:"authorized_grant_types"` // 授权类型
	WebServerRedirectURI  string `gorm:"column:web_server_redirect_uri;type:varchar(255);default:''" json:"web_server_redirect_uri"`        // 回调地址
	Authorities           string `gorm:"column:authorities;type:varchar(255);default:''" json:"authorities"`                                // 权限
	AccessTokenValidity   int    `gorm:"column:access_token_validity;type:int(11);not null;default:0" json:"access_token_validity"`         // 令牌过期秒数
	RefreshTokenValidity  int    `gorm:"column:refresh_token_validity;type:int(11);not null;default:0" json:"refresh_token_validity"`       // 刷新令牌过期秒数
	AdditionalInformation string `gorm:"column:additional_information;type:varchar(4096);default:''" json:"additional_information"`         // 附件说明
	Autoapprove           string `gorm:"column:autoapprove;type:varchar(255);default:''" json:"autoapprove"`                                // 自动授权
	CreateUser            string `gorm:"column:create_user;type:varchar(64);default:''" json:"create_user"`                                 // 创建人
	UpdateUser            string `gorm:"column:update_user;type:varchar(64);default:''" json:"update_user"`                                 // 修改人
	Status                int    `gorm:"column:status;type:int(2);not null;default:0" json:"status"`                                        // 状态
	IsDeleted             int    `gorm:"column:is_deleted;type:int(2);not null;default:0" json:"is_deleted"`                                // 是否已删除
}

// TableName, 获取应用表名称
func (m *App) TableName() string {
	return "sys_app"
}
