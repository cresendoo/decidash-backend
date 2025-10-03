package apiserver

import (
	"log/slog"

	"github.com/cresendoo/decidash-backend/pkg/config"
	"github.com/cresendoo/decidash-backend/pkg/db"
	"github.com/cresendoo/decidash-backend/pkg/utils"
)

type Config struct {
	Log struct {
		Level  slog.Level `yaml:"level"`
		Format string     `yaml:"format"`
		File   string     `yaml:"file"`
	} `yaml:"log"`

	Port      string `yaml:"port"`
	SentryDSN string `yaml:"sentry_dsn"`

	DB    db.Config `yaml:"db"`
	Redis struct {
		Addr string `yaml:"addr"`
		Pool int    `yaml:"pool"`
		DB   int    `yaml:"db"`
	} `yaml:"redis"`

	AptosAccounts struct {
		FeePayer string `yaml:"fee_payer"`
	} `yaml:"aptos_accounts"`
}

func (c *Config) Load() error {
	return config.LoadConfig(c)
}

func (c *Config) FileName() string {
	if utils.IsProductionPhase() {
		return "config-api-server-prod.yaml"
	}
	return "config-api-server-dev.yaml"
}
