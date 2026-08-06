<script setup>
import { useKlinik } from '../../stores/klinik'
import BottomSheet from '../ui/BottomSheet.vue'

const k = useKlinik()
const opts = [
  { label: 'Tidak tahu', val: null },
  { label: '30 mnt', val: 30 },
  { label: '1 jam', val: 60 },
  { label: '2 jam', val: 120 },
  { label: '3 jam', val: 180 },
]
</script>

<template>
  <BottomSheet :visible="k.modal === 'away'" @close="k.closeModal()">
    <div class="title">Atur status keluar</div>
    <div class="sub">Keterangan & perkiraan ini tampil di halaman web pasien</div>

    <label class="label">Keterangan · opsional</label>
    <input class="field note" v-model="k.awayDraft.note" placeholder="mis. Antar pasien ke RS" />

    <label class="label">Perkiraan kembali</label>
    <div class="chips">
      <button v-for="o in opts" :key="o.label" class="chip" :class="{ sel: k.awayDraft.minutes === o.val }" @click="k.awayDraft.minutes = o.val">
        {{ o.label }}
      </button>
    </div>

    <div class="actions">
      <button class="btn-ghost" @click="k.closeModal()">Batal</button>
      <button class="btn dark grow2" @click="k.confirmAway()">Set Keluar</button>
    </div>
  </BottomSheet>
</template>

<style scoped>
.title { font-size: 20px; font-weight: 800; color: var(--ink); margin-bottom: 3px; }
.sub { font-size: 12.5px; color: var(--muted); margin-bottom: 18px; }
.note { background: var(--fill2); margin: 7px 0 16px; }
.chips { display: flex; gap: 8px; margin: 8px 0 20px; flex-wrap: wrap; }
.chip { padding: 9px 14px; border-radius: 12px; background: var(--fill); color: var(--text-secondary); font-size: 13px; font-weight: 700; }
.chip.sel { background: var(--ink); color: #fff; }
.actions { display: flex; gap: 10px; }
.btn-ghost { flex: 1; padding: 15px; border-radius: 14px; background: var(--fill); color: var(--text-secondary); font-size: 15px; font-weight: 700; }
.btn.dark { background: var(--ink); }
.grow2 { flex: 2; }
</style>
