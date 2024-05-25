import { ref, watch } from "vue";

import { hoursToMilliseconds } from "@/utils";

class State {
  targetDate = ref(new Date());
  isReverseOn = ref(false);
  maxTime = ref(hoursToMilliseconds(24));

  *[Symbol.iterator]() {
    for (const property in this) {
      yield [this[property], property] as const;
    }
  }
}

const IDB_NAME = "flowey";
const OBJECT_STORE_NAME = "state";

class Database {
  idb: IDBDatabase | null = null;
  state: State = new State();

  async init() {
    this.idb = await this.openIndexedBD();

    for (const [ref, property] of this.state) {
      watch(ref, () => {
        void this.put(ref.value, property);
      });
    }
  };

  private async openIndexedBD() {
    return new Promise<IDBDatabase>((resolve) => {
      const request = window.indexedDB.open(IDB_NAME);

      request.onupgradeneeded = (event) => {
        const idb = (event.target as IDBOpenDBRequest).result;
        const objectStore = idb.createObjectStore(OBJECT_STORE_NAME);

        objectStore.transaction.oncomplete = () => {
          const objectStore = idb
            .transaction(OBJECT_STORE_NAME, "readwrite")
            .objectStore(OBJECT_STORE_NAME);

          for (const [ref, property] of this.state) {
            objectStore.put(ref.value, property);
          }
        };
      };

      request.onsuccess = (event) => {
        const idb = (event.target as IDBOpenDBRequest).result;
        const objectStore = idb
          .transaction(OBJECT_STORE_NAME, "readonly")
          .objectStore(OBJECT_STORE_NAME);

        for (const [ref, property] of this.state) {
          objectStore.get(property).onsuccess = (event) => {
            const value: unknown = (event.target as IDBRequest).result;
            if (value !== undefined) {
              if (typeof value === typeof ref.value) {
                ref.value = value as typeof ref.value;
                return;
              }
            }
            void this.put(ref.value, property);
          };
        }

        resolve(idb);
      };
    });
  }

  async put(value: unknown, property: string) {
    return new Promise<void>((resolve) => {
      this.idb!
        .transaction(OBJECT_STORE_NAME, "readwrite")
        .objectStore(OBJECT_STORE_NAME)
        .put(value, property)
        .onsuccess = () => { resolve(); };
    });
  }
}

const database = new Database();
await database.init();

export default database.state;
