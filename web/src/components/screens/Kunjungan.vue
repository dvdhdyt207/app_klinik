<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
const k = useKlinik()
const d = computed(() => k.derived)
</script>

<template>
  <div class="screen">
    <!-- Tombol "+" mengambang dilepas: mencatat kunjungan kini ada di navigasi,
         jadi tombol kedua di layar ini cuma mengulang tindakan yang sama. -->
    <header class="page-head">
      <h1 class="page-title">Kunjungan</h1>
      <p class="page-sub">{{ d.totalVisits }} tercatat · {{ d.todayCount }} hari ini</p>
    </header>

    <div v-if="d.visitCards.length" class="kartu">
      <ul class="rows">
        <li v-for="v in d.visitCards" :key="v.id" class="row">
          <span class="ava" aria-hidden="true">{{ v.initial }}</span>
          <div class="grow">
            <p class="r-name">
              {{ v.name }}<span class="r-age">{{ v.meta }}</span>
            </p>
            <p v-if="v.gejala" class="r-sub">{{ v.gejala }}</p>
            <!-- Obat ditulis sebagai satu baris teks, bukan deretan chip:
                 chip menambah kotak di dalam kotak, dan yang dicari bidan saat
                 memindai daftar ini adalah namanya, bukan bentuknya. -->
            <p v-if="v.chips.length" class="r-obat">{{ v.chips.join(' · ') }}</p>
          </div>
          <span class="r-side">{{ v.dateLabel }}</span>
        </li>
      </ul>
    </div>
    <div v-else class="kartu">
      <p class="kosong">Belum ada kunjungan tercatat. Ketuk <b>Catat Kunjungan</b> untuk memulai.</p>
    </div>
  </div>
</template>

<style scoped>
.screen { padding: 16px 16px 24px; }
.page-head { margin-bottom: 14px; }
.page-title { margin: 0; font-size: 20px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); }
.page-sub { margin: 3px 0 0; font-size: 13px; color: var(--muted); }

.kartu { background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg); }
.rows { list-style: none; margin: 0; padding: 0; }
.row { display: flex; align-items: flex-start; gap: 12px; padding: 13px 16px; }
.row + .row { border-top: 1px solid var(--hair); }
.grow { flex: 1; min-width: 0; }

.ava {
  width: 32px; height: 32px; border-radius: var(--ra-md); flex-shrink: 0; margin-top: 1px;
  background: var(--fill); color: var(--text-secondary);
  font-size: 12.5px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
}
.r-name { margin: 0; font-size: 14.5px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.r-age { font-size: 12.5px; font-weight: 500; color: var(--muted); margin-left: 8px; }
.r-sub { margin: 2px 0 0; font-size: 13px; color: var(--text-secondary); line-height: 1.45; }
.r-obat { margin: 4px 0 0; font-size: 12.5px; color: var(--muted); line-height: 1.45; }
.r-side { font-size: 12px; color: var(--muted); font-weight: 500; white-space: nowrap; margin-top: 2px; }
.kosong { margin: 0; padding: 26px 18px; text-align: center; font-size: 13.5px; color: var(--muted2); }
.kosong b { color: var(--text-secondary); font-weight: 700; }

@media (min-width: 920px) {
  .screen { padding: 4px 4px 28px; }
  .page-head { margin-bottom: 16px; }
  .page-title { font-size: 23px; }
  /* Daftar rekam medis dibaca dengan memindai ke bawah, jadi ia tetap satu
     kolom penuh — bukan kisi kartu yang memaksa mata melompat kiri-kanan. */
  .row { padding: 14px 20px; gap: 14px; }
}
</style>
