package db

import (
	"time"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"github.com/go-sql-driver/mysql"
)

func (c *Config) MySQLConfig() ([]*mysql.Config, error) {
	if c.Timezone == "" {
		c.Timezone = "UTC"
	}
	location, err := time.LoadLocation(c.Timezone)
	if err != nil {
		return nil, err
	}
	var cfgs []*mysql.Config
	for _, a := range c.Addresses {
		cfg := mysql.NewConfig()
		cfg.User = a.User
		cfg.Passwd = a.Passwd
		cfg.Net = "tcp"
		cfg.Addr = a.Addr
		cfg.Collation = "utf8mb4_unicode_ci"
		cfg.DBName = c.DBName

		cfg.ParseTime = true
		cfg.Loc = location
		cfgs = append(cfgs, cfg)
	}
	return cfgs, nil
}

func NewMySQLDB(cfg Config) (*DB, error) {
	var dsns []string
	configs, err := cfg.MySQLConfig()
	if err != nil {
		return nil, errorx.Wrap(err)
	}
	for i := range configs {
		dsns = append(dsns, configs[i].FormatDSN())
	}

	db, err := open("mysql", dsns)
	if err != nil {
		return nil, errorx.Wrap(err)
	}

	if cfg.MaxIdleConns != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	connMaxLifetime := 6 * time.Hour
	if cfg.ConnMaxLifetime != "" {
		connMaxLifetime, err = time.ParseDuration(cfg.ConnMaxLifetime)
		if err != nil {
			return nil, errorx.Wrap(err)
		}
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	return db, nil
}
