tinywasm:
	tinygo build -target wasm -o wshttp.wasm -no-debug -panic=trap && make wasm

wasm:
	GOOS=js GOARCH=wasm go build -o wshttp-go.wasm

build: wasm tinywasm

cp_execjs:
	cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js wasm_exec_tinygo.js
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js wasm_exec.js

test: wasm cp_execjs
	node test.js

prepack: build cp_execjs
	rm -f ./*.tgz
