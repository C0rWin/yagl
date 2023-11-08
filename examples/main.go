package main

import "github.com/c0rwin/yagl"

func main() {
	// Create a new logger with default options
	// by default logger is initialized with
	// Info level, hence Debug information won't be available
	logger := yagl.New()

	// Output a log message with Info level
	// The output message should be similar to this:

	// "[2023-11-09 00:56:33] [INFO]: This is a log message"
	logger.Logf(yagl.Info, "This is a log message")

	// Output a log message with Debug level
	// no message should be printed since the default level is Info
	logger.Logf(yagl.Debug, "This is a debug message")
}
