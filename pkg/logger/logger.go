package logger

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ANSI color codes
const (
	red     = 31
	yellow  = 33
	blue    = 34
	magenta = 35
	cyan    = 36
	white   = 37
)

// getColor returns the ANSI color code for the given log level.
func getColor(level logrus.Level) int {
	switch level {
	case logrus.FatalLevel, logrus.PanicLevel, logrus.ErrorLevel:
		return red
	case logrus.WarnLevel:
		return yellow
	case logrus.InfoLevel:
		return blue
	case logrus.DebugLevel:
		return magenta
	case logrus.TraceLevel:
		return cyan
	default:
		return white
	}
}

var logger = logrus.New()

func init() {
	env := os.Getenv("APP_ENV")

	logger.Level = logrus.InfoLevel
	if env == "production" {
		logger.SetFormatter(&CustomJSONFormatter{
			JSONFormatter: logrus.JSONFormatter{},
		})
	} else {
		logger.Formatter = &formatter{}
	}

	logger.SetReportCaller(false)
}

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

type Fields logrus.Fields

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debugf(format, args...)
	}
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Infof(format, args...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warnf(format, args...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Errorf(format, args...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatalf(format, args...)
	}
}

// Formatter implements logrus.Formatter interface.
type formatter struct {
	prefix string
}

// Format building log message.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	// Get caller information
	if entry.HasCaller() {
		frame := getFrame(10)                          // Adjust the frame number as needed
		sb.WriteString(fmt.Sprintf("\x1b[%dm", white)) // White color for file & line
		sb.WriteString(path.Base(frame.File))
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%d", frame.Line))
		sb.WriteString(" ")
	}

	// Apply color based on log level
	levelColor := getColor(entry.Level)
	sb.WriteString(fmt.Sprintf("\x1b[%dm", levelColor)) // Color start
	sb.WriteString(strings.ToUpper(entry.Level.String()))
	sb.WriteString(" ")
	sb.WriteString(entry.Time.Format(time.RFC3339))
	sb.WriteString(" ")
	sb.WriteString(f.prefix)
	sb.WriteString(entry.Message)
	sb.WriteString("\x1b[0m") // Reset color
	sb.WriteString("\n")

	return sb.Bytes(), nil
}

// Custom JSON Formatter
type CustomJSONFormatter struct {
	logrus.JSONFormatter
}

func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Adjust skipFrames based on your call stack
	frame := getFrame(10)

	// Adding custom fields for the corrected caller information
	entry.Data["location"] = fmt.Sprintf("%s:%d", frame.File, frame.Line)

	return f.JSONFormatter.Format(entry)
}

func getFrame(skipFrames int) runtime.Frame {
	pc := make([]uintptr, 15)
	runtime.Callers(skipFrames, pc)

	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()

		// Check if the frame is from your controllers and not from logger
		if strings.HasPrefix(frame.File, "/app") &&
			!strings.Contains(frame.File, "/app/infra/logger/") {
			return frame
		}

		if !more {
			break
		}
	}

	return runtime.Frame{} // Handle the case where no matching frame is found
}
