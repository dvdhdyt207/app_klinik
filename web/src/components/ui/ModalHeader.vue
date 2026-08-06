<script setup>
// Header modal layar-penuh: tombol back (‹) + judul + slot kanan.
defineProps({ title: String })
defineEmits(['back'])
</script>

<template>
  <div class="mh">
    <!-- Ikon SVG, bukan karakter "‹": tinggi dan garis dasar glif mengikuti
         font, jadi ia tidak pernah benar-benar di tengah kotaknya dan dulu
         diakali dengan padding-bottom. -->
    <button class="back" aria-label="Kembali" @click="$emit('back')">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor"
        stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <polyline points="15 18 9 12 15 6" />
      </svg>
    </button>
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
  color: var(--text-secondary); flex-shrink: 0;
  transition: background .15s ease;
}
.back:hover { background: var(--accent-soft); color: var(--accent-ink); }
.title { flex: 1; font-size: 15px; font-weight: 700; letter-spacing: -.01em; color: var(--ink); }
</style>
