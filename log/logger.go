// Package log describe logger
package log

// Level type for error levels
type Level int8

const (
	// DebugLevel lever for print debug messages
	DebugLevel Level = iota
	// InfoLevel lever for print info messages
	InfoLevel
	// WarningLevel lever for print warning messages
	WarningLevel
	// ErrorLevel lever for print error messages
	ErrorLevel
)

// Logger base logger interface
//
//nolint:interfacebloat // normal count methods to interface
type Logger interface {
	With(args ...any) Logger

	WithLevel(level Level) Logger
	Level(level Level) Logger

	Debug(message string, args ...any)
	Debugf(message string, args ...any)

	Info(message string, args ...any)
	Infof(message string, args ...any)

	Warning(message string, args ...any)
	Warningf(message string, args ...any)

	Error(message string, args ...any)
	Errorf(message string, args ...any)

	Panic(message string, args ...any)
	Panicf(message string, args ...any)

	Fatal(message string, args ...any)
	Fatalf(message string, args ...any)
}
