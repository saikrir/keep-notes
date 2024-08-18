package logger

import (
	"bytes"
	"log"
	"os"
)

const (
	DEBUG = iota
	INFO  = 1
	WARN  = 2
	ERROR = 3
)

const (
	DebugColor   = "\033[1;36m"
	InfoColor    = "\033[1;34m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
)

type logcolor = string
type loglevel = int

const logFlags = log.Ldate | log.Ltime | log.Lmicroseconds

func logMessage(level loglevel, messages ...any) {
	var logPrefix, logColor string

	switch level {
	case DEBUG:
		logPrefix = "DEBUG:\t"
		logColor = DebugColor
	case INFO:
		logPrefix = "INFO:\t"
		logColor = InfoColor
	case WARN:
		logPrefix = "WARN:\t"
		logColor = WarningColor
	case ERROR:
		logColor = "ERROR:\t"
		logColor = ErrorColor
	default:
		logPrefix = "UNKNOWN:\t"
		logColor = DebugColor
	}
	lgr := log.New(os.Stdout, logPrefix, logFlags)
	lgr.Printf(constructFormat(logColor, len(messages)), messages...)
}

func Debug(messages ...any) {
	logMessage(DEBUG, messages...)
}

func Info(messages ...any) {
	logMessage(INFO, messages...)
}

func Warning(messages ...any) {
	logMessage(WARN, messages...)
}

func Error(messages ...any) {
	logMessage(ERROR, messages...)
}

func constructFormat(color logcolor, nArgs int) string {
	buf := new(bytes.Buffer)
	buf.WriteString(color)
	for i := range nArgs {
		buf.WriteString("%v")
		if i != nArgs-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("\033[0m \n")
	return buf.String()
}
