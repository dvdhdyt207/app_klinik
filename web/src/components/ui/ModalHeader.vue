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
  padding: 14px 16px; background: #fff;
  border-bottom: 1px solid var(--hair);
}
.back {
  width: 32px; height: 32px; border-radius: var(--ra-md); background: var(--fill);
  display: flex; align-items: center; justify-content: center;
  font-size: 19px; color: var(--text-secondary); line-height: 1; padding-bottom: 2px;
  transition: background .15s ease;
}
.back:hover { background: var(--accent-soft); color: var(--accent-ink); }
.title { flex: 1; font-size: 15px; font-weight: 700; letter-spacing: -.01em; color: var(--ink); }
</style>
