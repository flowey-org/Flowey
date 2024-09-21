<script lang="ts" setup>
defineEmits<{ click: [] }>();

import Button from "@/components/Button.vue";
</script>

<template>
  <Button @click="$emit('click')">
    <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="M 17.5,6.5 H 6.5 v 11 h 11 z" fill="none" />
      <path d="M 9.375,9.375 H 8.5 v 5.25 h 0.875 z" fill="black" />
      <path d="M 9.375,8.5 V 9.375 H 11.125 V 8.5 Z" fill="black" />
      <path d="M 9.375,14.625 V 15.5 h 1.75 v -0.875 z" fill="black" />
      <path d="m 11.125,13.75 v 0.875 h 1.75 v -0.875 z" fill="black" />
      <path d="m 12.875,12.875 v 0.875 h 1.75 v -0.875 z" fill="black" />
      <path d="m 11.125,9.375 v 0.875 h 1.75 V 9.375 Z" fill="black" />
      <path d="m 12.875,10.25 v 0.875 h 1.75 v -0.875 z" fill="black" />
      <path d="M 14.625,12.875 H 15.5 v -1.75 h -0.875 z" fill="black" />
      <path d="m 6.5,5.5 v 1 h 11 v -1 z" fill="black" />
      <path d="m 18.5,6.5 h -1 v 11 h 1 z" fill="black" />
      <path d="m 6.5,17.5 v 1 h 11 v -1 z" fill="black" />
      <path d="m 5.5,17.5 h 1 v -11 h -1 z" fill="black" />
    </svg>
  </Button>
</template>
