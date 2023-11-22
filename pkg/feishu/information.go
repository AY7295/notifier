package feishu

import (
	"fmt"
	"github.com/AY7295/notifer/shared"
	"strings"
)

type notify func(*shared.App, shared.Information) error

func (n notify) Notify(app *shared.App, information shared.Information) error {
	return n(app, information)
}

func NewNotifyBuilder(config Config) (shared.NotifyBuilder, error) {
	err := config.Init()
	if err != nil {
		return nil, err
	}
	return &builder{&config}, nil
}

type builder struct {
	config *Config
}

func (b *builder) Build(level shared.Level) (shared.Notifier, error) {
	return notify(func(app *shared.App, info shared.Information) error {
		err := b.config.refreshToken()
		if err != nil {
			return err
		}

		ids, err := b.config.api.GetOpenIds(app.Notify.Phones, b.config.tenantAccessToken)
		if err != nil {
			return err
		}

		b.config.api.SendCard()

		return nil
	}), nil
}

func atBuilder(ids ...string) string {
	ats := make([]string, 0, len(ids))
	for _, id := range ids {
		ats = append(ats, fmt.Sprintf("<at id=%s><at/>", id))
	}
	return strings.Join(ats, " ")
}

func getTemplate(level shared.Level) string {
	switch level {
	case shared.Error:
		return "red"
	case shared.Warning:
		return "yellow"
	case shared.Info:
		return "green"
	case shared.Debug:
		return "gray"
	default:
		return "blue"
	}
}

const (
	card = `
			{
			  "header": {
			    "template": "%s",
			    "title": {
			      "content": "%s %s"
			    }
			  },
			  "elements": [
			    {
			      "tag": "div",
			      "text": {
			        "content": "%s \n%s",
			        "tag": "lark_md"
			      }
			    }
			  ]
			}
			`
)
