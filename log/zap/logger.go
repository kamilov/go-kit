// Package zap implementation logger with zap handler
package zap

import (
	"os"

	"github.com/kamilov/go-kit/log" //nolint:depguard // safe dependency
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	*options
	level  zap.AtomicLevel
	logger *zap.SugaredLogger
}

// New create new zap logger
func New(level log.Level, opts ...Option) log.Logger {
	//nolint:exhaustruct // use default values
	l := logger{
		level: zap.NewAtomicLevelAt(convertLevel(level)),
		options: &options{
			output:  os.Stdout,
			encoder: defaultEncoder(),
			zapOptions: []zap.Option{
				zap.AddCallerSkip(1),
				zap.AddCaller(),
			},
		},
	}

	for _, opt := range opts {
		opt.apply(l.options)
	}

	core := zapcore.NewCore(
		l.encoder,
		zapcore.AddSync(l.output),
		l.level,
	)

	l.logger = zap.New(core, l.zapOptions...).Sugar()

	return l
}

// NewNoop create and return logger for nullable output
func NewNoop() log.Logger {
	//nolint:exhaustruct // use default values
	return logger{
		level:  zap.NewAtomicLevel(),
		logger: zap.NewNop().Sugar(),
	}
}

func (l logger) With(args ...any) log.Logger {
	//nolint:exhaustruct // use default values
	return logger{
		level:  l.level,
		logger: l.logger.With(args...),
	}
}

func (l logger) WithLevel(level log.Level) log.Logger {
	zapLevel := zap.NewAtomicLevelAt(convertLevel(level))

	//nolint:exhaustruct // use default values
	return logger{
		level: zapLevel,
		logger: zap.New(
			zapcore.NewCore(
				l.encoder,
				zapcore.AddSync(l.output),
				zapLevel,
			),
			l.zapOptions...,
		).Sugar(),
	}
}

func (l logger) Level(level log.Level) log.Logger {
	l.level.SetLevel(convertLevel(level))

	return l
}

func (l logger) Debug(message string, args ...any) {
	l.logger.Debugw(message, args...)
}

func (l logger) Debugf(message string, args ...any) {
	l.logger.Debugf(message, args...)
}

func (l logger) Info(message string, args ...any) {
	l.logger.Infow(message, args...)
}

func (l logger) Infof(message string, args ...any) {
	l.logger.Infof(message, args...)
}

func (l logger) Warning(message string, args ...any) {
	l.logger.Warnw(message, args...)
}

func (l logger) Warningf(message string, args ...any) {
	l.logger.Warnf(message, args...)
}

func (l logger) Error(message string, args ...any) {
	l.logger.Errorw(message, args...)
}

func (l logger) Errorf(message string, args ...any) {
	l.logger.Errorf(message, args...)
}

func (l logger) Panic(message string, args ...any) {
	l.logger.Panicw(message, args...)
}

func (l logger) Panicf(message string, args ...any) {
	l.logger.Warnf(message, args...)
}

func (l logger) Fatal(message string, args ...any) {
	l.logger.Fatalw(message, args...)
}

func (l logger) Fatalf(message string, args ...any) {
	l.logger.Warnf(message, args...)
}

func convertLevel(level log.Level) zapcore.Level {
	switch level {
	case log.DebugLevel:
		return zapcore.DebugLevel
	case log.InfoLevel:
		return zapcore.InfoLevel
	case log.WarningLevel:
		return zapcore.WarnLevel
	case log.ErrorLevel:
		return zapcore.ErrorLevel
	}

	return zapcore.InfoLevel
}

func defaultEncoder() zapcore.Encoder {
	//nolint:exhaustruct // use default values
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "context",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}
