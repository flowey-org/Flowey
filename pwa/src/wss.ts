import { watch } from "vue";

import { state, store } from "@/store";

export class WebSocketService {
  ws: WebSocket | null;

  private connectAttempts = 0;
  private connectMaxAttempts = 5;
  private connectTimeout = 1000;
  private connectTimer: number | null = null;

  private statusListeners: (() => void)[] = [];

  constructor() {
    this.ws = null;
    this.setupNetworkListeners();
    this.setupWatchers();
  }

  async init() {
    await store.ready;
  }

  status() {
    if (!navigator.onLine) {
      return "offline";
    }
    if (!this.ws) {
      return "not initialized";
    }
    if (this.connectTimer) {
      return "connecting";
    }
    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return "connecting";
      case WebSocket.OPEN:
        return "connected";
      case WebSocket.CLOSING:
        return "closing";
      case WebSocket.CLOSED:
        return "closed";
      default:
        return "unknown";
    }
  }

  onStatusChange(callback: () => void) {
    this.statusListeners.push(callback);
  }

  private notifyStatusChange() {
    this.statusListeners.forEach((listener) => {
      listener();
    });
  }

  private resetAttempts() {
    this.connectAttempts = 0;
    this.connectTimeout = 1000;
  }

  private incrementAttempts() {
    this.connectAttempts++;
    this.connectTimeout = Math.min(this.connectTimeout * 2, 30000);
  }

  private setupNetworkListeners() {
    window.addEventListener("online", () => {
      console.log("[WebSocket] Network connection restored");
      this.resetAttempts();
      this.connect();
    });

    window.addEventListener("offline", () => {
      console.log("[WebSocket] Network connection lost");
      this.notifyStatusChange();
      this.cleanup();
    });
  }

  private sendUpdate() {
    const clientState = state.values();

    let clientStateString;
    try {
      clientStateString = JSON.stringify(clientState);
    } catch (error) {
      console.error("[WebSocket] Failed to stringify update:", error);
      return;
    }

    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(clientStateString);
    }
  }

  private receiveUpdate(data: unknown) {
    if (data == null || typeof data !== "string") {
      return;
    }

    const serverStateString = data;

    let serverState;
    try {
      serverState = JSON.parse(serverStateString) as Record<string, unknown>;
    } catch (error) {
      console.error("[WebSocket] Failed to parse update:", error);
      serverState = { version: 0 };
    }

    if (!("version" in serverState)) {
      serverState["version"] = 0;
    }

    const version = serverState["version"];
    if (version == state.version.value) {
      return;
    }

    console.log("[WebSocket] Applying update: ", serverState);
    state.apply(serverState);
  }

  private setupWatchers() {
    for (const ref of state) {
      watch(ref, () => {
        this.sendUpdate();
      });
    }
  }

  connect() {
    if (!navigator.onLine) {
      return;
    }

    let url: URL;

    try {
      url = new URL(state.endpoint.value + "ws/");
      url.protocol = url.protocol.replace("http", "ws");
    } catch {
      console.error("[WebSocket] Invalid endpoint URL:", state.endpoint.value);
      return;
    }

    console.log("[WebSocket] Attempting to connect...", {
      reconnectAttempt: this.connectAttempts,
      maxReconnectAttempts: this.connectMaxAttempts,
      reconnectTimeout: this.connectTimeout,
      url: url,
    });
    this.notifyStatusChange();

    try {
      this.ws = new WebSocket(url, ["flowey"]);

      this.ws.onopen = () => {
        console.log("[Websocket] Connected");
        this.notifyStatusChange();

        this.resetAttempts();
        this.sendUpdate();
      };

      this.ws.onclose = (event) => {
        console.log("[WebSocket] Connection closed", {
          code: event.code,
          reason: event.reason,
          wasClean: event.wasClean,
        });
        this.notifyStatusChange();

        if (!event.wasClean) {
          this.scheduleReconnect();
        }
      };

      this.ws.onerror = (error) => {
        console.error("[WebSocket] Encountered error: ", error);
        this.notifyStatusChange();
      };

      this.ws.onmessage = (event) => {
        this.receiveUpdate(event.data);
      };
    } catch (error) {
      console.error("[WebSocket] Connection error:", error);
      this.notifyStatusChange();

      this.scheduleReconnect();
    }
  }

  private scheduleReconnect() {
    if (!navigator.onLine) {
      return;
    }

    if (this.connectAttempts >= this.connectMaxAttempts) {
      console.log("[WebSocket] Max reconnection attempts reached");
      return;
    }

    if (this.connectTimer) {
      window.clearTimeout(this.connectTimer);
    }

    this.connectTimer = window.setTimeout(() => {
      if (navigator.onLine) {
        this.incrementAttempts();
        this.connect();
      }
    }, this.connectTimeout);
  }

  private cleanup() {
    if (this.ws) {
      if (this.ws.readyState === WebSocket.OPEN) {
        this.ws.close(1000, "connection cleaned up");
      }

      this.ws.onopen = null;
      this.ws.onclose = null;
      this.ws.onerror = null;
      this.ws.onmessage = null;

      this.ws = null;
    }

    if (this.connectTimer) {
      window.clearTimeout(this.connectTimer);
      this.connectTimer = null;
    }
  }

  disconnect() {
    console.log("[WebSocket] Disconnected");
    this.notifyStatusChange();

    this.cleanup();
    this.resetAttempts();
  }
}

export const wss = new WebSocketService();
void wss.init();
