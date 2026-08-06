<script setup>
import { computed } from 'vue'
import { useKlinik } from '../stores/klinik'
import { monogram } from '../lib/format'
const k = useKlinik()
// Nama klinik dulu ditulis mati "Bidan Pit" di sini. Bidan bisa menggantinya
// lewat Setelan → Profil klinik, dan sidebar yang tidak ikut berubah akan
// menyebut nama yang salah tepat di sebelah nama yang benar.
const namaKlinik = computed(() => k.profil.clinic || k.status.clinic || 'Klinik')
const inisial = computed(() => monogram(namaKlinik.value))
const tabs = [
  { key: 'beranda', label: 'Beranda' },
  { key: 'kunjungan', label: 'Kunjungan' },
  { key: 'stok', label: 'Stok' },
  { key: 'rekap', label: 'Rekap' },
  { key: 'akun', label: 'Setelan' },
]
// ikon garis sederhana (gaya Feather)
const icons = {
  beranda: '<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>',
  kunjungan: '<path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1"/>',
  stok: '<path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/>',
  rekap: '<line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>',
  akun: '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.6h.09A1.65 1.65 0 0 0 10 3.09V3a2 2 0 1 1 4 0v.09A1.65 1.65 0 0 0 15 4.6a1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9v.09a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>',
}
</script>

<template>
  <nav class="tabbar">
    <!-- brand (hanya desktop) -->
    <div class="brand">
      <div class="logo">{{ inisial }}</div>
      <!-- Tanpa label "KLINIK" di atasnya: nama yang diketik bidan hampir
           selalu sudah diawali kata itu ("Klinik Bidan Pit"), sehingga
           labelnya cuma mengulang kata pertama tepat di bawahnya. -->
      <div class="bwrap">
        <div class="bname">{{ namaKlinik }}</div>
      </div>
    </div>

    <div class="tabs">
      <button v-for="t in tabs" :key="t.key" class="tab" :class="{ active: k.screen === t.key }" @click="k.goScreen(t.key)">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor"
          stroke-width="2" stroke-linecap="round" stroke-linejoin="round" v-html="icons[t.key]" />
        <span class="tlabel">{{ t.label }}</span>
      </button>
    </div>
  </nav>
</template>

<style scoped>
/* ---- mobile: bottom tab bar (default) ---- */
.tabbar { display: flex; height: 66px; background: #fff; border-top: 1px solid var(--border); padding: 8px 6px 6px; flex-shrink: 0; }
.brand { display: none; }
.tabs { display: flex; flex: 1; }
.tab { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px; color: var(--muted2); }
.tab.active { color: var(--accent); }
.tlabel { font-size: 10.5px; font-weight: 700; }

/* ---- desktop: sidebar kiri ---- */
@media (min-width: 920px) {
  .tabbar {
    order: 1;
    flex-direction: column;
    height: 100%;
    width: 244px;
    padding: 22px 16px;
    border-top: none;
    border-right: 1px solid var(--border);
    gap: 20px;
  }
  .brand { display: flex; align-items: center; gap: 12px; padding: 4px 8px 6px; }
  .logo {
    width: 42px; height: 42px; border-radius: 13px; flex-shrink: 0;
    background: linear-gradient(145deg, var(--accent) 0%, var(--accent-press) 100%);
    color: #fff; font-weight: 800; font-size: 15px;
    display: flex; align-items: center; justify-content: center;
    box-shadow: 0 5px 14px rgba(18,128,106,.26);
  }
  .bwrap { min-width: 0; }
  /* Nama klinik bisa panjang ("Praktik Mandiri Bidan Pit Suryani"); tanpa ini
     ia mendorong lebar sidebar dan menggeser seluruh isi layar. */
  .bname { font-size: 17px; font-weight: 800; letter-spacing: -.02em; color: var(--ink); line-height: 1.2; overflow-wrap: anywhere; }

  .tabs { flex-direction: column; flex: none; gap: 4px; }
  .tab {
    flex-direction: row; justify-content: flex-start; align-items: center;
    gap: 13px; width: 100%; padding: 12px 14px; border-radius: 12px;
    color: var(--label);
  }
  .tab:hover { background: var(--fill2); }
  .tab.active { background: var(--accent-soft); color: var(--accent-ink); }
  .tab svg { width: 20px; height: 20px; }
  .tlabel { font-size: 14.5px; font-weight: 700; }
}
</style>
