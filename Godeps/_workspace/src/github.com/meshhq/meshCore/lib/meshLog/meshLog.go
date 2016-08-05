package meshLog

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/meshhq/meshCore/lib/gohttp/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/meshhq/meshCore/lib/gohttp/Godeps/_workspace/src/github.com/getsentry/raven-go"
	"github.com/meshhq/meshCore/lib/gohttp/Godeps/_workspace/src/github.com/gorilla/mux"
)

// LogLevel indicates the level of logging
type LogLevel int

const (
	debug LogLevel = iota
	info
	warn
	fatal
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

// Date / Time format
const dtFormat = "Jan 2 15:04:05"

// Logger is the wrapper around a given logging message
type Logger struct {
	envVar   string
	messages []loggedMessage
}

type loggedMessage struct {
	level   LogLevel
	message string
}

/**
 * Init Sentry
 */
func init() {
	raven.SetDSN(os.Getenv("SENTRY_KEY"))
}

/**
 * String Logging
 */

// Write allows the logger to conform to io.Writer. Assume
// info level logging
func Write(data []byte) (int, error) {
	Info(string(data))
	return len(data), nil
}

// Debug is a convenience method appending a debug message to the logger
func Debug(obj interface{}) {
	msg := fmt.Sprintf("%+v\n", obj)
	formattedMessage := formattedLogMessage("DEBUG", msg)
	color.Green(formattedMessage)
}

// Info is a convenience method appending a info style message to the logger
func Info(obj interface{}) {
	msg := fmt.Sprintf("%+v\n", obj)
	formattedMessage := formattedLogMessage("INFO", msg)
	color.White(formattedMessage)
}

// Warn is a convenience method appending a warning message to the logger
func Warn(obj interface{}) {
	msg := fmt.Sprintf("%+v\n", obj)
	formattedMessage := formattedLogMessage("WARN", msg)
	color.Yellow(formattedMessage)
}

// Fatal is a convenience method appending a fatal message to the logger
func Fatal(obj interface{}) {
	msg := fmt.Sprintf("%+v\n", obj)
	formattedMessage := formattedLogMessage("ERROR", msg)
	color.Red(formattedMessage)
}

/**
 * Formatted Strings
 */

// Debugf is a convenience method appending a debug message to the logger
func Debugf(msg string, a ...interface{}) {
	msg = fmt.Sprintf(msg, a...)
	formattedMessage := formattedLogMessage("DEBUG", msg)
	color.Green(formattedMessage)
}

// Infof is a convenience method appending a info style message to the logger
func Infof(msg string, a ...interface{}) {
	msg = fmt.Sprintf(msg, a...)
	formattedMessage := formattedLogMessage("INFO", msg)
	color.White(formattedMessage)
}

// Warnf is a convenience method appending a warning message to the logger
func Warnf(msg string, a ...interface{}) {
	msg = fmt.Sprintf(msg, a...)
	formattedMessage := formattedLogMessage("WARN", msg)
	color.Yellow(formattedMessage)
}

// Fatalf is a convenience method appending a fatal message to the logger
func Fatalf(msg string, a ...interface{}) {
	msg = fmt.Sprintf(msg, a...)
	formattedMessage := formattedLogMessage("ERROR", msg)
	color.Red(formattedMessage)
}

/**
 * Internal Formatting
 */

func formattedLogMessage(level string, logMessage string) string {
	env := "LOCAL"
	if len(os.Getenv("ENV")) > 0 {
		env = strings.ToUpper(os.Getenv("ENV"))
	}

	return fmt.Sprintf("[%s] - %s: %s", env, level, logMessage)
}

func formatColoredMessage(message string, level LogLevel) string {
	var levelColor int
	switch level {
	case debug:
		levelColor = yellow
	case info:
		levelColor = gray
	case warn:
		levelColor = green
	case fatal:
		levelColor = red
	}

	// levelText := strings.ToUpper(message)[0:4]
	return fmt.Sprintf("\x1b[%dm%s\x1b", levelColor, message)
}

func stringValueForLogLevel(level LogLevel) string {
	switch level {
	case debug:
		return "DEBUG"
	case info:
		return "INFO"
	case warn:
		return "WARN"
	case fatal:
		return "FATAL"
	}
	return "INFO"
}

/**
 * Convenience for panic / err
 */

// Perror is Syntax Sugga for panicing on error
func Perror(err error) {
	if err != nil {
		Fatal(err)
		panic(err)
	}
}

/**
 * Request Logging
 */

// RecordPanicWithRequest records an error, and parses the
// accompanying request
func RecordPanicWithRequest(capturedPanic interface{}, req *http.Request) {

	// Record Req Specifics
	tags := map[string]string{
		"env":             os.Getenv("ENV"),
		"req-url":         req.URL.String(),
		"req-method":      req.Method,
		"req-headers":     fmt.Sprintf("%+v", req.Header),
		"req-queryParams": fmt.Sprintf("%+v", mux.Vars(req)),
		"req-remoteAddr":  req.RemoteAddr,
	}

	// Attempt to report the error
	if err, ok := capturedPanic.(error); ok {
		// Send to reporting service
		Info(err)
		raven.CaptureError(err, tags)
	} else {
		errMsg := fmt.Sprintf("%+v\n", capturedPanic)
		err = errors.New(errMsg)
		Info(err)
		raven.CaptureError(err, tags)
	}
}

// RecoveryHandler wraps a http Handler to watch for exceptions
func RecoveryHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return raven.RecoveryHandler(handler)
}

// CapturePanic is a wrapper around a method that could panic
func CapturePanic(f func(), tags map[string]string, interfaces ...raven.Interface) (interface{}, string) {
	return raven.CapturePanic(f, tags, interfaces...)
}
