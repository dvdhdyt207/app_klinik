<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import ModalHeader from '../ui/ModalHeader.vue'

const k = useKlinik()
const d = computed(() => k.derived)
</script>

<template>
  <div v-if="k.modal === 'jadwal'" class="fullmodal">
    <ModalHeader title="Jadwal & Agenda" @back="k.closeModal()">
      <template #right>
        <button class="add" @click="k.openEventNew()"><span class="plus">+</span> Tambah</button>
      </template>
    </ModalHeader>
    <div class="body">
      <div v-if="d.allEvents.length" class="kartu">
        <button v-for="e in d.allEvents" :key="e.id" class="erow" @click="k.openEventForm(e)">
          <span class="dot" :style="{ background: e.dot }" aria-hidden="true" />
          <span class="grow">
            <span class="e-title">{{ e.title }}</span>
            <span class="e-when">{{ e.when }}</span>
          </span>
          <svg class="chev" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>
      </div>
      <p v-else class="kosong">Belum ada agenda. Ketuk <b>Tambah</b> untuk membuat.</p>
    </div>
  </div>
</template>

<style scoped>
.add { min-height: 34px; padding: 0 13px; border-radius: var(--ra-md); background: var(--accent); color: #fff; font-size: 13px; font-weight: 600; display: flex; align-items: center; gap: 5px; }
.add:hover { background: var(--accent-hover); }
.plus { font-size: 16px; }
.body { flex: 1; overflow-y: auto; padding: 16px; }
/* Satu permukaan, baris dipisah garis rambut — sama dengan daftar di layar. */
.kartu { background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg); }
.erow { display: flex; align-items: center; gap: 12px; text-align: left; width: 100%; padding: 12px 14px; }
.erow + .erow { border-top: 1px solid var(--hair); }
.erow:hover { background: var(--fill2); }
.dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; margin: 0 1px; }
.grow { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.e-title { font-size: 14px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.e-when { font-size: 12px; color: var(--muted); line-height: 1.35; }
.chev { color: var(--muted2); flex-shrink: 0; }
.kosong { margin: 0; padding: 26px 18px; text-align: center; font-size: 13.5px; color: var(--muted2); }
.kosong b { color: var(--text-secondary); font-weight: 700; }
</style>
