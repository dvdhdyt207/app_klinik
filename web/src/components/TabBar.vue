<script setup>
import { useKlinik } from '../stores/klinik'
const k = useKlinik()
const tabs = [
  { key: 'beranda', label: 'Beranda' },
  { key: 'kunjungan', label: 'Kunjungan' },
  { key: 'stok', label: 'Stok' },
  { key: 'rekap', label: 'Rekap' },
]
// ikon garis sederhana (gaya Feather)
const icons = {
  beranda: '<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>',
  kunjungan: '<path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1"/>',
  stok: '<path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/>',
  rekap: '<line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>',
}
</script>

<template>
  <nav class="tabbar">
    <button v-for="t in tabs" :key="t.key" class="tab" :class="{ active: k.screen === t.key }" @click="k.goScreen(t.key)">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor"
        stroke-width="2" stroke-linecap="round" stroke-linejoin="round" v-html="icons[t.key]" />
      <span class="tlabel">{{ t.label }}</span>
    </button>
  </nav>
</template>

<style scoped>
.tabbar { display: flex; height: 66px; background: #fff; border-top: 1px solid var(--border); padding: 8px 6px 6px; flex-shrink: 0; }
.tab { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px; color: var(--muted2); }
.tab.active { color: var(--accent); }
.tlabel { font-size: 10.5px; font-weight: 700; }
</style>
