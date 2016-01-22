package utils

import (
	"html/template"
	"path"

	"github.com/Felamande/filesync/settings"
	"github.com/beego/compress"
)

func AssetJS(src string) string {
	return path.Join(settings.Static, "js", src)
}

func DefaultFuncs() template.FuncMap {
	_, err := compress.LoadJsonConf(settings.CompressSetting, true, settings.Host)
	if err != nil {
		panic(err)

	}
    return template.FuncMap{
		"AssetJs": AssetJS,
	}
}
