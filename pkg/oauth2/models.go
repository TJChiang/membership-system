package oauth2

import "gorm.io/datatypes"

type Client struct {
	Id                     string         `gorm:"primary_key"`
	ClientName             string         `gorm:"not null"`
	ClientSecret           string         `gorm:"not null"`
	Scope                  string         `gorm:"not null"`
	GrantTypes             datatypes.JSON `gorm:"not null,default:'[]'"`
	Audience               datatypes.JSON `gorm:"not null,default:'[]'"`
	PostLogoutRedirectUris datatypes.JSON `gorm:"not null,default:'[]'"`
	BackchannelLogoutUri   string         `gorm:"not null,default:''"`
	RedirectUris           datatypes.JSON `gorm:"not null,default:'[]'"`
	CreatedAt              int64          `gorm:"autoCreateTime"`
	UpdatedAt              int64          `gorm:"autoUpdateTime"`
}
