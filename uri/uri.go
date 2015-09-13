package uri

import (
	"runtime"
)

type Uri struct {
	scheme string
	path   string
	host   string
}

func (u *Uri) Scheme() string {
	return u.scheme
}

type ParseError struct {
	Uri     string
	Message string
}

func (e ParseError) Error() string {
	return "Error parsing " + e.Uri + " on " + runtime.GOOS + ": " + e.Message
}

type ProtocolError struct {
	SchemeCurrent string
	SchemeNeeded  string
}

func (e ProtocolError) Error() string {
	return "Need protocol " + e.SchemeNeeded + ", current protocol is " + e.SchemeCurrent
}
