package logger

import (
	"fmt"
	"os"
	"strings"
)

type LogLevel int

const (
	NONE = iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
)

func (ll LogLevel) String() string {
	switch ll {
	case INFO:
		return "INFO "
	case WARN:
		return "WARN "
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}

	return "NONE"
}

func LogLevelFromString(s string) LogLevel {
	switch strings.ToUpper(s) {
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "DEBUG":
		return DEBUG
	case "ERROR":
		return ERROR
	}

	return NONE
}

type Logger struct {
	logLevel LogLevel
	prefix   string
}

func (l Logger) Prefix() string {
	ws := []rune("                         ")
	prefix := []rune(l.prefix)

	spacesToAdd := len(ws) - len(prefix)

	return fmt.Sprint(string(prefix) + string(ws[:spacesToAdd]))
}

func NewLogger(level string, destination *os.File) Logger {
	logLevel := LogLevelFromString(level)

	return Logger{logLevel, ""}
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

var EmptyLogger = NewLogger("NONE", nil)

func (l *Logger) write(msgLevel LogLevel, v ...any) {
	if l.logLevel >= msgLevel {
		fmt.Println(l.Prefix(), msgLevel.String()+": "+fmt.Sprint(v...))
	}
}

func (l *Logger) writef(msgLevel LogLevel, format string, v ...any) {
	if l.logLevel >= msgLevel {
		fmt.Printf(l.Prefix()+" "+msgLevel.String()+": "+format, v...)
	}
}

func (l *Logger) Info(v ...any)  { l.write(INFO, v...) }
func (l *Logger) Warn(v ...any)  { l.write(WARN, v...) }
func (l *Logger) Debug(v ...any) { l.write(DEBUG, v...) }
func (l *Logger) Error(v ...any) { l.write(ERROR, v...) }
func (l *Logger) Fatal(v ...any) {
	l.write(FATAL, v...)
	os.Exit(1)
}
func (l *Logger) Infof(format string, v ...any)  { l.writef(INFO, format, v...) }
func (l *Logger) Warnf(format string, v ...any)  { l.writef(WARN, format, v...) }
func (l *Logger) Debuf(format string, v ...any)  { l.writef(DEBUG, format, v...) }
func (l *Logger) Errorf(format string, v ...any) { l.writef(ERROR, format, v...) }
func (l *Logger) Fatalf(format string, v ...any) {
	l.writef(FATAL, format, v...)
	os.Exit(1)
}
