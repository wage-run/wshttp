package main

import (
	"context"
	"math"
	"net"
	"net/http"
	"sync"
	"syscall/js"
	"time"

	promise "github.com/nlepage/go-js-promise"
	"github.com/shynome/wahttp"
	"github.com/wage-run/wshttp"
	"github.com/xtaci/smux"
	"nhooyr.io/websocket"
)

func main() {
	var ep = js.Global().Get("WageEndpoint").String()
	var GoFetchExportName = "GoFetch"
	if v := js.Global().Get("WageFetchExport"); v.Type() == js.TypeString {
		GoFetchExportName = v.String()
	}
	js.Global().Set(GoFetchExportName, GoFetch(ep))
	<-make(chan any)
}

func getJsMaxRetry() int {
	if re := js.Global().Get("WageMaxRetry"); re.Type() == js.TypeNumber {
		return re.Int()
	} else {
		return 10
	}
}

func GoFetch(endpoint string) js.Func {

	if endpoint == "" {
		panic("env Endpoint is not set")
	}

	var session *smux.Session
	var locker = &sync.RWMutex{}
	var maxRetry = getJsMaxRetry()
	var connect = func() (err error) {
		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, endpoint, nil)
		if err != nil {
			return
		}
		rwc := wshttp.NewWSConn(conn)
		session, err = smux.Client(rwc, nil)
		if err != nil {
			return
		}
		return
	}
	var retryConnect = func(max int) {
		locker.Lock()
		defer locker.Unlock()
		if max == 0 {
			max = math.MaxInt
		}
		for i := 0; i < max; i++ {
			if err := connect(); err == nil {
				return
			}
			time.Sleep(time.Second)
		}
	}
	go retryConnect(0)

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				locker.RLock()
				defer locker.RUnlock()
				conn, err := session.OpenStream()
				if err != nil {
					locker.RUnlock()
					retryConnect(maxRetry)
					locker.RLock()
					conn, err = session.OpenStream()
				}
				return conn, err
			},
		},
	}

	return js.FuncOf(func(this js.Value, args []js.Value) any {
		p, resolve, reject := promise.New()
		go func() {
			var err error
			defer func() {
				if err != nil {
					reject(err.Error())
				}
			}()
			req, err := wahttp.JsRequest(args[0])
			if err != nil {
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			resolve(wahttp.GoResponse(resp))
		}()
		return p
	})
}
