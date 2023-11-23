package shared

import "strings"

type Notifier interface {
	Notify(*App, Information) error
}

type NotifyBuilder interface {
	Build(Level) Notifier
}

type Information interface {
	Format() string
}
type Option func(*information)

func NewInformation(message string, opts ...Option) Information {
	i := &information{Message: message}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

type information struct {
	Message string  `json:"message"`
	Errors  []error `json:"errors"`
}

func (i *information) Format() string {
	info := strings.Builder{}
	info.WriteString(i.Message)
	for _, err := range i.Errors {
		info.WriteString("\\n") // must be \\n, the '\\' will be escaped to '\' and '\n' will be escaped to '\n'
		info.WriteString(err.Error())
	}
	return info.String()
}

func WithErrors(err ...error) Option {
	return func(i *information) {
		if len(err) == 0 {
			return
		}

		i.Errors = append(i.Errors, err...)
	}
}
