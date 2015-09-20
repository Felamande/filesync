package syncer

import (
	"testing"

	"github.com/Felamande/filesync/uri"
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
