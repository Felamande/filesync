package uri

import (
	"testing"
)

func TestUriParse(t *testing.T) {
	UriString := "local:///D:/dir/file"
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

	UriString = "local://D:/dir/file"
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

func TestGetDrive(t *testing.T) {
	UriString := "local://D:/dir/test"
	u, e := Parse(UriString)
	if e != nil {
		t.Error("Error with " + UriString)
	}

	d, e := u.Drive()
	if e != nil {
		t.Error(e.Error())
	}
	if len(d) == 0 {
		t.Error("Error getting drive of path " + UriString)
	}

	UriString = "ftp://D:/dir/test"
	u, e = Parse(UriString)
	if e != nil {
		t.Error(e.Error())
	}

	d, e = u.Drive()

	if e == nil {
		t.Error("Cannot get drive of " + UriString + ", the protocol is not local.")
	}
}
