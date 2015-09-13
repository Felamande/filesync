package uri

import (
	"net/url"
	"strings"
)

func Parse(u string) (*Uri, error) {
	url, err := url.Parse(u)
	if err != nil {
		return nil, ParseError{u, err.Error()}
	}
	if url.Scheme == "local" && !strings.Contains(url.Host, ":") {
		return nil, ParseError{u, "not a local path."}
	}
	if len(url.Host) == 0 {
		return nil, ParseError{u, "Contains no domain"}
	}
	return &Uri{url.Scheme, url.Path, url.Host}, nil

}

func (u *Uri) Drive() (string, error) {
	if u.Scheme() == "local" {
		return u.host, nil
	}
	return "", ProtocolError{"local", u.scheme}

}

func (u *Uri) Full() string {
	if u.scheme != "local" {
		return u.scheme + "://" + u.host + u.path
	}
	return u.host + u.path
}
