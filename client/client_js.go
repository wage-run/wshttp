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
	js.Global().Set("wshttpGen", js.FuncOf(wshttpGen))
	<-make(chan any)
}

var defaultConfig js.Value

func init() {
	c := js.Global().Get("Object").New()
	c.Set("max_retry", 10)
	defaultConfig = c
}

func ifValDo(val js.Value, fn func(js.Value)) {
	if val.IsUndefined() {
		return
	}
	fn(val)
}

func wshttpGen(this js.Value, args []js.Value) any {
	endpoint := args[0].String()
	if endpoint == "" {
		panic("endpoint is not set")
	}
	config := js.Global().Get("Object").New()
	iconfig := js.Undefined()
	if len(args) == 2 {
		iconfig = args[1]
	}
	js.Global().Get("Object").Call("assign", config, defaultConfig, iconfig)

	smuxConfig := smux.DefaultConfig()
	ifValDo(config.Get("Version"), func(v js.Value) { smuxConfig.Version = v.Int() })
	ifValDo(config.Get("KeepAliveDisabled"), func(v js.Value) { smuxConfig.KeepAliveDisabled = v.Bool() })
	ifValDo(config.Get("KeepAliveInterval"), func(v js.Value) {
		if d, err := time.ParseDuration(v.String()); err == nil {
			smuxConfig.KeepAliveInterval = d
		}
	})
	ifValDo(config.Get("KeepAliveTimeout"), func(v js.Value) {
		if d, err := time.ParseDuration(v.String()); err == nil {
			smuxConfig.KeepAliveTimeout = d
		}
	})
	ifValDo(config.Get("MaxFrameSize"), func(v js.Value) { smuxConfig.MaxFrameSize = v.Int() })
	ifValDo(config.Get("MaxReceiveBuffer"), func(v js.Value) { smuxConfig.MaxReceiveBuffer = v.Int() })
	ifValDo(config.Get("MaxStreamBuffer"), func(v js.Value) { smuxConfig.MaxStreamBuffer = v.Int() })

	var session *smux.Session
	var locker = &sync.RWMutex{}
	var maxRetry = config.Get("max_retry").Int()
	var connect = func() (err error) {
		ctx := context.Background()
		conn, _, err := websocket.Dial(ctx, endpoint, nil)
		if err != nil {
			return
		}
		rwc := wshttp.NewWSConn(conn)
		session, err = smux.Client(rwc, smuxConfig)
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
