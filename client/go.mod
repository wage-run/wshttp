module github.com/wage-run/wshttp/client

go 1.20

replace github.com/wage-run/wshttp => ../

require (
	github.com/nlepage/go-js-promise v1.1.0
	github.com/shynome/wahttp v0.0.3
	github.com/wage-run/wshttp v0.0.0-00010101000000-000000000000
	github.com/xtaci/smux v1.5.24
	nhooyr.io/websocket v1.8.7
)

require github.com/klauspost/compress v1.10.3 // indirect
