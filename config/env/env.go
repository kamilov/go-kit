package env

import (
	"context"
	"os"

	"github.com/kamilov/go-kit/config"
	"github.com/kamilov/go-kit/utils/structure"
)

type (
	LookupFunc func(string) (string, bool)
	contextKey int
)

const (
	TagName contextKey = iota
	Prefix
	defaultTagName = "env"
)

type EnvData struct {
	prefix string
}

func (d EnvData) Get(key string) string {
	return os.Getenv(d.prefix + key)
}

//nolint:gochecknoinits // used for automatic adding read config function
func init() {
	config.RegisterReader(reader)
}

func reader(ctx context.Context, data any) error {
	tagName := getTagName(ctx)
	prefix := getPrefix(ctx)

	env := EnvData{prefix}
	dec := structure.NewDecoder(tagName)

	return dec.Decode(env, data)
}

func getTagName(ctx context.Context) string {
	if value, ok := ctx.Value(TagName).(string); ok {
		return value
	}

	return defaultTagName
}

func getPrefix(ctx context.Context) string {
	if value, ok := ctx.Value(Prefix).(string); ok {
		return value
	}

	return ""
}
