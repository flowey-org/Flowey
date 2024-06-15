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
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

button > svg {
  height: var(--timer-font-size);
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
