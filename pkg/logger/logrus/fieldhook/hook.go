package fieldhook

import (
	"github.com/sirupsen/logrus"
)

type DefaultFieldHook struct {
	Fields logrus.Fields
}

func (h *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *DefaultFieldHook) Fire(e *logrus.Entry) error {
	e.Data = h.Fields
	return nil
}
