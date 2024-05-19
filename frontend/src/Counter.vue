<script setup lang="ts">
import { computed, ref, watchEffect } from "vue";

const targetDate = ref(new Date());
const nowTime = ref(Date.now());

const isGameOn = computed(() => {
  return targetDate.value.getTime() - nowTime.value > 500;
});

function pad(number: number): string {
  return String(number).padStart(2, "0");
}

const counter = computed(() => {
  const elapsedMilliseconds = Math.abs(targetDate.value.getTime() - nowTime.value);
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
  targetDate.value = new Date(nowTime.value + 1000 * 10);
}
</script>

<template>
  <time :datetime="counter" class="counter">{{ counter }}</time>
  <button :disabled="isGameOn" @click="startGame">Start game</button>
</template>
