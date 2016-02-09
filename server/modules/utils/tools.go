package utils

import (
    "path/filepath"
    "github.com/Felamande/filesync/settings"
    "hash/adler32"

    "fmt"
    
)

func Abs(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(settings.Folder, path)
	}
	return path
}

func Adler32(s string)string{
    return fmt.Sprintf("%d",adler32.Checksum([]byte(s)))
    // hash.
}

