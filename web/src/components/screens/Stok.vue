<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
const k = useKlinik()
const d = computed(() => k.derived)
</script>

<template>
  <div class="screen">
    <header class="page-head">
      <div>
        <h1 class="page-title">Stok Obat</h1>
        <p class="page-sub">{{ d.medCount }} jenis terdaftar · {{ d.lowCount }} menipis</p>
      </div>
      <button class="aksi" @click="k.openPickStock()">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" aria-hidden="true"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        Tambah stok
      </button>
    </header>

    <div v-if="d.medsList.length" class="kartu">
      <ul class="rows">
        <li v-for="m in d.medsList" :key="m.id" class="row">
          <!-- Barisnya tombol, bukan <li> ber-@click: hanya elemen tombol yang
               bisa dicapai lewat Tab dan ditekan dengan Enter. -->
          <button class="rowbtn" @click="k.openMedEdit(m)">
            <span class="dot" :style="{ background: m.dot }" aria-hidden="true" />
            <span class="grow">
              <span class="r-name">{{ m.name }}</span>
              <span class="r-sub">{{ m.cat }}</span>
            </span>
            <span class="r-qty" :style="{ color: m.qtyColor }">
              {{ m.qty }}<i class="r-unit">{{ m.unit }}</i>
            </span>
            <svg class="chev" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
              stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 6 15 12 9 18"/></svg>
          </button>
        </li>
      </ul>
    </div>
    <div v-else class="kartu">
      <p class="kosong">Belum ada obat terdaftar. Ketuk <b>Tambah stok</b> untuk memulai.</p>
    </div>
  </div>
</template>

<style scoped>
.screen { padding: 16px 16px 24px; }
.page-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 14px; margin-bottom: 14px; }
.page-title { margin: 0; font-size: 20px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); }
.page-sub { margin: 3px 0 0; font-size: 13px; color: var(--muted); }
.aksi {
  display: flex; align-items: center; gap: 7px; flex-shrink: 0;
  min-height: 38px; padding: 0 14px; border-radius: var(--ra-md);
  background: var(--accent); color: #fff; font-size: 13px; font-weight: 600;
  transition: background .15s ease;
}
.aksi:hover { background: var(--accent-hover); }

.kartu { background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg); }
.rows { list-style: none; margin: 0; padding: 0; }
.row + .row { border-top: 1px solid var(--hair); }
/* Padding tinggal di tombol, bukan di <li>: kalau di luar, daerah yang bisa
   diketuk berhenti sebelum tepi baris — dan yang meleset justru sentuhan di
   pinggir, yang paling sering terjadi saat aplikasi dipakai satu tangan. */
.rowbtn {
  display: flex; align-items: center; gap: 12px;
  width: 100%; padding: 12px 16px; text-align: left;
  border-radius: inherit; transition: background .15s ease;
}
.rowbtn:hover { background: var(--fill2); }
.grow { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; margin: 0 1px; }
.r-name { font-size: 14.5px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.r-sub { font-size: 12px; color: var(--muted); line-height: 1.35; }
.chev { color: var(--muted2); flex-shrink: 0; }
.rowbtn:hover .chev { color: var(--accent); }
/* Jaraknya diatur di sini, bukan lewat spasi di dalam template: kompilator
   Vue merapatkan spasi tepi, dan hasilnya terbaca "8butir". */
.r-qty { font-size: 15px; font-weight: 700; white-space: nowrap; }
.r-unit { font-style: normal; font-size: 11px; font-weight: 500; color: var(--muted); margin-left: 5px; }
.kosong { margin: 0; padding: 26px 18px; text-align: center; font-size: 13.5px; color: var(--muted2); }
.kosong b { color: var(--text-secondary); font-weight: 700; }

@media (min-width: 920px) {
  .screen { padding: 4px 4px 28px; }
  .page-head { margin-bottom: 16px; }
  .page-title { font-size: 23px; }
  /* Dua kolom: daftar obat pendek dan sempit isinya, jadi satu kolom penuh
     menyisakan ruang kosong yang lebar di kanan. */
  .rows { display: grid; grid-template-columns: 1fr 1fr; }
  .row + .row { border-top: none; }
  .rows > .row { border-top: 1px solid var(--hair); }
  .rows > .row:nth-child(1), .rows > .row:nth-child(2) { border-top: none; }
  .rows > .row:nth-child(odd) { border-right: 1px solid var(--hair); }
  .rowbtn { padding: 13px 20px; }
}
</style>
