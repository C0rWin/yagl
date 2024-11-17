// Package yagl is a Yet Another Go Logger, a simple logger for go.
// Description: This file contains the interface for the logger. It is
// used to define the log levels and the log level strings.
package yagl

// LogLevel is the level of logging.
type LogLevel int

const (
	// Debug log level
	Debug LogLevel = 1 << iota
	// Info log level
	Info
	// Warn log level
	Warn
	// Error log level
	Error
	// Panic log level
	Panic
	// Fatal log level
	Fatal
)

var logLevelStrings = map[LogLevel]string{
	Debug: "DEBUG",
	Warn:  "WARN",
	Info:  "INFO",
	Error: "ERROR",
	Panic: "PANIC",
	Fatal: "FATAL",
}

// String returns the string representation of the LogLevel.
func (l LogLevel) String() string {
	s, exists := logLevelStrings[l]
	if exists {
		return s
	}
	return "INFO"
}

// AllLevels is a list of all log levels.
var AllLevels = []LogLevel{Debug, Info, Warn, Error}
