package notifier

import (
	"errors"
	"github.com/AY7295/notifer/model"
	"github.com/AY7295/tsmap"
	"sync"
)

var Global *hub

func Init(opts ...Option) {
	Global.Init(opts...)
}

type Option func(*hub)

type hub struct {
	sync.Once
	notifiers tsmap.TSMap[model.Level, []model.Notifier]
}

func (h *hub) Init(opts ...Option) {
	h.Do(func() {
		h.notifiers = tsmap.New[model.Level, []model.Notifier]()
		for _, opt := range opts {
			opt(h)
		}
	})
}

func (h *hub) Notify(level model.Level, info model.Information) error {
	notifiers, ok := h.notifiers.Get(level)
	if !ok {
		return errors.New("no notifier for level " + level.String())
	}

	errs := make([]error, 0, len(notifiers))
	for _, notifier := range notifiers {
		errs = append(errs, notifier.Notify(info))
	}

	return errors.Join(errs...)
}

func WithNotifier(level model.Level, notifiers []model.Notifier) Option {
	return func(h *hub) {
		if len(notifiers) != 0 {
			h.notifiers.Set(level, notifiers)
		}
	}
}
