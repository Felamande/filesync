package uri

import (
	"testing"
)

func TestUri(t *testing.T) {
	UriString := "local://D:/dir/file"
	_, e := Parse(UriString)
	if e == nil {
		t.Error("Error with " + UriString)
	}

	UriString = "local://D/dir/test"
	_, e = Parse(UriString)
	if e == nil {
		t.Error("Error with " + UriString)
	}

	UriString = "ftp:///D/dir/test"
	_, e = Parse(UriString)
	if e == nil {
		t.Error("Error with " + UriString)
	}

	UriString = "smb:///D/dir/test"
	_, e = Parse(UriString)
	if e == nil {
		t.Error("Error with " + UriString)
	}

	UriString = "local:///home/dir/test"
	_, e = Parse(UriString)
	if e != nil {
		t.Error("Error with " + UriString)
	}

	UriString = "proto://domain/dir/test"
	_, e = Parse(UriString)
	if e != nil {
		t.Error("Error with " + UriString)
	}

}
