import { watch } from "vue";

import { State } from "@/state";

const IDB_NAME = "flowey";

class StateStore {
  storeName = "state";
  state: State = new State();

  ready: Promise<void>;
  private resolveReady!: () => void;

  constructor() {
    this.ready = new Promise<void>((resolve) => {
      this.resolveReady = resolve;
    });
  }

  async init() {
    await this.loadState();

    for (const [ref, property] of this.state) {
      watch(ref, () => {
        void this.putState(ref.value, property);
      });
    }
  };

  private async openDB() {
    return new Promise<IDBDatabase>((resolve, reject) => {
      const request = indexedDB.open(IDB_NAME);

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result;
        if (!db.objectStoreNames.contains(this.storeName)) {
          const store = db.createObjectStore(this.storeName);

          store.transaction.oncomplete = () => {
            const transaction = db.transaction(this.storeName, "readwrite");
            const store = transaction.objectStore(this.storeName);

            for (const [ref, property] of this.state) {
              store.put(ref.value, property);
            }
          };
        }
      };

      request.onsuccess = () => {
        resolve(request.result);
      };
      request.onerror = () => {
        reject(new Error(request.error?.message ?? `[IndexedDB] [${this.storeName}] Failed to open a connection`));
      };
    });
  }

  private async loadState() {
    const db = await this.openDB();
    return new Promise<void>((resolve, reject) => {
      const transaction = db.transaction(this.storeName, "readonly");
      const store = transaction.objectStore(this.storeName);

      let pending = 0;

      for (const [ref, property] of this.state) {
        const request = store.get(property);
        pending++;

        request.onsuccess = () => {
          const value: unknown = request.result;
          if (value !== undefined && typeof value === typeof ref.value) {
            ref.value = value;
          } else {
            void this.putState(ref.value, property);
          }

          pending--;
          if (pending === 0) {
            this.resolveReady();
          }
        };
        request.onerror = () => {
          reject(new Error(request.error?.message ?? `[IndexedDB] [${this.storeName}] Failed to get state`));
        };
      }

      resolve();
    });
  }

  private async putState(value: unknown, property: string) {
    const db = await this.openDB();
    return new Promise<void>((resolve, reject) => {
      const transaction = db.transaction(this.storeName, "readwrite");
      const store = transaction.objectStore(this.storeName);
      const request = store.put(value, property);

      request.onsuccess = () => {
        resolve();
      };
      request.onerror = () => {
        reject(new Error(request.error?.message ?? `[IndexedDB] [${this.storeName}] Failed to put state`));
      };
    });
  }
}

export const stateStore = new StateStore();
void stateStore.init();

export const state = stateStore.state;
