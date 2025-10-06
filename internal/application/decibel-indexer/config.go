package decibelindexer

import (
	"log/slog"

	"github.com/cresendoo/decidash-backend/pkg/config"
	"github.com/cresendoo/decidash-backend/pkg/utils"
)

type Config struct {
	Log struct {
		Level  slog.Level `yaml:"level"`
		Format string     `yaml:"format"`
		File   string     `yaml:"file"`
	} `yaml:"log"`

	SentryDSN string `yaml:"sentry_dsn"`

	DB    string `yaml:"db"`
	Redis struct {
		Addr string `yaml:"addr"`
		Pool int    `yaml:"pool"`
		DB   int    `yaml:"db"`
	} `yaml:"redis"`
}

func (c *Config) Load() error {
	return config.LoadConfig(c)
}

func (c *Config) FileName() string {
	if utils.IsProductionPhase() {
		return "config-decibel-indexer-prod.yaml"
	}
	return "config-decibel-indexer-dev.yaml"
}
