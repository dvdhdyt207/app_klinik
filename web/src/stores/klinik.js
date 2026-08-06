import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../api/client'
import { CAT, baseUnit, isLow, isDanger } from '../lib/catalog'
import { eventWhen, fmtDate, dayKey, dayLabel } from '../lib/format'

const emptyDraft = () => ({ name: '', age: '', gejala: '', items: [] })
const pad = (n) => String(n).padStart(2, '0')
const toDate = (ts) => { const d = new Date(ts); return d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate()) }
const toTime = (ts) => { const d = new Date(ts); return pad(d.getHours()) + ':' + pad(d.getMinutes()) }

// Warna tokens yang dibutuhkan derivasi (nilai mentah agar bisa diikat inline).
const COL = {
  danger: '#d64545', dangerBg: '#fdeaea',
  warning: '#e0a52a', warningBg: '#fdf3df',
  healthy: '#3fbf8f', ink: '#16202e', accent: '#2f6ce0', past: '#c0c8d4',
}

export const useKlinik = defineStore('klinik', () => {
  // ---- data server ----
  const data = ref({ status: null, events: [], medicines: [], visits: [] })
  const loading = ref(true)
  const error = ref(null)
  const now = ref(Date.now()) // di-refresh berkala agar countdown/label waktu reaktif

  // ---- UI state ----
  const screen = ref('beranda')
  const modal = ref(null) // 'visit'|'pick'|'addstock'|'away'|'jadwal'|'event'
  const rekapTab = ref('kunjungan')
  const pickContext = ref('visit')
  const query = ref('')
  const draft = ref(emptyDraft())
  const awayDraft = ref({ note: '', minutes: null })
  const eventDraft = ref(null)
  const stockTarget = ref(null)
  const stockUnitIdx = ref(0)
  const stockCount = ref('1')

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

  // ---- pick / stok ----
  function openPickVisit() { pickContext.value = 'visit'; query.value = ''; modal.value = 'pick' }
  function openPickStock() { pickContext.value = 'stock'; query.value = ''; modal.value = 'pick' }
  function backFromPick() { modal.value = pickContext.value === 'visit' ? 'visit' : null }

  function pick(c) {
    if (pickContext.value === 'visit') {
      const items = [...draft.value.items]
      const idx = items.findIndex((i) => i.name.toLowerCase() === c.name.toLowerCase())
      if (idx >= 0) items[idx] = { ...items[idx], qty: items[idx].qty + 1 }
      else items.push({ name: c.name, cat: c.cat, qty: 1, unit: baseUnit(c.cat) })
      draft.value = { ...draft.value, items }; modal.value = 'visit'
    } else {
      const existing = (data.value.medicines || []).find((m) => m.name.toLowerCase() === c.name.toLowerCase())
      stockTarget.value = existing || { id: null, name: c.name, cat: c.cat, qty: 0 }
      stockUnitIdx.value = 0; stockCount.value = '1'; modal.value = 'addstock'
    }
  }
  async function confirmAddStock() {
    const st = stockTarget.value; if (!st) return
    const cfg = CAT[st.cat] || CAT.Tablet
    const mult = (cfg.units[stockUnitIdx.value] || cfg.units[0]).mult
    const amount = (parseInt(stockCount.value, 10) || 0) * mult
    if (amount <= 0) return
    await api.addStock({ id: st.id, name: st.name, cat: st.cat, amount })
    stockTarget.value = null; modal.value = null; screen.value = 'stok'; refresh()
  }
  function closeAddStock() { stockTarget.value = null; modal.value = null }

  // ---- visit ----
  function openVisit() { draft.value = emptyDraft(); modal.value = 'visit' }
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
    const meds = data.value.medicines || []
    const visits = data.value.visits || []
    const events = data.value.events || []
    const nowMs = Date.now()

    const low = meds.filter((m) => isLow(m.cat, m.qty)).sort((a, b) => a.qty - b.qty).map((m) => {
      const danger = isDanger(m.cat, m.qty)
      return { ...m, unit: baseUnit(m.cat), color: danger ? COL.danger : COL.warning,
        tint: danger ? COL.dangerBg : COL.warningBg, status: danger ? 'segera dibeli' : 'perlu diperhatikan' }
    })
    const medsList = meds.map((m) => {
      const lo = isLow(m.cat, m.qty)
      return { ...m, unit: baseUnit(m.cat), dot: lo ? COL.danger : COL.healthy, qtyColor: lo ? COL.danger : COL.ink }
    })
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
    const totalMedsOut = visits.reduce((s, v) => s + (v.items || []).reduce((a, i) => a + i.qty, 0), 0)

    const allEvents = [...events].sort((a, b) => a.startTs - b.startTs).map((e) => ({
      ...e, when: eventWhen(e), past: e.endTs < nowMs,
      dot: e.endTs < nowMs ? COL.past : (e.allDay ? COL.warning : COL.accent),
    }))
    const upcoming = allEvents.filter((e) => !e.past)

    return { low, medsList, visitCards, todayCount, dailyRecap, totalMedsOut, allEvents, upcoming,
      lowCount: low.length, totalVisits: visits.length, medCount: meds.length }
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
    screen, modal, rekapTab, pickContext, query,
    draft, awayDraft, eventDraft, stockTarget, stockUnitIdx, stockCount,
    derived, refresh,
    setHadir, openAway, confirmAway, extend,
    openEventForm, saveEvent, deleteEvent, openJadwal, openEventNew, backJadwal,
    openPickVisit, openPickStock, backFromPick, pick, confirmAddStock, closeAddStock,
    openVisit, saveVisit, goScreen, closeModal,
  }
})
