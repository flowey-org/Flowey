<script lang="ts" setup>
import Counter from "@/components/Counter.vue";
</script>

<template>
  <Counter />
</template>

<style>
@font-face {
  font-family: "Minecraft";
  font-style: normal;
  font-weight: 400;
  font-display: block;
  src: url("/assets/fonts/Minecraft.woff2") format("woff2");
}

html, body {
  margin: 0;
  padding: 0;
}

#app {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100vh;
  width: 100vw;
}

* {
  font-family: Minecraft, serif;
}

:root {
  --timer-font-size: clamp(1em, 1em + 10vw, 6em);
}
</style>
