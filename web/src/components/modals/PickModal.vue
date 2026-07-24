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
      <button v-for="r in results" :key="r.name" class="rrow" @click="k.pick(r)">
        <span class="ava" :style="{ background: r.tint, color: r.ink }">{{ r.letter }}</span>
        <span class="grow">
          <span class="r-name">{{ r.name }}</span>
          <span class="r-sub">{{ r.sub }}</span>
        </span>
        <span class="plus">+</span>
      </button>
      <button v-if="canCreate" class="createrow" @click="k.pick({ name: k.query.trim(), cat: 'Tablet' })">
        <span class="ava" style="background:var(--tablet-tint);color:var(--accent);font-size:20px">+</span>
        <span class="grow">
          <span class="c-name">Tambah "{{ k.query }}" sebagai obat baru</span>
          <span class="r-sub">Obat tablet</span>
        </span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.fullmodal { position: absolute; inset: 0; z-index: 21; background: var(--app-bg); display: flex; flex-direction: column; animation: modal-in .2s ease; }
.searchbar { padding: 14px 16px 10px; background: #fff; }
.search { background: var(--fill2); }
.body { flex: 1; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 7px; }
.rrow, .createrow { display: flex; align-items: center; gap: 12px; text-align: left; width: 100%; background: #fff; border: 1px solid var(--border); border-radius: 12px; padding: 12px 14px; }
.createrow { border: 1.5px dashed #b9cdf0; background: transparent; margin-top: 4px; }
.ava { width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 15px; font-weight: 800; flex-shrink: 0; }
.grow { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.r-name { font-size: 14.5px; font-weight: 600; color: var(--ink); }
.c-name { font-size: 14px; font-weight: 700; color: var(--accent); }
.r-sub { font-size: 12px; color: var(--muted); }
.plus { font-size: 20px; color: var(--accent); }
</style>
