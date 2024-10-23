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
      <path class="background" d="M 22.153846,1.8461538 H 1.8461538 V 22.153846 H 22.153846 Z" fill="none" />
      <!-- eslint-disable-next-line max-len -->
      <path d="m 11.076923,7.3846154 h 1.846154 v 3.6923076 h 3.692308 v 1.846154 h -3.692308 v 3.692308 H 11.076923 V 12.923077 H 7.3846154 v -1.846154 h 3.6923076 z" fill="#000000" />
      <path d="M 1.8461538,0 V 1.8461538 H 22.153846 V 0 Z" fill="#000000" />
      <path d="M 24,1.8461538 H 22.153846 V 22.153846 H 24 Z" fill="#000000" />
      <path d="M 1.8461538,22.153846 V 24 H 22.153846 v -1.846154 z" fill="#000000" />
      <path d="M 0,22.153846 H 1.8461538 V 1.8461538 H 0 Z" fill="#000000" />
    </svg>
  </Button>
</template>
