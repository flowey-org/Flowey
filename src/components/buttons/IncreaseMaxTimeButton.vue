<script lang="ts" setup>
defineEmits<{ click: [] }>();

import Button from "@/components/Button.vue";

import { state } from "@/store";
</script>

<template>
  <Button
    :disabled="state.maxTime.value === state.maxTime.maximum"
    @click="$emit('click')"
  >
    <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="M 17.5,6.5 H 6.5 v 11 h 11 z" fill="none" />
      <path d="m 11.5,9.5 h 1 v 2 h 2 v 1 h -2 v 2 h -1 v -2 H 9.5 v -1 h 2 z" fill="black" />
      <path d="m 6.5,5.5 v 1 h 11 v -1 z" fill="black" />
      <path d="m 18.5,6.5 h -1 v 11 h 1 z" fill="black" />
      <path d="m 6.5,17.5 v 1 h 11 v -1 z" fill="black" />
      <path d="m 5.5,17.5 h 1 v -11 h -1 z" fill="black" />
    </svg>
  </Button>
</template>
