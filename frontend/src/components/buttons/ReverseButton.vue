<script lang="ts" setup>
defineEmits<{ click: [] }>();

import Button from "@/components/Button.vue";

import { state } from "@/store";
</script>

<template>
  <Button
    :class="!state.isReverseOn.value && 'suggested'"
    @click="$emit('click')"
  >
    <template v-if="state.isReverseOn.value">
      <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path class="background" d="M 17.5,6.5 H 6.5 v 11 h 11 z" fill="none" />
        <path d="M 9.33333 10.3333 H 10.3333 V 11.3333 H 9.33333 V 10.3333 Z" fill="black" />
        <path d="M 9.33333 12.3333 H 10.3333 V 13.3333 H 9.33333 V 12.3333 Z" fill="black" />
        <path d="M 10.3333 9.33333 H 11.3333 V 10.3333 H 10.3333 V 9.33333 Z" fill="black" />
        <path d="M 10.3333 13.3333 H 11.3333 V 14.3333 H 10.3333 V 13.3333 Z" fill="black" />
        <path d="M 13.5 10.3333 H 14.5 V 11.3333 H 13.5 V 10.3333 Z" fill="black" />
        <path d="M 13.5 12.3333 H 14.5 V 13.3333 H 13.5 V 12.3333 Z" fill="black" />
        <path d="M 14.5 9.33333 H 15.5 V 10.3333 H 14.5 V 9.33333 Z" fill="black" />
        <path d="M 14.5 13.3333 H 15.5 V 14.3333 H 14.5 V 13.3333 Z" fill="black" />
        <path d="M 8.33333 11.3333 H 9.33333 V 12.3333 H 8.33333 V 11.3333 Z" fill="black" />
        <path d="M 12.5 11.3333H13.5 V 12.3333 H 12.5 V 11.3333 Z" fill="black" />
        <path d="m 6.5,5.5 v 1 h 11 v -1 z" fill="black" />
        <path d="m 18.5,6.5 h -1 v 11 h 1 z" fill="black" />
        <path d="m 6.5,17.5 v 1 h 11 v -1 z" fill="black" />
        <path d="m 5.5,17.5 h 1 v -11 h -1 z" fill="black" />
      </svg>
    </template>
    <template v-else>
      <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path class="background" d="M 17.5,6.5 H 6.5 v 11 h 11 z" fill="none" />
        <path d="M 14.6667 10.3333 H 13.6667 V 11.3333 H 14.6667 V 10.3333 Z" fill="black" />
        <path d="M 14.6667 12.3333 H 13.6667 V 13.3333 H 14.6667 V 12.3333 Z" fill="black" />
        <path d="M 13.6667 9.33333 H 12.6667 V 10.3333 H 13.6667 V 9.33333 Z" fill="black" />
        <path d="M 13.6667 13.3333 H 12.6667 V 14.3333 H 13.6667 V 13.3333 Z" fill="black" />
        <path d="M 10.5 10.3333 H 9.5 V 11.3333 H 10.5 V 10.3333 Z" fill="black" />
        <path d="M 10.5 12.3333 H 9.5 V 13.3333 H 10.5 V 12.3333 Z" fill="black" />
        <path d="M 9.5 9.33333 H 8.5 V 10.3333 H 9.5 V 9.33333 Z" fill="black" />
        <path d="M 9.5 13.3333 H 8.5 V 14.3333 H 9.5 V 13.3333 Z" fill="black" />
        <path d="M 15.6667 11.3333 H 14.6667 V 12.3333 H 15.6667 V 11.3333 Z" fill="black" />
        <path d="M 11.5 11.3333 H 10.5 V 12.3333 H 11.5 V 11.3333 Z" fill="black" />
        <path d="m 6.5,5.5 v 1 h 11 v -1 z" fill="black" />
        <path d="m 18.5,6.5 h -1 v 11 h 1 z" fill="black" />
        <path d="m 6.5,17.5 v 1 h 11 v -1 z" fill="black" />
        <path d="m 5.5,17.5 h 1 v -11 h -1 z" fill="black" />
      </svg>
    </template>
  </Button>
</template>
