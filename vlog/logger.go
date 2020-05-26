package vlog

import (
	"fmt"
	"log"
	"os"
)

var (
	Log = NewLogger("Default Logger", InfoLevel)
)

const (
	ErrorLevel int = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

func InitLog(prefix string, logLevel int) {
	Log = NewLogger(prefix, logLevel)
}

func NewLogger(prefix string, logLevel int) *Logger {
	return &Logger{
		OutLogger: log.New(os.Stdout, "", log.LstdFlags),
		ErrLogger: log.New(os.Stderr, "", log.LstdFlags),
		LogLevel:  logLevel,
		Prefix:    prefix,
	}
}

type Logger struct {
	OutLogger *log.Logger
	ErrLogger *log.Logger
	LogLevel  int
	Prefix    string
}

func getLogPrefixHeader(prefix, level string) string {
	return fmt.Sprintf("[%s][%s] %s", prefix, level, "%s")
}

func (l Logger) Error(s string, i ...interface{}) {
	if l.LogLevel >= ErrorLevel {
		header := getLogPrefixHeader(l.Prefix, "ERROR")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.ErrLogger.Println(message)
	}
}

func (l Logger) Info(s string, i ...interface{}) {
	if l.LogLevel >= InfoLevel {
		header := getLogPrefixHeader(l.Prefix, "INFO")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.OutLogger.Println(message)
	}
}

func (l Logger) Debug(s string, i ...interface{}) {
	if l.LogLevel >= DebugLevel {
		header := getLogPrefixHeader(l.Prefix, "DEBUG")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.OutLogger.Println(message)
	}
}

func (l Logger) Warn(s string, i ...interface{}) {
	if l.LogLevel >= WarnLevel {
		header := getLogPrefixHeader(l.Prefix, "WARNING")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.OutLogger.Println(message)
	}
}

func (l Logger) Fatal(s string, i ...interface{}) {
	header := getLogPrefixHeader(l.Prefix, "FATAL")
	message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
	panic(message)
}

func Error(s string, i ...interface{}) {
	Log.Error(s, i...)
}

func Info(s string, i ...interface{}) {
	Log.Info(s, i...)
}

func Debug(s string, i ...interface{}) {
	Log.Debug(s, i...)
}

func Warn(s string, i ...interface{}) {
	Log.Warn(s, i...)
}

func Fatal(s string, i ...interface{}) {
	Log.Fatal(s, i...)
}
