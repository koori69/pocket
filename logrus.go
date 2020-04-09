// Package pocket K·J Create at 2020-04-09 21:30
package pocket

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// ContextHook for log the call context
type ContextHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

// NewContextHook use to make an hook
func NewContextHook(levels ...logrus.Level) logrus.Hook {
	hook := ContextHook{
		Field:  "source",
		Skip:   5,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook ContextHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
