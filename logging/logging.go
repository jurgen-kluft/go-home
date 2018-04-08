package logging

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log     *logrus.Logger
	context map[string]*logrus.Entry
}

func New() *Logger {
	logger := &Logger{}
	logger.log = logrus.New()

	return logger
}

func (log *Logger) AddEntry(context string) {
	log.context[context] = logrus.WithFields(logrus.Fields{
		"category": context,
	})
}

func (log *Logger) LogInfo(context string, line string) {
	entry := log.context[context]
	entry.Info(line)
}

func (log *Logger) LogError(context string, line string) {
	entry := log.context[context]
	entry.Error(line)
}
