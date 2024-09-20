import process from "process";
import repl from "repl";

import WebSocket from "ws";

const args = process.argv;

if (args.length < 3) {
  console.error("no address specified");
  process.exit(1);
}

const ws = new WebSocket("ws://" + args[2]);

function cleanup() {
  ws.close();
  process.exit();
}

function respond(value) {
  replServer.setPrompt("< ");
  replServer.displayPrompt();
  console.log(value);
  replServer.setPrompt("> ");
  replServer.displayPrompt();
}

const replServer = repl.start();
replServer.context.ws = ws;
replServer.ignoreUndefined = true;

replServer.on("SIGINT", () => {
  replServer.once("SIGINT", () => {
    cleanup();
  });
});

ws.onmessage = (event) => {
  respond(event.data);
};

ws.onclose = (event) => {
  let code = event.code;
  let reason = event.reason;
  let clean = event.wasClean ? "clean" : "not clean";
  respond(`Connection closed (${code}, "${reason}", ${clean})`);
};

ws.onerror = (event) => {
  console.log(`Error: ${event.error.code}`);
  cleanup();
};
