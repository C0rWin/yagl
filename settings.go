package yagl

import (
	"html/template"
	"io"
	"os"
)

// Setting is a function that sets a logger option
type Setting func(*Logger)

// Format sets the logger format to the custom format
func Format(format string) Setting {
	return func(l *Logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.format = format
		l.tmpl = template.Must(template.New("log").Parse(l.format))
	}
}

var (
	// StdFormat sets the logger format to the default format
	StdFormat = Format(Std)
	// DebugFormat sets the logger format to the debug format
	DebugFormat = Format(Dbg)
)

// Level sets the logger level
func Level(level LogLevel) Setting {
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
// if no log level is provided, all log levels will be set to the
// provided writer. If a log level is provided, only that log level
// will be set to the provided writer.
func CustomLogOut(out io.Writer, levels ...LogLevel) Setting {
	return func(l *Logger) {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		if len(levels) == 0 {
			levels = []LogLevel{Debug, Info, Warn, Error}
		}
		for _, level := range levels {
			l.levelOuts[level] = out
		}
	}
}
