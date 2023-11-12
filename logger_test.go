package yagl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerHierarchy(t *testing.T) {
	for _, logLevel := range []LogLevel{Debug, Info, Warn, Error} {
		for _, messageLogLevel := range []LogLevel{Debug, Info, Warn, Error} {
			t.Run(logLevel.String()+"-"+messageLogLevel.String(), func(t *testing.T) {
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
