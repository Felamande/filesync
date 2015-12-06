package uri

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var protocolRegistry map[string]reflect.Type

type Handler func(root Uri, u Uri) error

func init() {
	protocolRegistry = make(map[string]reflect.Type, 4)
	protocolRegistry["local"] = reflect.TypeOf((*UriLocal)(nil)).Elem()
	protocolRegistry["ftp"] = reflect.TypeOf((*UriFtp)(nil)).Elem()
}

type Uri interface {
	Scheme() string
	Host() string
	Path() string
	Uri() string
	Abs() string
	Parent() (Uri, error)

	Create(bool, os.FileMode) error
	OpenRead() (io.ReadCloser, error)
	OpenWrite() (io.WriteCloser, error)
	Remove() error
	Walk(dh, fh Handler) error //dh is the handler to handle a directory,and fh is the handler to handle a file.

	IsDir() bool
	Exist() bool
	IsAbs() bool

	Mode() os.FileMode
	ModTime() time.Time

	setHost(h string)
	setPath(p string)
	setScheme(s string)
}

func Parse(u string) (Uri, error) {
	urlp, err := url.Parse(u)
	if err != nil {
		return nil, ParseError{u, err.Error()}
	}

	UriType, exist := protocolRegistry[urlp.Scheme]
	if !exist {
		return nil, ProtocolError{urlp.Scheme, "protocol not supported."}
	}

	UriVal := reflect.New(UriType)
	i := UriVal.Interface()
	Urip, ok := i.(Uri)
	if !ok {
		return nil, ProtocolError{urlp.Scheme, "protocol not fully supported."}
	}
	Urip.setScheme(urlp.Scheme)
	Urip.setHost(urlp.Host)
	Urip.setPath(urlp.Path)

	return Urip, nil

}

type UriLocal struct {
	scheme string
	path   string
	host   string
}

func (u *UriLocal) Host() string {
	return u.host

}

func (u *UriLocal) Uri() string {
	return u.scheme + "://" + u.host + u.path
}

func (u *UriLocal) Abs() string {
	return u.host + u.path
}

func (u *UriLocal) IsAbs() bool {
	return filepath.IsAbs(u.Abs())
}

func (u *UriLocal) Scheme() string {
	return u.scheme
}

func (u *UriLocal) Mode() os.FileMode {
	fi, err := os.Stat(u.host + u.path)
	if err != nil {
		return os.ModePerm
	}
	return fi.Mode()
}

func (u *UriLocal) Exist() bool {
	fi, _ := os.Stat(u.Abs())
	if fi != nil {
		return true
	}
	return false
}

func (u *UriLocal) ModTime() time.Time {
	fi, _ := os.Stat(u.host + u.path)
	if fi == nil {
		return time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	return fi.ModTime()
}

func (u *UriLocal) Create(IsDir bool, m os.FileMode) (err error) {

	if u.Exist() {
		return nil
	}

	if IsDir {
		err = os.Mkdir(u.Abs(), m)
		if err != nil {
			return
		}
	} else {
		var fd *os.File
		fd, err = os.OpenFile(u.Abs(), os.O_CREATE, m)
		defer fd.Close()
		if err != nil {
			return
		}
	}
	return nil
}

func (u *UriLocal) OpenRead() (io.ReadCloser, error) {
	AbsPath := u.host + u.path

	if !filepath.IsAbs(AbsPath) {
		return nil, OpenError{u.Uri(), "is not an absolute path."}
	}
	return os.OpenFile(AbsPath, os.O_RDONLY, u.Mode())

}
func (u *UriLocal) OpenWrite() (io.WriteCloser, error) {
	AbsPath := u.Abs()

	if !filepath.IsAbs(AbsPath) {
		return nil, OpenError{u.Uri(), "is not an absolute path."}
	}

	return os.OpenFile(AbsPath, os.O_WRONLY, u.Mode())

}

func (u *UriLocal) Remove() error {
	return os.Remove(u.Abs())
}

func (u *UriLocal) Walk(dh, fh Handler) error {

	if !u.IsDir() {
		return errors.New("walk " + u.Abs() + ": is not a directory")
	}

	err := filepath.Walk(u.Abs(),
		func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}
			path = strings.Replace(path, "\\", "/", -1)
			urip, err := Parse(u.Scheme() + "://" + path)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			if urip.IsDir() {
				err = dh(u, urip)
			} else {
				err = fh(u, urip)
			}
			if err != nil {
				fmt.Println(err)
				return nil
			}
			return nil
		},
	)

	return err
}

func (u *UriLocal) Path() string {

	return u.path
}

func (u *UriLocal) Parent() (Uri, error) {
	p := filepath.Dir(u.Abs())
	p = strings.Replace(p, "\\", "/", -1)
	return Parse(u.Scheme() + "://" + p)
}

func (u *UriLocal) IsDir() bool {
	fi, _ := os.Stat(u.Abs())
	if fi == nil {
		return false
		//TODO
	}

	return fi.IsDir()
}

func (u *UriLocal) setHost(h string) {
	u.host = h
}
func (u *UriLocal) setPath(p string) {
	u.path = p
}
func (u *UriLocal) setScheme(s string) {
	u.scheme = s
}

type ParseError struct {
	Uri     string
	Message string
}

func (e ParseError) Error() string {
	return "Error parsing " + e.Uri + " on " + runtime.GOOS + ": " + e.Message
}

type ProtocolError struct {
	Protocol string
	Message  string
}

func (e ProtocolError) Error() string {
	return "When handling protocol " + e.Protocol + ": " + e.Message
}

type OpenError struct {
	Uri     string
	Message string
}

func (e OpenError) Error() string {
	return "Open " + e.Uri + "error: " + e.Message
}
