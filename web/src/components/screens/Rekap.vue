<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
const k = useKlinik()
const d = computed(() => k.derived)
const isK = computed(() => k.rekapTab === 'kunjungan')
</script>

<template>
  <div class="screen">
    <h1 class="title">Rekap</h1>
    <div class="seg">
      <button :class="['seg-btn', { active: isK }]" @click="k.rekapTab = 'kunjungan'">Kunjungan</button>
      <button :class="['seg-btn', { active: !isK }]" @click="k.rekapTab = 'stok'">Stok Menipis</button>
    </div>

    <template v-if="isK">
      <div class="stat-row">
        <!-- Penekanan mengikuti arti: jumlah kunjungan itu angka pokok layar
             ini, obat yang keluar adalah turunannya. -->
        <div class="big-card solid">
          <div class="big-num">{{ d.totalVisits }}</div>
          <div class="big-lbl">Total kunjungan</div>
        </div>
        <div class="big-card soft">
          <div class="big-num">{{ d.totalMedsOut }}</div>
          <div class="big-lbl">Obat dikeluarkan</div>
        </div>
      </div>
      <div class="recap-head">Ringkasan per hari</div>
      <div class="list">
        <div v-for="(row, i) in d.dailyRecap" :key="i" class="rline">
          <div class="grow">
            <div class="rlabel">{{ row.label }}</div>
            <div class="rnames">{{ row.patients }}</div>
          </div>
          <div class="rcount">{{ row.count }}</div>
        </div>
      </div>
    </template>

    <template v-else>
      <div class="subtitle">{{ d.lowCount }} obat perlu segera dibeli · urut dari terkecil</div>
      <div class="list">
        <div v-for="m in d.low" :key="m.id" class="lowline">
          <div class="qtybox" :style="{ background: m.tint }">
            <div class="qb-num" :style="{ color: m.color }">{{ m.qty }}</div>
            <div class="qb-unit" :style="{ color: m.color }">{{ m.unit }}</div>
          </div>
          <div class="grow">
            <div class="low-name">{{ m.name }}</div>
            <div class="low-sub">{{ m.cat }} · {{ m.status }}</div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.screen { padding: 20px 18px 24px; }
.title { margin: 0 0 16px; font-size: 22px; font-weight: 800; letter-spacing: -.4px; color: var(--ink); }
.seg { display: flex; background: var(--fill); border-radius: 14px; padding: 4px; margin-bottom: 18px; }
.seg-btn { flex: 1; text-align: center; padding: 9px 0; border-radius: 10px; font-size: 13px; font-weight: 700; color: var(--muted); }
.seg-btn.active { background: #fff; color: var(--ink); box-shadow: var(--shadow-sm); }

.stat-row { display: flex; gap: 10px; margin-bottom: 14px; }
.big-card { flex: 1; border-radius: var(--r-md); padding: 15px; }
/* Dulu kartu kiri hitam pekat dan kanan biru — dua warna yang tidak punya arti
   apa pun, cuma supaya berbeda. Sekarang keduanya teal, bedanya tingkat
   penekanan: yang berisi angka utama lebih pekat. */
.big-card.soft { background: var(--accent-soft); color: var(--accent-ink); }
.big-card.solid { background: var(--accent); color: #fff; }
.big-num { font-size: 26px; font-weight: 800; letter-spacing: -.5px; }
.big-lbl { font-size: 11.5px; opacity: .85; font-weight: 600; }
.recap-head { font-size: 14px; font-weight: 700; color: var(--ink); margin: 8px 2px; }

.list { display: flex; flex-direction: column; gap: 9px; }
.rline { display: flex; align-items: center; gap: 12px; background: #fff; border: 1px solid var(--border); border-radius: 14px; padding: 13px 15px; }
.grow { flex: 1; min-width: 0; }
.rlabel { font-size: 14px; font-weight: 700; color: var(--ink); }
.rnames { font-size: 12px; color: var(--muted); }
.rcount { font-size: 18px; font-weight: 800; color: var(--accent); }

.subtitle { font-size: 13px; color: var(--muted); margin-bottom: 12px; }
.lowline { display: flex; align-items: center; gap: 14px; background: #fff; border: 1px solid var(--border); border-radius: 14px; padding: 13px 15px; }
.qtybox { min-width: 52px; height: 52px; border-radius: 12px; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.qb-num { font-size: 20px; font-weight: 800; line-height: 1; }
.qb-unit { font-size: 10px; font-weight: 600; }
.low-name { font-size: 14.5px; font-weight: 700; color: var(--ink); }
.low-sub { font-size: 12px; color: var(--muted); }

@media (min-width: 920px) {
  .screen { padding: 8px 4px 24px; }
  .title { font-size: 26px; }
  /* Pemilih tab dan kartu angka tidak ikut melebar sampai ujung layar —
     tombol selebar 500px tidak lebih mudah ditekan, hanya lebih aneh. */
  .seg, .stat-row { max-width: 620px; }
  .list { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 10px; align-items: start; }
}
</style>
