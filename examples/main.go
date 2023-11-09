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

	// Update log level to Debug
	logger.SetOptions(yagl.LogLevelOption(yagl.Debug))

	// After, updating the log level to Debug
	// the following message should be printed:
	// "[2023-11-09 00:56:33] [DEBUG]: This is a debug message"
	logger.Logf(yagl.Debug, "This is a debug message")

	// Now we can also change to more detailed format to contain
	// package name and function name
	logger.SetOptions(yagl.DebugFormatOption)

	// The output should be similar to this:
	// "[2023-11-09 00:56:33] [DEBUG]: [main.main]: This is a debug message"
	logger.Logf(yagl.Debug, "This is a debug message")
}
