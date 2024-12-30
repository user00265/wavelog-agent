package main

import (
	"log"
	"os"
)

type loggerTypes struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

var loggerFlags = log.LstdFlags | log.Lshortfile

var infoTLogger = log.New(os.Stdout, "INFO: ", loggerFlags)
var warnTLogger = log.New(os.Stdout, "WARN: ", loggerFlags)
var errorTLogger = log.New(os.Stdout, "ERROR: ", loggerFlags)

func (l *loggerTypes) info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l *loggerTypes) warn(v ...interface{}) {
	l.warnLogger.Println(v...)
}

func (l *loggerTypes) error(v ...interface{}) {
	l.errorLogger.Fatal(v...)
}

var logger = loggerTypes{
	infoLogger:  infoTLogger,
	warnLogger:  warnTLogger,
	errorLogger: errorTLogger,
}
