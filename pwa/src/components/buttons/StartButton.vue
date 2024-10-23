<script lang="ts" setup>
defineEmits<{ click: [] }>();

import Button from "@/components/Button.vue";
</script>

<template>
  <Button @click="$emit('click')">
    <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="M 22.153846,1.8461538 H 1.8461538 V 22.153846 H 22.153846 Z" fill="none" />
      <path d="M 7.1538462,7.1538462 H 5.5384615 v 9.6923078 h 1.6153847 z" fill="#000000" />
      <path d="M 7.1538462,5.5384615 V 7.1538462 H 10.384615 V 5.5384615 Z" fill="#000000" />
      <path d="m 7.1538462,16.846154 v 1.615384 h 3.2307688 v -1.615384 z" fill="#000000" />
      <path d="m 10.384615,15.230769 v 1.615385 h 3.23077 v -1.615385 z" fill="#000000" />
      <path d="m 13.615385,13.615385 v 1.615384 h 3.230769 v -1.615384 z" fill="#000000" />
      <path d="m 10.384615,7.1538462 v 1.6153846 h 3.23077 V 7.1538462 Z" fill="#000000" />
      <path d="m 13.615385,8.7692308 v 1.6153842 h 3.230769 V 8.7692308 Z" fill="#000000" />
      <path d="m 16.846154,13.615385 h 1.615384 v -3.23077 h -1.615384 z" fill="#000000" />
      <path d="M 1.8461538,0 V 1.8461538 H 22.153846 V 0 Z" fill="#000000" />
      <path d="M 24,1.8461538 H 22.153846 V 22.153846 H 24 Z" fill="#000000" />
      <path d="M 1.8461538,22.153846 V 24 H 22.153846 v -1.846154 z" fill="#000000" />
      <path d="M 0,22.153846 H 1.8461538 V 1.8461538 H 0 Z" fill="#000000" />
    </svg>
  </Button>
</template>
