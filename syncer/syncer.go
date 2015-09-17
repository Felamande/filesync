package syncer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Felamande/filesync/uri"

	"github.com/Felamande/filesync/log"
	fsnotify "gopkg.in/fsnotify.v1"
)

var logger log.Logger

func init() {
	fi, _ := os.Stat(".")
	os.Mkdir(".log", fi.Mode())
	logger = log.NewFileLogger("./.log/sync.log")
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
	SyncPairs []SyncPair
}

func New() *Syncer {
	return &Syncer{
		SyncPairs: []SyncPair{},
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

func (s *Syncer) readConfig() {

	configFile, err := os.Open("config.json")
	if err != nil {
		logger.Warn("*Syncer.readConfig", err.Error())
		return
	}

	JsonBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		logger.Warn("*Syncer.readConfig", err.Error())
		return
	}

	var config SavedConfig = SavedConfig{}

	err = json.Unmarshal(JsonBytes, &config.Pairs)

	if err != nil {
		logger.Warn("*Syncer.readConfig", err.Error())
	}

	if len(config.Pairs) == 0 {
		logger.Info("&SyncPair.readConfig", "no pairs in config.json.")
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
			SyncPair{
				Left:   LeftUri,
				Right:  RightUri,
				Config: pair.Config,
			},
		)
		fmt.Println(pair.Config)
	}

}

func (s *Syncer) Run() {

	s.readConfig()

	for _, pair := range s.SyncPairs {

		if !pair.Left.IsAbs() || !pair.Left.IsAbs() {
			logger.Warn("*Syncer.Run", PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Pair Uris not absolute"}.Error())
			continue
		}
		if !pair.Left.Exist() || !pair.Left.Exist() {
			logger.Warn("*Syncer.Run", PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Res of pair Uris not exist"}.Error())
			continue
		}

		func(p SyncPair) {
			go p.BeginWatch()
		}(pair)

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

	pair := SyncPair{
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

	go pair.BeginWatch()

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

			err = p.watcher.Add(lDir.Abs())
			if err != nil {
				logger.Error("*SyncPair.BeginWatch@walkDir", err.Error())
				return nil
			}
			logger.Info("*SyncPair.BeginWatch@walkDir", "Add to watcher: "+lDir.Abs())
			p.handleCreate(lDir)
			return nil

		},
		func(root, lFile uri.Uri) error {
			p.handleCreate(lFile)
			return nil
		},
	)

	if err != nil {
		logger.Info("*SyncPair.BeginWatch", err.Error())
		return
	}

	logger.Info("*SyncPair.BeginWatch", "Add pair : "+p.Left.Uri()+" ==> "+p.Right.Uri())
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
		logger.Error("*SyncPair.handleCreateRight", err.Error())
		return
	}

	if !lName.ModTime().After(rName.ModTime()) {
		fmt.Println(rName.Abs(), "The left is older or the right exist.")
		return
	} else {
		fmt.Println(rName.Abs(), "the left is newer.")
	}

	for {
		err := rName.Create(lName.IsDir(), lName.Mode())
		if err == nil {

			break
		} else {
			fmt.Println(err.Error())
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
	Uris := strings.Replace(u.Uri(), p.Left.Uri(), p.Right.Uri(), -1)
	return uri.Parse(Uris)

}

func (p *SyncPair) WatchLeft(left uri.Uri) error {
	for {
		err := p.watcher.Add(left.Abs())
		if err != nil {
			fmt.Println("file occupied: " + err.Error())
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
