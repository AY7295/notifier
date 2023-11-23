package notifier

import (
	"github.com/AY7295/notifer/pkg/feishu"
	"github.com/AY7295/notifer/shared"
)

func NewFeishuBuilder(config feishu.Config) (shared.NotifyBuilder, error) {
	return feishu.NewNotifyBuilder(config)
}

func NewFeishuConfig(appId, appSecret string, groupNotify bool) feishu.Config {
	return feishu.Config{
		Lark: feishu.Lark{
			AppId:     appId,
			AppSecret: appSecret,
		},
		GroupNotify: groupNotify,
	}
}
