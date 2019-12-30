package astilog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type sourceHook struct{}

func (h *sourceHook) Fire(e *logrus.Entry) error {
	// Skip logrus and asticode callers
	i := 0
	_, file, line, ok := runtime.Caller(i)
	for ok && (strings.Contains(file, "asticode/go-astilog") || strings.Contains(file, "sirupsen/logrus")) {
		i++
		_, file, line, ok = runtime.Caller(i)
	}

	// Process file
	if !ok {
		file = "<???>"
		line = 1
	} else {
		file = filepath.Base(file)
	}
	e.Data["source"] = fmt.Sprintf("%s:%d", file, line)
	return nil
}

func (h *sourceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
