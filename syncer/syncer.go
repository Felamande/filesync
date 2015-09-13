package syncer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
	Left    string     `json:"left"`
	Right   string     `json:"right"`
	Config  SyncConfig `json:"config"`
	watcher *fsnotify.Watcher
}

type SyncMsg struct {
	Op       fsnotify.Op
	FullName string
	IsDir    bool
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

	config, err := os.Open("config.json")
	if err != nil {
		Log.Warn("*Syncer.readConfig", err.Error())
		return
	}

	JsonBytes, err := ioutil.ReadAll(config)
	if err != nil {
		Log.Warn("*Syncer.readConfig", err.Error())
		return
	}

	err = json.Unmarshal(JsonBytes, &s.SyncPairs)
	if err != nil {
		Log.Warn("*Syncer.readConfig", err.Error())
	}

	fmt.Println(len(s.SyncPairs))

}

func (s *Syncer) Run() {

	s.readConfig()

	if len(s.SyncPairs) == 0 {
		fmt.Println(len(s.SyncPairs))
		Log.Info("*Syncer.Run", "No sync pair in config file.")
		return
	}

	for _, pair := range s.SyncPairs {

		if !pair.IsValid() {
			Log.Warn("*Syncer.Run", PairNotValidError{pair.Left, pair.Right}.Error())
			continue
		}
		fmt.Println(pair)

		err := pair.ExistOrCreate()
		if err != nil {
			Log.Error("*Syncer.Run@*SyncPair.ExistOrCreate", err.Error())
			continue
		}

		func(p SyncPair) {
			go p.BeginWatch()
		}(pair)

	}

}

func (s *Syncer) NewPair(config SyncConfig, source, target string) (err error) {

	p := SyncPair{
		Left:   source,
		Right:  target,
		Config: config,
	}
	if !p.IsValid() {
		return PairNotValidError{source, target}
	}

	p.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		Log.Error("*Syncer.NewPair", err.Error())
		return nil
	}

	err = p.ExistOrCreate()
	if err != nil {
		Log.Error("*Syncer.NewPair@*SyncPair.ExistOrCreate", err.Error())
		return
	}

	s.SyncPairs = append(s.SyncPairs, p)

	go p.BeginWatch()

	return nil
}

func (p *SyncPair) IsValid() bool {
	return filepath.IsAbs(p.Left) && filepath.IsAbs(p.Right)
}

func (p *SyncPair) ExistOrCreate() error {
	if !exist(p.Left) {
		err := os.MkdirAll(p.Left, 0777)
		return err
	}
	if !exist(p.Right) {
		err := os.MkdirAll(p.Right, 0777)
		return err
	}
	return nil
}

func (p *SyncPair) BeginWatch() {
	var err error

	p.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		Log.Error("*SyncPair.BeginWatch", err.Error())
		return
	}

	fmt.Println("Start Walking: " + p.Left)

	walkDir(p.Left,
		func(lDir string) {
			err := p.watcher.Add(lDir)
			if err != nil {
				Log.Error("*SyncPair.BeginWatch@walkDir", err.Error()+" "+lDir)
				return
			}
			Log.Info("*SyncPair.BeginWatch@walkDir", "Add to watcher: "+lDir)

			p.handleCreateDir(lDir)

		},
		func(lFile string) {
			rFile := strings.Replace(lFile, p.Left, p.Right, -1)
			lFi, _ := os.Stat(lFile)
			rFi, err := os.Stat(rFile)
			if os.IsNotExist(err) {
				p.handleFile(lFile, true)
				return
			}

			if lFi.ModTime().After(rFi.ModTime()) {
				p.handleFile(lFile, true)
				return
			}

		},
	)

	Log.Info("*SyncPair.BeginWatch", "Add pair : "+p.Left+" ==> "+p.Right)
	p.loopMsg()

}

func (p *SyncPair) loopMsg() {
	tokens := make(chan bool, runtime.NumCPU())

	for {
		select {
		case e := <-p.watcher.Events:
			fi, err := os.Stat(e.Name)
			if os.IsNotExist(err) {
				continue
			}
			go p.handle(SyncMsg{e.Op, e.Name, fi.IsDir()}, tokens)
		case e := <-p.watcher.Errors:
			Log.Info("*SyncPair.BeginWatch", e.Error())

		}

	}
}

func (p *SyncPair) handle(msg SyncMsg, tokens chan bool) {
	tokens <- true
	defer func() { <-tokens }()

	switch msg.Op {
	case fsnotify.Create:
		if msg.IsDir {
			p.watcher.Add(msg.FullName)
			Log.Info("*SyncPair.handle@switch", "Create dir and add to Watcher: "+msg.FullName)
			p.handleCreateDir(msg.FullName)

		} else {
			Log.Info("*SyncPair.handle@switch", "Create file: "+msg.FullName)
			p.handleFile(msg.FullName, false)
		}
	case fsnotify.Write:
		if msg.IsDir {
			return
		}
		p.handleFile(msg.FullName, true)
	case fsnotify.Remove:
		if !p.Config.SyncDelete {
			return
		}
		p.handleRemove(msg.FullName, msg.IsDir)
	}
}

func (p *SyncPair) handleFile(lFile string, ModWrite bool) {
	var err error
	fi, err := os.Stat(lFile)
	if err != nil {
		Log.Warn("*SyncPair.handleCreateFile@os.Stat", err.Error()+" "+lFile)
	}
	rFile := strings.Replace(lFile, p.Left, p.Right, -1)

	if rFile == lFile {
		Log.Error("*SyncPair.handleCreateFile@strings.Replace", lFile+" is not in directory "+p.Left)
		return
	}

	if !ModWrite {
		if exist(rFile) {
			Log.Info("*SyncPair.handleCreateFile@exist", "file "+rFile+" alread exists")
			return
		}
	}
	//BUG unknown bug here causing panic.
	lFd, err := os.OpenFile(lFile, os.O_RDWR|os.O_CREATE, fi.Mode())
	for {
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
		lFd, err = os.OpenFile(lFile, os.O_RDWR|os.O_CREATE, fi.Mode())

	}
	defer lFd.Close()

	rFd, err := os.OpenFile(rFile, os.O_RDWR|os.O_CREATE, fi.Mode())
	for {
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
		rFd, err = os.OpenFile(rFile, os.O_RDWR|os.O_CREATE, fi.Mode())

	}
	defer rFd.Close()

	io.Copy(rFd, lFd)
	Log.Info("*SyncPair.handleFile@io.Copy", "Sync file succesfully: "+lFile+" ==> "+rFile)
}

func (p *SyncPair) handleCreateDir(lDir string) {
	fi, err := os.Stat(lDir)
	if err != nil {
		Log.Warn("*SyncPair.handleCreateFile@os.Stat", err.Error()+" "+lDir)
	}

	rDir := strings.Replace(lDir, p.Left, p.Right, -1)
	if exist(rDir) {
		return
	}

	for {
		err := os.MkdirAll(rDir, fi.Mode())
		if err == nil {
			break
		} else {
			time.Sleep(time.Second * 1)
		}
	}

	Log.Info("*SyncPair.handleCreateFile@io.Copy", "Sync dir succesfully: "+lDir+" ==> "+rDir)

}

func (p *SyncPair) handleRemove(lName string, IsDir bool) {
	rName := strings.Replace(lName, p.Left, p.Right, -1)
	_, err := os.Stat(rName)

	if os.IsNotExist(err) {
		return
	}

	for {
		err := os.Remove(rName)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}

	Log.Info("*SyncPair.handleRemove", "Remove successfully: "+rName)
	return

}

type PairNotValidError struct {
	Left  string
	Right string
}

func (e PairNotValidError) Error() string {
	return "Pair not valid :" + e.Left + " ==> " + e.Right
}

func walkDirTranverse(Name string, dh DirHandler, fh FileHandler) {

	fis, _ := ioutil.ReadDir(Name)

	if len(fis) == 0 {
		fmt.Println("no dir in: " + Name)
		return
	}
	for _, fi := range fis {
		if !fi.IsDir() {
			fh(Name + "\\" + fi.Name())
			continue
		}

		FullDirName := Name + "\\" + fi.Name()
		dh(FullDirName)
		walkDirTranverse(FullDirName, dh, fh)
	}

}

func walkDir(Name string, dh DirHandler, fh FileHandler) {
	walkDirTranverse(Name, dh, fh)
	dh(Name)
}

func exist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
