package syncer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Felamande/filesync/uri"
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

}

func TestToRight(t *testing.T) {
	ln, _ := uri.Parse("local://D:/pictures/")
	rn, _ := uri.Parse("local://E:/pictures")
	p := &SyncPair{
		Left:  ln,
		Right: rn,
	}
	lnn, _ := uri.Parse("local://D:/pictures/2015-2.jpg")
	rnn, err := p.ToRight(lnn)
	if err != nil {
		t.Error(err.Error())
	}
	if rnn.Uri() != "local://E:/pictures/2015-2.jpg" {
		t.Error("ToRight: " + rnn.Uri())
	}

	ln, _ = uri.Parse("local://D:/pictures/")
	rn, _ = uri.Parse("local://E:/pictures/")
	p = &SyncPair{
		Left:  ln,
		Right: rn,
	}
	lnn, _ = uri.Parse("local://D:/pictures/2015-2.jpg")
	rnn, err = p.ToRight(lnn)
	if err != nil {
		t.Error(err.Error())
	}
	if rnn.Uri() != "local://E:/pictures/2015-2.jpg" {
		t.Error("ToRight: " + rnn.Uri())
	}

	ln, _ = uri.Parse("local://D:/pictures")
	rn, _ = uri.Parse("local://E:/pictures")
	p = &SyncPair{
		Left:  ln,
		Right: rn,
	}
	lnn, _ = uri.Parse("local://D:/pictures/2015-2.jpg")
	rnn, err = p.ToRight(lnn)
	if err != nil {
		t.Error(err.Error())
	}
	if rnn.Uri() != "local://E:/pictures/2015-2.jpg" {
		t.Error("ToRight: " + rnn.Uri())
	}

	ln, _ = uri.Parse("local://D:/pictures")
	rn, _ = uri.Parse("local://E:/pictures/")
	p = &SyncPair{
		Left:  ln,
		Right: rn,
	}
	lnn, _ = uri.Parse("local://D:/pictures/2015-2.jpg")
	rnn, err = p.ToRight(lnn)
	if err != nil {
		t.Error(err.Error())
	}
	if rnn.Uri() != "local://E:/pictures/2015-2.jpg" {
		t.Error("ToRight: " + rnn.Uri())
	}

}
