<template>
  <div class="text">
    <slot />
  </div>
</template>

<style>
.text {
  font-size: calc(var(--timer-font-size) / 3);
  display: flex;
  align-items: center;
}
</style>
