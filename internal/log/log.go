package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Level = string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

type Fields = logrus.Fields

type Logger interface {
	WithErr(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Config struct {
	Out   io.Writer
	Level Level
}

type wrappedLogger struct {
	entry *logrus.Entry
}

func NewDiscard() Logger {
	return NewLogger(Config{Out: io.Discard})
}

func NewLogger(config Config) Logger {
	if config.Out == nil {
		config.Out = os.Stderr
	}
	if config.Level == "" {
		config.Level = LevelDebug
	}

	l, err := logrus.ParseLevel(config.Level)
	if err != nil {
		panic("invalid log level: " + config.Level)
	}

	logger := logrus.New()
	logger.SetLevel(l)
	logger.SetOutput(config.Out)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return &wrappedLogger{entry: logrus.NewEntry(logger)}
}

func (l *wrappedLogger) WithErr(err error) Logger {
	return &wrappedLogger{entry: l.entry.WithError(err)}
}

func (l *wrappedLogger) WithField(key string, value interface{}) Logger {
	return &wrappedLogger{entry: l.entry.WithField(key, value)}
}

func (l *wrappedLogger) WithFields(fields Fields) Logger {
	return &wrappedLogger{entry: l.entry.WithFields(fields)}
}

func (l *wrappedLogger) Debug(args ...interface{}) { l.entry.Debug(args...) }
func (l *wrappedLogger) Info(args ...interface{})  { l.entry.Info(args...) }
func (l *wrappedLogger) Warn(args ...interface{})  { l.entry.Warn(args...) }
func (l *wrappedLogger) Error(args ...interface{}) { l.entry.Error(args...) }
func (l *wrappedLogger) Fatal(args ...interface{}) { l.entry.Fatal(args...) }

func (l *wrappedLogger) Debugf(format string, args ...interface{}) { l.entry.Debugf(format, args...) }
func (l *wrappedLogger) Infof(format string, args ...interface{})  { l.entry.Infof(format, args...) }
func (l *wrappedLogger) Warnf(format string, args ...interface{})  { l.entry.Warnf(format, args...) }
func (l *wrappedLogger) Errorf(format string, args ...interface{}) { l.entry.Errorf(format, args...) }
func (l *wrappedLogger) Fatalf(format string, args ...interface{}) { l.entry.Fatalf(format, args...) }
