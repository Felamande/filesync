package syncer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Felamande/filesync/log"
	"github.com/Felamande/filesync/uri"
	"github.com/kardianos/osext"
	fsnotify "gopkg.in/fsnotify.v1"
	yaml "gopkg.in/yaml.v2"
)

var logger log.Logger

func init() {
	f, _ := osext.ExecutableFolder()
	logger = log.NewFileLogger(filepath.Join(f, ".log/sync.log"))
}

type DirHandler func(string)
type FileHandler func(string)

type SyncConfig struct {
	CoverSameName bool `json:"cover_same_name"`
	SyncDelete    bool `json:"sync_delete"`
	SyncRename    bool `json:"sync_rename"`
}

type SyncPair struct {
	Left    uri.Uri
	Right   uri.Uri
	Config  SyncConfig
	watcher *fsnotify.Watcher
	tokens  chan bool
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
		PairMap:   make(map[string]string, 4),
	}
}

// @param1 interface log.Logger
// package github.com/Felamande/filesync/log
// type Logger interface {
//	Info(source, messsage string)
//	Debug(source, messsage string)
//	Warn(source, messsage string)
//	Error(source, messsage string)
//	Critical(source, messsage string)
//	Panic(source, messsage string)
//	Close() error
//}
func SetLogger(logger log.Logger) {
	logger = logger
}

func (s *Syncer) Run(config SavedConfig) {

	if len(config.Pairs) == 0 {
		logger.Info("*SyncPair.Run", "no pairs in config.json.")
		return
	}

	for _, pair := range config.Pairs {
		LeftUri, err := uri.Parse(pair.Left)
		if err != nil {
			logger.Error("*SyncPair.readConfig", err.Error())
			return
		}
		RightUri, err := uri.Parse(pair.Right)
		if err != nil {
			logger.Error("*SyncPair.readConfig", err.Error())
			return
		}
		s.SyncPairs = append(s.SyncPairs,
			&SyncPair{
				Left:   LeftUri,
				Right:  RightUri,
				Config: pair.Config,
			},
		)
	}

	for _, pair := range s.SyncPairs {

		if !pair.Left.IsAbs() || !pair.Left.IsAbs() {
			logger.Warn("*Syncer.Run", PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Pair Uris not absolute"}.Error())
			continue
		}
		if !pair.Left.Exist() || !pair.Left.Exist() {
			logger.Warn("*Syncer.Run", PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Res of pair Uris not exist"}.Error())
			continue
		}

		s.PairMap[pair.Left.Uri()] = pair.Right.Uri()

		go func(p *SyncPair) {
			p.BeginWatch()
		}(pair)

	}

}

func (s *Syncer) WatchConfigChange(file string) error {
	fmt.Println("Im here.")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = watcher.Add(file)
	if err != nil {
		return err
	}
	fmt.Println("start watching config file: " + file)
	for {
		select {
		case e := <-watcher.Events:
			fmt.Println(e.Op, e.Name)
			if e.Op == fsnotify.Remove || e.Op == fsnotify.Rename {
				continue
			}
			go s.updateConfig(file)
		case e := <-watcher.Errors:
			fmt.Println(e.Error())
		}
	}

}

func (s *Syncer) updateConfig(file string) {
	var data []byte = make([]byte, 4096)
	var err error
	var timeout int = 0
	for {
		if timeout > 100 {
			return
		}
		data, err = ioutil.ReadFile(file)
		if err != nil {
			timeout++
		} else {
			break
		}

	}

	var config *SavedConfig = &SavedConfig{
		Pairs: []SyncPairConfig{},
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Println("unmarshal failed: ", err.Error())
		logger.Info("*Syncer.updateConfig@Unmarshall", err.Error())
		return
	}

	for _, p := range config.Pairs {
		right, exist := s.PairMap[p.Left]
		if !exist {
			err = s.NewPair(p.Config, p.Left, p.Right)
			if err != nil {
				logger.Error("*Syncer.updateConfig", err.Error())
				continue
			}
			fmt.Println("update: ", p.Config, p.Left, p.Right)
			logger.Info("*Syncer.updateConfig", "config updated.")
		} else {
			if right == p.Right {
				continue
			}

			err = s.NewPair(p.Config, p.Left, p.Right)
			if err != nil {
				logger.Error("*Syncer.updateConfig", err.Error())
			}
			fmt.Println("update: ", p.Config, p.Left, p.Right)
			logger.Info("*Syncer.updateConfig", "config updated.")
		}
	}
}

func (s *Syncer) NewPair(config SyncConfig, source, target string) error {

	lUri, err := uri.Parse(source)
	if err != nil {
		logger.Error("*Syncer.NewPair", "Parse source: "+err.Error())
		return err
	}
	rUri, err := uri.Parse(target)
	if err != nil {
		logger.Error("*Syncer.NewPair", "Parse target: "+err.Error())
		return err
	}

	pair := &SyncPair{
		Left:  lUri,
		Right: rUri,
	}
	pair.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Error("*SyncPair.NewPair", err.Error())
		return err
	}

	if !pair.Left.IsAbs() || !pair.Left.IsAbs() {
		err = PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Pair Uris not absolute"}
		logger.Warn("*Syncer.Run", err.Error())
		return err
	}
	if !pair.Left.Exist() || !pair.Left.Exist() {
		err = PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Res of pair Uris not exist"}
		logger.Warn("*Syncer.Run", err.Error())
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
		logger.Error("*SyncPair.BeginWatch", err.Error())
		return
	}

	fmt.Println("Start Walking: " + p.Left.Abs())

	err = p.Left.Walk(
		func(root, lDir uri.Uri) error {
			//fmt.Println("now: " + lDir.Abs())
			err = p.WatchLeft(lDir)
			if err != nil {
				logger.Error("*SyncPair.BeginWatch@walkDir", err.Error())
				return nil
			}
			logger.Info("*SyncPair.BeginWatch@walkDir", "Add to watcher: "+lDir.Abs())
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
		fmt.Println("walk: ", err)
		logger.Info("*SyncPair.BeginWatch", err.Error())
		return
	}

	logger.Info("*SyncPair.BeginWatch", "Add pair : "+p.Left.Uri()+" ==> "+p.Right.Uri())
	fmt.Println("end Walking: " + p.Left.Uri())
	p.loopMsg()
}

func (p *SyncPair) loopMsg() {
	p.tokens = make(chan bool, runtime.NumCPU())

	for {
		select {
		case e := <-p.watcher.Events:
			UriString := p.formatLeftUriString(e.Name)
			u, err := uri.Parse(UriString)
			if err != nil {
				logger.Error("*SyncPair.loopMsg@events", err.Error())
				continue
			}
			go p.handle(SyncMsg{e.Op, u})
		case e := <-p.watcher.Errors:
			logger.Info("*SyncPair.looMsg@errors", e.Error())

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
				logger.Error("*SyncPair.handle@watcher.Add", err.Error())
				return
			}
			logger.Info("*SyncPair.handle@switch", "Add to Watcher: "+msg.Left.Abs())
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
		logger.Error("*SyncPair.handleWrite@ToRight", err.Error())
		return
	}

	fmt.Println("start copying:", lFile.Abs(), "==>", rFile.Abs())

	lFd, rFd := copyFile(rFile, lFile)
	defer lFd.Close()
	defer rFd.Close()

	fmt.Println("finish copying:", lFile.Abs(), "==>", rFile.Abs())
	logger.Info("*SyncPair.handleFile@io.Copy", "Sync file succesfully: "+lFile.Abs()+" ==> "+rFile.Abs())
}

func (p *SyncPair) handleCreate(lName uri.Uri) {
	rName, err := p.ToRight(lName)
	if err != nil {
		logger.Error("*SyncPair.handleCreate@ToRight", err.Error())
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

	logger.Info("*SyncPair.handleCreate@io.Copy", "Sync  succesfully: "+lName.Abs()+" ==> "+rName.Abs())

}

func (p *SyncPair) formatLeftUriString(name string) string {
	name = strings.Replace(name, "\\", "/", -1)
	return p.Left.Scheme() + "://" + name
}

func (p *SyncPair) handleRemove(lName uri.Uri) {
	fmt.Println(lName.Abs(), "removed")
	rName, err := p.ToRight(lName)
	if err != nil {
		logger.Error("*SyncPair.handleRemove", err.Error())
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
		logger.Warn("*syncPair.handleRemove", err.Error())
	}

	logger.Info("*SyncPair.handleRemove", "Remove successfully: "+rName.Abs())
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

func exist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
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
