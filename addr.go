package wshttp

import "net"

type TCPAddr string

var _ net.Addr = (*TCPAddr)(nil)

func (addr TCPAddr) Network() string { return "tcp" }
func (addr TCPAddr) String() string  { return string(addr) }
