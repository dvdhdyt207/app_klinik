<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import { hhmm, untilLabel, monogram, tanggalPanjang } from '../../lib/format'
import Toggle from '../ui/Toggle.vue'

const k = useKlinik()

const namaKlinik = computed(() => k.profil.clinic || k.status.clinic || 'Klinik')
const inisial = computed(() => monogram(namaKlinik.value))
const st = computed(() => k.status)
const away = computed(() => !st.value.bidanHadir)
const d = computed(() => k.derived)
const untilText = computed(() => { void k.now; return st.value.awayUntil ? untilLabel(st.value.awayUntil) : '' })
const hariIni = computed(() => { void k.now; return tanggalPanjang(Date.now()) })
</script>

<template>
  <div class="screen">
    <header class="b-head">
      <div class="head">
        <div class="head-id"><p class="clinic">{{ namaKlinik }}</p></div>
        <div class="head-desk">
          <h1 class="page-title">Beranda</h1>
          <p class="page-date">{{ hariIni }}</p>
        </div>
        <button class="avatar" type="button" title="Setelan" @click="k.goScreen('akun')">{{ inisial }}</button>
      </div>
    </header>

    <!-- ===== status ===== -->
    <section class="b-status">
      <div class="kartu status" :class="away ? 'is-away' : 'is-here'">
        <div class="s-row">
          <p class="s-title">
            <span class="s-dot" aria-hidden="true" />{{ away ? 'Tidak di klinik' : 'Sedang di klinik' }}
          </p>
          <button v-if="!away" class="s-toggle" aria-label="Atur status keluar" @click="k.openAway()">
            <Toggle :on="true" variant="small" />
          </button>
          <button v-else class="s-toggle" aria-label="Tandai sudah kembali" @click="k.setHadir()">
            <Toggle :on="false" variant="small" />
          </button>
        </div>

        <template v-if="away">
          <p v-if="st.awayNote" class="s-note">{{ st.awayNote }}</p>
          <p class="s-meta">
            <template v-if="st.awayUntil">Perkiraan kembali <b>{{ hhmm(st.awayUntil) }}</b> · {{ untilText }}</template>
            <template v-else>Waktu kembali belum dipastikan</template>
          </p>
          <div class="s-actions">
            <button class="btn-ghost" @click="k.extend(15)">+15 mnt</button>
            <button class="btn-ghost" @click="k.extend(30)">+30 mnt</button>
            <button class="btn-ghost" @click="k.openAway()">Ubah</button>
            <button class="btn-solid" @click="k.setHadir()">Saya sudah kembali</button>
          </div>
        </template>
        <p v-else class="s-meta">Pasien melihat status ini terbuka. Ketuk sakelar untuk mengatur status keluar.</p>

        <router-link to="/" class="s-web">
          Lihat halaman pasien
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M7 17L17 7M17 7H8M17 7v9"/></svg>
        </router-link>
      </div>
    </section>

    <!-- ===== ringkasan: kartu tersendiri, tidak menyatu dgn status =====
         Baris "Obat menipis" dilepas bersama pengelolaan stok. Penggantinya
         bukan angka obat yang lain, melainkan angka kunjungan — itu satu-satunya
         hitungan yang masih dipegang aplikasi ini. -->
    <section class="b-angka">
      <div class="kartu ringkas">
        <p class="rk-label">Ringkasan</p>
        <div class="rk-rows">
          <button class="rk-row is-link" @click="k.goScreen('kunjungan')">
            <span class="rk-num">{{ d.todayCount }}</span>
            <span class="rk-lbl">Kunjungan hari ini</span>
            <svg class="rk-chev" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"/></svg>
          </button>
          <button class="rk-row is-link" @click="k.goScreen('rekap')">
            <span class="rk-num">{{ d.totalVisits }}</span>
            <span class="rk-lbl">Total kunjungan</span>
            <svg class="rk-chev" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"/></svg>
          </button>
        </div>
      </div>
    </section>

    <!-- ===== dua daftar sederajat ===== -->
    <section class="b-recent">
      <div class="sec-head">
        <h2 class="sec-title">Kunjungan terakhir</h2>
        <button class="sec-link" @click="k.goScreen('kunjungan')">Semua</button>
      </div>
      <div class="kartu isi">
        <ul v-if="d.visitCards.length" class="rows">
          <li v-for="v in d.visitCards.slice(0, 4)" :key="v.id" class="row">
            <span class="ava" aria-hidden="true">{{ v.initial }}</span>
            <span class="grow">
              <span class="r-name">{{ v.name }}</span>
              <span class="r-sub">{{ v.meta }}</span>
            </span>
            <span class="r-side">{{ v.dateLabel }}</span>
          </li>
        </ul>
        <p v-else class="kosong">Belum ada kunjungan tercatat</p>
      </div>
    </section>

    <section class="b-sched">
      <div class="sec-head">
        <h2 class="sec-title">Jadwal &amp; agenda</h2>
        <button class="sec-link" @click="k.openJadwal()">Kelola</button>
      </div>
      <div class="kartu isi">
        <ul v-if="d.upcoming.length" class="rows">
          <li v-for="e in d.upcoming.slice(0, 4)" :key="e.id" class="row">
            <span class="dot" :style="{ background: e.dot }" aria-hidden="true" />
            <span class="grow">
              <span class="r-name">{{ e.title }}</span>
              <span class="r-sub">{{ e.when }}</span>
            </span>
          </li>
        </ul>
        <p v-else class="kosong">Belum ada agenda mendatang</p>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* ================= HP (bawaan) ================= */
.screen { padding: 16px 16px 24px; display: flex; flex-direction: column; gap: 20px; }

.head { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.head-id { min-width: 0; }
.head-desk { display: none; }
.clinic { font-size: 19px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); margin: 0; line-height: 1.25; overflow-wrap: anywhere; }
.avatar {
  width: 40px; height: 40px; border-radius: 20px; flex-shrink: 0;
  background: var(--fill); color: var(--text-secondary);
  font-size: 13px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
}
.avatar:hover { background: var(--accent-soft); color: var(--accent-ink); }

.kartu { background: var(--card); border: 1px solid var(--line); border-radius: var(--ra-lg); }

/* ---- status ---- */
.status { padding: 16px 18px 14px; }
.s-row { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.s-title {
  display: inline-flex; align-items: center; gap: 9px; margin: 0;
  font-size: 18px; font-weight: 700; letter-spacing: -.015em; color: var(--ink); line-height: 1.35;
}
.s-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.is-here .s-dot { background: var(--accent); }
.is-away .s-dot { background: var(--away); }
.s-toggle { display: flex; align-items: center; min-height: 40px; }

.s-note { font-size: 14px; font-weight: 500; color: var(--text-secondary); margin: 10px 0 0; line-height: 1.5; }
.s-meta { font-size: 13px; color: var(--muted); margin: 6px 0 0; line-height: 1.5; }
.s-meta b { color: var(--ink); font-weight: 700; }

.s-actions { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 14px; }
.btn-ghost {
  min-height: 38px; padding: 0 14px; border-radius: var(--ra-md);
  background: var(--fill); color: var(--text-secondary);
  font-size: 13px; font-weight: 600;
}
.btn-ghost:hover { color: var(--ink); }
.btn-solid {
  flex: 1; min-width: 150px; min-height: 38px; padding: 0 16px; border-radius: var(--ra-md);
  background: var(--accent); color: #fff; font-size: 13px; font-weight: 600;
}
.btn-solid:hover { background: var(--accent-hover); }

.s-web {
  display: inline-flex; align-items: center; gap: 5px; margin-top: 14px;
  font-size: 12.5px; font-weight: 600; color: var(--muted);
}
.s-web:hover { color: var(--accent); }

/* ---- ringkasan hari ini (kartu sendiri) ---- */
.ringkas { padding: 14px 18px 6px; display: flex; flex-direction: column; }
.rk-label {
  margin: 0 0 4px; font-size: var(--micro); font-weight: 700;
  letter-spacing: .09em; text-transform: uppercase; color: var(--muted);
}
.rk-rows { display: flex; flex-direction: column; flex: 1; }
.rk-row {
  display: flex; align-items: center; gap: 12px; padding: 12px 0; text-align: left;
  border-radius: var(--ra-sm);
}
.rk-row + .rk-row { border-top: 1px solid var(--hair); }
.rk-num { font-size: 22px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); line-height: 1; min-width: 30px; }
.rk-lbl { flex: 1; font-size: 13px; color: var(--text-secondary); }
.rk-chev { color: var(--muted2); flex-shrink: 0; }
.rk-row.is-link:hover .rk-lbl { color: var(--ink); }
.rk-row.is-link:hover .rk-chev { color: var(--accent); }

/* ---- bagian daftar ---- */
.sec-head { display: flex; align-items: baseline; justify-content: space-between; gap: 12px; margin: 0 2px 8px; }
.sec-title {
  margin: 0; font-size: var(--micro); font-weight: 700;
  letter-spacing: .09em; text-transform: uppercase; color: var(--muted);
}
.sec-link { font-size: 12.5px; font-weight: 600; color: var(--muted); min-height: 30px; }
.sec-link:hover { color: var(--accent); }

.rows { list-style: none; margin: 0; padding: 0; }
.row { display: flex; align-items: center; gap: 12px; padding: 11px 16px; }
.row + .row { border-top: 1px solid var(--hair); }
.grow { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.r-name { font-size: 14px; font-weight: 600; color: var(--ink); line-height: 1.35; }
.r-sub { font-size: 12px; color: var(--muted); line-height: 1.35; }
.r-side { font-size: 12px; color: var(--muted); font-weight: 500; white-space: nowrap; }
.ava {
  width: 30px; height: 30px; border-radius: var(--ra-md); flex-shrink: 0;
  background: var(--fill); color: var(--text-secondary);
  font-size: 12.5px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
}
.dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; margin: 0 1px; }
.kosong { font-size: 13px; color: var(--muted2); margin: 0; padding: 22px 16px; text-align: center; }

/* ================= DESKTOP ================= */
@media (min-width: 920px) {
  .screen {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 20px;
    padding: 4px 4px 28px;
  }
  /* Baris ditulis eksplisit supaya penempatannya pasti — properti `order` dan
     penempatan otomatis pernah melempar satu kartu ke baris sendiri. */
  .b-head   { grid-column: 1 / -1; grid-row: 1; }
  .b-status { grid-column: 1 / 3;  grid-row: 2; }   /* dua kolom */
  .b-angka  { grid-column: 3;      grid-row: 2; }   /* kartu terpisah, ada jarak */
  /* Kunjungan terakhir mengambil dua kolom yang ditinggalkan daftar stok
     menipis: barisnya berisi nama, umur, dan tanggal sekaligus, dan di satu
     kolom sempit ketiganya mulai berdesakan. */
  .b-recent { grid-column: 1 / 3; grid-row: 3; }
  .b-sched  { grid-column: 3; grid-row: 3; }

  /* Kartu status dan kartu ringkasan disejajarkan tingginya — keduanya sebaris
     dan perbedaan tinggi di situ terlihat seperti kesalahan.
     Ketiga daftar di bawahnya TIDAK dipaksa sama tinggi: isinya berbeda-beda
     panjang, dan menyamakannya hanya menambah ruang putih kosong di dalam kartu
     yang isinya sedikit. */
  .screen { align-items: start; }
  .b-status, .b-angka { align-self: stretch; display: flex; flex-direction: column; }
  .status, .ringkas { flex: 1; }

  .head { align-items: flex-end; margin-bottom: 2px; }
  .head-id, .avatar { display: none; }
  .head-desk { display: flex; align-items: baseline; justify-content: space-between; gap: 16px; width: 100%; }
  .page-title { margin: 0; font-size: 23px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); }
  .page-date { margin: 0; font-size: 13px; color: var(--muted); }

  .status { padding: 20px 22px; }
  .s-title { font-size: 19px; }
  .ringkas { padding: 16px 20px 10px; }
  /* Kedua baris berbagi tinggi kartu. Tanpa ini, saat status sedang panjang
     (bidan tidak di klinik) kartu ini ikut memanjang dan menyisakan ruang
     kosong menganggur di bawahnya. */
  .rk-row { flex: 1; padding: 14px 0; }
  .rk-num { font-size: 24px; min-width: 34px; }
}
</style>
