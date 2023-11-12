package yagl

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerHierarchy(t *testing.T) {
	t.Parallel()
	for _, logLevel := range []LogLevel{Debug, Info, Warn, Error} {
		for _, messageLogLevel := range []LogLevel{Debug, Info, Warn, Error} {
			messageLogLevel := messageLogLevel
			t.Run(logLevel.String()+"-"+messageLogLevel.String(), func(t *testing.T) {
				t.Parallel()
				buffer := bytes.NewBuffer(nil)
				logger := New(CustomLogOut(buffer), StdFormat, Level(logLevel))
				logger.Logf(messageLogLevel, "Hello World")
				if messageLogLevel >= logLevel {
					require.Contains(t, buffer.String(), "Hello World")
				} else {
					require.Empty(t, buffer.String())
				}
			})
		}
	}
}

func TestLoggerConcurrentUsage(t *testing.T) {
	t.Parallel()
	buffer := bytes.NewBuffer([]byte{})
	logger := New(CustomLogOut(buffer), StdFormat, Level(Debug))
	wg := &sync.WaitGroup{}
	wg.Add(4)
	for _, logLevel := range []LogLevel{Debug, Info, Warn, Error} {
		go func(l *Logger, level LogLevel) {
			defer wg.Done()
			l.Setup(Level(level))
			for _, lvl := range []LogLevel{Debug, Info, Warn, Error} {
				// Since the setting the log level and logging
				// actually happens concurrently, we can't be sure
				// that the log level will be set before the log
				// and since these operations are not atomic, there
				// is no point to assert that message is logged

				// This test is just to ensure that the logger
				// could be used concurrently and we can access it's
				// method to change internal state concurrently
				l.Logf(lvl, "Hello World")
			}
		}(logger, logLevel)
	}

	wg.Wait()
}
