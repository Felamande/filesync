package settings

import (
	"io/ioutil"
	"path/filepath"
	"sync"
    "time"

	"github.com/go-fsnotify/fsnotify"
	"github.com/go-ini/ini"
	"github.com/gosexy/yaml"
	ymlread "gopkg.in/yaml.v2"

	"github.com/kardianos/osext"
)

type cfgMgr struct {
	cfg    SavedConfig
	writer *yaml.Yaml
}

func(m *cfgMgr) Save() error{
    return m.writer.Save()
}

func(m *cfgMgr) Add(c SyncPairConfig)error{
    m.cfg.Pairs = append(m.cfg.Pairs,c)
    return m.writer.Set(m.cfg)
}

func(m *cfgMgr) Cfg() SavedConfig{
    return m.cfg
}

type SyncConfig struct {
	CoverSameName bool `json:"cover_same_name"`
	SyncDelete    bool `json:"sync_delete"`
	SyncRename    bool `json:"sync_rename"`
}

type GlobalConfig struct {
}

type SavedConfig struct {
	Pairs     []SyncPairConfig `json:"pairs"`
	LogPath   string           `json:"log_path"`
	IgnoreExt []string         `yaml:"ignore_ext"`
}

type SyncPairConfig struct {
	Left   string     `json:"left"`
	Right  string     `json:"right`
	Config SyncConfig `json:"config"`
}

const (
	staticSec   = "static"
	serverSec   = "server"
	fileSyncSec = "filesync"
	tplSec      = "templates"
    AdminSec    ="admin"
    LogSec      = "log"
)

var (
	Host            string
	Static          string
	LocalStatic     string
	CompressSetting string
	fileSyncCfgFile string
	FsCfgMgr        *cfgMgr
	Port            string
	Folder          string
	TplHome         string
	DelimesLeft     string
	DelimesRight    string
	TplCharset      string
	TplReload       bool
	CfgFile         string
	AdminPwd        string
	LogFile         string
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

	fileSyncCfgFile = getAbs(get(fileSyncSec, "config"))
	fileSyncCfg, err := readConfig(fileSyncCfgFile)

	fsCfgWriter,err := yaml.Open(fileSyncCfgFile)
    if err!=nil{
        panic(err)
    }

	FsCfgMgr = &cfgMgr{fileSyncCfg, fsCfgWriter}
    
	TplHome = getAbs(get(tplSec, "home"))

	DelimesLeft = get(tplSec, "ldelime")
	DelimesRight = get(tplSec, "rdelime")

	TplCharset = get(tplSec, "charset")

	TplReload = cfg.Section(tplSec).Key("reload").MustBool(false)
	AdminPwd = get(AdminSec, "passwd")

    LogFile = filepath.Join(
        getAbs(get(LogSec,"path")),
        time.Now().Format(get(LogSec,"format")),
        )
	go Watch()
}

func Reload() {
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
	for {
		select {
		case e := <-w.Events:
			if e.Op == fsnotify.Remove {
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

func readConfig(ConfigFile string) (SavedConfig, error) {

	config := SavedConfig{
		Pairs: []SyncPairConfig{},
	}
	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return SavedConfig{}, err
	}

	err = ymlread.Unmarshal(data, &config)
	if err != nil {
		return SavedConfig{}, err
	}
    
	return config, nil

}
