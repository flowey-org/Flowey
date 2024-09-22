<script lang="ts" setup>
defineProps<{
  time: string;
}>();
</script>

<template>
  <time id="timer" :datetime="time">{{ time }}</time>
</template>

<style>
#timer {
  font-size: var(--timer-font-size);
}
</style>
