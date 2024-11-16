package config_test

import (
	"embed"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/syntaxfa/syntax-backend/config"
	"github.com/syntaxfa/syntax-backend/logger"
	"os"
	"reflect"
	"strings"
	"testing"
)

//go:embed test.yml
var configYML embed.FS

type dbConfig struct {
	Host     string `koanf:"host"`
	Password string `koanf:"password"`
	User     string `koanf:"user"`
}

type structConfig struct {
	Debug    bool     `koanf:"debug"`
	FilePath string   `koanf:"file_path"`
	DB       dbConfig `koanf:"db"`
}

const (
	prefix    = "SYNTAX_"
	delimiter = "."
	separator = "__"
)

func callBackEnv(source string) string {
	base := strings.ToLower(strings.TrimPrefix(source, prefix))

	return strings.ReplaceAll(base, separator, delimiter)
}

func TestLoadingDefaultConfigFromStruct(t *testing.T) {
	k := koanf.New(delimiter)

	testStruct := structConfig{
		Debug:    true,
		FilePath: "config.yml",
		DB: dbConfig{
			Host: "localhost",
		},
	}

	if err := k.Load(structs.Provider(testStruct, "koanf"), nil); err != nil {
		t.Fatalf("error loading default config. error: %s", err.Error())
	}

	var c structConfig

	if err := k.Unmarshal("", &c); err != nil {
		t.Fatalf("error marshalling config. error: %s", err.Error())
	}

	if !reflect.DeepEqual(c, testStruct) {
		t.Fatalf("expected: %+v, got: %+v", testStruct, c)
	}
}

func TestLoadingConfigFromFile(t *testing.T) {
	k := koanf.New(delimiter)

	tmpFile, err := os.CreateTemp("", "test.yml")
	defer tmpFile.Close()

	data, err := configYML.ReadFile("test.yml")
	if err != nil {
		t.Fatalf("error load config from embedded file, error: %s", err.Error())
	}

	if _, err := tmpFile.Write(data); err != nil {
		t.Fatalf("error write to tmpFile, error: %s", err.Error())
	}

	if err := k.Load(file.Provider(tmpFile.Name()), yaml.Parser()); err != nil {
		t.Fatalf("error load config from file. error: %s", err.Error())
	}

	var c config.Config

	if err := k.Unmarshal("", &c); err != nil {
		t.Fatalf("error unmarshalling from file. error: %s", err.Error())
	}

	want := config.Config{
		Logger: logger.Config{FilePath: "logs/logs.json"},
	}

	if !reflect.DeepEqual(want, c) {
		t.Fatalf("expected: %+v, got: %+v", want, c)
	}
}

func TestLoadConfigFromEnvironmentVariable(t *testing.T) {
	k := koanf.New(".")

	os.Setenv("SYNTAX_LOGGER__FILE_PATH", "new_file_path.yml")

	if err := k.Load(env.Provider(prefix, delimiter, callBackEnv), nil); err != nil {
		t.Fatalf("error environment variables, error: %s", err.Error())
	}

	var instance config.Config
	if err := k.Unmarshal("", &instance); err != nil {
		t.Fatalf("error unmarshalling config, error: %s", err.Error())
	}

	want := config.Config{
		Logger: logger.Config{FilePath: "new_file_path.yml"},
	}

	if !reflect.DeepEqual(want, instance) {
		t.Fatalf("expected: %+v, got: %+v", want, instance)
	}
}
