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
)

var logLevelStrings = map[LogLevel]string{
	Debug: "DEBUG",
	Warn:  "WARN",
	Info:  "INFO",
	Error: "ERROR",
}

// String returns the string representation of the LogLevel.
func (l LogLevel) String() string {
	s, exists := logLevelStrings[l]
	if exists {
		return s
	}
	return "INFO"
}
