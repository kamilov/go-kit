package tracer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tracer, err := New(
		context.Background(),
		WithHost("localhost"),
		WithPort(9999),
		WithName("test"),
	)

	assert.Nil(t, err)
	assert.NotPanics(t, tracer.Sync)
}

func TestNewNoop(t *testing.T) {
	tracer := NewNoop("noop")

	assert.NotNil(t, tracer)
	assert.NotPanics(t, tracer.Sync)
}
