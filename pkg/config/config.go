package config

import (
	"os"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"github.com/cresendoo/decidash-backend/pkg/utils"
	"gopkg.in/yaml.v3"
)

type IConfig interface {
	FileName() string
}

func LoadConfigWithFilename(fileName string, cfg IConfig) error {
	path, err := utils.FindFilePath(fileName)
	if err != nil {
		return err
	}
	data, err := os.Open(path)
	if err != nil {
		return errorx.Wrap(err).With("path", path)
	}
	err = yaml.NewDecoder(data).Decode(cfg)
	if err != nil {
		return errorx.Wrap(err).With("path", path)
	}
	return nil
}

func LoadConfig(cfg IConfig) error {
	fileName := cfg.FileName()
	return LoadConfigWithFilename(fileName, cfg)
}
