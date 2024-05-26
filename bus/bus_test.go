package bus_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/kamilov/go-kit/bus"
)

const testTopic = "topic"

func TestNew(t *testing.T) {
	bus := bus.New(uint(runtime.NumCPU()))

	if bus == nil {
		t.Fatal("could not create bus")
	}
}

func TestBus_Subscribe(t *testing.T) {
	bus := bus.New(uint(runtime.NumCPU()))

	if bus.Subscribe(testTopic, func() {}) != nil {
		t.Fatal("could not subscribe")
	}

	if bus.Subscribe(testTopic, testTopic) == nil {
		t.Fatal("expected error")
	}
}

func TestBus_Unsubscribe(t *testing.T) {
	bus := bus.New(uint(runtime.NumCPU()))
	handler := func() {}
	handler2 := func(x int) int { return x }

	if err := bus.Subscribe(testTopic, handler); err != nil {
		t.Fatal(err)
	}

	if err := bus.Subscribe(testTopic, handler2); err != nil {
		t.Fatal(err)
	}

	if err := bus.Unsubscribe(testTopic, handler); err != nil {
		t.Fatal(err)
	}

	if err := bus.Unsubscribe(testTopic, handler2); err != nil {
		t.Fatal(err)
	}

	if bus.Unsubscribe(testTopic, handler) == nil {
		t.Fatal("expected error")
	}

	if bus.Unsubscribe(testTopic, testTopic) == nil {
		t.Fatal("expected error")
	}
}

func TestBus_Close(t *testing.T) {
	bus := bus.New(uint(runtime.NumCPU()))
	handler := func() {}

	if err := bus.Subscribe(testTopic, handler); err != nil {
		t.Fatal(err)
	}

	bus.Close(testTopic)
	bus.Close(testTopic + "-non")

	if err := bus.Unsubscribe(testTopic, handler); err == nil {
		t.Fatal("expected error")
	}
}

func TestBus_Publish(t *testing.T) {
	bus := bus.New(uint(runtime.NumCPU()))

	var (
		first, second bool
		wg            sync.WaitGroup
	)
	wg.Add(2)

	handler := func(x *bool) func(v bool) {
		return func(v bool) {
			wg.Done()
			*x = v
		}
	}

	if err := bus.Subscribe(testTopic, handler(&first)); err != nil {
		t.Fatal(err)
	}

	if err := bus.Subscribe(testTopic, handler(&second)); err != nil {
		t.Fatal(err)
	}

	bus.Publish(testTopic, true)
	bus.Publish(testTopic+"-non", true)

	wg.Wait()

	if first == false || second == false {
		t.Fatal("expected true")
	}
}
