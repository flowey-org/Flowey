<template>
  <div class="block">
    <slot />
  </div>
</template>

<style>
.block {
  display: flex;
  justify-content: center;
  width: 100%;
  height: var(--timer-font-size);
}
</style>
