<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
const k = useKlinik()
const d = computed(() => k.derived)
const isK = computed(() => k.rekapTab === 'kunjungan')
</script>

<template>
  <div class="screen">
    <header class="page-head">
      <h1 class="page-title">Rekap</h1>
      <!-- Tab bergaris bawah, bukan tombol pil di dalam kotak abu: satu kotak
           lebih sedikit, dan yang aktif tetap terbaca jelas. -->
      <div class="tabs" role="tablist">
        <button class="tb" :class="{ on: isK }" role="tab" :aria-selected="isK" @click="k.rekapTab = 'kunjungan'">Kunjungan</button>
        <button class="tb" :class="{ on: !isK }" role="tab" :aria-selected="!isK" @click="k.rekapTab = 'stok'">Stok menipis</button>
      </div>
    </header>

    <template v-if="isK">
      <div class="kartu angka">
        <div class="ang">
          <span class="ang-num">{{ d.totalVisits }}</span>
          <span class="ang-lbl">Total kunjungan</span>
        </div>
        <div class="ang">
          <span class="ang-num">{{ d.totalMedsOut }}</span>
          <span class="ang-lbl">Obat dikeluarkan</span>
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
    </template>

    <template v-else>
      <p class="sec-title">{{ d.lowCount }} obat perlu diperhatikan · urut dari terkecil</p>
      <div class="kartu">
        <ul v-if="d.low.length" class="rows">
          <li v-for="m in d.low" :key="m.id" class="row">
            <span class="dot" :style="{ background: m.color }" aria-hidden="true" />
            <div class="grow">
              <p class="r-name">{{ m.name }}</p>
              <p class="r-sub">{{ m.cat }} · {{ m.status }}</p>
            </div>
            <span class="r-qty" :style="{ color: m.color }">{{ m.qty }}<i class="r-unit">{{ m.unit }}</i></span>
          </li>
        </ul>
        <p v-else class="kosong">Semua stok aman</p>
      </div>
    </template>
  </div>
</template>

<style scoped>
.screen { padding: 16px 16px 24px; }
.page-head { margin-bottom: 16px; }
.page-title { margin: 0 0 12px; font-size: 20px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); }

.tabs { display: flex; gap: 20px; border-bottom: 1px solid var(--line); }
.tb {
  padding: 0 0 9px; font-size: 13.5px; font-weight: 600; color: var(--muted);
  border-bottom: 2px solid transparent; margin-bottom: -1px; min-height: 34px;
}
.tb:hover { color: var(--text-secondary); }
.tb.on { color: var(--accent-ink); border-bottom-color: var(--accent); }

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
.dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; margin: 0 1px; }
.r-name { margin: 0; font-size: 14.5px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.r-sub { margin: 1px 0 0; font-size: 12px; color: var(--muted); }
.r-count { font-size: 15px; font-weight: 700; color: var(--ink); white-space: nowrap; }
.r-qty { font-size: 15px; font-weight: 700; white-space: nowrap; }
.r-unit { font-style: normal; font-size: 11px; font-weight: 500; color: var(--muted); margin-left: 5px; }
.kosong { margin: 0; padding: 26px 18px; text-align: center; font-size: 13.5px; color: var(--muted2); }

@media (min-width: 920px) {
  .screen { padding: 4px 4px 28px; max-width: 900px; }
  .page-title { font-size: 23px; }
  .row { padding: 13px 20px; }
  .ang { padding: 16px 22px; }
  .ang-num { font-size: 26px; }
}
</style>
