package feishu

import (
	"time"
)

type Config struct {
	// Lark: the config of the feishu bot
	Lark Lark
	// GroupNotify: will send the notification in the chat group where the bot is in
	GroupNotify bool
	// api: the apis of feishu
	api *api
}

// Lark : the app must have the ability of bot and must have the permission of "send message to user/group"
// [https://open.feishu.cn] here you can create an app and config it
type Lark struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func (c *Config) Init() error {
	c.api = newApi()
	return c.refreshToken()
}

func (c *Config) refreshToken() error {
	if c.api.expire.Sub(time.Now()) > 5*time.Minute {
		return nil
	}
	return c.api.refreshToken(c.Lark)
}
