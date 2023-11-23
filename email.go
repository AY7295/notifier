package notifier

import (
	"github.com/AY7295/notifer/pkg/email"
	"github.com/AY7295/notifer/shared"
)

func NewEmailBuilder(config email.Config) (shared.NotifyBuilder, error) {
	return email.NewNotifyBuilder(config)
}

func NewEmailConfig() email.Config {
	return email.Config{}
}
