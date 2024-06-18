<script lang="ts" setup>
import { computed, ref, watchEffect } from "vue";

import Block from "@/components/Block.vue";
import Box from "@/components/Box.vue";
import Text from "@/components/Text.vue";
import Timer from "@/components/Timer.vue";

import AcceptButton from "@/components/buttons/AcceptButton.vue";
import BuffButton from "@/components/buttons/BuffButton.vue";
import CancelButton from "@/components/buttons/CancelButton.vue";
import DecreaseMaxTimeButton from "@/components/buttons/DecreaseMaxTimeButton.vue";
import IncreaseMaxTimeButton from "@/components/buttons/IncreaseMaxTimeButton.vue";
import ReverseButton from "@/components/buttons/ReverseButton.vue";
import StartButton from "@/components/buttons/StartButton.vue";
import StopButton from "@/components/buttons/StopButton.vue";

import state from "@/state";
import { formatSeconds, hoursToMilliseconds } from "@/utils";

const nowTime = ref(Date.now());
const isStopping = ref(false);

const currentDifference = computed(() => {
  const difference = state.refs.targetDate.value - nowTime.value;
  let distance = Math.abs(difference);
  if (difference < 0) {
    distance *= state.computed.buff.value;
  }
  return Math.sign(difference) * Math.min(distance, state.refs.maxTime.value);
});

const isGameOn = computed(() => {
  return state.refs.isReverseOn.value || currentDifference.value > 500;
});

const view = computed(() => {
  if (isStopping.value) return "stopping";
  return isGameOn.value ? "gameOn" : "gameOff";
});

const time = computed(() => {
  if (isGameOn.value) {
    const elapsedMilliseconds = Math.abs(currentDifference.value);
    const elapsedSeconds = Math.round(elapsedMilliseconds / 1000);
    return formatSeconds(elapsedSeconds);
  } else {
    const seconds = state.refs.maxTime.value / 1000;
    return formatSeconds(seconds);
  }
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

function increaseMaxTime() {
  state.refs.maxTime.value = Math.min(state.refs.maxTime.maximum, state.refs.maxTime.value + hoursToMilliseconds(1));
}

function decreaseMaxTime() {
  state.refs.maxTime.value = Math.max(state.refs.maxTime.minimum, state.refs.maxTime.value - hoursToMilliseconds(1));
}

function reverseTime() {
  let difference = currentDifference.value;
  if (difference > 0) {
    difference /= state.computed.buff.value;
  }
  state.refs.targetDate.value = nowTime.value - difference;

  state.refs.isReverseOn.value = !state.refs.isReverseOn.value;
}

function startGame() {
  nowTime.value = Date.now();
  state.refs.targetDate.value = nowTime.value + state.refs.maxTime.value;
}

function stopGame() {
  isStopping.value = false;
  state.refs.isReverseOn.value = false;
  state.refs.targetDate.value = nowTime.value;
}

function toggleStopping() {
  isStopping.value = !isStopping.value;
}

function nextBuff() {
  const [buff, prevBuff] = state.nextBuff();
  const difference = state.refs.targetDate.value - nowTime.value;
  let distance = Math.abs(difference);
  if (difference < 0) {
    distance = distance * prevBuff / buff;
    state.refs.targetDate.value = nowTime.value - distance;
  }
}
</script>

<template>
  <Box>
    <Timer :time />
    <Block v-if="view==='gameOff'">
      <StartButton @click="startGame" />
      <BuffButton @click="nextBuff" />
      <IncreaseMaxTimeButton @click="increaseMaxTime" />
      <DecreaseMaxTimeButton @click="decreaseMaxTime" />
    </Block>
    <Block v-else-if="view==='gameOn'">
      <ReverseButton @click="reverseTime" />
      <BuffButton @click="nextBuff" />
      <StopButton @click="toggleStopping" />
    </Block>
    <Block v-else-if="view==='stopping'">
      <Text>You sure?</Text>
      <AcceptButton @click="stopGame" />
      <CancelButton @click="toggleStopping" />
    </Block>
  </Box>
</template>
