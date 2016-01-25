package utils

import (
	"html/template"
	"path"
    "fmt"

	"github.com/Felamande/filesync/settings"
    
	"github.com/beego/compress"
)

func AssetJS(src string) template.HTML {
    return template.HTML(fmt.Sprintf(`<script src="%s"></script>`,path.Join(settings.Static.VirtualRoot, "js", src)))
	// return path.Join(settings.Static.VirtualRoot, "js", src)
}

func AssetCss(src string)template.HTML{
    return template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s" type="text/css" />`,path.Join(settings.Static.VirtualRoot, "css", src)))
}

func DefaultFuncs() template.FuncMap {
	s, err := compress.LoadJsonConf(Abs(settings.Static.CompressDef), true, settings.Server.Host)
	if err != nil {
		panic(err)

	}
    return template.FuncMap{
		"AssetJs": AssetJS,
        "AssetCss":AssetCss,
        "CompressCss":s.Css.CompressCss,
        "CompressJs":s.Js.CompressJs,
	}
}
