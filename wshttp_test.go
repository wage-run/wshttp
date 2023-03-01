package wshttp

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/xtaci/smux"
	"nhooyr.io/websocket"
)

var l net.Listener

var smuxConfig *smux.Config

func TestMain(m *testing.M) {

	l = try.To1(net.Listen("tcp", "127.0.0.1:0"))

	smuxConfig = smux.DefaultConfig()
	smuxConfig.Version = 2

	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			http.Error(w, err.Error(), 500)
		})

		conn := try.To1(websocket.Accept(w, r, nil))
		rwc := NewWSConn(conn)
		rwc.RAddr = TCPAddr(r.RemoteAddr)
		session := try.To1(smux.Server(rwc, smuxConfig))
		defer session.Close()

		l := &SmuxListener{Session: session}

		mux := http.NewServeMux()
		mux.HandleFunc("/hello1", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "hello1")
		})
		mux.HandleFunc("/hello2", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "hello2")
		})

		try.To(http.Serve(l, mux))

	}))
	m.Run()
}

func TestClient(t *testing.T) {
	ctx := context.Background()
	conn, _ := try.To2(websocket.Dial(ctx, fmt.Sprintf("ws://%s", l.Addr().String()), nil))
	rwc := NewWSConn(conn)
	session := try.To1(smux.Client(rwc, smuxConfig))
	defer session.Close()

	client := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return session.OpenStream()
			},
		},
	}

	testCases := []struct {
		path   string
		status int
		body   string
	}{
		{"/hello1", 200, "hello1"},
		{"/hello2", 200, "hello2"},
		{"/hello3", 404, ""},
	}

	for _, item := range testCases {
		t.Run("client request", func(t *testing.T) {

			defer err2.Catch(func(err error) {
				t.Error(err)
			})

			r := try.To1(client.Get(fmt.Sprintf("http://127.0.0.1%s", item.path)))

			assert.Equal(r.StatusCode, item.status)

			if item.body == "" {
				return
			}

			b := try.To1(io.ReadAll(r.Body))

			assert.Equal(string(b), item.body)

		})
	}
}
