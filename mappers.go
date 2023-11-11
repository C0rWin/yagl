package yagl

// Mapper is a function that maps a string to another string
// might be used for example to mutate string and remove
// sensitive data
type Mapper func(string) string

// NoOpMapper is a mapper that does nothing
var noOpMapper Mapper = func(s string) string { return s }
