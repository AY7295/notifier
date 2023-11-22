package notifier

import (
	"errors"
	"github.com/AY7295/notifer/shared"
	"github.com/AY7295/tsmap"
	"sync"
)

var Global *hub

func Init(app shared.App, opts ...Option) {
	Global.Init(app, opts...)
}

type Option func(*hub)

type hub struct {
	sync.Once
	app       *shared.App
	notifiers tsmap.TSMap[shared.Level, []shared.Notifier]
}

func (h *hub) Init(app shared.App, opts ...Option) {
	h.Do(func() {
		h.app = &app
		h.notifiers = tsmap.New[shared.Level, []shared.Notifier]()
		for _, opt := range opts {
			opt(h)
		}
	})
}

func (h *hub) Notify(level shared.Level, info shared.Information) error {
	notifiers, ok := h.notifiers.Get(level)
	if !ok {
		return errors.New("no notifier for level: " + level.String())
	}

	errs := make([]error, 0, len(notifiers))
	for _, notifier := range notifiers {
		errs = append(errs, notifier.Notify(h.app, info))
	}

	return errors.Join(errs...)
}

func WithNotifier(level shared.Level, notifiers []shared.Notifier) Option {
	return func(h *hub) {
		if len(notifiers) != 0 {
			h.notifiers.Set(level, notifiers)
		}
	}
}
