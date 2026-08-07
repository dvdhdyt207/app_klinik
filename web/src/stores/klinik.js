import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../api/client'
import { SATUAN_BAWAAN, satuanAtau } from '../lib/obat'
import { eventWhen, fmtDate, dayKey, dayLabel } from '../lib/format'

const emptyDraft = () => ({ name: '', age: '', gejala: '', items: [] })
const pad = (n) => String(n).padStart(2, '0')
const toDate = (ts) => { const d = new Date(ts); return d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate()) }
const toTime = (ts) => { const d = new Date(ts); return pad(d.getHours()) + ':' + pad(d.getMinutes()) }

// Warna yang dipakai derivasi. Ditulis sebagai var() — bukan hex — supaya
// tokens.css tetap satu-satunya sumber kebenaran warna. Nilai ini berakhir di
// atribut style inline (`:style="{ background: m.color }"`), dan var() ikut
// diselesaikan di sana persis seperti di stylesheet.
//
// Sebelumnya di sini ada hex tersalin, termasuk biru lama #2f6ce0. Akibatnya
// mengganti palet di tokens.css tidak mengubah titik agenda dan angka stok —
// warna yang hidup di JavaScript tidak ikut terbawa.
const COL = {
  warning: 'var(--warning)', accent: 'var(--accent)', past: 'var(--disabled)',
}

export const useKlinik = defineStore('klinik', () => {
  // ---- data server ----
  const data = ref({ status: null, events: [], visits: [] })
  const loading = ref(true)
  const error = ref(null)
  const now = ref(Date.now()) // di-refresh berkala agar countdown/label waktu reaktif

  // ---- UI state ----
  const screen = ref('beranda')
  const modal = ref(null) // 'visit'|'pick'|'away'|'jadwal'|'event'
  const query = ref('')
  const draft = ref(emptyDraft())
  const awayDraft = ref({ note: '', minutes: null })
  const eventDraft = ref(null)
  // Satuan untuk obat yang namanya baru diketik di layar Cari Obat. Dipegang di
  // sini, bukan di komponennya, supaya pilihannya tidak ikut hilang setiap kali
  // modal dirender ulang.
  const pickSatuan = ref(SATUAN_BAWAAN)

  // ---- load & polling ----
  async function refresh() {
    try {
      data.value = await api.getState()
      error.value = null
    } catch (e) {
      error.value = e.message || 'Gagal memuat data'
    } finally {
      loading.value = false
    }
  }
  refresh()
  setInterval(() => { now.value = Date.now() }, 30000)        // refresh label waktu
  setInterval(() => { if (!modal.value) refresh() }, 15000)   // sinkron ringan saat idle

  // ---- status ----
  async function setHadir() {
    await api.setStatus({ bidanHadir: true, awayNote: '', awayUntil: null })
    modal.value = null; refresh()
  }
  function openAway() {
    awayDraft.value = { note: data.value.status?.awayNote || '', minutes: null }
    modal.value = 'away'
  }
  async function confirmAway() {
    const au = awayDraft.value.minutes ? Date.now() + awayDraft.value.minutes * 60000 : null
    await api.setStatus({ bidanHadir: false, awayNote: (awayDraft.value.note || '').trim(), awayUntil: au })
    modal.value = null; refresh()
  }
  async function extend(mins) {
    const cur = data.value.status?.awayUntil || Date.now()
    await api.setStatus({ bidanHadir: false, awayNote: data.value.status?.awayNote || '', awayUntil: cur + mins * 60000 })
    refresh()
  }

  // ---- events ----
  function openEventForm(ev) {
    let dr
    if (ev) {
      dr = { id: ev.id, title: ev.title, allDay: !!ev.allDay,
        startDate: toDate(ev.startTs), startTime: toTime(ev.startTs),
        endDate: toDate(ev.endTs), endTime: toTime(ev.endTs) }
    } else {
      const d = new Date()
      dr = { id: null, title: '', allDay: false,
        startDate: toDate(d.getTime()), startTime: pad(d.getHours()) + ':00',
        endDate: toDate(d.getTime()), endTime: pad(Math.min(23, d.getHours() + 1)) + ':00' }
    }
    eventDraft.value = dr; modal.value = 'event'
  }
  async function saveEvent() {
    const e = eventDraft.value
    if (!e || !e.title.trim() || !e.startDate) return
    const mk = (date, time) => new Date(date + 'T' + time).getTime()
    let startTs = mk(e.startDate, e.allDay ? '00:00' : (e.startTime || '00:00'))
    let endTs = mk(e.endDate || e.startDate, e.allDay ? '23:59' : (e.endTime || '23:59'))
    if (endTs < startTs) endTs = startTs
    const body = { title: e.title.trim(), allDay: e.allDay, startTs, endTs }
    if (e.id) await api.updateEvent(e.id, body); else await api.createEvent(body)
    eventDraft.value = null; modal.value = 'jadwal'; refresh()
  }
  async function deleteEvent() {
    if (!eventDraft.value?.id) return
    await api.deleteEvent(eventDraft.value.id)
    eventDraft.value = null; modal.value = 'jadwal'; refresh()
  }

  // ---- pilih obat ----
  // Satu konteks saja sekarang: obat selalu dipilih untuk catatan kunjungan.
  // Dulu layar ini melayani dua tujuan (kunjungan & tambah stok) lewat
  // pickContext; menambah stok sudah tidak ada.
  function openPickVisit() { query.value = ''; pickSatuan.value = SATUAN_BAWAAN; modal.value = 'pick' }
  function backFromPick() { modal.value = 'visit' }

  function pick(c) {
    const nama = (c.name || '').trim()
    if (!nama) return
    const items = [...draft.value.items]
    const idx = items.findIndex((i) => i.name.toLowerCase() === nama.toLowerCase())
    // Obat yang sama dipilih dua kali menambah jumlahnya, bukan membuat baris
    // kedua dengan nama yang sama.
    if (idx >= 0) items[idx] = { ...items[idx], qty: items[idx].qty + 1 }
    else items.push({ name: nama, qty: 1, unit: satuanAtau(c.unit) })
    draft.value = { ...draft.value, items }; modal.value = 'visit'
  }

  // ---- visit ----
  function openVisit() { draft.value = emptyDraft(); modal.value = 'visit' }

  // Mengisi form dari pasien yang pernah datang. Umur ikut terisi karena itu
  // satu-satunya keterangan lain yang tersimpan tentang orangnya — dan mengetik
  // ulang umur pasien yang sama tiap kunjungan hanya menambah peluang salah.
  // Nilainya tetap bisa diubah: umur berubah, dan yang tercatat harus umur saat
  // kunjungan ini.
  function pakaiPasien(p) {
    draft.value = { ...draft.value, name: p.name, age: p.age ? String(p.age) : '' }
  }
  async function saveVisit() {
    const dr = draft.value
    if (!dr.name.trim() || dr.items.length === 0) return
    await api.createVisit({
      name: dr.name.trim(), age: parseInt(dr.age, 10) || 0, gejala: (dr.gejala || '').trim(),
      items: dr.items.map((i) => ({ name: i.name, qty: i.qty, unit: i.unit })),
    })
    draft.value = emptyDraft(); modal.value = null; screen.value = 'kunjungan'; refresh()
  }

  // ---- navigasi ----
  function goScreen(s) { screen.value = s; modal.value = null }
  function closeModal() { modal.value = null }
  function openJadwal() { modal.value = 'jadwal' }
  function openEventNew() { openEventForm(null) }
  function backJadwal() { eventDraft.value = null; modal.value = 'jadwal' }

  // ---- derivasi (mirror useDerived) ----
  const derived = computed(() => {
    const _ = now.value // dep agar recompute berkala
    const visits = data.value.visits || []
    const events = data.value.events || []
    const nowMs = Date.now()

    // Daftar obat & daftar pasien dibentuk dari RIWAYAT KUNJUNGAN, bukan dari
    // tabel tersendiri: tidak ada lagi daftar obat yang harus didata lebih dulu,
    // dan tidak ada daftar yang bisa melenceng dari apa yang benar-benar pernah
    // dipakai. Sekali dicatat, nama itu bisa dipakai ulang.
    //
    // `visits` datang urut terbaru dulu (ORDER BY ts DESC), jadi kemunculan
    // PERTAMA sebuah nama adalah pemakaian terakhirnya — itu yang disimpan.
    const obatMap = new Map()
    const pasienMap = new Map()
    visits.forEach((v) => {
      const nm = (v.name || '').trim()
      const key = nm.toLowerCase()
      if (nm && !pasienMap.has(key)) pasienMap.set(key, { name: nm, age: v.age, lastTs: v.ts, count: 0 })
      if (nm) pasienMap.get(key).count++
      ;(v.items || []).forEach((it) => {
        const on = (it.name || '').trim()
        const ok = on.toLowerCase()
        if (!on || obatMap.has(ok)) return
        obatMap.set(ok, { name: on, unit: satuanAtau(it.unit), lastTs: v.ts })
      })
    })
    const obatDipakai = [...obatMap.values()].sort((a, b) => b.lastTs - a.lastTs)
    const pasienDipakai = [...pasienMap.values()].sort((a, b) => b.lastTs - a.lastTs)

    const visitCards = visits.map((v) => ({
      ...v, initial: (v.name || '?').trim().charAt(0).toUpperCase(),
      meta: v.age + ' tahun', dateLabel: fmtDate(v.ts),
      chips: (v.items || []).map((it) => it.name + ' · ' + it.qty + ' ' + it.unit),
    }))
    const todayKey = new Date().toDateString()
    const todayCount = visits.filter((v) => dayKey(v.ts) === todayKey).length

    const byDay = {}
    visits.forEach((v) => {
      const key = dayKey(v.ts)
      if (!byDay[key]) byDay[key] = { ts: v.ts, count: 0, names: [] }
      byDay[key].count++; byDay[key].names.push((v.name || '').split(' ')[0])
      if (v.ts > byDay[key].ts) byDay[key].ts = v.ts
    })
    const dailyRecap = Object.values(byDay).sort((a, b) => b.ts - a.ts).map((row) => ({
      label: dayLabel(row.ts), count: row.count,
      patients: row.names.slice(0, 3).join(', ') + (row.names.length > 3 ? ' +' + (row.names.length - 3) : ''),
    }))
    const allEvents = [...events].sort((a, b) => a.startTs - b.startTs).map((e) => ({
      ...e, when: eventWhen(e), past: e.endTs < nowMs,
      dot: e.endTs < nowMs ? COL.past : (e.allDay ? COL.warning : COL.accent),
    }))
    const upcoming = allEvents.filter((e) => !e.past)

    return { obatDipakai, pasienDipakai, visitCards, todayCount, dailyRecap, allEvents, upcoming,
      totalVisits: visits.length, hariTercatat: dailyRecap.length }
  })

  const status = computed(() => data.value.status || {})
  const profil = computed(() => data.value.profil || { clinic: '', alamat: '', whatsapp: '' })

  // ---- profil klinik ----
  // Server mengembalikan bentuk yang sudah dinormalkan (nomor jadi 62…), jadi
  // hasilnya dipakai langsung alih-alih menebak apa yang tersimpan.
  async function simpanProfil(body) {
    const baru = await api.setProfil(body)
    data.value = { ...data.value, profil: baru, status: { ...data.value.status, clinic: baru.clinic } }
    return baru
  }

  return {
    data, loading, error, now, status, profil, simpanProfil,
    screen, modal, query,
    draft, awayDraft, eventDraft, pickSatuan,
    derived, refresh,
    setHadir, openAway, confirmAway, extend,
    openEventForm, saveEvent, deleteEvent, openJadwal, openEventNew, backJadwal,
    openPickVisit, backFromPick, pick,
    openVisit, saveVisit, pakaiPasien, goScreen, closeModal,
  }
})
