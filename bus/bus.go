package bus

import (
	"fmt"
	"reflect"
	"sync"
)

type (
	handler struct {
		callback reflect.Value
		queue    chan []reflect.Value
	}

	// Bus implements publish/subscribe messaging paradigm.
	Bus interface {
		// Publish publishes arguments to the given topic subscribers
		// Publish block only when the buffer of one of the subscribers is full.
		Publish(topic string, args ...any)
		// Close unsubscribe all handlers from given topic
		Close(topic string)
		// Subscribe subscribes to the given topic
		Subscribe(topic string, callback any) error
		// Unsubscribe unsubscribes handler from the given topic
		Unsubscribe(topic string, callback any) error
	}

	bus struct {
		queueSize uint
		handlers  map[string][]*handler
		mutex     sync.RWMutex
	}
)

func New(queueSize uint) Bus {
	return &bus{
		queueSize: queueSize,
		handlers:  make(map[string][]*handler),
	}
}

func (b *bus) Publish(topic string, args ...any) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	handlers, ok := b.handlers[topic]
	if !ok {
		return
	}

	callbackArgs := buildArgs(args)

	for _, h := range handlers {
		h.queue <- callbackArgs
	}
}

func (b *bus) Close(topic string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.handlers[topic]; !ok {
		return
	}

	for _, h := range b.handlers[topic] {
		close(h.queue)
	}

	delete(b.handlers, topic)
}

func (b *bus) Subscribe(topic string, callback any) error {
	if err := isValidHandler(callback); err != nil {
		return err
	}

	h := &handler{
		callback: reflect.ValueOf(callback),
		queue:    make(chan []reflect.Value, b.queueSize),
	}

	go func() {
		for args := range h.queue {
			h.callback.Call(args)
		}
	}()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[topic] = append(b.handlers[topic], h)

	return nil
}

func (b *bus) Unsubscribe(topic string, callback any) error {
	if err := isValidHandler(callback); err != nil {
		return err
	}

	rv := reflect.ValueOf(callback)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.handlers[topic]; !ok {
		return fmt.Errorf("handler for topic %s not found", topic)
	}

	for i, h := range b.handlers[topic] {
		//nolint:govet // using this compare
		if h.callback == rv {
			close(h.queue)

			if len(b.handlers[topic]) == 1 {
				delete(b.handlers, topic)
			} else {
				b.handlers[topic] = append(b.handlers[topic][:i], b.handlers[topic][i+1:]...)
			}
		}
	}

	return nil
}

func isValidHandler(callback any) error {
	if rt := reflect.TypeOf(callback); rt.Kind() != reflect.Func {
		return fmt.Errorf("%s must be a function", rt)
	}

	return nil
}

func buildArgs(args []any) []reflect.Value {
	result := make([]reflect.Value, len(args))

	for i, arg := range args {
		result[i] = reflect.ValueOf(arg)
	}

	return result
}
