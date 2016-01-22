package syncer

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"
    "github.com/Felamande/filesync/settings"
	"github.com/Felamande/filesync/log"
	"github.com/Felamande/filesync/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

var defaultSyncer *Syncer

func Default()*Syncer{
    return defaultSyncer
}

func init(){
    defaultSyncer = New()
}

var logger *log.Logger

type DirHandler func(string)
type FileHandler func(string)

type SyncPair struct {
	Left      uri.Uri
	Right     uri.Uri
	Config    settings.SyncConfig
	watcher   *fsnotify.Watcher
	tokens    chan bool
	IgnoreMap map[string]bool
}

type SyncMsg struct {
	Op   fsnotify.Op
	Left uri.Uri
}

type Syncer struct {
	SyncPairs []*SyncPair
	PairMap   map[string]string
}

func New() *Syncer {
	return &Syncer{
		SyncPairs: []*SyncPair{},
		PairMap:   make(map[string]string,32),
	}
}

func SetLogger(l *log.Logger) {
	logger = l
}

func (s *Syncer) Run() {
    config := settings.FsCfgMgr.Cfg()
	if logger == nil {
		panic("logger is nil")
	}

	if len(config.Pairs) == 0 {
		logger.Info("no pairs in config.json.")
		return
	}

	for _, pair := range config.Pairs {
		err := s.newPair(pair.Config, pair.Left, pair.Right, config.IgnoreExt)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

func (s *Syncer) NewPair(config settings.SyncConfig, source, target string, IgnoreRules []string) error {
    err :=s.newPair(config,source,target,IgnoreRules)
    if err!=nil{
        return err
    }
    err = settings.FsCfgMgr.Add(settings.SyncPairConfig{
        Left:source,
        Right:target,
        Config:config,
    })
    if err!=nil{
        return err
    }
    // return settings.FsCfgMgr.Save()
    return nil
}

func (s *Syncer) newPair(config settings.SyncConfig, source, target string, IgnoreRules []string) error {
	lURI, err := uri.Parse(source)
	if err != nil {
		return err
	}
	rURI, err := uri.Parse(target)
	if err != nil {
		return err
	}
	var m = make(map[string]bool, 6)
	for _, ignore := range IgnoreRules {
		m[ignore] = true
	}

	pair := &SyncPair{
		Left:      lURI,
		Right:     rURI,
		Config:    config,
		IgnoreMap: m,
	}
	fmt.Println("s.ignore", IgnoreRules)
	pair.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if !pair.Left.IsAbs() || !pair.Right.IsAbs() {
		err = PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Pair Uris not absolute"}
		return err
	}
	if !pair.Left.Exist() {
		err = PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Res of left Uri not exist"}
		return err
	}

	s.SyncPairs = append(s.SyncPairs, pair)
	s.PairMap[pair.Left.Uri()] = pair.Right.Uri()
	go func(p *SyncPair) {
		p.BeginWatch()
	}(pair)

	return nil
}

func (p *SyncPair) BeginWatch() {
	var err error

	p.watcher, err = fsnotify.NewWatcher()
	p.tokens = make(chan bool, runtime.NumCPU())
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = p.Left.Walk(
		func(root, lDir uri.Uri) error {
			//fmt.Println("now: " + lDir.Abs())
			err = p.WatchLeft(lDir)
			if err != nil {
				logger.Error(err.Error())
				return nil
			}
			logger.Info("Add to watcher: " + lDir.Abs())
			p.handleCreate(lDir)
			return nil

		},
		func(root, lFile uri.Uri) error {
			//fmt.Println("now: " + lFile.Abs())
			p.handleCreate(lFile)
			return nil
		},
	)

	if err != nil {
		logger.Info(err.Error())
		return
	}

	logger.Info("Add pair : " + p.Left.Uri() + " ==> " + p.Right.Uri())
	p.loopMsg()
}

func (p *SyncPair) loopMsg() {
	p.tokens = make(chan bool, runtime.NumCPU())

	for {
		select {
		case e := <-p.watcher.Events:
			URIString := p.formatLeftUriString(e.Name)
			u, err := uri.Parse(URIString)
			ext := filepath.Ext(u.Abs())
			if(p.IgnoreMap[ext]){
				logger.Info("ignored",u.Uri())
				
				continue
			}

			if err != nil {
				logger.Error(err.Error())
				continue
			}
			go p.handle(SyncMsg{e.Op, u})
		case e := <-p.watcher.Errors:
			logger.Info(e.Error())

		}

	}
}

func (p *SyncPair) handle(msg SyncMsg) {
	p.tokens <- true
	defer func() { <-p.tokens }()

	switch msg.Op {
	case fsnotify.Create:
		if msg.Left.IsDir() {
			err := p.WatchLeft(msg.Left)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			logger.Info("Add to Watcher: " + msg.Left.Abs())
		}

		p.handleCreate(msg.Left)

	case fsnotify.Write:
		if msg.Left.IsDir() {
			return
		}
		p.handleWrite(msg.Left)
	case fsnotify.Remove:
		if !p.Config.SyncDelete {
			return
		}
		p.handleRemove(msg.Left)
	case fsnotify.Rename:
		if !p.Config.SyncRename {
			return
		}
		p.handleRename(msg.Left)
	}
}

func (p *SyncPair) handleWrite(lFile uri.Uri) {
	var err error

	rFile, err := p.ToRight(lFile)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	lFd, rFd := copyFile(rFile, lFile)
	defer lFd.Close()
	defer rFd.Close()

	logger.Info("Sync file succesfully: " + lFile.Abs() + " ==> " + rFile.Abs())
}

func (p *SyncPair) handleCreate(lName uri.Uri) {
	rName, err := p.ToRight(lName)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if !lName.ModTime().After(rName.ModTime()) {
		return
	}

	for {
		err := rName.Create(lName.IsDir(), lName.Mode())
		if err == nil {
			break
		} else {
			time.Sleep(time.Second * 1)
		}
	}

	if rName.IsDir() {
		return
	}

	//first time copy file.
	lFd, rFd := copyFile(rName, lName)
	defer lFd.Close()
	defer rFd.Close()

	logger.Info("Sync  succesfully: " + lName.Abs() + " ==> " + rName.Abs())

}

func (p *SyncPair) formatLeftUriString(name string) string {
	name = strings.Replace(name, "\\", "/", -1)
	return p.Left.Scheme() + "://" + name
}

func (p *SyncPair) handleRemove(lName uri.Uri) {
	fmt.Println(lName.Abs(), "removed")
	rName, err := p.ToRight(lName)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if !rName.Exist() {
		return
	}

	for {
		err := rName.Remove()
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
	fmt.Println(rName.Abs(), "removed")

	err = p.watcher.Remove(lName.Abs())
	if err != nil {
		logger.Warn(err.Error())
	}

	logger.Info("Remove successfully: " + rName.Abs())
	return

}

func (p *SyncPair) handleRename(lName uri.Uri) {
	fmt.Println(lName.Abs(), "rename")
}

func (p *SyncPair) ToRight(u uri.Uri) (uri.Uri, error) {
	lTmp := p.Left.Uri()
	rTmp := p.Right.Uri()
	lTmplen := len(lTmp)
	rTmplen := len(rTmp)
	if lTmp[lTmplen-1] == '/' {
		lTmp = lTmp[0 : lTmplen-1]
	}
	if rTmp[rTmplen-1] == '/' {
		rTmp = rTmp[0 : rTmplen-1]
	}
	Uris := strings.Replace(u.Uri(), lTmp, rTmp, -1)
	return uri.Parse(Uris)

}

func (p *SyncPair) WatchLeft(left uri.Uri) error {
	for {
		err := p.watcher.Add(left.Abs())
		if err != nil {
			err = p.watcher.Add(left.Abs())
		} else {
			break
		}
	}
	return nil

}

type PairNotValidError struct {
	Left    string
	Right   string
	Message string
}

func (e PairNotValidError) Error() string {
	return e.Message + ": " + e.Left + " ==> " + e.Right
}

func copyFile(rFile, lFile uri.Uri) (io.ReadCloser, io.WriteCloser) {
	var err error

	lFd, err := lFile.OpenRead()
	for {
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
		lFd, err = lFile.OpenRead()

	}

	rFd, err := rFile.OpenWrite()
	for {
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {

			break
		}
		rFd, err = rFile.OpenWrite()

	}

	io.Copy(rFd, lFd)
	return lFd, rFd
}

// func (s *Syncer) WatchConfigChange(file string) error {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		return err
// 	}
// 	err = watcher.Add(file)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("start watching config file: " + file)
// 	for {
// 		select {
// 		case e := <-watcher.Events:
// 			fmt.Println(e.Op, e.Name)
// 			if e.Op == fsnotify.Remove || e.Op == fsnotify.Rename {
// 				continue
// 			}
// 			go s.updateConfig(file)
// 		case e := <-watcher.Errors:
// 			logger.Error(e.Error())
// 		}
// 	}

// }

// func (s *Syncer) updateConfig(file string) {
// 	data := make([]byte, 4096)
// 	var err error
// 	timeout := 0
// 	for {
// 		if timeout > 100 {
// 			return
// 		}
// 		data, err = ioutil.ReadFile(file)
// 		if err != nil {
// 			timeout++
// 		} else {
// 			break
// 		}

// 	}

// 	config := &SavedConfig{
// 		Pairs: []SyncPairConfig{},
// 	}

// 	err = yaml.Unmarshal(data, config)
// 	if err != nil {
// 		fmt.Println("unmarshal failed: ", err.Error())
// 		logger.Info(err.Error())
// 		return
// 	}

// 	for _, p := range config.Pairs {
// 		right, exist := s.PairMap[p.Left]
// 		if !exist {
// 			err = s.NewPair(p.Config, p.Left, p.Right)
// 			if err != nil {
// 				logger.Error(err.Error())
// 				continue
// 			}
// 			fmt.Println("update: ", p.Config, p.Left, p.Right)
// 			logger.Info("config updated.")
// 		} else {
// 			if right == p.Right {
// 				continue
// 			}

// 			err = s.NewPair(p.Config, p.Left, p.Right)
// 			if err != nil {
// 				logger.Error(err.Error())
// 			}
// 			fmt.Println("update: ", p.Config, p.Left, p.Right)
// 			logger.Info("config updated.")
// 		}
// 	}
// }
