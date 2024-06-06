<script lang="ts" setup>
import { computed, ref, watchEffect } from "vue";

import Block from "@/components/Block.vue";
import Box from "@/components/Box.vue";
import Button from "@/components/Button.vue";
import Text from "@/components/Text.vue";
import Timer from "@/components/Timer.vue";

import AcceptIcon from "@/icons/accept.svg";
import BackwardIcon from "@/icons/backward.svg";
import CancelIcon from "@/icons/cancel.svg";
import ForwardIcon from "@/icons/forward.svg";
import MinusIcon from "@/icons/minus.svg";
import PlayIcon from "@/icons/play.svg";
import PlusIcon from "@/icons/plus.svg";
import StopIcon from "@/icons/stop.svg";

import state from "@/state";
import { formatSeconds, hoursToMilliseconds } from "@/utils";

const nowTime = ref(Date.now());
const isStopping = ref(false);

const currentTime = computed(() => {
  const difference = state.targetDate.value.getTime() - nowTime.value;
  return Math.sign(difference) * Math.min(Math.abs(difference), state.maxTime.value);
});

const isGameOn = computed(() => {
  return state.isReverseOn.value || currentTime.value > 500;
});

const view = computed(() => {
  if (isStopping.value) return "stopping";
  return isGameOn.value ? "gameOn" : "gameOff";
});

const time = computed(() => {
  if (isGameOn.value) {
    const elapsedMilliseconds = Math.abs(currentTime.value);
    const elapsedSeconds = Math.round(elapsedMilliseconds / 1000);
    return formatSeconds(elapsedSeconds);
  } else {
    const seconds = state.maxTime.value / 1000;
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

const MIN_MAX_TIME = hoursToMilliseconds(24);
const MAX_MAX_TIME = hoursToMilliseconds(72);

function increaseMaxTime() {
  state.maxTime.value = Math.min(MAX_MAX_TIME, state.maxTime.value + hoursToMilliseconds(1));
}

function decreaseMaxTime() {
  state.maxTime.value = Math.max(MIN_MAX_TIME, state.maxTime.value - hoursToMilliseconds(1));
}

function reverseTime() {
  state.isReverseOn.value = !state.isReverseOn.value;
  state.targetDate.value = new Date(nowTime.value - currentTime.value);
}

function startGame() {
  nowTime.value = Date.now();
  state.targetDate.value = new Date(nowTime.value + state.maxTime.value);
}

function stopGame() {
  isStopping.value = false;
  state.isReverseOn.value = false;
  state.targetDate.value = new Date(nowTime.value);
}

function toggleStopping() {
  isStopping.value = !isStopping.value;
}
</script>

<template>
  <Box>
    <Timer :time />
    <Block v-if="view==='gameOff'">
      <Button @click="startGame">
        <PlayIcon />
      </Button>
      <Button
        :disabled="state.maxTime.value === MAX_MAX_TIME"
        @click="increaseMaxTime"
      >
        <PlusIcon />
      </Button>
      <Button
        :disabled="state.maxTime.value === MIN_MAX_TIME"
        @click="decreaseMaxTime"
      >
        <MinusIcon />
      </Button>
    </Block>
    <Block v-else-if="view==='gameOn'">
      <Button
        :class="!state.isReverseOn.value && 'suggested'"
        @click="reverseTime"
      >
        <BackwardIcon v-if="state.isReverseOn.value" />
        <ForwardIcon v-else />
      </Button>
      <Button @click="toggleStopping">
        <StopIcon />
      </Button>
    </Block>
    <Block v-else-if="view==='stopping'">
      <Text>You sure?</Text>
      <Button @click="stopGame">
        <AcceptIcon />
      </Button>
      <Button @click="toggleStopping">
        <CancelIcon />
      </Button>
    </Block>
  </Box>
</template>
