package endpoint_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kamilov/go-kit/endpoint"
)

type Counter int

func testEndpoint(_ context.Context, counter Counter) (Counter, error) {
	counter++
	if counter == 0 {
		return counter, errors.New("counter is zero")
	}

	return counter, nil
}

func testMiddleware(
	next func(ctx context.Context, counter Counter) (Counter, error),
) func(ctx context.Context, counter Counter) (Counter, error) {
	return func(ctx context.Context, counter Counter) (Counter, error) {
		counter++
		return next(ctx, counter)
	}
}

func TestEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint endpoint.Endpoint[Counter, Counter]
		input    Counter
		output   Counter
		isError  bool
	}{
		{
			"simple endpoint",
			testEndpoint,
			Counter(1),
			Counter(2),
			false,
		},
		{
			"endpoint with middleware",
			testMiddleware(testEndpoint),
			Counter(1),
			Counter(3),
			false,
		},
		{
			"endpoint with middleware and expected error",
			testMiddleware(testEndpoint),
			Counter(-2),
			Counter(0),
			true,
		},
	}

	ctx := context.Background()

	var (
		output Counter
		err    error
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			output, err = test.endpoint(ctx, test.input)

			if test.isError {
				if err == nil {
					t.Fatal("error is not received")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				} else if output != test.output {
					t.Fatalf("got %v, want %v", output, test.output)
				}
			}
		})
	}
}
