const path = require("path");
const wasmUrl = path.join(module.path, "wshttp-go.wasm");
require("./wasm_exec_node");
globalThis["WebSocket"] = require("websocket").w3cwebsocket;
const b = fs.readFileSync(wasmUrl);
const go = new Go();

exports.GoFetchInit = Promise.resolve(1)
  .then(() => WebAssembly.instantiate(b, go.importObject))
  .then(({ instance }) => instance)
  .then((inst) => {
    return { process: go.run(inst) };
  });
