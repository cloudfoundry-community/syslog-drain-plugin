package command

// Logger is used for outputting results and errors
type Logger interface {
	Fatalf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Printf(format string, args ...interface{})
}
