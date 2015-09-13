//+build darwin
//+build unix
package uri

import (
	"net/url"
)

func Parse(u string) (*Uri, error) {
	urlp, err := url.Parse(u)
	if err != nil {
		return nil, ParseError{u, err.Error()}
	}

	if urlp.Scheme == "local" && len(urlp.Host) != 0 {
		return nil, ParseError{u, "not a local path."}
	} else if urlp.Scheme != "local" && len(urlp.Host) == 0 {
		return nil, ParseError{u, "need a protocol."}
	}

	return &Uri{urlp.Scheme, urlp.Path, urlp.Host}, nil

}

func (u *Uri) Full() string {
	if u.scheme != "local" {
		return u.scheme + "://" + u.host + u.path
	}
	return u.path
}
