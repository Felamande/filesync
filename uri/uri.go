package uri

import (
	"errors"
	"io"
	"net/url"
	"reflect"
	"runtime"
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
	Parent() Uri

	Open() (io.ReadWriteCloser, error)
	Remove() error
	Walk(Handler) error
	IsDir() bool
	IsFile() bool

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
		return nil, ProtocolError{urlp.Scheme, "protocol not supported."}
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
	return u.scheme + "://" + u.Host() + u.path
}

func (u *UriLocal) Scheme() string {
	return u.scheme
}

func (u *UriLocal) Open() (io.ReadWriteCloser, error) {

	return nil, errors.New("not implmented")
}

func (u *UriLocal) Remove() error {
	return errors.New("not implemented")
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

func (u *UriLocal) IsFile() bool {
	return !u.IsDir()
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
