package feishu

import (
	"time"
)

type Config struct {
	// Lark: the config of the feishu bot
	Lark Lark
	// NeedNotifyInGroup: will send the notification in the chat group where the bot is in
	NeedNotifyInGroup bool
	// api: the apis of feishu
	api *api
}

type Lark struct {
	ID     string `json:"app_id"`
	Secret string `json:"app_secret"`
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
