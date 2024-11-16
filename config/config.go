package config

import "github.com/syntaxfa/syntax-backend/logger"

type Config struct {
	Logger logger.Config `koanf:"logger"`
}
