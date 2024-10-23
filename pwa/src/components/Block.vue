<template>
  <div class="block">
    <slot />
  </div>
</template>

<style>
.block {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: var(--timer-font-size);
  gap: calc(var(--timer-font-size) / 2.25);
}
</style>
