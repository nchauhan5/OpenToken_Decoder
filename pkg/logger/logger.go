package logger

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// DefaultLoggerProperties would be used to setup the logger if environment variables are not set.
type DefaultLoggerProperties struct {
	formatter log.Formatter
	output    io.Writer
	logLevel  log.Level
}

var defaultLogProps = DefaultLoggerProperties{
	formatter: &log.TextFormatter{},
	output:    os.Stdout,
	logLevel:  log.DebugLevel,
}

// Logger to be used across the auth service
var Logger = log.New()

func init() {

	fmt.Println("Inside logger init() *************")
	// Check for formatter type
	if os.Getenv("formatter") == "JSON" {
		Logger.SetFormatter(&log.JSONFormatter{})
	} else {
		Logger.SetFormatter(defaultLogProps.formatter)
	}

	// Check for Output type either os.StdOut or file
	if os.Getenv("output") == "File" {
		logfile, err := os.OpenFile("../../logs/access.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err == nil {
			Logger.SetOutput(logfile)
		} else {
			Logger.SetOutput(defaultLogProps.output)
		}
	} else {
		Logger.SetOutput(defaultLogProps.output)
	}

	// Check the log level type among DEBUG, INFO, WARN, ERROR, FATAL, PANIC
	if os.Getenv("logLevel") == "INFO" {
		Logger.SetLevel(log.InfoLevel)
	} else if os.Getenv("logLevel") == "WARN" {
		Logger.SetLevel(log.WarnLevel)
	} else if os.Getenv("logLevel") == "ERROR" {
		Logger.SetLevel(log.ErrorLevel)
	} else if os.Getenv("logLevel") == "FATAL" {
		Logger.SetLevel(log.FatalLevel)
	} else if os.Getenv("logLevel") == "PANIC" {
		Logger.SetLevel(log.PanicLevel)
	} else {
		Logger.SetLevel(defaultLogProps.logLevel)
	}

}
