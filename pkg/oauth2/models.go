package oauth2

import "gorm.io/datatypes"

type Client struct {
	Id                     string         `gorm:"primary_key" json:"client_id"`
	ClientName             string         `gorm:"not null" json:"client_name"`
	ClientSecret           string         `gorm:"not null" json:"-"`
	Scope                  string         `gorm:"not null" json:"scope"`
	GrantTypes             datatypes.JSON `gorm:"not null;default:'[]'" json:"grant_types"`
	Audience               datatypes.JSON `gorm:"not null;default:'[]'" json:"audience"`
	PostLogoutRedirectUris datatypes.JSON `gorm:"not null" json:"post_logout_redirect_uris"`
	BackchannelLogoutUri   string         `gorm:"not null;default:''" json:"backchannel_logout_uri"`
	RedirectUris           datatypes.JSON `gorm:"not null;default:'[]'" json:"redirect_uris"`
	CreatedAt              int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              int64          `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Client) TableName() string {
	return "oauth2_client"
}
