<script setup>
// Header modal layar-penuh: tombol back (‹) + judul + slot kanan.
defineProps({ title: String })
defineEmits(['back'])
</script>

<template>
  <div class="mh">
    <button class="back" @click="$emit('back')">‹</button>
    <span class="title">{{ title }}</span>
    <slot name="right" />
  </div>
</template>

<style scoped>
.mh {
  display: flex; align-items: center; gap: 12px;
  padding: 16px 16px 12px; background: #fff;
  border-bottom: 1px solid var(--border);
}
.back {
  width: 36px; height: 36px; border-radius: 18px; background: var(--fill);
  display: flex; align-items: center; justify-content: center;
  font-size: 20px; color: var(--text-secondary); line-height: 1; padding-bottom: 2px;
}
.title { flex: 1; font-size: 17px; font-weight: 800; color: var(--ink); }
</style>
