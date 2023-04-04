const { GoFetchInit } = require("./index");

GoFetchInit.then(async () => {
  const fetch = wshttpGen("ws://127.0.0.1:3000/", { Version: 2 });
  let req = new Request("http://127.0.0.1:3000/hello1");
  let r = await fetch(req).then((r) => r.text());
  console.log(r);
  process.exit(0);
});
