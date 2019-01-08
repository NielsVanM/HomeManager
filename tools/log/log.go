package log

import (
	"fmt"
	"os"
	"strings"
)

const (
	// Information is a log level for nice to know information
	Information = "INFO"

	// Warning is a log level for recoverable errors that should
	// not occur
	Warning = "WARNING"

	// Error is a log level for product breaking errors
	Error = "ERROR"

	// Fatality is a log level that is identifical to error but
	// quits the application by panicking
	Fatality = "FATAL"
)

// Log is a global function for outputting information to stdout
// and log.txt
func Log(level, module string, message ...string) {

	compactMessage := strings.Join(message, " ")
	outMessage := fmt.Sprintf("[%s/%s] %s", module, level, compactMessage)

	fmt.Println(outMessage)

	if level == Fatality {
		panic(message)
	}

	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		//REVIEW Should we panic or do something else when we can't store logs?
		panic(err)
	}
	defer f.Close()

	if _, err := f.WriteString(outMessage + "\n"); err != nil {
		panic(err)
	}
}

// Info logs message at information level
func Info(module string, message ...string) {
	Log(Information, module, message...)
}

// Warn logs message at information level
func Warn(module string, message ...string) {
	Log(Warning, module, message...)
}

// Err logs message at information level
func Err(module string, message ...string) {
	Log(Error, module, message...)
}

// Fatal logs message at information level
func Fatal(module string, message ...string) {
	Log(Fatality, module, message...)
}
