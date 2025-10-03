package utils

import (
	"os"
	"path/filepath"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func FindFile(filename string, startDir string) (string, error) {
	foundPath := ""
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errorx.Wrap(err)
		}
		if !info.IsDir() && info.Name() == filename {
			foundPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if foundPath == "" {
		return "", os.ErrNotExist
	}
	return foundPath, nil
}

func FindFileInParentDirs(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		filePath := filepath.Join(dir, filename)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}

		// 상위 디렉토리로 이동
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// 루트 디렉토리에 도달
			break
		}
		dir = parentDir
	}

	return "", os.ErrNotExist
}
