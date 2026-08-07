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

// Tiga tujuan harian. "Setelan" TIDAK ada di sini: ia dibuka lewat avatar di
// Beranda (HP) atau baris di dasar sidebar (desktop). Sebelumnya ia jadi tab
// kelima, sederajat dengan pekerjaan yang dilakukan bidan tiap hari — padahal
// profil klinik dan ganti sandi disentuh beberapa kali setahun.
//
// "Stok" dulu tab ketiga. Ia dilepas bersama seluruh pengelolaan persediaan:
// yang dicatat klinik ini cukup obat apa yang diberikan ke pasien.
const tabs = [
  { key: 'beranda', label: 'Beranda' },
  { key: 'kunjungan', label: 'Kunjungan' },
  { key: 'rekap', label: 'Rekap' },
]
// ikon garis sederhana (gaya Feather)
const icons = {
  beranda: '<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>',
  kunjungan: '<path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1"/>',
  rekap: '<line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>',
  setelan: '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.6h.09A1.65 1.65 0 0 0 10 3.09V3a2 2 0 1 1 4 0v.09A1.65 1.65 0 0 0 15 4.6a1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9v.09a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>',
}
</script>

<template>
  <nav class="tabbar" aria-label="Navigasi utama">
    <!-- brand + tindakan utama (hanya desktop) -->
    <div class="sidebar-top">
      <div class="brand">
        <div class="logo" aria-hidden="true">{{ inisial }}</div>
        <div class="bwrap"><div class="bname">{{ namaKlinik }}</div></div>
      </div>
      <button class="aksi-desk" @click="k.openVisit()">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" aria-hidden="true"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        Catat Kunjungan
      </button>
    </div>

    <div class="tabs">
      <button v-for="(t, i) in tabs" :key="t.key" class="tab" :class="{ active: k.screen === t.key }"
        :aria-current="k.screen === t.key ? 'page' : undefined" @click="k.goScreen(t.key)">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor"
          stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" v-html="icons[t.key]" />
        <span class="tlabel">{{ t.label }}</span>
      </button>

      <!-- Tindakan utama duduk DI DALAM navigasi, bukan sebagai tombol di dasar
           satu layar. Mencatat kunjungan adalah alasan aplikasi ini ada; dulu ia
           hanya terjangkau dari Beranda. Di HP tombol ini disisipkan di antara
           tab lewat `order`, dan menempati satu petak selebar tab supaya jarak
           antar-ikon tetap rata. -->
      <button class="aksi-hp" aria-label="Catat kunjungan baru" @click="k.openVisit()">
        <span class="aksi-bulat" aria-hidden="true">
          <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </span>
      </button>
    </div>

    <!-- Setelan dipisahkan: bukan pekerjaan harian, jadi tidak sederajat -->
    <button class="tab setelan" :class="{ active: k.screen === 'akun' }"
      :aria-current="k.screen === 'akun' ? 'page' : undefined" @click="k.goScreen('akun')">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"
        stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" v-html="icons.setelan" />
      <span class="tlabel">Setelan</span>
    </button>
  </nav>
</template>

<style scoped>
/* ---- HP: bar bawah (bawaan) ---- */
.tabbar { display: flex; height: 68px; background: #fff; border-top: 1px solid var(--border); padding: 6px 2px; flex-shrink: 0; }
/* Selektor anak-langsung, bukan `.setelan` polos: aturan `.tab { display: flex }`
   di bawah punya kekhususan yang sama dan ditulis belakangan, jadi ia menang dan
   Setelan ikut muncul di bar bawah HP. */
.tabbar > .sidebar-top, .tabbar > .setelan { display: none; }
.tabs { display: flex; flex: 1; align-items: center; }
.tab {
  flex: 1; min-width: 0;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  gap: 3px; padding: 0 2px; color: var(--muted2);
  min-height: 44px;               /* target sentuh 44px — bidan memakai satu tangan */
  border-radius: 12px;
}
.tab.active { color: var(--accent); }
.tlabel { font-size: 10px; font-weight: 700; letter-spacing: -.01em; max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* Tombol tindakan disisipkan setelah dua tab pertama: Beranda, Kunjungan,
   tombol, Rekap. Ia ikut `flex: 1` seperti tab supaya keempat petak sama lebar
   dan jarak antar-ikon rata — dengan lebar tetap, tombolnya berdiri lebih rapat
   ke tab sebelahnya dan barisnya terbaca miring ke satu sisi. */
.tabs .tab:nth-child(1) { order: 1; }
.tabs .tab:nth-child(2) { order: 2; }
.aksi-hp {
  order: 3; flex: 1; min-width: 0;
  display: flex; align-items: center; justify-content: center;
}
.aksi-bulat {
  width: 50px; height: 50px; border-radius: 16px;
  display: flex; align-items: center; justify-content: center;
  background: var(--accent); color: #fff;
  transition: background .15s ease;
}
.aksi-hp:active .aksi-bulat { background: var(--accent-press); }
.tabs .tab:nth-child(3) { order: 4; }

/* ---- desktop: sidebar kiri ---- */
@media (min-width: 920px) {
  .tabbar {
    order: 1;
    flex-direction: column;
    height: 100%;
    width: 250px;
    padding: 20px 14px 16px;
    border-top: none;
    border-right: 1px solid var(--border);
    gap: 18px;
  }
  .tabbar > .sidebar-top { display: block; }
  .brand { display: flex; align-items: center; gap: 12px; padding: 2px 6px 0; }
  .logo {
    width: 36px; height: 36px; border-radius: var(--ra-md); flex-shrink: 0;
    background: var(--accent);
    color: #fff; font-weight: 700; font-size: 13px;
    display: flex; align-items: center; justify-content: center;
  }
  .bwrap { min-width: 0; }
  /* Nama klinik bisa panjang ("Praktik Mandiri Bidan Pit Suryani"); tanpa ini
     ia mendorong lebar sidebar dan menggeser seluruh isi layar. */
  .bname { font-size: 15px; font-weight: 700; letter-spacing: -.015em; color: var(--ink); line-height: 1.25; overflow-wrap: anywhere; }

  .aksi-desk {
    display: flex; align-items: center; justify-content: center; gap: 8px;
    width: 100%; margin-top: 16px; min-height: 40px; padding: 0 12px;
    border-radius: var(--ra-md); font-size: 13.5px; font-weight: 600; color: #fff;
    background: var(--accent);
    transition: background .15s ease;
  }
  .aksi-desk:hover { background: var(--accent-hover); }
  .aksi-hp { display: none; }

  .tabs { flex-direction: column; flex: 1; gap: 3px; align-items: stretch; }
  .tab {
    flex: none; flex-direction: row; justify-content: flex-start; align-items: center;
    gap: 12px; width: 100%; min-height: 38px; padding: 0 12px; border-radius: var(--ra-md);
    color: var(--text-secondary);
  }
  .tab:hover { background: var(--fill2); }
  .tab.active { background: var(--accent-soft); color: var(--accent-ink); }
  .tab svg { width: 18px; height: 18px; }
  .tlabel { font-size: 13.5px; font-weight: 600; }

  /* Setelan menempel di dasar sidebar, dipisah garis: terjangkau, tapi jelas
     bukan bagian dari pekerjaan harian. */
  .tabbar > .setelan {
    display: flex; flex: none; flex-direction: row;
    margin-top: auto; padding-top: 14px;
    border-top: 1px solid var(--border);
    border-radius: 0;
  }
  .setelan .tlabel { font-size: 13.5px; font-weight: 600; }
  .setelan.active { background: none; color: var(--accent-ink); }
}
</style>
