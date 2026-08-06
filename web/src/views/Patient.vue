<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { api } from '../api/client'
import { hhmm, eventWhen, relTime, countdown, monogram as buatMonogram } from '../lib/format'

const st = ref(null)
const now = ref(Date.now())
let pollTimer, tickTimer

async function refresh() {
  try { st.value = await api.publicStatus() } catch { /* offline: pertahankan render terakhir */ }
}

const upcoming = computed(() => {
  void now.value
  if (!st.value) return []
  const n = Date.now()
  return (st.value.events || []).filter((e) => e.endTs >= n).sort((a, b) => a.startTs - b.startTs).slice(0, 4)
})
// Ketiganya datang dari /api/public/status, tidak lagi ditulis mati di sini.
// Fallback dipakai hanya sebelum permintaan pertama selesai — bukan sebagai
// nilai cadangan permanen, karena nilai cadangan yang salah persis itulah yang
// dulu tampil ke pasien.
const namaKlinik = computed(() => st.value?.clinic || 'Klinik')
const alamat = computed(() => st.value?.alamat || '')
const whatsapp = computed(() => st.value?.whatsapp || '')
const updated = computed(() => { void now.value; return st.value ? relTime(st.value.ts) : '' })
const cd = computed(() => { void now.value; return st.value?.awayUntil ? countdown(st.value.awayUntil) : '' })

const monogram = computed(() => buatMonogram(namaKlinik.value))

onMounted(() => {
  refresh()
  pollTimer = setInterval(refresh, 10000)              // polling status 10 dtk
  tickTimer = setInterval(() => { now.value = Date.now() }, 30000) // refresh label waktu
})
onUnmounted(() => { clearInterval(pollTimer); clearInterval(tickTimer) })
</script>

<template>
  <div class="page">
    <div class="shell">
      <!-- brand header -->
      <header class="brand">
        <div class="logo" aria-hidden="true">{{ monogram }}</div>
        <div class="bwrap">
          <h1 class="bname">{{ namaKlinik }}</h1>
          <p class="btag">Praktik Bidan Mandiri · melayani 24 jam</p>
        </div>
      </header>

      <!-- ---------- hero: jawaban atas satu-satunya pertanyaan pasien ---------- -->
      <section v-if="st && st.hadir" class="hero hero-open">
        <div class="hero-in">
          <div class="badge badge-open"><span class="pulse" />SEDANG BUKA</div>
          <p class="hero-title">Bidan sedang di klinik</p>
          <p class="hero-line">Silakan datang langsung untuk pemeriksaan atau konsultasi.</p>
          <p class="hero-updated">Diperbarui {{ updated }}</p>
        </div>
      </section>

      <section v-else-if="st" class="hero hero-away">
        <div class="hero-in">
          <div class="badge badge-away"><span class="graydot" />SEDANG TIDAK ADA</div>
          <p class="hero-title">Bidan sedang tidak di tempat</p>
          <!-- Catatan bidan tetap paragraf, bukan judul: isinya teks bebas yang
               panjangnya tak terduga, dan alasan sungguhan ("menolong persalinan")
               justru perlu terbaca utuh, bukan terpotong demi kerapian judul. -->
          <p v-if="st.awayNote" class="away-note">{{ st.awayNote }}</p>
          <div class="return-box">
            <template v-if="st.awayUntil">
              <p class="rb-label">Perkiraan kembali</p>
              <p class="rb-time">Pukul {{ hhmm(st.awayUntil) }}</p>
              <p class="rb-cd">{{ cd }}</p>
            </template>
            <p v-else class="rb-none">Waktu kembali belum dapat dipastikan</p>
          </div>
          <p class="hero-updated">Diperbarui {{ updated }}</p>
        </div>
      </section>

      <!-- Sebelum permintaan pertama selesai. Kotak diam lebih baik daripada
           menebak "buka" atau "tutup" — tebakan yang salah membuat orang
           berangkat ke klinik yang kosong. -->
      <section v-else class="hero hero-load" aria-busy="true">
        <div class="hero-in"><p class="loading">Memuat status klinik…</p></div>
      </section>

      <!-- ---------- dua kartu: di HP bertumpuk, di layar lebar berdampingan ---------- -->
      <div class="grid">
        <section class="card2">
          <h2 class="c-title">Jadwal Bidan</h2>
          <p class="c-sub">Agenda mendatang saat bidan berhalangan</p>
          <ul v-if="upcoming.length" class="events">
            <li v-for="e in upcoming" :key="e.id" class="ev">
              <span class="ev-dot" :class="e.allDay ? 'is-allday' : 'is-timed'" />
              <div class="ev-grow">
                <p class="ev-title">{{ e.title }}</p>
                <p class="ev-when">{{ eventWhen(e) }}</p>
              </div>
            </li>
          </ul>
          <p v-else class="empty">Tidak ada agenda mendatang — bidan siap melayani.</p>
        </section>

        <!-- Seluruh kartunya hilang kalau bidan belum mengisi apa pun. Sebelumnya
             di sini tertulis alamat dan nomor karangan; halaman yang jujur
             menampilkan lebih sedikit lebih baik daripada halaman yang mengarang. -->
        <section v-if="alamat || whatsapp" class="card2">
          <h2 class="c-title">Lokasi &amp; Kontak</h2>
          <p class="c-sub">Alamat praktik dan nomor yang bisa dihubungi</p>
          <p v-if="alamat" class="addr">{{ alamat }}</p>
          <a v-if="whatsapp" class="wa" :href="'https://wa.me/' + whatsapp" target="_blank" rel="noopener">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M12 2a10 10 0 00-8.6 15l-1.3 4.7 4.8-1.3A10 10 0 1012 2zm0 2a8 8 0 11-4.2 14.8l-.3-.2-2.8.8.8-2.7-.2-.3A8 8 0 0112 4z"/></svg>
            Hubungi via WhatsApp
          </a>
        </section>
      </div>

      <p class="foot">Status &amp; jadwal diperbarui otomatis oleh {{ namaKlinik }}</p>
    </div>
  </div>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: var(--patient-bg);
  padding: 20px 16px 44px;
}
/* clamp: satu aturan untuk HP sampai monitor lebar. Halaman ini dulu terkunci
   460px, jadi di 1440px dua pertiga layarnya kosong. */
.shell {
  width: 100%;
  max-width: 1120px;
  margin: 0 auto;
  display: flex; flex-direction: column;
  gap: clamp(14px, 1.6vw, 22px);
}

/* ---------------- brand ---------------- */
.brand { display: flex; align-items: center; gap: 14px; padding: 4px 2px; }
.logo {
  width: 50px; height: 50px; border-radius: 16px; flex-shrink: 0;
  background: linear-gradient(145deg, var(--accent) 0%, var(--accent-press) 100%);
  color: #fff; font-weight: 800; font-size: 17px; letter-spacing: .01em;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 6px 16px rgba(18,128,106,.28);
}
.bwrap { min-width: 0; }
.bname { font-size: clamp(19px, 2vw, 23px); font-weight: 800; letter-spacing: -.02em; line-height: 1.15; margin: 0; }
.btag { font-size: 13px; color: var(--muted); font-weight: 600; margin: 2px 0 0; }

/* ---------------- hero ---------------- */
.hero { border-radius: var(--r-xl); padding: clamp(26px, 3.4vw, 44px) clamp(20px, 3vw, 40px); }
/* Isinya dibatasi lebar & dipusatkan: pita boleh selebar layar, kalimatnya
   tidak — baris sepanjang 1100px tidak terbaca, ia terpindai. */
.hero-in { max-width: 640px; margin: 0 auto; text-align: center; }

.hero-open { background: linear-gradient(160deg, var(--accent) 0%, var(--accent-press) 100%); color: #fff; box-shadow: var(--shadow-accent); }
.hero-away { background: var(--away-soft); border: 1px solid var(--away-line); box-shadow: var(--shadow-md); }
.hero-load { background: var(--card); border: 1px solid var(--border); }

.badge {
  display: inline-flex; align-items: center; gap: 9px;
  padding: 8px 16px; border-radius: 30px;
  font-size: 12.5px; font-weight: 700; letter-spacing: .07em;
}
.badge-open { background: rgba(255,255,255,.18); }
.badge-away { background: rgba(184,121,28,.13); color: var(--away-ink); }
.pulse { width: 10px; height: 10px; border-radius: 50%; background: #fff; animation: pit-pulse 1.8s infinite; }
.graydot { width: 10px; height: 10px; border-radius: 50%; background: var(--away); }

.hero-title {
  font-size: clamp(28px, 3.6vw, 42px); font-weight: 800;
  letter-spacing: -.03em; line-height: 1.13;
  margin: 18px 0 0; text-wrap: balance;
}
.hero-away .hero-title { color: var(--ink); font-size: clamp(24px, 2.8vw, 34px); }
.hero-line { font-size: clamp(14.5px, 1.15vw, 16.5px); opacity: .92; margin: 12px 0 0; line-height: 1.5; }
.away-note { font-size: clamp(15px, 1.2vw, 17px); color: var(--away-ink); font-weight: 600; margin: 12px 0 0; line-height: 1.5; }

.return-box {
  margin: 22px auto 0; max-width: 340px;
  background: var(--card); border: 1px solid var(--away-line);
  border-radius: var(--r-lg); padding: 18px;
}
.rb-label { font-size: 12px; font-weight: 700; color: var(--label); letter-spacing: .06em; text-transform: uppercase; margin: 0; }
.rb-time { font-size: clamp(26px, 2.4vw, 32px); font-weight: 800; color: var(--ink); margin: 5px 0 0; letter-spacing: -.02em; }
.rb-cd { font-size: 13.5px; color: var(--away-ink); font-weight: 700; margin: 3px 0 0; }
.rb-none { font-size: 14.5px; color: var(--label); font-weight: 600; margin: 0; }

.hero-updated { font-size: 12.5px; font-weight: 600; margin: 22px 0 0; padding-top: 16px; }
.hero-open .hero-updated { border-top: 1px solid rgba(255,255,255,.22); opacity: .85; }
.hero-away .hero-updated { border-top: 1px solid var(--away-line); color: var(--muted); }
.loading { color: var(--muted); font-size: 15px; font-weight: 600; margin: 0; }

/* ---------------- kartu ---------------- */
/* auto-fit + minmax: satu kolom di HP, dua kolom begitu muat, tanpa breakpoint
   tersendiri. Kartu Lokasi bisa hilang seluruhnya (belum diisi) — dengan
   auto-fit, Jadwal lalu melebar sendiri alih-alih menyisakan kolom kosong. */
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(310px, 1fr)); gap: clamp(14px, 1.6vw, 22px); }
/* Kartu Lokasi hilang selama bidan belum mengisi alamat & nomor — dan itu
   keadaan produksi hari ini, bukan kasus pinggiran. Tanpa batas ini, Jadwal
   melar selebar 1120px untuk memuat satu baris kalimat. */
.grid > .card2:only-child { max-width: 640px; margin: 0 auto; width: 100%; }
.card2 {
  background: var(--card); border: 1px solid var(--border);
  border-radius: var(--r-lg); padding: clamp(18px, 1.8vw, 26px);
  box-shadow: var(--shadow-sm);
  display: flex; flex-direction: column;
}
.c-title { font-size: 16px; font-weight: 800; margin: 0; letter-spacing: -.01em; }
.c-sub { font-size: 12.5px; color: var(--muted); margin: 5px 0 16px; }

.events { display: flex; flex-direction: column; gap: 14px; list-style: none; margin: 0; padding: 0; }
.ev { display: flex; gap: 13px; align-items: flex-start; }
.ev-dot { width: 9px; height: 9px; border-radius: 50%; margin-top: 6px; flex-shrink: 0; }
.ev-dot.is-timed { background: var(--accent); }
.ev-dot.is-allday { background: var(--away); }
.ev-grow { flex: 1; min-width: 0; }
.ev-title { font-size: 14.5px; font-weight: 700; line-height: 1.35; margin: 0; }
.ev-when { font-size: 12.5px; color: var(--muted); margin: 2px 0 0; }
.empty { font-size: 13.5px; color: var(--muted); font-weight: 600; margin: 0; }

/* pre-line: alamat diketik bidan di textarea, jadi enter-nya harus tetap jadi
   baris baru. Sebelumnya alamatnya <br /> yang ditulis tangan di template. */
.addr { font-size: 14.5px; color: var(--text-secondary); line-height: 1.6; margin: 0 0 18px; white-space: pre-line; }
.wa {
  display: flex; align-items: center; justify-content: center; gap: 9px;
  margin-top: auto;                     /* tombol rata bawah walau kartu kiri lebih tinggi */
  background: var(--whatsapp); color: #fff;
  border-radius: var(--r-md); padding: 15px;
  font-size: 14.5px; font-weight: 700;
  transition: background .15s ease;
}
.wa:hover { background: var(--accent-hover); }

.foot { text-align: center; font-size: 12px; color: var(--muted2); padding: 8px 0 0; font-weight: 600; margin: 0; }

/* ---------------- layar lebar ---------------- */
@media (min-width: 900px) {
  .page { padding: 34px 28px 56px; }
  .brand { padding: 0 2px 4px; }
}
</style>
