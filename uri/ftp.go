package uri

import (
	"net"
)

type UriFtp struct {
	conn   net.Conn
	scheme string
	path   string
	host   string
}

func (u *UriFtp) Scheme() string {
	return "ftp"
}

func (u *UriFtp) Host() string {
	return ""
}

