package logger

import (
	"log"
	"os"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})

	Infof(format string, args ...interface{})
	Info(args ...interface{})
}

type LogLevelSetter interface {
	SetLogLevel(level Level)
}

var Log Logger = &defaultLogger{
	Logger: log.New(os.Stderr, "", log.LstdFlags),
	Level:  LevelInfo,
}

func SetLogger(l Logger) {
	Log = l
}

func SetLogLevel(level Level) {
	if l, ok := Log.(LogLevelSetter); ok {
		l.SetLogLevel(level)
	}
}

type defaultLogger struct {
	*log.Logger
	Level Level
}

func Debugf(format string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Debug(args...)
}

func Infof(format string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Infof(format, args...)
}

func Info(args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Info(args...)
}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	if l.Level > LevelDebug {
		return
	}
	l.Printf(format, args...)
}

func (l *defaultLogger) Debug(args ...interface{}) {
	if l.Level > LevelDebug {
		return
	}
	l.Print(args...)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func (l *defaultLogger) Info(args ...interface{}) {
	l.Print(args...)
}

func (l *defaultLogger) SetLogLevel(lvl Level) {
	l.Level = lvl
}
