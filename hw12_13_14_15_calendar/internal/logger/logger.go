package logger

import (
	"fmt"
)

// error / warn / info / debug
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

type Logger struct {
	level string
	lev   int
}

func New(level string) *Logger {
	var lev int
	switch level {
	case DEBUG:
		lev = 1
	case INFO:
		lev = 2
	case WARN:
		lev = 3
	case ERROR:
		lev = 4
	default:
		lev = 0 // print Any
	}
	return &Logger{
		level: level,
		lev:   lev,
	}
}

func (l *Logger) print(num int, msg string) {
	if l.lev <= num {
		fmt.Printf("%s: %s", l.level, msg)
	}
}

func (l Logger) Debug(msg string) {
	l.print(1, msg)
}

func (l Logger) Info(msg string) {
	l.print(2, msg)
}

func (l Logger) Warn(msg string) {
	l.print(3, msg)
}

func (l Logger) Error(msg string) {
	l.print(4, msg)
}
