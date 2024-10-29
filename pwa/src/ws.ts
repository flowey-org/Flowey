import { state, store } from "@/store";

class WebSocketService {
  ws: WebSocket | null;

  constructor() {
    this.ws = null;
  }

  async init() {
    await store.ready;
    this.connect(new URL("ws", state.endpoint.value));
  }

  connect(url: URL) {
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log("[Websocket] Connected");
    };

    this.ws.onclose = (event) => {
      const code = event.code;
      const reason = event.reason;
      const clean = event.wasClean ? "clean" : "not clean";
      console.log(`[WebSocket] Connection closed (${code}, "${reason}", ${clean})`);
    };

    this.ws.onerror = (error) => {
      console.error("[WebSocket] Encountered error: ", error);
    };

    this.ws.onmessage = (event) => {
      console.log("[WebSocket] Received message: ", event.data);
    };
  }
}

export const ws = new WebSocketService();
void ws.init();
