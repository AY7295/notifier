package shared

import (
	"fmt"
	"github.com/AY7295/tsmap"
)

type App struct {
	Name    string
	Mobiles []string
	Emails  []string
}

type Level int

var (
	levels = func() tsmap.TSMap[Level, string] {
		ts := tsmap.New[Level, string]()
		ts.Set(Error, "Error")
		ts.Set(Warning, "Warning")
		ts.Set(Info, "Info")
		ts.Set(Debug, "Debug")
		return ts
	}()
)

func RegisterLevel(level Level, name string) {
	levels.Set(level, name)
}

func (l Level) String() string {
	if s, ok := levels.Get(l); ok {
		return s
	}
	return fmt.Sprintf("Unknown level %d", l)
}

const (
	Error Level = iota
	Warning
	Info
	Debug
)
