package settings

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
	"github.com/go-ini/ini"
	// "github.com/gosexy/yaml"
	ymlread "gopkg.in/yaml.v2"

	"github.com/kardianos/osext"
)

type staticCfg struct {
	VirtualRoot string `ini:"vstatic"`
	LocalRoot   string `ini:"lstatic"`
	CompressDef string `ini:"compress"`
}

type serverCfg struct {
	Port string `ini:"port"`
	Host string `ini:"host"`
}
type fileSyncCfg struct {
	CfgFile string `ini:"config"`
}

type templateCfg struct {
	Home         string `ini:"home"`
	DelimesLeft  string `ini:"ldelime"`
	DelimesRight string `ini:"rdelime"`
	Charset      string `ini:"charset"`
	Reload       bool   `ini:"reload"`
}
type defaultVar struct {
	AppName string `ini:"appname"`
}

type adminCfg struct {
	Passwd string `ini:"passwd"`
}

type logCfg struct {
	Path   string `ini:"path"`
	Format string `ini:"format"`
	File   string `ini:"-"`
}

type setting struct {
	Static      staticCfg   `ini:"static"`
	Server      serverCfg   `ini:"server"`
	Filesync    fileSyncCfg `ini:"filesync"`
	Template    templateCfg `ini:"template"`
	DefaultVars defaultVar  `ini:"defaultvars"`
	Admin       adminCfg    `ini:"admin"`
	Log         logCfg      `ini:"log"`
}

var (
	FsCfgMgr      *cfgMgr
	Folder        string
	settingStruct = new(setting)

	//GlobalSettings
	Static      staticCfg
	Server      serverCfg
	Filesync    fileSyncCfg
	Template    templateCfg
	DefaultVars defaultVar
	Admin       adminCfg
	Log         logCfg
)

var lock = new(sync.Mutex)

func init() {
	var err error
	Folder, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
}

func Init() {
	cfgFile := getAbs("./settings/settings.ini")

	cfg := new(ini.File)
	cfg.BlockMode = false
	cfg, err := ini.Load(cfgFile)
	if err != nil {
		panic(err)
	}

	cfg.MapTo(&settingStruct)
	settingStruct.Log.File = filepath.Join(
		getAbs(settingStruct.Log.Path),
		time.Now().Format(settingStruct.Log.Format),
	)

	//map to global
	{
		Static = settingStruct.Static
		Server = settingStruct.Server
		Filesync = settingStruct.Filesync
		Template = settingStruct.Template
		DefaultVars = settingStruct.DefaultVars
		Admin = settingStruct.Admin
		Log = settingStruct.Log
	}

	// FsCfgMgr = new(cfgMgr)
	// FsCfgMgr.Init()

	go watch()
}

func reload() {
	lock.Lock()
	defer lock.Unlock()
	Init()
}

func watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = w.Add(getAbs(settingStruct.Filesync.CfgFile))
	if err != nil {
		panic(err)
	}
	for {
		select {
		case e := <-w.Events:
			if e.Op == fsnotify.Remove {
				break
			}
			reload()
		}
	}
}

func getAbs(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(Folder, path)
	}
	return path
}

func readConfig(ConfigFile string) *SavedConfig {
	fmt.Println(ConfigFile)
	config := new(SavedConfig)
	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic(err)
	}

	err = ymlread.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	return config

}
