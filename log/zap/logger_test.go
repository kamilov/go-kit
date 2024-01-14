package zap

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/kamilov/go-kit/log"
	"github.com/stretchr/testify/assert"
)

type clock time.Time

func (c clock) Now() time.Time {
	return time.Time(c).UTC()
}

func (c clock) NewTicker(duration time.Duration) *time.Ticker {
	return time.NewTicker(duration)
}

func newTestClock(sec int64) clock {
	return clock(time.Unix(sec, 0))
}

func Test_New(t *testing.T) {
	l := New(log.DebugLevel)

	assert.NotNil(t, l)
}

func Test_NewNoop(t *testing.T) {
	l := NewNoop()

	assert.NotNil(t, l)
}

func TestLogger_Log(t *testing.T) {
	calls := []struct {
		name string
		fn   func(l log.Logger, message string, args ...any)
	}{
		{
			name: "DEBUG",
			fn: func(l log.Logger, message string, args ...any) {
				l.Debug(message, args...)
			},
		},
		{
			name: "INFO",
			fn: func(l log.Logger, message string, args ...any) {
				l.Info(message, args...)
			},
		},
		{
			name: "WARN",
			fn: func(l log.Logger, message string, args ...any) {
				l.Warning(message, args...)
			},
		},
		{
			name: "ERROR",
			fn: func(l log.Logger, message string, args ...any) {
				l.Error(message, args...)
			},
		},
	}

	tests := []struct {
		name    string
		message string
		args    []any
		expect  string
	}{
		{
			"Without Args",
			"test message",
			[]any{},
			"{\"level\":\"LEVEL\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message\",\"app\":\"test\"}\n",
		},
		{
			"With Args",
			"test message",
			[]any{
				"arg",
				"value",
			},
			"{\"level\":\"LEVEL\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message\",\"app\":\"test\",\"arg\":\"value\"}\n",
		},
	}

	var buffer bytes.Buffer

	l := New(
		log.DebugLevel,
		WithOutput(&buffer),
		WithClock(newTestClock(0)),
		WithCaller(false),
		WithField("app", "test"),
	)

	for _, test := range tests {
		for _, call := range calls {
			t.Run(call.name+" - "+test.name, func(t *testing.T) {
				call.fn(l, test.message, test.args...)

				assert.NotEmpty(t, buffer.String())
				assert.Equal(t, strings.Replace(test.expect, "LEVEL", call.name, 1), buffer.String())

				buffer.Reset()
			})
		}
	}
}

func TestLogger_Logf(t *testing.T) {
	calls := []struct {
		name string
		fn   func(l log.Logger, message string, args ...any)
	}{
		{
			name: "DEBUG",
			fn: func(l log.Logger, message string, args ...any) {
				l.Debugf(message, args...)
			},
		},
		{
			name: "INFO",
			fn: func(l log.Logger, message string, args ...any) {
				l.Infof(message, args...)
			},
		},
		{
			name: "WARN",
			fn: func(l log.Logger, message string, args ...any) {
				l.Warningf(message, args...)
			},
		},
		{
			name: "ERROR",
			fn: func(l log.Logger, message string, args ...any) {
				l.Errorf(message, args...)
			},
		},
	}

	tests := []struct {
		name    string
		message string
		args    []any
		expect  string
	}{
		{
			"Without Args",
			"test message %s %d",
			[]any{"a", 1},
			"{\"level\":\"LEVEL\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message a 1\",\"app\":\"test\"}\n",
		},
		{
			"With Args",
			"test message %s",
			[]any{
				"test",
				"arg",
				"value",
			},
			"{\"level\":\"LEVEL\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message test%!(EXTRA string=arg, string=value)\",\"app\":\"test\"}\n",
		},
	}

	var buffer bytes.Buffer

	l := New(
		log.DebugLevel,
		WithOutput(&buffer),
		WithClock(newTestClock(0)),
		WithCallerSkip(1),
		WithCaller(false),
		WithField("app", "test"),
	)

	for _, test := range tests {
		for _, call := range calls {
			t.Run(call.name+" - "+test.name, func(t *testing.T) {
				call.fn(l, test.message, test.args...)

				assert.NotEmpty(t, buffer.String())
				assert.Equal(t, strings.Replace(test.expect, "LEVEL", call.name, 1), buffer.String())

				buffer.Reset()
			})
		}
	}
}

func TestLogger_With(t *testing.T) {
	var buffer bytes.Buffer

	l := New(
		log.DebugLevel,
		WithOutput(&buffer),
		WithClock(newTestClock(0)),
		WithCaller(false),
	)

	assert.NotNil(t, l)

	l.Info("test message")

	assert.Equal(t, "{\"level\":\"INFO\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message\"}\n", buffer.String())

	buffer.Reset()

	l = l.With("key", "value")

	l.Warning("test message")

	assert.Equal(t, "{\"level\":\"WARN\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message\",\"key\":\"value\"}\n", buffer.String())

	buffer.Reset()

	l.With("struct", struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}{"a", 100}).Error("test message")

	assert.Equal(t, "{\"level\":\"ERROR\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message\",\"key\":\"value\",\"struct\":{\"key\":\"a\",\"value\":100}}\n", buffer.String())
}

func TestLogger_Level(t *testing.T) {
	var buffer bytes.Buffer

	l := New(
		log.DebugLevel,
		WithOutput(&buffer),
		WithClock(newTestClock(0)),
		WithCaller(false),
	)

	assert.NotNil(t, l)

	l.Debug("test message")

	assert.Equal(t, "{\"level\":\"DEBUG\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test message\"}\n", buffer.String())

	buffer.Reset()

	l = l.WithLevel(log.InfoLevel)

	l.Debug("test debug message")
	l.Infof("test info message")

	assert.Equal(t, "{\"level\":\"INFO\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test info message\"}\n", buffer.String())

	buffer.Reset()

	l.Level(log.ErrorLevel)

	l.Debug("test debug message")
	l.Infof("test info message")
	l.Warning("test warn message")
	l.Error("test error message")

	assert.Equal(t, "{\"level\":\"ERROR\",\"time\":\"1970-01-01T00:00:00Z\",\"message\":\"test error message\"}\n", buffer.String())
}
