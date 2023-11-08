# YAGL (Yet Another Go Logger)

YAGL is a lightweight, powerful, and flexible logging library for Go. It provides an alternative implementation for application logging, focusing on simplicity and customizability.

## Features

- **Simple API**: Easy to use and understand.
- **Highly Configurable**: Customize logging levels, formats, and destinations.
- **Asynchronous Logging**: Improve performance by logging in the background.
- **Structured Logging**: Support for structured logging formats like JSON.
- **Extensible**: Easily extend with custom handlers and formatters.

## Getting Started

### Prerequisites

- Go version 1.15 or higher

### Installation

To start using YAGL, install Go and run `go get`:

```sh
go get -u github.com/C0rWin/yagl
```

This will retrieve the library.

## Usage

Here's a simple example of how to use YAGL:

```go
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
	logger.SetOptions(yagl.WithDebug)

	// The output should be similar to this:
	// "[2023-11-09 00:56:33] [DEBUG]: [main.main]: This is a debug message"
	logger.Logf(yagl.Debug, "This is a debug message")
}
```

## Examples

You can find more examples in the `/examples` directory within this repository.

## License

YAGL is released under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Authors

- Artem Barger - *Initial work* - [YourUsername](https://github.com/C0rWin)

See also the list of [contributors](https://github.com/C0rWin/yagl/contributors) who participated in this project.


For any additional questions or feedback, please contact the maintainers directly or open an issue in the GitHub repository.
```
