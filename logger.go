package yagl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"text/template"
	"time"
)

// loginfo is the log info struct, represents a log message
// and information to be printed along aside with the message
type loginfo struct {
	DateTime time.Time `json:"datetime"`
	Level    LogLevel  `json:"level"`
	Message  string    `json:"message"`
	PkgName  string    `json:"package_name"`
	FuncName string    `json:"func_name"`
}

// ToJSON marshals the loginfo struct to json
func (l *loginfo) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}

// Logger is the logger struct
type Logger struct {
	// format is the logger format
	format string
	// level is the logger level
	level     LogLevel
	tmpl      *template.Template
	mtx       sync.Mutex
	levelOuts map[LogLevel]io.Writer
	// mapper is the message mapper function that is
	// used to mutate the message before printing it
	mapper Mapper
	// isJSON is a flag that indicates whether the logger
	// should print the log message as json or not
	isJSON bool

	buffersPool sync.Pool
}

// Defaults list of the default logger settings
var Defaults = []Setting{StdFormat, Level(Info), DefaultStd, WithMapper(noOpMapper)}

// New creates a new logger
func New(opts ...Setting) *Logger {
	if len(opts) == 0 {
		opts = Defaults
	}

	l := &Logger{
		levelOuts: make(map[LogLevel]io.Writer),
		buffersPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	// apply the options
	for _, opt := range opts {
		opt(l)
	}

	// ensure that the logger has a mapper
	if l.mapper == nil {
		l.mapper = noOpMapper
	}

	// if no format is defined, use the default format
	if l.format == "" {
		StdFormat(l)
	}

	// check whenver all levels have corresponding writers or not
	for _, level := range AllLevels {
		if _, exists := l.levelOuts[level]; !exists {
			// if not, use the info writer
			fmt.Fprintln(os.Stderr, "No writers defined, trying to fallback to the info writer", level)
			l.levelOuts[level] = os.Stdout
		}
	}
	return l
}

// Logf logs a message with given arguments and log level
func (l *Logger) Logf(level LogLevel, msg string, args ...interface{}) {
	// Ensure logger could be used concurrently
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if level < l.level {
		return
	}

	// Reduce allocations by using a buffer pool
	buffer := l.buffersPool.Get().(*bytes.Buffer)
	defer l.buffersPool.Put(buffer)

	info := l.loginfo(level, msg, args...)

	if l.isJSON {
		// if jsonMessage is enabled, marshal the loginfo struct to jsonMessage
		jsonMessage, err := info.ToJSON()
		if err != nil {
			m := fmt.Sprintf(msg, args...)
			info = l.loginfo(Error, "Failed to marshal log message [%s] to json, %+v", m, err)
			jsonMessage, err = info.ToJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to marshal log messsage [%s] to json %+v", m, err)
			}
		}
		buffer.Write(jsonMessage)
		buffer.WriteString("\n")
	} else {
		if err := l.tmpl.Execute(buffer, info); err != nil {
			m := fmt.Sprintf(msg, args...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to formate log messsage [%s] according to format [%s], %+v",
					m, l.format, err)
			}
		}
	}

	bOut := bytes.Join([][]byte{buffer.Bytes(), []byte("\n")}, []byte(""))
	// Write to the appropriate writer
	if out, exists := l.levelOuts[level]; exists {
		_, err := out.Write(bOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log message [%s] to writer level [%s],  %+v",
				buffer.String(), level, err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "No writers defined, trying to fallback to the info writer")
		infoOut, exists := l.levelOuts[Info]
		if !exists {
			fmt.Fprintf(os.Stderr, "No writers defined, failed to output log message, %s", buffer.String())
			return
		}
		if _, err := infoOut.Write(bOut); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log message [%s] to writer level [%s],  %+v",
				buffer.String(), level, err)
		}
	}

}

// Setup sets the logger options
func (l *Logger) Setup(opt ...Setting) {
	for _, o := range opt {
		o(l)
	}
}

// loginfo creates a loginfo struct for a message with given arguments
func (l *Logger) loginfo(level LogLevel, msg string, args ...interface{}) *loginfo {
	// mutate the message if needed
	message := l.mapper(msg)
	if len(args) > 0 {
		message = l.mapper(fmt.Sprintf(msg, args...))
	}
	if l.format == dbgFormat {
		// in case there is debug formatting is enabled
		// there is a need to get the caller info to
		// extract the package name and function name
		pkgName, funcName, _ := getCallerInfo() // get package name and function name
		return &loginfo{
			DateTime: time.Now(),
			Level:    level,
			Message:  message,
			PkgName:  pkgName,
			FuncName: funcName,
		}
	}
	return &loginfo{
		DateTime: time.Now(),
		Level:    level,
		Message:  message,
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
