<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import ModalHeader from '../ui/ModalHeader.vue'
import Toggle from '../ui/Toggle.vue'

const k = useKlinik()
const e = computed(() => k.eventDraft)
const canSave = computed(() => e.value && e.value.title.trim() && e.value.startDate)
</script>

<template>
  <div v-if="k.modal === 'event'" class="fullmodal">
    <ModalHeader :title="e && e.id ? 'Ubah Agenda' : 'Agenda Baru'" @back="k.backJadwal()" />
    <div v-if="e" class="body">
      <label class="label">Judul agenda</label>
      <input class="field mt mb" v-model="e.title" placeholder="mis. Keluar kota / Posyandu" />

      <button class="allday" @click="e.allDay = !e.allDay">
        <span class="ad-label">Seharian / beberapa hari</span>
        <Toggle :on="e.allDay" variant="small" />
      </button>

      <label class="label">Mulai</label>
      <div class="drow mt mb">
        <input type="date" class="field grow" v-model="e.startDate" />
        <input v-if="!e.allDay" type="time" class="field time" v-model="e.startTime" />
      </div>

      <label class="label">Selesai</label>
      <div class="drow mt mb">
        <input type="date" class="field grow" v-model="e.endDate" />
        <input v-if="!e.allDay" type="time" class="field time" v-model="e.endTime" />
      </div>

      <button v-if="e.id" class="del" @click="k.deleteEvent()">Hapus agenda</button>
    </div>
    <div class="footer">
      <button class="btn" :class="{ 'is-disabled': !canSave }" :disabled="!canSave" @click="k.saveEvent()">Simpan Agenda</button>
    </div>
  </div>
</template>

<style scoped>
.fullmodal { position: absolute; inset: 0; z-index: 22; background: var(--app-bg); display: flex; flex-direction: column; animation: modal-in .2s ease; }
.body { flex: 1; overflow-y: auto; padding: 18px; }
.mt { margin-top: 7px; } .mb { margin-bottom: 16px; }
.allday { display: flex; align-items: center; justify-content: space-between; width: 100%; background: #fff; border: 1px solid var(--input-border); border-radius: 12px; padding: 13px 14px; margin-bottom: 16px; }
.ad-label { font-size: 14px; font-weight: 600; color: var(--ink); }
.drow { display: flex; gap: 8px; }
.grow { flex: 1; min-width: 0; }
.time { width: 110px; }
.del { width: 100%; padding: 13px; border-radius: 12px; background: var(--del-bg); color: var(--del); font-size: 14px; font-weight: 700; margin-top: 6px; }
.footer { padding: 16px; background: #fff; border-top: 1px solid var(--border); }
</style>
