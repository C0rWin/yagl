package yagl

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"time"
)

var (
	// Default logger format template
	StdFormat = `{{.DateTime.Format "2006-01-02 15:04:05"}} [{{.Level}}]: {{.Message}}`
	// Default logger format template with package name and function name, mainly for debug purpose
	DebugFormat = `{{.DateTime.Format "2006-01-02 15:04:05"}} ({{.PkgName}}/{{.FuncName}}) [{{.Level}}]: {{.Message}}`
)

type Option func(*Logger)

// StdFormatOption sets the logger format to the default format
func StdFormatOption(l *Logger) {
	l.format = StdFormat
}

// DebugFormatOption sets the logger format to the debug format
func DebugFormatOption(l *Logger) {
	l.format = DebugFormat
}

// CustomFormanOption sets the logger format to the custom format
func CustomFormatOption(format string) Option {
	return func(l *Logger) {
		l.format = format
	}
}

// LogLevelOption sets the logger level
func LogLevelOption(level LogLevel) Option {
	return func(l *Logger) {
		l.level = level
	}
}

// DefaultLogPrinterOption sets the logger printer to the default printer
func DefaultStd(l *Logger) {
	l.stdOut = os.Stdout
}

// CustomLogOut(stdOut io.Writer) sets the logger printer to the custom printer
func CustomLogOut(stdOut io.Writer) Option {
	return func(l *Logger) {
		l.stdOut = stdOut
	}
}

// loginfo is the log info struct, represents a log message
// and information to be printed along aside with the message
type loginfo struct {
	DateTime time.Time
	Level    LogLevel
	Message  string
	PkgName  string
	FuncName string
}

// Logger is the logger struct
type Logger struct {
	format string
	level  LogLevel
	tmpl   *template.Template
	stdOut io.Writer
}

// New creates a new logger
func New(opts ...Option) *Logger {
	if len(opts) == 0 {
		opts = append(opts, StdFormatOption, LogLevelOption(Info), DefaultStd)
	}
	l := &Logger{}
	for _, opt := range opts {
		opt(l)
	}
	l.tmpl = template.Must(template.New("log").Parse(l.format))
	return l
}

// Log logs a message
func (l *Logger) Logf(level LogLevel, msg string, args ...interface{}) {
	if level >= l.level {
		buffer := bytes.NewBuffer(nil)
		info := l.logi(level, msg, args...)

		if err := l.tmpl.Execute(buffer, info); err != nil {
			panic(err)
		}

		l.stdOut.Write(buffer.Bytes())
		l.stdOut.Write([]byte("\n"))
	}
}

// logi creates a loginfo struct for a message with given arguments
func (l *Logger) logi(level LogLevel, msg string, args ...interface{}) *loginfo {
	return &loginfo{
		DateTime: time.Now(),
		Level:    level,
		Message:  msg,
	}
}
