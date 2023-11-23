package feishu

import (
	"encoding/json"
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

func (b *builder) Build(level shared.Level) shared.Notifier {
	var (
		userIds []string
		chatIds []string
		err     error
	)
	return notify(func(app *shared.App, info shared.Information) error {
		err = b.config.refreshToken()
		if err != nil {
			return err
		}

		userIds, err = b.config.api.getOpenIds(app.Mobiles)
		if err != nil {
			return err
		}

		if b.config.GroupNotify {
			chatIds, err = b.config.api.getChatIds()
			if err != nil {
				return err
			}
		}

		return b.config.api.sendCard(
			newCard(getTemplate(level), app.Name, level.String(), info.Format(), builderAt(userIds...)),
			userIds, chatIds...,
		)
	})
}

func builderAt(ids ...string) string {
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

func newCard(templateColor, appName, level, content, ats string) string {
	return fmt.Sprintf(cardFormat, templateColor, appName, level, content, ats)
}

// init: just to escape the useless characters
func init() {
	var temp any
	_ = json.Unmarshal([]byte(cardFormat), &temp)
	bytes, _ := json.Marshal(temp)
	cardFormat = string(bytes)
}

var (
	cardFormat = `
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
