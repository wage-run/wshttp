package wshttp

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"nhooyr.io/websocket"
)

type WSConn struct {
	conn *websocket.Conn

	reader       io.Reader
	readerLocker sync.Locker

	LAddr net.Addr
	RAddr net.Addr
}

var _ io.ReadWriteCloser = (*WSConn)(nil)

func NewWSConn(conn *websocket.Conn) (ws *WSConn) {
	return &WSConn{
		conn: conn,

		readerLocker: &sync.Mutex{},
	}
}

func (ws *WSConn) Read(p []byte) (n int, err error) {

	ws.readerLocker.Lock()
	defer ws.readerLocker.Unlock()

	if ws.reader == nil {
		ctx := context.Background()
		_, ws.reader, err = ws.conn.Reader(ctx)
		if err != nil {
			return
		}
	}
	n, err = ws.reader.Read(p)
	if errors.Is(err, io.EOF) {
		ws.reader = nil
		err = nil
	}
	return
}

func (ws *WSConn) Write(p []byte) (n int, err error) {
	ctx := context.Background()
	writer, err := ws.conn.Writer(ctx, websocket.MessageBinary)
	if err != nil {
		return
	}
	defer writer.Close()
	return writer.Write(p)
}

func (ws *WSConn) Close() (err error) {
	return ws.conn.Close(websocket.StatusNormalClosure, "")
}

func (ws *WSConn) LocalAddr() net.Addr  { return ws.LAddr }
func (ws *WSConn) RemoteAddr() net.Addr { return ws.RAddr }
