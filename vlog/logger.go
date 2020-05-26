package vlog

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Log *Logger
)

const (
	ErrorLevel int = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

func InitLog(prefix string, logLevel int) {
	if Log != nil {
		return
	}
	Log = NewLogger(prefix, logLevel)
}

func NewLogger(prefix string, logLevel int) *Logger {
	return &Logger{
		outLogger: log.New(os.Stdout, prefix, log.LstdFlags),
		errLogger: log.New(os.Stderr, prefix, log.LstdFlags),
		LogLevel:  logLevel,
	}
}

type Logger struct {
	outLogger *log.Logger
	errLogger *log.Logger
	LogLevel  int
}

func getLogPrefixHeader(level string) string {
	logTimestamp := time.Now().Format("20130519:23:12:00")
	return fmt.Sprintf(" %s [%s] %s", logTimestamp, level, "%s")
}

func (l Logger) Error(s string, i ...interface{}) {
	if l.LogLevel >= ErrorLevel {
		header := getLogPrefixHeader("ERROR")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.errLogger.Println(message)
	}
}

func (l Logger) Info(s string, i ...interface{}) {
	if l.LogLevel >= InfoLevel {
		header := getLogPrefixHeader("INFO")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.errLogger.Println(message)
	}
}

func (l Logger) Debug(s string, i ...interface{}) {
	if l.LogLevel >= DebugLevel {
		header := getLogPrefixHeader("DEBUG")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.errLogger.Println(message)
	}
}

func (l Logger) Warn(s string, i ...interface{}) {
	if l.LogLevel >= WarnLevel {
		header := getLogPrefixHeader("WARNING")
		message := fmt.Sprintf(header, fmt.Sprintf(s, i...))
		l.errLogger.Println(message)
	}
}

func (l Logger) Fatal(v ...interface{}) {
	l.errLogger.Println(v...)
	os.Exit(1)
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

func Fatal(v ...interface{}) {
	Log.Fatal(v...)
}
