package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/utils/reflect"
)

type (
	Reader               = func(ctx context.Context, data any) error
	ReaderWithoutContext = func(data any) error

	ReaderConstraint interface {
		Reader | ReaderWithoutContext
	}
)

//nolint:gochecknoglobals // used to register readers
var readers []Reader

func RegisterReader[T ReaderConstraint](reader T) {
	switch r := any(reader).(type) {
	case ReaderWithoutContext:
		readers = append(readers, func(_ context.Context, data any) error {
			return r(data)
		})

	case Reader:
		readers = append(readers, r)
	}
}

func Read(ctx context.Context, data any) error {
	for _, reader := range readers {
		if err := reader(ctx, data); err != nil {
			return err
		}
	}

	return nil
}

func ReadFile(ctx context.Context, path string, data any) error {
	if err := reflect.ValidateStructPointer(data); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		return fmt.Errorf("file open error: %w", err)
	}

	defer file.Close()

	extension := strings.TrimLeft(strings.ToLower(filepath.Ext(path)), ".")
	decoder := coder.GetDecoder(extension)
	if decoder == nil {
		return fmt.Errorf("not found decoder for file extension: %s", extension)
	}

	err = decoder(ctx, file, data)
	if err != nil {
		return fmt.Errorf("decoder error: %w", err)
	}

	return Read(ctx, data)
}
