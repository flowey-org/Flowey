<script lang="ts" setup>
defineEmits<{ click: [] }>();

import { computed } from "vue";

import Button from "@/components/Button.vue";

import { state } from "@/store";

const text = computed(() => {
  return `${state.buff.value.toFixed(2)}x`;
});
</script>

<template>
  <Button @click="$emit('click')">
    <svg viewBox="0 0 36 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="m 29.5,6.515 h -23 v 11 h 23 z" fill="none" />
      <text x="9.7037077" y="14.074073">
        {{ text }}
      </text>
      <path d="m 6.5,5.515 v 1 h 23 v -1 z" fill="black" />
      <path d="m 6.405,17.485 v 1 H 29.405 v -1 z" fill="black" />
      <path d="m 30.5,6.515 h -1 v 11 h 1 z" fill="black" />
      <path d="m 5.5,17.515 h 1 v -11 h -1 z" fill="black" />
    </svg>
  </Button>
</template>

<style>
svg > text {
  font-size: 5px;
}
</style>
