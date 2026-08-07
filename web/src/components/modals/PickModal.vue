<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import { SATUAN } from '../../lib/obat'
import { fmtDate } from '../../lib/format'
import ModalHeader from '../ui/ModalHeader.vue'

const k = useKlinik()
const q = computed(() => k.query.trim().toLowerCase())

// Daftarnya adalah obat yang PERNAH DIBERIKAN pada kunjungan sebelumnya —
// bukan daftar stok, dan bukan katalog tertulis di kode. Tidak ada lagi yang
// perlu didata lebih dulu: obat yang sekali dicatat langsung bisa dipakai ulang
// pada kunjungan berikutnya, lengkap dengan satuan terakhirnya.
//
// Query kosong berarti seluruhnya tampil, terbaru dulu — membuka layar ini
// tanpa mengetik apa pun sudah memperlihatkan obat yang paling sering dipakai
// belakangan.
const results = computed(() => k.derived.obatDipakai
  .filter((m) => m.name.toLowerCase().includes(q.value))
  .map((m) => ({
    name: m.name, unit: m.unit,
    sub: m.unit + ' · terakhir ' + fmtDate(m.lastTs),
    letter: m.name.charAt(0).toUpperCase(),
  })))

const exact = computed(() => k.derived.obatDipakai.some((m) => m.name.toLowerCase() === q.value))
const bisaBaru = computed(() => q.value.length >= 2 && !exact.value)
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
          <span class="ava" aria-hidden="true">{{ r.letter }}</span>
          <span class="grow">
            <span class="r-name">{{ r.name }}</span>
            <span class="r-sub">{{ r.sub }}</span>
          </span>
          <svg class="plus" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
            stroke-width="2.4" stroke-linecap="round" aria-hidden="true"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
      </div>
      <p v-else-if="!bisaBaru" class="kosong">
        {{ k.derived.obatDipakai.length
          ? 'Tidak ada obat yang cocok. Ketik namanya untuk menambahkan yang baru.'
          : 'Belum ada obat yang pernah dicatat. Ketik nama obatnya untuk menambahkan.' }}
      </p>

      <!-- Obat baru: nama diketik sendiri, satuannya dipilih di sini juga.
           Satuan ditanyakan di depan karena sesudah ini tidak ada layar mana pun
           untuk memperbaikinya — daftar obat tidak lagi berdiri sendiri, jadi
           yang salah satuan akan terbawa ke rekam medis kunjungan ini. -->
      <div v-if="bisaBaru" class="baru">
        <p class="b-judul">Belum pernah dicatat</p>
        <div class="b-satuan" role="group" aria-label="Satuan obat">
          <button v-for="s in SATUAN" :key="s" class="chip" :class="{ on: k.pickSatuan === s }"
            @click="k.pickSatuan = s">{{ s }}</button>
        </div>
        <button class="b-aksi" @click="k.pick({ name: k.query.trim(), unit: k.pickSatuan })">
          <span class="ava is-new" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </span>
          <span class="grow">
            <span class="c-name">Tambahkan “{{ k.query.trim() }}”</span>
            <span class="r-sub">Dihitung dalam {{ k.pickSatuan }}</span>
          </span>
        </button>
      </div>
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

.baru { margin-top: 14px; padding: 12px 14px 14px; border: 1px dashed var(--line); border-radius: var(--ra-lg); }
.b-judul {
  margin: 0 0 9px; font-size: var(--micro); font-weight: 700;
  letter-spacing: .09em; text-transform: uppercase; color: var(--muted);
}
.b-satuan { display: flex; gap: 8px; margin-bottom: 12px; }
.chip {
  min-height: 34px; padding: 0 14px; border-radius: var(--ra-md);
  border: 1px solid var(--input-border); background: var(--fill2);
  font-size: 13px; font-weight: 600; color: var(--text-secondary);
}
.chip:hover { color: var(--ink); }
.chip.on { border-color: var(--accent); background: var(--accent-soft); color: var(--accent-ink); }
.b-aksi { display: flex; align-items: center; gap: 12px; text-align: left; width: 100%; }

/* Avatar seragam. Dulu warnanya mengikuti kategori obat (Tablet biru, Sirup
   teal, Sachet ungu); kategori sudah tidak ada, dan mewarnai per satuan hanya
   akan jadi hiasan — di aplikasi ini warna dipakai sebagai isyarat. */
.ava {
  width: 32px; height: 32px; border-radius: var(--ra-md); flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  background: var(--fill); color: var(--text-secondary);
  font-size: 13px; font-weight: 700;
}
.ava.is-new { background: var(--accent-soft); color: var(--accent-ink); }
.grow { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.r-name { font-size: 14px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.c-name { font-size: 13.5px; font-weight: 600; color: var(--accent-ink); line-height: 1.35; }
.r-sub { font-size: 12px; color: var(--muted); line-height: 1.35; }
.plus { color: var(--muted2); flex-shrink: 0; }
.rrow:hover .plus { color: var(--accent); }
.kosong { margin: 0; padding: 26px 18px; text-align: center; font-size: 13.5px; color: var(--muted2); line-height: 1.5; }
</style>
