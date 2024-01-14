// Package env environment values parser
package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/kamilov/go-kit/utils" //nolint:depguard // safe dependency
)

type (
	// LookupFunc func definition for lookup environment variables
	LookupFunc func(name string) (string, bool)

	// Loader configuration structure for loading environment values
	Loader struct {
		prefix string
		lookup LookupFunc
	}
)

const (
	tagName      = "env"
	secretSuffix = ",secret"
)

var (
	// ErrStructPointer error for unpoint struct cases
	ErrStructPointer = errors.New("must be a pointer to struct")
	// ErrNilPointer error for nullable point cases
	ErrNilPointer   = errors.New("the pointer should be nil")
	camelCaseRegexp = regexp.MustCompile("([^A-Z_])([A-Z])")
)

// New create new environment loader
func New(opts ...Option) *Loader {
	//nolint:exhaustruct // use default values, and apply with options
	loader := &Loader{
		lookup: os.LookupEnv,
	}

	for _, opt := range opts {
		opt.apply(loader)
	}

	return loader
}

// Load parse value from environment
//
//nolint:cyclop // normal cyclomatic
func (l *Loader) Load(target any) error {
	reflectValue := reflect.ValueOf(target)

	if reflectValue.Kind() != reflect.Ptr || !reflectValue.IsNil() && reflectValue.Elem().Kind() != reflect.Struct {
		return ErrStructPointer
	}

	if reflectValue.IsNil() {
		return ErrNilPointer
	}

	reflectValue = reflectValue.Elem()
	reflectType := reflectValue.Type()

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)

		if !field.CanSet() {
			continue
		}

		fieldType := reflectType.Field(i)

		if fieldType.Anonymous {
			field = utils.Indirect(field)

			if field.Kind() == reflect.Struct {
				if err := l.Load(field.Addr().Interface()); err != nil {
					return err
				}
			}

			continue
		}

		name, _ := getName(fieldType.Tag.Get(tagName), fieldType.Name)

		if name == "-" {
			continue
		}

		name = l.prefix + name

		if value, ok := l.lookup(name); ok {
			if err := utils.SetValue(field, value); err != nil {
				return fmt.Errorf("error reading \"%v\": %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

func getName(tag, field string) (string, bool) {
	name := strings.TrimSuffix(tag, secretSuffix)
	secret := strings.HasSuffix(tag, secretSuffix)

	if name == "" {
		name = strings.ToUpper(camelCaseRegexp.ReplaceAllString(field, "${1}_${2}"))
	}

	return name, secret
}
