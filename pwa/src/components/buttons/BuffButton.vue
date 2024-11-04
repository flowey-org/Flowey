<script lang="ts" setup>
defineEmits<{ click: [] }>();

import { computed } from "vue";

import Button from "@/components/Button.vue";

import { state } from "@/state";

const text = computed(() => {
  return `${state.buff.value.toFixed(2)}x`;
});
</script>

<template>
  <Button @click="$emit('click')">
    <svg viewBox="0 0 36 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="M 34.56,1.8504241 H 1.44 V 22.20509 h 33.12 z" fill="none" />
      <path d="M 1.44,0 V 1.8504241 H 34.56 V 0 Z" fill="#000000" />
      <path d="M 1.3032,22.149576 V 24 h 33.12 v -1.850424 z" fill="#000000" />
      <path d="M 36,1.8504241 H 34.56 V 22.20509 H 36 Z" fill="#000000" />
      <path d="M 0,22.20509 H 1.44 V 1.8504241 H 0 Z" fill="#000000" />
      <text>
        <tspan x="18.249996" y="15.5" style="text-align:center;text-anchor:middle">
          {{ text }}
        </tspan>
      </text>
    </svg>
  </Button>
</template>

<style>
svg > text {
  font-size: 9px;
}
</style>
