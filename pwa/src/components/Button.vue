<script lang="ts" setup>
defineEmits<{ click: [] }>();
</script>

<template>
  <button @click="$emit('click')">
    <slot />
  </button>
</template>

<style>
:root {
  --color-button: #a6dd8e;
  --color-button-disabled: #a2aebb;
  --color-button-suggested: #ff8c42;
}

button {
  all: unset;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 1px;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

button > svg {
  box-sizing: border-box;
  height: 100%;
}

button > svg > .background {
  fill: var(--color-button);
}

button.suggested > svg > .background {
  fill: var(--color-button-suggested);
}

button:disabled > svg > .background {
  fill: var(--color-button-disabled);
  cursor: not-allowed;
}

button:not(button:disabled):is(:focus-visible, :hover) > svg > .background {
  filter: brightness(0.9);
}
</style>
