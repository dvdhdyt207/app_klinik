<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { api } from '../api/client'
import { hhmm, eventWhen, relTime, countdown } from '../lib/format'

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

onMounted(() => {
  refresh()
  pollTimer = setInterval(refresh, 10000)              // polling status 10 dtk
  tickTimer = setInterval(() => { now.value = Date.now() }, 30000) // refresh label waktu
})
onUnmounted(() => { clearInterval(pollTimer); clearInterval(tickTimer) })
</script>

<template>
  <div class="page">
    <div class="col">
      <!-- brand header -->
      <div class="brand">
        <div class="logo">P</div>
        <div>
          <div class="bname">{{ namaKlinik }}</div>
          <div class="btag">Praktik Bidan Mandiri · melayani 24 jam</div>
        </div>
      </div>

      <!-- hero -->
      <div v-if="st && st.hadir" class="hero-open">
        <div class="badge open"><span class="pulse" />SEDANG BUKA</div>
        <div class="hero-title">Bidan sedang<br />di klinik</div>
        <div class="hero-line">Silakan datang langsung ke klinik untuk pemeriksaan atau konsultasi.</div>
        <div class="hero-updated open-updated">Diperbarui {{ updated }}</div>
      </div>

      <div v-else-if="st" class="hero-away">
        <div class="badge away"><span class="graydot" />SEDANG TIDAK ADA</div>
        <div class="hero-title dark">Bidan sedang<br />tidak di tempat</div>
        <div v-if="st.awayNote" class="away-note">{{ st.awayNote }}</div>
        <div class="return-box">
          <template v-if="st.awayUntil">
            <div class="rb-label">Perkiraan kembali</div>
            <div class="rb-time">Pukul {{ hhmm(st.awayUntil) }}</div>
            <div class="rb-cd">{{ cd }}</div>
          </template>
          <div v-else class="rb-none">Waktu kembali belum dapat dipastikan</div>
        </div>
        <div class="hero-updated away-updated">Diperbarui {{ updated }}</div>
      </div>

      <!-- jadwal -->
      <div class="card2">
        <div class="c-title">Jadwal Bidan</div>
        <div class="c-sub">Agenda mendatang saat bidan berhalangan</div>
        <div v-if="upcoming.length" class="events">
          <div v-for="e in upcoming" :key="e.id" class="ev">
            <span class="ev-dot" :style="{ background: e.allDay ? '#e0a52a' : '#2f6ce0' }" />
            <div class="ev-grow">
              <div class="ev-title">{{ e.title }}</div>
              <div class="ev-when">{{ eventWhen(e) }}</div>
            </div>
          </div>
        </div>
        <div v-else class="empty">Tidak ada agenda mendatang — bidan siap melayani.</div>
      </div>

      <!-- lokasi & kontak -->
      <!-- Seluruh kartunya hilang kalau bidan belum mengisi apa pun. Sebelumnya
           di sini tertulis alamat dan nomor karangan; halaman yang jujur
           menampilkan lebih sedikit lebih baik daripada halaman yang mengarang. -->
      <div v-if="alamat || whatsapp" class="card2">
        <div class="c-title mb12">Lokasi &amp; Kontak</div>
        <div v-if="alamat" class="addr">{{ alamat }}</div>
        <a v-if="whatsapp" class="wa" :href="'https://wa.me/' + whatsapp" target="_blank" rel="noopener">
          <svg width="19" height="19" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2a10 10 0 00-8.6 15l-1.3 4.7 4.8-1.3A10 10 0 1012 2zm0 2a8 8 0 11-4.2 14.8l-.3-.2-2.8.8.8-2.7-.2-.3A8 8 0 0112 4z"/></svg>
          Hubungi via WhatsApp
        </a>
      </div>

      <div class="foot">Status &amp; jadwal diperbarui otomatis oleh {{ namaKlinik }}</div>
    </div>
  </div>
</template>

<style scoped>
.page { min-height: 100vh; display: flex; justify-content: center; padding: 22px 16px 40px; background: var(--patient-bg); }
.col { width: 100%; max-width: 460px; display: flex; flex-direction: column; gap: 16px; }

.brand { display: flex; align-items: center; gap: 13px; padding: 6px 2px; }
.logo { width: 46px; height: 46px; border-radius: 14px; background: var(--accent); color: #fff; font-weight: 800; font-size: 20px; display: flex; align-items: center; justify-content: center; }
.bname { font-size: 19px; font-weight: 800; letter-spacing: -.02em; line-height: 1.1; }
.btag { font-size: 12.5px; color: var(--muted); font-weight: 600; }

.hero-open { background: var(--green); border-radius: 24px; padding: 30px 24px; color: #fff; box-shadow: 0 18px 40px rgba(31,157,87,.28); text-align: center; }
.hero-away { background: #fff; border: 1px solid var(--border2); border-radius: 24px; padding: 30px 24px; text-align: center; box-shadow: 0 12px 30px rgba(22,32,46,.06); }
.badge { display: inline-flex; align-items: center; gap: 9px; padding: 7px 15px; border-radius: 30px; font-size: 12.5px; font-weight: 700; letter-spacing: .06em; }
.badge.open { background: rgba(255,255,255,.18); }
.badge.away { background: #f0f2f6; color: var(--label); }
.pulse { width: 10px; height: 10px; border-radius: 50%; background: #fff; animation: pit-pulse 1.8s infinite; }
.graydot { width: 10px; height: 10px; border-radius: 50%; background: #b6bdc9; }
.hero-title { font-size: 32px; font-weight: 800; letter-spacing: -.03em; margin-top: 18px; line-height: 1.1; }
.hero-title.dark { color: var(--ink); }
.hero-line { font-size: 14.5px; opacity: .92; margin-top: 10px; line-height: 1.45; }
.away-note { font-size: 14.5px; color: var(--text-secondary); margin-top: 10px; line-height: 1.45; font-weight: 600; }
.return-box { margin-top: 18px; background: var(--fill2); border-radius: 16px; padding: 16px; }
.rb-label { font-size: 12px; font-weight: 700; color: var(--muted); letter-spacing: .06em; text-transform: uppercase; }
.rb-time { font-size: 24px; font-weight: 800; color: var(--ink); margin-top: 4px; letter-spacing: -.02em; }
.rb-cd { font-size: 13px; color: var(--accent); font-weight: 700; margin-top: 2px; }
.rb-none { font-size: 14px; color: var(--label); font-weight: 600; }
.hero-updated { font-size: 12.5px; font-weight: 600; }
.open-updated { margin-top: 20px; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.22); opacity: .85; }
.away-updated { margin-top: 18px; padding-top: 16px; border-top: 1px solid #eef2f7; color: var(--muted2); }

.card2 { background: #fff; border: 1px solid var(--border); border-radius: 20px; padding: 20px; }
.c-title { font-size: 15px; font-weight: 800; margin-bottom: 6px; }
.c-title.mb12 { margin-bottom: 12px; }
.c-sub { font-size: 12.5px; color: var(--muted); margin-bottom: 14px; }
.events { display: flex; flex-direction: column; gap: 12px; }
.ev { display: flex; gap: 13px; align-items: flex-start; }
.ev-dot { width: 10px; height: 10px; border-radius: 50%; margin-top: 5px; flex-shrink: 0; }
.ev-grow { flex: 1; min-width: 0; }
.ev-title { font-size: 14.5px; font-weight: 700; line-height: 1.3; }
.ev-when { font-size: 12.5px; color: var(--muted); margin-top: 1px; }
.empty { font-size: 13.5px; color: var(--muted); font-weight: 600; }
/* pre-line: alamat diketik bidan di textarea, jadi enter-nya harus tetap jadi
   baris baru. Sebelumnya alamatnya <br /> yang ditulis tangan di template. */
.addr { font-size: 14px; color: var(--text-secondary); line-height: 1.55; margin-bottom: 16px; white-space: pre-line; }
.wa { display: flex; align-items: center; justify-content: center; gap: 9px; background: var(--whatsapp); color: #fff; border-radius: 14px; padding: 14px; font-size: 14.5px; font-weight: 700; }
.foot { text-align: center; font-size: 11.5px; color: #a3adba; padding: 6px 0 2px; font-weight: 600; }
</style>
