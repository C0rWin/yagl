package yagl

import (
	"bytes"
	"fmt"
	"regexp"
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

func TestLoggerHidingSecrets(t *testing.T) {
	t.Parallel()
	buffer := bytes.NewBuffer(nil)
	logger := New(CustomLogOut(buffer), StdFormat, Level(Debug), WithMapper(func(s string) string {
		return string(bytes.ReplaceAll([]byte(s), []byte("Secret"), []byte("********")))
	}))
	logger.Logf(Debug, "Hello Secret World")
	require.Contains(t, buffer.String(), "Hello ******** World")
	// Ensure that the original message is not logged and secret is wiped out
	require.NotContains(t, buffer.String(), "Secret")

	logger.Setup(WithMapper(func(s string) string {
		emailRegEx := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}`)
		return emailRegEx.ReplaceAllString(s, "********")
	}))

	logger.Logf(Debug, "User alice@gmail.com entered the system")
	require.Contains(t, buffer.String(), "User ******** entered the system")
	require.NotContains(t, buffer.String(), "alice@gmail.com")
}

func TestLoggerLogsWithArgs(t *testing.T) {
	t.Parallel()
	buffer := bytes.NewBuffer(nil)
	logger := New(CustomLogOut(buffer), StdFormat, Level(Debug))
	logger.Logf(Debug, "Hello %s", "World")
	require.Contains(t, buffer.String(), "Hello World")
}

func FuzzLoggerInputs(f *testing.F) {
	f.Add("yagl")

	f.Fuzz(func(t *testing.T, s string) {
		fmt.Println("XXX", s)
		buffer := bytes.NewBuffer(nil)
		logger := New(CustomLogOut(buffer), StdFormat, Level(Debug))
		logger.Logf(Debug, s)
		require.Contains(t, buffer.String(), s)
	})
}
