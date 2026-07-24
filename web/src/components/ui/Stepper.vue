<script setup>
defineProps({ value: { type: Number, default: 0 } })
defineEmits(['dec', 'inc'])
</script>

<template>
  <div class="stepper">
    <button class="sbtn" @click="$emit('dec')">−</button>
    <span class="val">{{ value }}</span>
    <button class="sbtn" @click="$emit('inc')">+</button>
  </div>
</template>

<style scoped>
.stepper { display: flex; align-items: center; border: 1px solid var(--border2); border-radius: 10px; overflow: hidden; }
.sbtn { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; font-size: 18px; color: #3b4a5e; background: var(--fill2); }
.val { width: 34px; text-align: center; font-size: 14px; font-weight: 700; }
</style>
