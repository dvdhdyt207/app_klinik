// API client. URL relatif '/api/...' — di dev di-proxy Vite ke :4000,
// di produksi di-serve origin yang sama oleh server Go.

// Dipanggil saat server menjawab 401 pada request yang seharusnya sudah punya
// sesi (mis. sesi kedaluwarsa karena aplikasi dibiarkan terbuka semalaman).
//
// Didaftarkan dari main.js, bukan dengan mengimpor router di sini: router
// mengimpor view, view mengimpor berkas ini, dan lingkaran impor itu membuat
// salah satunya masih kosong saat dipakai.
let onTanpaSesi = null
export function setPenanganTanpaSesi(fn) { onTanpaSesi = fn }

// Jalur yang 401-nya wajar dan tidak berarti "sesi habis": login (sandi salah)
// dan me (pemeriksaan awal — memang belum tentu punya sesi).
const diamSaat401 = ['/auth/login', '/auth/me']

async function req(path, options) {
  const res = await fetch('/api' + path, {
    headers: { 'Content-Type': 'application/json' },
    // Cookie sesi ikut terkirim. same-origin, bukan include: frontend dan API
    // berada di origin yang sama, dan mengirim kredensial ke origin lain tidak
    // pernah kita perlukan.
    credentials: 'same-origin',
    ...options,
  })
  if (!res.ok) {
    if (res.status === 401 && !diamSaat401.includes(path) && onTanpaSesi) onTanpaSesi()
    let msg = res.status + ''
    try { const j = await res.json(); msg = j.error || msg } catch {}
    const err = new Error(msg)
    err.status = res.status
    throw err
  }
  if (res.status === 204) return null
  return res.json()
}

export const api = {
  // ---- auth ----
  login: (body) => req('/auth/login', { method: 'POST', body: JSON.stringify(body) }),
  logout: () => req('/auth/logout', { method: 'POST' }),
  me: () => req('/auth/me'),
  gantiSandi: (body) => req('/auth/password', { method: 'POST', body: JSON.stringify(body) }),

  getState: () => req('/state'),
  publicStatus: () => req('/public/status'),
  setStatus: (body) => req('/status', { method: 'PUT', body: JSON.stringify(body) }),
  setProfil: (body) => req('/profil', { method: 'PUT', body: JSON.stringify(body) }),
  createEvent: (body) => req('/events', { method: 'POST', body: JSON.stringify(body) }),
  updateEvent: (id, body) => req('/events/' + id, { method: 'PUT', body: JSON.stringify(body) }),
  deleteEvent: (id) => req('/events/' + id, { method: 'DELETE' }),
  addStock: (body) => req('/medicines/stock', { method: 'POST', body: JSON.stringify(body) }),
  createVisit: (body) => req('/visits', { method: 'POST', body: JSON.stringify(body) }),
}
