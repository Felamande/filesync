package utils

import (
    "path/filepath"
    "github.com/Felamande/filesync/settings"
    "hash/adler32"
    "crypto/md5"
    "io"
    "encoding/hex"
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

func Md5(source ...interface{}) string {
	ctx := md5.New()
	for _, s := range source {
		switch ss := s.(type) {
		case io.Reader:
			io.Copy(ctx, ss)
		case string:
			ctx.Write([]byte(ss))
		case []byte:
			ctx.Write(ss)

		}
	}

	return hex.EncodeToString(ctx.Sum(nil))
}