package syncer

import (
	"io/ioutil"
	"testing"

	"github.com/Felamande/filesync/uri"
	yaml "gopkg.in/yaml.v2"
)

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

func TestWriteConfig(t *testing.T) {
	s := SavedConfig{
		Pairs: []SyncPairConfig{
			{
				Left:   "local://E:",
				Right:  "local://D",
				Config: SyncConfig{true, true, true},
			},
		},
		Port:      30000,
		LogPath:   "./.log",
		IgnoreExt: []string{"tmp", "temp", "download"},
	}
	b, _ := yaml.Marshal(&s)
	ioutil.WriteFile("testdata/test.config.yaml", b, 0777)

}
