package logging

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log     *logrus.Logger
	process string
	context map[string]*logrus.Entry
}

func New(process string) *Logger {
	logger := &Logger{}
	logger.log = logrus.New()
	logger.log.Formatter.(*logrus.TextFormatter).DisableColors = false
	logger.log.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	logger.log.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	logger.log.Formatter.(*logrus.TextFormatter).TimestampFormat = "Mon 15:04:05"
	logger.process = process
	logger.context = map[string]*logrus.Entry{}
	return logger
}

func (log *Logger) AddEntry(context string) {
	logentry := logrus.WithFields(logrus.Fields{
		"process":  log.process,
		"category": context,
	})

	logentry.Logger.Formatter.(*logrus.TextFormatter).DisableColors = false
	logentry.Logger.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	logentry.Logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	logentry.Logger.Formatter.(*logrus.TextFormatter).TimestampFormat = "Mon 15:04:05"

	log.context[context] = logentry
}

func (log *Logger) LogInfo(context string, line string) {
	entry := log.context[context]
	entry.Info(line)
}

func (log *Logger) LogError(context string, line string) {
	entry := log.context[context]
	entry.Error(line)
}
