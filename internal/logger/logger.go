package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

var (
	InfoLogger  Logger
	ErrorLogger Logger
)

type stdLogger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func (l *stdLogger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

func (l *stdLogger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

func Init() {
	logger := &stdLogger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	InfoLogger = logger
	ErrorLogger = logger
}

func Info(format string, v ...interface{}) {
	InfoLogger.Info(format, v...)
}

func Error(format string, v ...interface{}) {
	ErrorLogger.Error(format, v...)
}
