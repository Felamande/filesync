package uri

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"
)

var protocolRegistry map[string]reflect.Type

type Handler func(u Uri) error

func init() {
	protocolRegistry = make(map[string]reflect.Type, 4)
	protocolRegistry["local"] = reflect.TypeOf((*UriLocal)(nil)).Elem()
}

type Uri interface {
	Scheme() string
	Host() string
	Path() string
	Uri() string
	Abs() string
	Parent() Uri

	Create(bool, os.FileMode) error
	OpenRead() (io.ReadCloser, error)
	OpenWrite() (io.WriteCloser, error)
	Remove() error
	Walk(Handler) error

	IsDir() bool
	Exist() bool

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
	fi, err := os.Stat(u.host + u.path)
	if err != nil {
		return time.Now()
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
	if u.IsDir() {
		return nil, OpenError{u.Uri(), "is a directory."}
	}
	AbsPath := u.host + u.path

	if !filepath.IsAbs(AbsPath) {
		return nil, OpenError{u.Uri(), "is not an absolute path."}
	}
	return os.OpenFile(AbsPath, os.O_RDONLY, u.Mode())

}
func (u *UriLocal) OpenWrite() (io.WriteCloser, error) {
	if u.IsDir() {
		return nil, OpenError{u.Uri(), "is a directory."}
	}
	AbsPath := u.Abs()

	if !filepath.IsAbs(AbsPath) {
		return nil, OpenError{u.Uri(), "is not an absolute path."}
	}

	return os.OpenFile(AbsPath, os.O_WRONLY, u.Mode())

}

func (u *UriLocal) Remove() error {
	return os.Remove(u.host + u.path)
}

func (u *UriLocal) Walk(h Handler) error {
	return nil
}

func (u *UriLocal) Path() string {
	return u.path
}

func (u *UriLocal) Parent() Uri {
	return nil
}

func (u *UriLocal) IsDir() bool {
	return false
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
