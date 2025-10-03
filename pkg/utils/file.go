package utils

import (
	"io/fs"
	"path/filepath"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func FindFilePath(fileName string) (string, error) {
	var filePath string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if d.Name() == fileName {
			filePath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", errorx.Wrap(err).With("file_name", fileName)
	}
	return filePath, nil
}
