<script setup lang="ts">
import { computed, ref, watchEffect } from "vue";

import state from "@/state";
import { hoursToMilliseconds } from "@/utils";

const nowTime = ref(Date.now());

const currentTime = computed(() => {
  const difference = state.targetDate.value.getTime() - nowTime.value;
  return Math.sign(difference) * Math.min(Math.abs(difference), state.maxTime.value);
});

const isGameOn = computed(() => {
  return state.isReverseOn.value || currentTime.value > 500;
});

function pad(number: number): string {
  return String(number).padStart(2, "0");
}

const counter = computed(() => {
  if (!isGameOn.value) {
    return "00:00:00";
  }
  const elapsedMilliseconds = Math.abs(currentTime.value);
  const elapsedSeconds = Math.round(elapsedMilliseconds / 1000);
  const hours = Math.floor(elapsedSeconds / 3600);
  const minutes = Math.floor(elapsedSeconds % 3600 / 60);
  const seconds = elapsedSeconds % 60;
  return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
});

// Doing it this way allows us to avoid time drift
//
// References:
// - https://gist.github.com/jakearchibald/cb03f15670817001b1157e62a076fe95
// - https://youtu.be/MCi6AZMkxcU
function setAnimation(callback: () => void, interval: number) {
  const start = document.timeline.currentTime as number;

  function frame(time: number) {
    if (!isGameOn.value) return;
    callback();
    scheduleFrame(time);
  }

  function scheduleFrame(time: number) {
    const elapsed = time - start;
    const roundedElapsed = Math.round(elapsed / interval) * interval;
    const targetNext = start + roundedElapsed + interval;
    const delay = targetNext - performance.now();
    setTimeout(() => requestAnimationFrame(frame), delay);
  }

  scheduleFrame(start);
}

watchEffect(() => {
  if (isGameOn.value) {
    setAnimation(() => {
      nowTime.value = Date.now();
    }, 1000);
  }
});

function startGame() {
  nowTime.value = Date.now();
  state.targetDate.value = new Date(nowTime.value + state.maxTime.value);
}

function reverseTime() {
  state.isReverseOn.value = !state.isReverseOn.value;
  state.targetDate.value = new Date(nowTime.value - currentTime.value);
}

function stopGame() {
  state.isReverseOn.value = false;
  state.targetDate.value = new Date(nowTime.value);
}
</script>

<template>
  <time :datetime="counter" class="counter">{{ counter }}</time>
  <br>
  <label for="maxTimeSelect">Start with</label>
  <select id="maxTimeSelect" v-model="state.maxTime.value" :disabled="isGameOn">
    <option :value="10 * 1000">10s</option>
    <option :value="hoursToMilliseconds(24)">24h</option>
    <option :value="hoursToMilliseconds(48)">48h</option>
    <option :value="hoursToMilliseconds(72)">72h</option>
  </select>
  <button :disabled="isGameOn" @click="startGame">Start game</button>
  <button :disabled="!isGameOn" @click="stopGame">Stop game</button>
  <button :disabled="!isGameOn" @click="reverseTime">
    {{ state.isReverseOn.value ? "Disable reverse": "Enable reverse" }}
  </button>
</template>
