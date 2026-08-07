<script setup>
import { computed, ref } from 'vue'
import { useKlinik } from '../../stores/klinik'
import { fmtDate } from '../../lib/format'
import ModalHeader from '../ui/ModalHeader.vue'
import Stepper from '../ui/Stepper.vue'

const k = useKlinik()
const dr = computed(() => k.draft)
const canSave = computed(() => dr.value.name.trim() && dr.value.items.length > 0)

function ageInput(e) { k.draft.age = e.target.value.replace(/[^0-9]/g, '') }
function dec(it) { it.qty = Math.max(1, it.qty - 1) }
function inc(it) { it.qty = it.qty + 1 }
function remove(i) { k.draft.items.splice(i, 1) }

// ---- pasien yang pernah datang ----
// Nama tetap diketik bebas — kolomnya kolom biasa, bukan daftar pilihan. Yang
// ditambahkan hanya saran dari kunjungan sebelumnya, supaya pasien yang datang
// lagi tidak perlu diketik ulang (dan tidak tercatat dengan dua ejaan berbeda).
const fokusNama = ref(false)
const saran = computed(() => {
  const q = dr.value.name.trim().toLowerCase()
  if (!q) return []
  return k.derived.pasienDipakai
    .filter((p) => p.name.toLowerCase().includes(q) && p.name.toLowerCase() !== q)
    .slice(0, 5)
    .map((p) => ({ ...p, sub: (p.age ? p.age + ' th · ' : '') + 'terakhir ' + fmtDate(p.lastTs) }))
})
const saranTampil = computed(() => fokusNama.value && saran.value.length > 0)

function pakai(p) { k.pakaiPasien(p); fokusNama.value = false }

// Esc menutup daftar saran lebih dulu, bukan langsung membuang seluruh form.
// Penjaga Esc global ada di BidanApp; propagasinya dihentikan di sini supaya
// tidak dua-duanya jalan sekaligus.
function escNama(e) {
  if (!saranTampil.value) return
  e.stopPropagation()
  fokusNama.value = false
}
</script>

<template>
  <div v-if="k.modal === 'visit'" class="fullmodal">
    <ModalHeader title="Catat Kunjungan" @back="k.closeModal()" />
    <div class="body">
      <label class="label">Nama pasien</label>
      <div class="namawrap mt mb">
        <input class="field" v-model="k.draft.name" placeholder="mis. Siti Aminah"
          autocomplete="off" @focus="fokusNama = true" @blur="fokusNama = false" @keydown.esc="escNama" />
        <!-- mousedown.prevent: tanpa ini kolom nama kehilangan fokus lebih dulu,
             daftarnya keburu tertutup, dan ketukan jatuh ke tempat kosong. -->
        <div v-if="saranTampil" class="saran">
          <button v-for="p in saran" :key="p.name" class="s-row" @mousedown.prevent="pakai(p)">
            <span class="s-name">{{ p.name }}</span>
            <span class="s-sub">{{ p.sub }}</span>
          </button>
        </div>
      </div>

      <label class="label">Umur</label>
      <div class="agerow mt mb">
        <input class="field age" :value="k.draft.age" @input="ageInput" inputmode="numeric" placeholder="0" />
        <span class="suffix">tahun</span>
      </div>

      <label class="label">Gejala / keluhan · opsional</label>
      <textarea class="field area mt mb" v-model="k.draft.gejala" rows="3" placeholder="mis. Demam 2 hari, batuk, pilek" />

      <div class="obat-head">
        <label class="label">Obat yang diberikan</label>
        <button class="link" @click="k.openPickVisit()">+ Tambah obat</button>
      </div>

      <div v-if="dr.items.length" class="list">
        <div v-for="(it, i) in dr.items" :key="i" class="item">
          <div class="grow">
            <div class="it-name">{{ it.name }}</div>
            <div class="it-unit">{{ it.unit }}</div>
          </div>
          <Stepper :value="it.qty" @dec="dec(it)" @inc="inc(it)" />
          <button class="rm" :aria-label="'Hapus ' + it.name" @click="remove(i)">×</button>
        </div>
      </div>
      <div v-else class="dashed">Belum ada obat. Ketuk "Tambah obat".</div>
    </div>
    <div class="footer">
      <button class="btn" :class="{ 'is-disabled': !canSave }" :disabled="!canSave" @click="k.saveVisit()">Simpan Kunjungan</button>
    </div>
  </div>
</template>

<style scoped>
.body { flex: 1; overflow-y: auto; padding: 18px; }
.mt { margin-top: 7px; } .mb { margin-bottom: 16px; }

/* Daftar saran mengambang di atas isi form, bukan mendorongnya ke bawah:
   kolom Umur yang bergeser tiap kali satu huruf diketik membuat form terasa
   goyah, dan yang sedang mengetik nama tidak sedang melihat ke bawah sana. */
.namawrap { position: relative; }
.saran {
  position: absolute; z-index: 5; top: calc(100% + 4px); left: 0; right: 0;
  background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg);
  box-shadow: var(--shadow-md); overflow: hidden;
}
.s-row { display: flex; flex-direction: column; gap: 1px; width: 100%; text-align: left; padding: 9px 13px; }
.s-row + .s-row { border-top: 1px solid var(--hair); }
.s-row:hover { background: var(--fill2); }
.s-name { font-size: 14px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.s-sub { font-size: 12px; color: var(--muted); line-height: 1.35; }
.agerow { display: flex; align-items: center; gap: 8px; }
.age { width: 90px; }
.suffix { font-size: 14px; color: var(--muted); font-weight: 600; }
.area { min-height: 84px; resize: none; }
.obat-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.link { font-size: 12.5px; font-weight: 600; color: var(--accent); min-height: 32px; }
.link:hover { color: var(--accent-hover); }
/* Satu permukaan, baris dipisah garis rambut — bukan tiap obat jadi kartu
   ber-border sendiri di dalam kartu modal. */
.list { background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg); }
.item { display: flex; align-items: center; gap: 10px; padding: 10px 12px; }
.item + .item { border-top: 1px solid var(--hair); }
.grow { flex: 1; min-width: 0; }
.it-name { font-size: 14px; font-weight: 600; color: var(--ink); }
.it-unit { font-size: 11.5px; color: var(--muted); }
.rm { width: 32px; height: 32px; border-radius: var(--ra-sm); color: var(--muted2); font-size: 17px; }
.rm:hover { background: var(--danger-bg); color: var(--danger); }
.footer { padding: 14px 16px; background: #fff; border-top: 1px solid var(--hair); }
</style>
