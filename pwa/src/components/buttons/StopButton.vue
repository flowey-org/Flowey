<script lang="ts" setup>
defineEmits<{ click: [] }>();

import Button from "@/components/Button.vue";
</script>

<template>
  <Button @click="$emit('click')">
    <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="M 22.153846,1.8461538 H 1.8461538 V 22.153846 H 22.153846 Z" fill="none" />
      <path d="M 15.692308,6.4615385 V 8.3076923 H 8.3076923 V 6.4615385 Z" fill="#000000" />
      <path d="m 15.692308,15.692308 v 1.846154 H 8.3076923 v -1.846154 z" fill="#000000" />
      <path d="M 17.538462,15.692308 H 15.692308 V 8.3076923 h 1.846154 z" fill="#000000" />
      <path d="M 8.3076923,15.692308 H 6.4615385 V 8.3076923 h 1.8461538 z" fill="#000000" />
      <path d="M 1.8461538,0 V 1.8461538 H 22.153846 V 0 Z" fill="#000000" />
      <path d="M 24,1.8461538 H 22.153846 V 22.153846 H 24 Z" fill="#000000" />
      <path d="M 1.8461538,22.153846 V 24 H 22.153846 v -1.846154 z" fill="#000000" />
      <path d="M 0,22.153846 H 1.8461538 V 1.8461538 H 0 Z" fill="#000000" />
    </svg>
  </Button>
</template>
