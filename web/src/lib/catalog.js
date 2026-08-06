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

// Daftar kategori untuk pilihan di layar ubah obat. Urutannya menentukan
// urutan tampil.
export const KATEGORI = Object.keys(CAT)

// TIDAK ADA daftar obat di berkas ini. Dulu ada 14 obat contoh dari design
// handoff, dan itu keliru: daftar obat adalah data milik klinik, bukan aturan
// hitungan. Obat karangan ikut muncul di pencarian dan tidak bisa dihapus
// tanpa deploy ulang.
//
// Sumbernya sekarang tabel `medicines` di database — apa yang benar-benar
// distok bidan. Lihat PickModal.vue.

// Tint kategori untuk avatar/kotak (bg + warna teks). var(), bukan hex — lihat
// alasannya di komentar COL pada stores/klinik.js.
export const catTint = (cat) =>
  cat === 'Sirup' ? { bg: 'var(--sirup-tint)', ink: 'var(--sirup-ink)' }
  : cat === 'Sachet' ? { bg: 'var(--sachet-tint)', ink: 'var(--sachet-ink)' }
  : { bg: 'var(--tablet-tint)', ink: 'var(--tablet-ink)' }
