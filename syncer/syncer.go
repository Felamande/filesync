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

var Log log.Logger

func init() {
	Log = log.NewFileLogger("sync.log")
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
	Log = logger
}

func (s *Syncer) readConfig() {

	configFile, err := os.Open("config.json")
	if err != nil {
		Log.Warn("*Syncer.readConfig", err.Error())
		return
	}

	JsonBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		Log.Warn("*Syncer.readConfig", err.Error())
		return
	}

	var config SavedConfig = SavedConfig{}

	err = json.Unmarshal(JsonBytes, &config.Pairs)

	if err != nil {
		Log.Warn("*Syncer.readConfig", err.Error())
	}

	if len(config.Pairs) == 0 {
		Log.Info("&SyncPair.readConfig", "no pairs in config.json.")
		return
	}

	for _, pair := range config.Pairs {
		LeftUri, err := uri.Parse(pair.Left)
		if err != nil {
			Log.Error("*SyncPair.readConfig", err.Error())
			return
		}
		RightUri, err := uri.Parse(pair.Right)
		if err != nil {
			Log.Error("*SyncPair.readConfig", err.Error())
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
			Log.Warn("*Syncer.Run", PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Pair Uris not absolute"}.Error())
			continue
		}
		if !pair.Left.Exist() || !pair.Left.Exist() {
			Log.Warn("*Syncer.Run", PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Res of pair Uris not exist"}.Error())
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
		Log.Error("*Syncer.NewPair", "Parse source: "+err.Error())
		return err
	}
	rUri, err := uri.Parse(target)
	if err != nil {
		Log.Error("*Syncer.NewPair", "Parse target: "+err.Error())
		return err
	}

	pair := SyncPair{
		Left:  lUri,
		Right: rUri,
	}
	pair.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		Log.Error("*SyncPair.NewPair", err.Error())
		return err
	}

	if !pair.Left.IsAbs() || !pair.Left.IsAbs() {
		err = PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Pair Uris not absolute"}
		Log.Warn("*Syncer.Run", err.Error())
		return err
	}
	if !pair.Left.Exist() || !pair.Left.Exist() {
		err = PairNotValidError{pair.Left.Abs(), pair.Right.Abs(), "Res of pair Uris not exist"}
		Log.Warn("*Syncer.Run", err.Error())
		return err
	}

	s.SyncPairs = append(s.SyncPairs, pair)

	go pair.BeginWatch()

	return nil
}

func (p *SyncPair) BeginWatch() {
	var err error

	p.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		Log.Error("*SyncPair.BeginWatch", err.Error())
		return
	}

	fmt.Println("Start Walking: " + p.Left.Abs())

	err = p.Left.Walk(
		func(root, lDir uri.Uri) error {
			rDir, err := p.ToRight(lDir)
			if err != nil {
				Log.Error("*SyncPair.BeginWatch@walkDir", err.Error())
				return nil
			}
			err = p.watcher.Add(lDir.Abs())
			if err != nil {
				Log.Error("*SyncPair.BeginWatch@walkDir", err.Error())
				return nil
			}
			Log.Info("*SyncPair.BeginWatch@walkDir", "Add to watcher: "+lDir.Abs())
			rDir.Create(true, lDir.Mode())
			return nil

		},
		func(root, lFile uri.Uri) error {
			rFile, err := p.ToRight(lFile)
			if err != nil {
				Log.Error("*SyncPair.BeginWatch@walkDir", err.Error())
				return nil
			}
			if !rFile.Exist() {
				rFile.Create(false, lFile.Mode())
				p.handleWrite(rFile)
			}
			if lFile.ModTime().After(rFile.ModTime()) {
				p.handleWrite(rFile)
			}
			return nil
		},
	)

	if err != nil {
		Log.Info("*SyncPair.BeginWatch", err.Error())
		return
	}

	Log.Info("*SyncPair.BeginWatch", "Add pair : "+p.Left.Uri()+" ==> "+p.Right.Uri())
	p.loopMsg()
}

func (p *SyncPair) loopMsg() {
	tokens := make(chan bool, runtime.NumCPU())

	for {
		select {
		case e := <-p.watcher.Events:
			UriString := p.formatLeftUriString(e.Name)
			u, err := uri.Parse(UriString)
			if err != nil {
				Log.Error("*SyncPair.loopMsg@events", err.Error())
				continue
			}
			go p.handle(SyncMsg{e.Op, u}, tokens)
		case e := <-p.watcher.Errors:
			Log.Info("*SyncPair.looMsg@errors", e.Error())

		}

	}
}

func (p *SyncPair) handle(msg SyncMsg, tokens chan bool) {
	tokens <- true
	defer func() { <-tokens }()

	switch msg.Op {
	case fsnotify.Create:
		if msg.Left.IsDir() {
			err := p.WatchLeft(msg.Left)
			if err != nil {
				Log.Error("*SyncPair.handle@watcher.Add", err.Error())
				return
			}
			Log.Info("*SyncPair.handle@switch", "Add to Watcher: "+msg.Left.Abs())
		}

		p.handleCreateRight(msg.Left)

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
	}
}

func (p *SyncPair) handleWrite(lFile uri.Uri) {
	var err error

	rFile, err := p.ToRight(lFile)
	if err != nil {
		Log.Error("*SyncPair.handleWrite@ToRight", err.Error())
		return
	}

	//BUG unknown bug here causing panic.
	lFd, err := lFile.OpenRead()
	for {
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
		lFd, err = lFile.OpenRead()

	}
	defer lFd.Close()

	rFd, err := rFile.OpenWrite()
	for {
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
		rFd, err = rFile.OpenWrite()

	}
	defer rFd.Close()
	fmt.Println("start copying:", lFile.Abs(), "==>", rFile.Abs())
	io.Copy(rFd, lFd)
	fmt.Println("finish copying:", lFile.Abs(), "==>", rFile.Abs())
	Log.Info("*SyncPair.handleFile@io.Copy", "Sync file succesfully: "+lFile.Abs()+" ==> "+rFile.Abs())
}

func (p *SyncPair) handleCreateRight(lName uri.Uri) {
	rName, err := p.ToRight(lName)
	if err != nil {
		Log.Error("*SyncPair.handleCreateRight", err.Error())
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

	Log.Info("*SyncPair.handleCreateFile@io.Copy", "Sync  succesfully: "+lName.Abs()+" ==> "+rName.Abs())

}

func (p *SyncPair) formatLeftUriString(name string) string {
	name = strings.Replace(name, "\\", "/", -1)
	return p.Left.Scheme() + "://" + name
}

func (p *SyncPair) handleRemove(lName uri.Uri) {
	fmt.Println(lName.Abs(), "removed")
	rName, err := p.ToRight(lName)
	if err != nil {
		Log.Error("*SyncPair.handleRemove", err.Error())
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

	Log.Info("*SyncPair.handleRemove", "Remove successfully: "+rName.Abs())
	return

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
