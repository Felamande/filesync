package syncer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	fsnotify "gopkg.in/fsnotify.v1"
)

type Test struct {
	Pc   uintptr
	File string
	Line int
	Ok   bool
}

func TestSync(t *testing.T) {

	//	s := New()
	//	s.NewPair(SyncConfig{true, true, false}, "D:\\pictures", "E:\\pictures")
	//	s.NewPair(SyncConfig{true, true, false}, "D:\\video", "E:\\pictures")
	//	b, _ := json.Marshal(s.SyncPairs)
	//	ioutil.WriteFile("config.json", b, 777)
	_, err := os.Stat(`E:\\Picturesl`)
	fmt.Println(!os.IsNotExist(err))
}

func TestMarshalConfig(t *testing.T) {
	b, e := ioutil.ReadFile(`D:\Dev\gopath\src\github.com\Felamande\filesync\testdata\config.json`)
	if e != nil {
		t.Log(e.Error())
	}

	var config SavedConfig = SavedConfig{}

	e = json.Unmarshal(b, &config.Pairs)
	if e != nil {
		t.Log(e.Error())
	}
	t.Log(config)

}

func TestRename(t *testing.T) {
	w, _ := fsnotify.NewWatcher()
	w.Add("./testdata")
	for {
		select {
		case e := <-w.Events:
			fmt.Println(e.Op, e.Name)
			switch e.Op {
			case fsnotify.Rename:
				fmt.Println("rename: ", e.Name)
			}
		}
	}
}
