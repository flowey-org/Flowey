import { ref, shallowReactive } from "vue";

import { hoursToMilliseconds } from "@/utils";

class Buff {
  value = 1.00;

  constructor() {
    return shallowReactive(this);
  }

  next(): [number, number] {
    const prev = this.value;
    this.value = Math.round((this.value + 0.05) * 100) / 100;
    (this.value > 2.00) && (this.value = 1.00);
    return [this.value, prev];
  }
}

class MaxTime {
  value = hoursToMilliseconds(24);
  minimum = hoursToMilliseconds(24);
  maximum = hoursToMilliseconds(72);

  constructor() {
    return shallowReactive(this);
  }

  increment() {
    this.value = Math.min(this.maximum, this.value + hoursToMilliseconds(1));
  }

  decrement() {
    this.value = Math.max(this.minimum, this.value - hoursToMilliseconds(1));
  }
}

export class State {
  buff = new Buff();
  isReverseOn = ref(false);
  maxTime = new MaxTime();
  targetDate = ref(Date.now());

  endpoint = ref("");
  username = ref("");
  isLoggedIn = ref(false);

  version = ref(0);

  * [Symbol.iterator]() {
    for (const property in this) {
      const ref = this[property as keyof this];
      if (ref && typeof ref === "object" && "value" in ref) {
        yield [ref, property] as const;
      }
    }
  }

  values() {
    const object: Record<string, unknown> = {};

    for (const [ref, property] of this) {
      object[property] = ref.value;
    }

    return object;
  }

  apply(object: Record<string, unknown>) {
    for (const [ref, property] of this) {
      if (property in object) {
        ref.value = object[property];
      }
    }
  }
}
