<template>
  <div id="box">
    <slot />
  </div>
</template>

<style>
#box {
  width: fit-content;
  height: fit-content;
  margin: auto;
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
