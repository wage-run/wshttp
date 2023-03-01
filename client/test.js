globalThis["WageEndpoint"] = "ws://127.0.0.1:3000/live";
const { GoFetchInit } = require("./index");

GoFetchInit.then(async () => {
  let req = new Request("http://127.0.0.1:3000/live/events/");
  let r = await GoFetch(req).then((r) => r.json());
  console.log(r);
  process.exit(0);
});
