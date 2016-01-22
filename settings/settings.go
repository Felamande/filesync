package settings

import (
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"
	"github.com/go-ini/ini"
	"github.com/kardianos/osext"
    "sync"
)

const (
	staticSec   = "static"
	serverSec   = "server"
	fileSyncSec = "filesync"
	tplSec      = "templates"
)

var (
	Host            string
	Static          string
	LocalStatic     string
	CompressSetting string
	FileSyncConfig  string
	Port            string
	Folder          string
	TplHome         string
	DelimesLeft     string
	DelimesRight    string
	TplCharset      string
	TplReload       bool
	CfgFile         string
)

var cfg *ini.File
var lock = new(sync.Mutex)
func init() {
	var err error
	Folder, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
}

func Init() {
	var err error

	CfgFile = getAbs("./settings/settings.ini")

	cfg, err = ini.Load(CfgFile)
	if err != nil {
		panic(err)
	}

	Port = get(serverSec, "port")
	Host = get(serverSec, "host")

	Static = get(staticSec, "vstatic")

	LocalStatic = getAbs(get(staticSec, "lstatic"))
	CompressSetting = getAbs(get(staticSec, "compress"))

	FileSyncConfig = getAbs(get(fileSyncSec, "config"))

	TplHome = getAbs(get(tplSec, "home"))

	DelimesLeft = get(tplSec, "ldelime")
	DelimesRight = get(tplSec, "rdelime")

	TplCharset = get(tplSec, "charset")

	TplReload = cfg.Section(tplSec).Key("reload").MustBool(false)
    
    go Watch()
}

func Reload(){
    lock.Lock()
    defer lock.Unlock()
    Init()
}

func Watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
    
    err = w.Add(CfgFile)
    if err != nil {
		panic(err)
	}
    for{
        select{
            case e:=<-w.Events:
               if e.Op == fsnotify.Remove{
                   break
               }
               Reload()
        }
    }
}

func getAbs(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(Folder, path)
	}
	return path
}

func get(sec, key string) string {
	return cfg.Section(sec).Key(key).Value()
}
