# Changelog

## [0.0.7] - 2023-03-01

### Improve

- 现在可以使用多个不同的实例而不互相冲突, 通过全局变量`WageFetchExport`设置不一样的函数名称

## [0.0.6] - 2023-03-01

### Fix

- 将 `websocket-polyfill` 替换为 `globalThis["WebSocket"] = require("websocket").w3cwebsocket`.
  `websocket-polyfill` 的实现不规范, 不可用

## [0.0.5] - 2023-03-01

### Improve

- 可通过 `import "@shynome/wshttp"` 来引入 tinygo wasm_exec browser

## [0.0.4] - 2023-02-28

### Change

- WebSocket 重连默认 10 次. 可通过 `WageMaxRetry` 设置

## [0.0.3] - 2023-02-28

### Improve

- 可直接在 node 中使用 GoFetch 了

## [0.0.2] - 2023-02-27

### Fix

- `DisableKeepAlives` 每次链接都打开一个新的 `smux stream`
