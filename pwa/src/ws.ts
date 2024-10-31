import { state, store } from "@/store";

class WebSocketService {
  ws: WebSocket | null;

  private reconnectAttempts = 0;
  private reconnectMaxAttempts = 5;
  private reconnectTimeout = 1000;
  private reconnectTimer: number | null = null;

  private resetAttempts() {
    this.reconnectAttempts = 0;
    this.reconnectTimeout = 1000;
  }

  private incrementAttempts() {
    this.reconnectAttempts++;
    this.reconnectTimeout = Math.min(this.reconnectTimeout * 2, 30000);
  }

  private setupNetworkListeners() {
    window.addEventListener("online", () => {
      console.log("[WebSocket] Network connection restored");
      this.resetAttempts();
      this.connect();
    });

    window.addEventListener("offline", () => {
      console.log("[WebSocket] Network connection lost");
      this.cleanup();
    });
  }

  constructor() {
    this.ws = null;
    this.setupNetworkListeners();
  }

  async init() {
    await store.ready;
    this.connect();
  }

  connect() {
    const url = new URL("ws", state.endpoint.value);

    const message = {
      reconnectAttempt: this.reconnectAttempts,
      maxReconnectAttempts: this.reconnectMaxAttempts,
      reconnectTimeout: this.reconnectTimeout,
      url: url.toString(),
    };

    console.log("[WebSocket] Attempting to connect...", message);

    try {
      this.ws = new WebSocket(url, ["flowey"]);

      this.ws.onopen = () => {
        console.log("[Websocket] Connected");
        this.reconnectAttempts = 0;
        this.reconnectTimeout = 1000;
      };

      this.ws.onclose = (event) => {
        const message = {
          code: event.code,
          reason: event.reason,
          wasClean: event.wasClean,
        };
        console.log("[WebSocket] Connection closed", message);

        if (!event.wasClean) {
          this.scheduleReconnect();
        }
      };

      this.ws.onerror = (error) => {
        console.error("[WebSocket] Encountered error: ", error);
      };

      this.ws.onmessage = (event) => {
        console.log("[WebSocket] Received message: ", event.data);
      };
    } catch (error) {
      console.error("[WebSocket] Connection error:", error);
      this.scheduleReconnect();
    }
  }

  private scheduleReconnect() {
    if (!navigator.onLine) {
      return;
    }

    if (this.reconnectAttempts >= this.reconnectMaxAttempts) {
      console.log("[WebSocket] Max reconnection attempts reached");
      return;
    }

    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer);
    }

    this.reconnectTimer = window.setTimeout(() => {
      if (navigator.onLine) {
        this.incrementAttempts();
        this.connect();
      }
    }, this.reconnectTimeout);
  }

  private cleanup() {
    if (this.ws) {
      this.ws.onopen = null;
      this.ws.onclose = null;
      this.ws.onerror = null;
      this.ws.onmessage = null;

      if (this.ws.readyState === WebSocket.OPEN) {
        this.ws.close();
      }

      this.ws = null;
    }

    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  disconnect() {
    this.cleanup();
    this.resetAttempts();
  }
}

export const ws = new WebSocketService();
void ws.init();
