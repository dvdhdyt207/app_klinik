// Helper format tanggal/waktu Bahasa Indonesia (dari prototipe .dc.html)
const BL = ['Jan','Feb','Mar','Apr','Mei','Jun','Jul','Agu','Sep','Okt','Nov','Des']
const BL_FULL = ['Januari','Februari','Maret','April','Mei','Juni','Juli','Agustus','September','Oktober','November','Desember']
const HR = ['Min','Sen','Sel','Rab','Kam','Jum','Sab']
const pad = (n) => String(n).padStart(2, '0')

export const hhmm = (ts) => { const d = new Date(ts); return pad(d.getHours()) + '.' + pad(d.getMinutes()) }
export const shortDate = (ts) => { const d = new Date(ts); return HR[d.getDay()] + ', ' + d.getDate() + ' ' + BL[d.getMonth()] }

export function eventWhen(e) {
  const s = new Date(e.startTs), en = new Date(e.endTs)
  const sameDay = s.toDateString() === en.toDateString()
  if (e.allDay) {
    if (sameDay) return shortDate(e.startTs) + ' · seharian'
    const days = Math.round((new Date(e.endTs).setHours(0,0,0,0) - new Date(e.startTs).setHours(0,0,0,0)) / 86400000) + 1
    return shortDate(e.startTs) + ' – ' + shortDate(e.endTs) + ' · ' + days + ' hari'
  }
  if (sameDay) return shortDate(e.startTs) + ' · ' + hhmm(e.startTs) + '–' + hhmm(e.endTs)
  return shortDate(e.startTs) + ' ' + hhmm(e.startTs) + ' – ' + shortDate(e.endTs) + ' ' + hhmm(e.endTs)
}

export function untilLabel(ts) {
  const diff = ts - Date.now()
  if (diff <= 0) return 'sudah lewat perkiraan'
  const m = Math.round(diff / 60000)
  if (m < 60) return '± ' + m + ' menit lagi'
  const h = Math.floor(m / 60), mm = m % 60
  return '± ' + h + ' jam' + (mm ? ' ' + mm + ' mnt' : '') + ' lagi'
}

// label tanggal kunjungan: "Hari ini HH.MM" / "Kemarin HH.MM" / "D Mon"
export function fmtDate(ts) {
  const d = new Date(ts), now = new Date()
  const sameDay = d.toDateString() === now.toDateString()
  const yest = new Date(now.getTime() - 86400000).toDateString() === d.toDateString()
  const hh = hhmm(ts)
  if (sameDay) return 'Hari ini ' + hh
  if (yest) return 'Kemarin ' + hh
  return d.getDate() + ' ' + BL[d.getMonth()]
}

export const dayKey = (ts) => new Date(ts).toDateString()
export function dayLabel(ts) {
  const d = new Date(ts), now = new Date()
  if (d.toDateString() === now.toDateString()) return 'Hari ini'
  if (new Date(now.getTime() - 86400000).toDateString() === d.toDateString()) return 'Kemarin'
  return d.getDate() + ' ' + BL_FULL[d.getMonth()]
}

// waktu relatif untuk halaman pasien ("baru saja", "5 menit lalu", ...)
export function relTime(ts) {
  const s = Math.max(0, Math.floor((Date.now() - ts) / 1000))
  if (s < 45) return 'baru saja'
  const m = Math.floor(s / 60); if (m < 60) return m + ' menit lalu'
  const h = Math.floor(m / 60); if (h < 24) return h + ' jam lalu'
  return Math.floor(h / 24) + ' hari lalu'
}

// countdown halaman pasien ("± X menit lagi")
export function countdown(ts) {
  const diff = ts - Date.now(); if (diff <= 0) return 'sekitar sekarang'
  const m = Math.round(diff / 60000); if (m < 60) return '± ' + m + ' menit lagi'
  const h = Math.floor(m / 60), mm = m % 60; return '± ' + h + ' jam' + (mm ? ' ' + mm + ' mnt' : '') + ' lagi'
}

const HARI = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
// "Rabu, 6 Agustus 2026" — dipakai di kepala Beranda.
export function tanggalPanjang(ts) {
  const d = new Date(ts)
  return HARI[d.getDay()] + ', ' + d.getDate() + ' ' + BL_FULL[d.getMonth()] + ' ' + d.getFullYear()
}

// Monogram logo dari nama klinik. Dipakai halaman pasien, sidebar Bidan App,
// dan avatar Beranda — satu fungsi, karena nama klinik bisa diubah bidan lewat
// Setelan dan logo yang tidak ikut berubah akan jadi satu-satunya bagian
// aplikasi yang menyebut klinik yang salah.
const KATA_UMUM = new Set(['klinik', 'praktik', 'bidan', 'mandiri', 'pmb'])
export function monogram(nama) {
  const kata = String(nama || '').split(/\s+/).filter(Boolean)
  const inti = kata.filter((w) => !KATA_UMUM.has(w.toLowerCase()))
  const dipakai = (inti.length ? inti : kata).slice(0, 2)
  return dipakai.map((w) => w[0]).join('').toUpperCase() || 'K'
}

// util konversi date/time input (untuk Form Agenda)
export const toDateStr = (ts) => { const d = new Date(ts); return d.getFullYear() + '-' + pad(d.getMonth()+1) + '-' + pad(d.getDate()) }
export const toTimeStr = (ts) => { const d = new Date(ts); return pad(d.getHours()) + ':' + pad(d.getMinutes()) }
export const mkTs = (dateStr, timeStr) => new Date(dateStr + 'T' + (timeStr || '00:00')).getTime()
