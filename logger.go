package yagl

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

var (
	// Default logger format template
	StdFormat = `[{{.DateTime.Format "2006-01-02 15:04:05"}}] [{{.Level}}]: {{.Message}}`
	// Default logger format template with package name and function name, mainly for debug purpose
	DebugFormat = `[{{.DateTime.Format "2006-01-02 15:04:05"}}] ({{.PkgName}}/{{.FuncName}}) [{{.Level}}]: {{.Message}}`
)

type Option func(*Logger)

// CustomFormanOption sets the logger format to the custom format
func CustomFormatOption(format string) Option {
	return func(l *Logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.format = format
		l.tmpl = template.Must(template.New("log").Parse(l.format))
	}
}

var (
	// StdFormatOption sets the logger format to the default format
	StdFormatOption = CustomFormatOption(StdFormat)
	// DebugFormatOption sets the logger format to the debug format
	DebugFormatOption = CustomFormatOption(DebugFormat)
)

// LogLevelOption sets the logger level
func LogLevelOption(level LogLevel) Option {
	return func(l *Logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.level = level
	}
}

// DefaultStd sets the logger output for different log levels
// default output writers. All but Error log level are set to
// os.Stdout and the Error level set to the os.Stder
//
// It's possible to overwrite the default output writers with
// CustomLogOut option, where for each log level a custom writer
// could be defined
func DefaultStd(l *Logger) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.levelOuts = map[LogLevel]io.Writer{
		Debug: os.Stdout,
		Info:  os.Stdout,
		Warn:  os.Stdout,
		Error: os.Stderr,
	}
}

// CustomLogOut sets the logger output for different log levels
func CustomLogOut(out io.Writer, levels ...LogLevel) Option {
	return func(l *Logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		for _, level := range levels {
			l.levelOuts[level] = out
		}
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
	format       string
	level        LogLevel
	tmpl         *template.Template
	debugEnabled bool
	mtx          sync.Mutex
	levelOuts    map[LogLevel]io.Writer
}

// New creates a new logger
func New(opts ...Option) *Logger {
	if len(opts) == 0 {
		opts = append(opts, StdFormatOption, LogLevelOption(Info), DefaultStd)
	}

	l := &Logger{
		levelOuts: make(map[LogLevel]io.Writer),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Log logs a message
func (l *Logger) Logf(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	buffer := bytes.NewBuffer(nil)
	info := l.logi(level, msg, args...)

	if err := l.tmpl.Execute(buffer, info); err != nil {
		panic(err)
	}

	// Ensure logger could be used concurrently
	l.mtx.Lock()
	defer l.mtx.Unlock()

	// Write to the appropriate writer
	if out, exists := l.levelOuts[level]; exists {
		out.Write(buffer.Bytes())
		out.Write([]byte("\n"))
	} else {
		panic(fmt.Sprintf("No writer for level %s", level.String()))
	}

}

// SetOptions sets the logger options
func (l *Logger) SetOptions(opt ...Option) {
	for _, o := range opt {
		o(l)
	}
}

// logi creates a loginfo struct for a message with given arguments
func (l *Logger) logi(level LogLevel, msg string, args ...interface{}) *loginfo {
	if l.format == DebugFormat {
		pkgName, funcName, _ := getCallerInfo()
		return &loginfo{
			DateTime: time.Now(),
			Level:    level,
			Message:  fmt.Sprintf(msg, args...),
			PkgName:  pkgName,
			FuncName: funcName,
		}
	}
	return &loginfo{
		DateTime: time.Now(),
		Level:    level,
		Message:  fmt.Sprintf(msg, args...),
	}
}

// getCallerInfo gets the package name and function name of the caller
func getCallerInfo() (pakageName, funcName string, line int) {
	pc, _, line, ok := runtime.Caller(3)
	if !ok {
		return "undefined", "undefined", 0
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "undefined", "undefined", 0
	}

	funcName = fn.Name()
	pakageName = path.Dir(funcName)

	return pakageName, path.Base(funcName), line
}
