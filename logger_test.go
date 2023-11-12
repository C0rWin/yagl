package yagl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		logLevel Setting
		log      func(l *Logger)
		assert   func(t *testing.T, actual, expected string)
	}{
		{
			name:     "Testing default logger trying to print a message with INFO level",
			message:  "Hello World",
			logLevel: Level(Info),
			log: func(l *Logger) {
				l.Logf(Info, "Hello World")
			},
			assert: func(t *testing.T, actual, expected string) {
				require.Contains(t, actual, expected)
				require.Contains(t, actual, Info.String())
			},
		},
		{
			name:     "Testing default logger trying to print a message with DEBUG level",
			message:  "Hello World",
			logLevel: Level(Info),
			log: func(l *Logger) {
				l.Logf(Debug, "Hello World")
			},
			assert: func(t *testing.T, actual, expected string) {
				require.Empty(t, actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(nil)
			logger := New(CustomLogOut(buffer), StdFormat, test.logLevel)
			test.log(logger)
			test.assert(t, buffer.String(), test.message)
		})
	}
}
