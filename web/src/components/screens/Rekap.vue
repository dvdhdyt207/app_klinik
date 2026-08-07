<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
const k = useKlinik()
const d = computed(() => k.derived)
</script>

<template>
  <div class="screen">
    <!-- Dulu layar ini bertab: "Kunjungan" dan "Stok menipis". Rekap stok
         dilepas bersama seluruh pengelolaan persediaan, dan satu tab yang
         berdiri sendiri bukan tab — ia hanya kotak yang tidak bisa diapa-apakan.
         Jadi tabnya ikut hilang, bukan disisakan satu. -->
    <header class="page-head">
      <h1 class="page-title">Rekap</h1>
      <p class="page-sub">Ringkasan kunjungan yang tercatat</p>
    </header>

    <div class="kartu angka">
      <div class="ang">
        <span class="ang-num">{{ d.totalVisits }}</span>
        <span class="ang-lbl">Total kunjungan</span>
      </div>
      <div class="ang">
        <span class="ang-num">{{ d.hariTercatat }}</span>
        <span class="ang-lbl">Hari ada kunjungan</span>
      </div>
    </div>

    <p class="sec-title">Ringkasan per hari</p>
    <div class="kartu">
      <ul v-if="d.dailyRecap.length" class="rows">
        <li v-for="(row, i) in d.dailyRecap" :key="i" class="row">
          <div class="grow">
            <p class="r-name">{{ row.label }}</p>
            <p class="r-sub">{{ row.patients }}</p>
          </div>
          <span class="r-count">{{ row.count }}</span>
        </li>
      </ul>
      <p v-else class="kosong">Belum ada kunjungan untuk diringkas</p>
    </div>
  </div>
</template>

<style scoped>
.screen { padding: 16px 16px 24px; }
.page-head { margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); }
.page-sub { margin: 3px 0 0; font-size: 13px; color: var(--muted); }

.kartu { background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg); }
.angka { display: flex; margin-bottom: 22px; }
.ang { flex: 1; display: flex; flex-direction: column; gap: 2px; padding: 14px 18px; }
.ang + .ang { border-left: 1px solid var(--hair); }
.ang-num { font-size: 23px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); line-height: 1.15; }
.ang-lbl { font-size: 12px; color: var(--muted); }

.sec-title {
  margin: 0 2px 8px; font-size: var(--micro); font-weight: 700;
  letter-spacing: .09em; text-transform: uppercase; color: var(--muted);
}

.rows { list-style: none; margin: 0; padding: 0; }
.row { display: flex; align-items: center; gap: 12px; padding: 12px 16px; }
.row + .row { border-top: 1px solid var(--hair); }
.grow { flex: 1; min-width: 0; }
.r-name { margin: 0; font-size: 14.5px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.r-sub { margin: 1px 0 0; font-size: 12px; color: var(--muted); }
.r-count { font-size: 15px; font-weight: 700; color: var(--ink); white-space: nowrap; }
.kosong { margin: 0; padding: 26px 18px; text-align: center; font-size: 13.5px; color: var(--muted2); }

@media (min-width: 920px) {
  .screen { padding: 4px 4px 28px; max-width: 900px; }
  .page-title { font-size: 23px; }
  .row { padding: 13px 20px; }
  .ang { padding: 16px 22px; }
  .ang-num { font-size: 26px; }
}
</style>
