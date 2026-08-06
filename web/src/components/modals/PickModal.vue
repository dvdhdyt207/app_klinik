<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import { CATALOG, baseUnit, catTint } from '../../lib/catalog'
import ModalHeader from '../ui/ModalHeader.vue'

const k = useKlinik()
const q = computed(() => k.query.trim().toLowerCase())
const meds = computed(() => k.data.medicines || [])

const results = computed(() => CATALOG.filter((c) => c.name.toLowerCase().includes(q.value)).map((c) => {
  const ex = meds.value.find((m) => m.name.toLowerCase() === c.name.toLowerCase())
  const tint = catTint(c.cat)
  return { ...c, sub: c.cat + (ex ? ' · stok ' + ex.qty + ' ' + baseUnit(c.cat) : ' · belum ada stok'),
    letter: c.name.charAt(0).toUpperCase(), tint: tint.bg, ink: tint.ink }
}))
const exact = computed(() => CATALOG.some((c) => c.name.toLowerCase() === q.value) || meds.value.some((m) => m.name.toLowerCase() === q.value))
const canCreate = computed(() => k.pickContext === 'stock' && q.value.length >= 3 && !exact.value)
</script>

<template>
  <div v-if="k.modal === 'pick'" class="fullmodal">
    <ModalHeader title="Cari Obat" @back="k.backFromPick()" />
    <div class="searchbar">
      <input class="field search" v-model="k.query" placeholder="Ketik nama obat…" autofocus />
    </div>
    <div class="body">
      <div v-if="results.length" class="kartu">
        <button v-for="r in results" :key="r.name" class="rrow" @click="k.pick(r)">
          <span class="ava" :style="{ background: r.tint, color: r.ink }" aria-hidden="true">{{ r.letter }}</span>
          <span class="grow">
            <span class="r-name">{{ r.name }}</span>
            <span class="r-sub">{{ r.sub }}</span>
          </span>
          <svg class="plus" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="2.4" stroke-linecap="round" aria-hidden="true"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
      </div>
      <p v-else-if="!canCreate" class="kosong">Tidak ada obat yang cocok.</p>

      <button v-if="canCreate" class="createrow" @click="k.pick({ name: k.query.trim(), cat: 'Tablet' })">
        <span class="ava is-new" aria-hidden="true">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </span>
        <span class="grow">
          <span class="c-name">Tambah “{{ k.query }}” sebagai obat baru</span>
          <span class="r-sub">Obat tablet</span>
        </span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.searchbar { padding: 12px 16px; background: #fff; border-bottom: 1px solid var(--hair); }
.body { flex: 1; overflow-y: auto; padding: 16px; }
/* Satu permukaan, baris dipisah garis rambut — daftar ini bisa panjang, dan
   kartu per baris membuatnya terbaca seperti tumpukan kotak. */
.kartu { background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg); }
.rrow { display: flex; align-items: center; gap: 12px; text-align: left; width: 100%; padding: 10px 14px; }
.rrow + .rrow { border-top: 1px solid var(--hair); }
.rrow:hover { background: var(--fill2); }
.createrow {
  display: flex; align-items: center; gap: 12px; text-align: left; width: 100%;
  margin-top: 10px; padding: 12px 14px;
  border: 1px dashed var(--line); border-radius: var(--ra-lg); background: transparent;
}
.createrow:hover { border-color: var(--accent); background: var(--fill2); }
.ava {
  width: 32px; height: 32px; border-radius: var(--ra-md); flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  font-size: 13px; font-weight: 700;
}
.ava.is-new { background: var(--accent-soft); color: var(--accent-ink); }
.grow { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.r-name { font-size: 14px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.c-name { font-size: 13.5px; font-weight: 600; color: var(--accent-ink); line-height: 1.35; }
.r-sub { font-size: 12px; color: var(--muted); line-height: 1.35; }
.plus { color: var(--muted2); flex-shrink: 0; }
.rrow:hover .plus { color: var(--accent); }
.kosong { margin: 0; padding: 26px 18px; text-align: center; font-size: 13.5px; color: var(--muted2); }
</style>
