package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/syntaxfa/syntax-backend/logger"
)

// TODO - defaultYamlFilePath has some problem in testing environment ut have normal behavior in build environment.
/*
	How can this problem be fixed?
		In the test environment, our binary file is located in a different directory.
		That's why it can't recognize the config.yml file. To solve this problem, we can embed the config file.
*/
const (
	defaultPrefix       = "SYNTAX_"
	defaultDelimiter    = "."
	defaultSeparator    = "__"
	defaultYamlFilePath = "config.yml"
)

var c Config

type Options struct {
	Prefix       string
	Delimiter    string
	Separator    string
	YamlFilePath string
	CallBackEnV  func(string) string
}

func defaultCallBackEnv(source string) string {
	fmt.Println("touch default call back env")
	base := strings.ToLower(strings.TrimPrefix(source, defaultPrefix))

	return strings.ReplaceAll(base, defaultSeparator, defaultDelimiter)
}

func init() {
	k := koanf.New(defaultDelimiter)

	// load default configuration from Default function
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		logger.L().Error("error loading default config", "error", err.Error())
	}

	// load configuration from yaml file
	if err := k.Load(file.Provider(defaultYamlFilePath), yaml.Parser()); err != nil {
		logger.L().Error("error loading config from 'config.yml'", "file path", defaultYamlFilePath, "error", err.Error())
	}

	// load from environment variable
	if err := k.Load(env.Provider(defaultPrefix, defaultDelimiter, defaultCallBackEnv), nil); err != nil {
		logger.L().Error("error loading environment variables", "error", err.Error())
	}

	if err := k.Unmarshal("", &c); err != nil {
		logger.L().Error("error unmarshalling config", "error", err.Error())
	}
}

func C() *Config {
	return &c
}

func New(opt Options) *Config {
	k := koanf.New(opt.Delimiter)

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		logger.L().Error("error loading config from Default()", "error", err.Error())
	}

	if err := k.Load(file.Provider(opt.YamlFilePath), yaml.Parser()); err != nil {
		logger.L().Error("error loading from `config.yml`", "error", err.Error())
	}

	if err := k.Load(env.Provider(opt.Prefix, opt.Delimiter, opt.CallBackEnV), nil); err != nil {
		logger.L().Error("error loading from environment variables", "error", err.Error())
	}

	if err := k.Unmarshal("", &c); err != nil {
		logger.L().Error("error unmarshalling config", "error", err.Error())
	}

	return &c
}
