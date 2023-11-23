package email

import "github.com/AY7295/notifer/shared"

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
	//TODO implement me
	panic("implement me")
}
