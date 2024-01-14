package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	testConfigWithoutTags struct {
		Host string
		Port int
		Pass string
	}

	testConfigWithSkippedName struct {
		host string `env:"HOST"`
		Port int    `env:"-"`
		Pass string `env:",secret"`
	}

	Embedded struct {
		H string `env:"HOST"`
		P int    `env:"PORT"`
	}

	testConfigWithEmbedded struct {
		Embedded
		embedded Embedded
	}
)

func testLookup(name string) (string, bool) {
	data := map[string]string{
		"HOST": "localhost",
		"PORT": "80",
		"PASS": "test",
	}

	value, ok := data[name]

	return value, ok
}

func testLookupPrefix(name string) (string, bool) {
	data := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "80",
		"APP_PASS": "test",
	}

	value, ok := data[name]

	return value, ok
}

func TestLoader_Load(t *testing.T) {
	var config testConfigWithoutTags

	env := New(
		WithLookup(testLookup),
	)

	err := env.Load(&config)

	assert.Nil(t, err)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 80, config.Port)
	assert.Equal(t, "test", config.Pass)

	assert.Equal(t, ErrStructPointer, env.Load(config))

	var configNilPointer *testConfigWithoutTags

	assert.Equal(t, ErrNilPointer, env.Load(configNilPointer))

	var configWithSkipped testConfigWithSkippedName

	err = env.Load(&configWithSkipped)

	assert.Nil(t, err)
	assert.Equal(t, "", configWithSkipped.host)
	assert.Equal(t, 0, configWithSkipped.Port)
	assert.Equal(t, "test", configWithSkipped.Pass)

	var configWithEmbedded testConfigWithEmbedded

	err = env.Load(&configWithEmbedded)

	assert.Nil(t, err)
	assert.Equal(t, "localhost", configWithEmbedded.Embedded.H)
	assert.Equal(t, 80, configWithEmbedded.Embedded.P)
	assert.Equal(t, "", configWithEmbedded.embedded.H)
	assert.Equal(t, 0, configWithEmbedded.embedded.P)
}

func Test_getName(t *testing.T) {
	tests := []struct {
		name   string
		tag    string
		field  string
		expect string
		secret bool
	}{
		{"without tag", "", "Name", "NAME", false},
		{"without tag snake", "", "NameName", "NAME_NAME", false},
		{"without tag secret", ",secret", "Name", "NAME", true},
		{"tag", "NaMe", "Name", "NaMe", false},
		{"tag secret", "NaMe,secret", "Name", "NaMe", true},
		{"more comma", "Name,Comma", "Name", "Name,Comma", false},
		{"more comma secret", "Name,Comma,secret", "Name", "Name,Comma", true},
	}

	for _, test := range tests {
		name, secret := getName(test.tag, test.field)

		assert.Equal(t, test.expect, name, test.name)
		assert.Equal(t, test.secret, secret, test.name)
	}
}
