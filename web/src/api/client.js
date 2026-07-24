// API client. URL relatif '/api/...' — di dev di-proxy Vite ke :4000,
// di produksi di-serve origin yang sama oleh server Go.
async function req(path, options) {
  const res = await fetch('/api' + path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    let msg = res.status + ''
    try { const j = await res.json(); msg = j.error || msg } catch {}
    throw new Error(msg)
  }
  if (res.status === 204) return null
  return res.json()
}

export const api = {
  getState: () => req('/state'),
  publicStatus: () => req('/public/status'),
  setStatus: (body) => req('/status', { method: 'PUT', body: JSON.stringify(body) }),
  createEvent: (body) => req('/events', { method: 'POST', body: JSON.stringify(body) }),
  updateEvent: (id, body) => req('/events/' + id, { method: 'PUT', body: JSON.stringify(body) }),
  deleteEvent: (id) => req('/events/' + id, { method: 'DELETE' }),
  addStock: (body) => req('/medicines/stock', { method: 'POST', body: JSON.stringify(body) }),
  createVisit: (body) => req('/visits', { method: 'POST', body: JSON.stringify(body) }),
}
