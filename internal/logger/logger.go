package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	Level        logrus.Level
	logrusLogger *logrus.Logger
}

func New(level string, loggerOut io.Writer) *Logger {
	log := logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	loggerLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatal(err)
	}
	log.Level = loggerLevel
	log.Out = loggerOut

	return &Logger{
		Level:        loggerLevel,
		logrusLogger: log,
	}
}

func (l *Logger) getEntry() *logrus.Entry {
	return l.logrusLogger.WithFields(logrus.Fields{})
}

func (l *Logger) Info(msg string) {
	l.logrusLogger.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.getEntry().Error(msg)
}

func (l *Logger) Fatal(msg string) {
	l.logrusLogger.Fatal(msg)
}
