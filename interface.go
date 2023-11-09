package yagl

// LogLevel is the level of logging.
type LogLevel int

const (
	Debug LogLevel = 1 << iota
	Warn
	Info
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
