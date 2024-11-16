package exampletwo

import "github.com/syntaxfa/syntax-backend/config"

func LoadConfig() *config.Config {
	return config.C()
}
