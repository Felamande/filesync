package utils

import (
    "path/filepath"
    "github.com/Felamande/filesync/settings"
)

func Abs(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(settings.Folder, path)
	}
	return path
}