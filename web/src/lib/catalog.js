// Aturan bisnis konversi satuan & katalog obat (dari README design handoff).
// qty obat SELALU disimpan dalam base unit. Mirror dari server Go (catalog.go).

export const CAT = {
  Tablet: { base: 'butir',  units: [ { label: 'Box (100)', mult: 100 }, { label: 'Strip (10)', mult: 10 }, { label: 'Butir', mult: 1 } ], thr: 20 },
  Sirup:  { base: 'botol',  units: [ { label: 'Botol', mult: 1 } ], thr: 5 },
  Sachet: { base: 'sachet', units: [ { label: 'Box (100)', mult: 100 }, { label: 'Sachet', mult: 1 } ], thr: 20 },
}

export const baseUnit = (cat) => (CAT[cat] || CAT.Tablet).base
export const threshold = (cat) => (CAT[cat] || CAT.Tablet).thr

// danger (merah) bila qty <= threshold*0.5, selain itu warning (amber)
export const isLow = (cat, qty) => qty <= threshold(cat)
export const isDanger = (cat, qty) => qty <= threshold(cat) * 0.5

// Katalog master obat untuk fitur "Cari Obat"
export const CATALOG = [
  { name: 'Paracetamol 500mg', cat: 'Tablet' },
  { name: 'Paracetamol sirup', cat: 'Sirup' },
  { name: 'Amoxicillin 500mg', cat: 'Tablet' },
  { name: 'Amoxicillin sirup', cat: 'Sirup' },
  { name: 'Ibuprofen 400mg', cat: 'Tablet' },
  { name: 'Asam Mefenamat 500mg', cat: 'Tablet' },
  { name: 'Cetirizine 10mg', cat: 'Tablet' },
  { name: 'CTM 4mg', cat: 'Tablet' },
  { name: 'Antasida', cat: 'Tablet' },
  { name: 'Vitamin B Complex', cat: 'Tablet' },
  { name: 'Vitamin C 500mg', cat: 'Tablet' },
  { name: 'Domperidone 10mg', cat: 'Tablet' },
  { name: 'Dexamethasone 0.5mg', cat: 'Tablet' },
  { name: 'Oralit', cat: 'Sachet' },
]

// Tint kategori untuk avatar/kotak (bg + warna teks). var(), bukan hex — lihat
// alasannya di komentar COL pada stores/klinik.js.
export const catTint = (cat) =>
  cat === 'Sirup' ? { bg: 'var(--sirup-tint)', ink: 'var(--sirup-ink)' }
  : cat === 'Sachet' ? { bg: 'var(--sachet-tint)', ink: 'var(--sachet-ink)' }
  : { bg: 'var(--tablet-tint)', ink: 'var(--tablet-ink)' }
