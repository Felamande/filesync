package uri

import "testing"

type Left interface {
	Uri
}

func TestUriParse(t *testing.T) {
	UriString := "local://D:/dir/file"
	ut, e := Parse(UriString)
	if e != nil {
		t.Log(e.Error())
		return
	}

	if ut.Uri() != UriString {
		t.Error(ut.Uri(), UriString)
	}
	if ut.Host() != "D:" {
		t.Error(UriString)
	}

	if ut.Scheme() != "local" {
		t.Error(UriString)
	}

	UriString = "noreg://test/test"
	_, e = Parse(UriString)
	if e == nil {
		t.Error(UriString)
	}

}

func TestUriLocalOpen(t *testing.T) {
	UriString := "local://D:/Dev/gopath/src/github.com/Felamande/filesync/testdata/testopen.txt"
	u, e := Parse(UriString)
	if e != nil {
		t.Error(e.Error())
		return
	}

	//	w, e := u.OpenWrite()
	//	if e == nil {
	//		t.Error("Not created: ", UriString)
	//		return
	//	}
	e = u.Create(false, 0777)
	if e != nil {
		t.Error(e.Error())
		return
	}

	w, e := u.OpenWrite()
	if e != nil {
		t.Error(e.Error())
		return
	}

	var TestData string = "testdfgsfafhgvfvgdfvdsgdf"
	_, e = w.Write([]byte(TestData))
	if e != nil {
		t.Error(e.Error())
		return
	}
	w.Close()

	var rb []byte = make([]byte, len(TestData))

	r, e := u.OpenRead()
	if e != nil {
		t.Error(e.Error())
		return
	}
	_, e = r.Read(rb)
	if e != nil {
		t.Error(e.Error())
		return
	}

	if string(rb) != TestData {
		t.Error(string(rb))
	}

	r.Close()

	e = u.Remove()
	if e != nil {
		t.Error(e.Error())
		return
	}

	if u.Exist() {
		t.Error(UriString, "not removed.")
	}

}

func TestOpenDir(t *testing.T) {
}
