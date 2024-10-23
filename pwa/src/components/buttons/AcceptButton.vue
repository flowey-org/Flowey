<script lang="ts" setup>
defineEmits<{ click: [] }>();

import Button from "@/components/Button.vue";
</script>

<template>
  <Button @click="$emit('click')">
    <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="M 22.153846,1.8461538 H 1.8461538 V 22.153846 H 22.153846 Z" fill="none" />
      <path d="m 14.769231,11.076923 h -1.846154 v 1.846154 h 1.846154 z" fill="#000000" />
      <path d="m 12.923077,12.923077 h -1.846154 v 1.846154 h 1.846154 z" fill="#000000" />
      <path d="m 16.615385,9.2307692 h -1.846154 v 1.8461538 h 1.846154 z" fill="#000000" />
      <path d="m 18.461538,7.3846154 h -1.846153 v 1.8461538 h 1.846153 z" fill="#000000" />
      <path d="M 11.076923,14.769231 H 9.2307692 v 1.846154 h 1.8461538 z" fill="#000000" />
      <path d="M 7.3846154,11.076923 H 5.5384615 v 1.846154 h 1.8461539 z" fill="#000000" />
      <path d="M 9.2307692,12.923077 H 7.3846154 v 1.846154 h 1.8461538 z" fill="#000000" />
      <path d="M 1.8461538,0 V 1.8461538 H 22.153846 V 0 Z" fill="#000000" />
      <path d="M 24,1.8461538 H 22.153846 V 22.153846 H 24 Z" fill="#000000" />
      <path d="M 1.8461538,22.153846 V 24 H 22.153846 v -1.846154 z" fill="#000000" />
      <path d="M 0,22.153846 H 1.8461538 V 1.8461538 H 0 Z" fill="#000000" />
    </svg>
  </Button>
</template>
