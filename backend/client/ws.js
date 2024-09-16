import process from "process";
import repl from "repl";

import WebSocket from "ws";

const args = process.argv;

if (args.length < 3) {
  console.error("no address specified");
  process.exit(1);
}

const ws = new WebSocket("ws://" + args[2]);

const replServer = repl.start();
replServer.context.ws = ws;
replServer.ignoreUndefined = true;

replServer.on("SIGINT", () => {
  replServer.once("SIGINT", () => {
    ws.close();
    process.exit();
  });
});

ws.onmessage = (event) => {
  replServer.setPrompt("< ");
  replServer.displayPrompt();
  console.log(event.data);
  replServer.setPrompt("> ");
  replServer.displayPrompt();
};