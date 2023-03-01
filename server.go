package wshttp

import (
	"net"

	"github.com/xtaci/smux"
)

type SmuxListener struct {
	*smux.Session
}

var _ net.Listener = (*SmuxListener)(nil)

func (l *SmuxListener) Accept() (net.Conn, error) {
	return l.AcceptStream()
}

func (l *SmuxListener) Addr() net.Addr { return nil }
