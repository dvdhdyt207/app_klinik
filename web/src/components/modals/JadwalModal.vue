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
      <button v-for="e in d.allEvents" :key="e.id" class="erow" @click="k.openEventForm(e)">
        <span class="dot" :style="{ background: e.dot, width: '10px', height: '10px' }" />
        <span class="grow">
          <span class="e-title">{{ e.title }}</span>
          <span class="e-when">{{ e.when }}</span>
        </span>
        <span class="chev">›</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.fullmodal { position: absolute; inset: 0; z-index: 20; background: var(--app-bg); display: flex; flex-direction: column; animation: modal-in .2s ease; }
.add { height: 36px; padding: 0 14px; border-radius: 18px; background: var(--accent); color: #fff; font-size: 13px; font-weight: 700; display: flex; align-items: center; gap: 5px; }
.plus { font-size: 18px; }
.body { flex: 1; overflow-y: auto; padding: 18px; display: flex; flex-direction: column; gap: 10px; }
.erow { display: flex; align-items: center; gap: 13px; text-align: left; width: 100%; background: #fff; border: 1px solid var(--border); border-radius: 14px; padding: 14px 15px; }
.grow { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.e-title { font-size: 14.5px; font-weight: 700; color: var(--ink); }
.e-when { font-size: 12.5px; color: var(--muted); }
.chev { font-size: 18px; color: #c0c8d4; }
</style>
