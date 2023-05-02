package utils

import (
	"os"
	"path/filepath"
)

func IsDir(path string) bool {
	if pathAbs, err := filepath.Abs(path); err == nil {
		if fileInfo, err := os.Stat(pathAbs); !os.IsNotExist(err) && fileInfo.IsDir() {
			return true
		}
	}
	return false
}