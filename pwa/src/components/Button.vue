<script lang="ts" setup>
defineEmits<{ click: [] }>();
</script>

<template>
  <button @click="$emit('click')">
    <slot />
  </button>
</template>

<style>
button {
  all: unset;
  display: flex;
  align-items: center;
  justify-content: center;
  height: calc(var(--timer-font-size) / 1.75);
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

button > svg {
  box-sizing: border-box;
  padding: 0;
  margin: 0;
  height: 100%;
  max-width: 100%;
  max-height: 100%;
}

button > svg > .background {
  fill: #a6dd8e
}

button.suggested > svg > .background {
  fill: #ff8c42
}

button:disabled > svg > .background {
  fill: #a2aebb
}
</style>
