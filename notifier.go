package notifier

import (
	"errors"
	"github.com/AY7295/notifer/shared"
	"github.com/AY7295/tsmap"
	"reflect"
	"sync"
)

// Global is the only global hub instance
var Global *hub

func Init(app shared.App, opts ...Option) {
	Global.Init(app, opts...)
}

type Option func(*hub)

type hub struct {
	once      sync.Once
	app       *shared.App
	notifiers tsmap.TSMap[shared.Level, []shared.Notifier]
}

func (h *hub) Init(app shared.App, opts ...Option) {
	h.once.Do(func() {
		h.app = &app
		h.notifiers = tsmap.New[shared.Level, []shared.Notifier]()
		for _, opt := range opts {
			opt(h)
		}
	})
}

// Notify : call all notifiers set for level
func (h *hub) Notify(level shared.Level, info shared.Information) error {
	if h.app == nil {
		panic("notifier not init") // must init first
	}
	if reflect.ValueOf(info).IsNil() {
		return errors.New("info must not be nil")
	}

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

// WithNotifier : set notifiers for level
func WithNotifier(level shared.Level, notifiers ...shared.Notifier) Option {
	return func(h *hub) {
		if len(notifiers) != 0 {
			h.notifiers.Set(level, notifiers)
		}
	}
}
