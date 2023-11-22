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
	// tenantAccessToken: the access-token of the lark-api, need to be refreshed every 2 hours
	tenantAccessToken string
	// expire: the expiring time of the tenantAccessToken
	expire time.Time
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
	if c.expire.Sub(time.Now()) > 5*time.Minute {
		return nil
	}
	token, expire, err := c.api.GetToken(c.Lark)
	if err != nil {
		return err
	}

	c.tenantAccessToken = "Bearer " + token
	c.expire = time.Now().Add(expire)

	return nil
}
