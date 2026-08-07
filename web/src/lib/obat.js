// Satuan obat yang ikut tercatat pada rekam medis.
//
// Berkas ini menggantikan lib/catalog.js, yang memuat aturan persediaan:
// kategori obat (Tablet/Sirup/Sachet), konversi kemasan (box 100 → butir), dan
// ambang "stok menipis". Pengelolaan stok sudah dilepas — aplikasi ini tidak
// lagi menghitung sisa obat di lemari — jadi yang tersisa hanya satu hal:
// label satuan, supaya "10" pada catatan kunjungan punya arti.
//
// Kategori sengaja TIDAK diganti nama jadi satuan lalu dipertahankan sebagai
// konsep tersendiri. Setelah stok hilang, kategori tidak menentukan apa pun
// selain satuannya sendiri — dua nama untuk satu hal yang sama.

export const SATUAN = ['butir', 'botol', 'sachet']
export const SATUAN_BAWAAN = 'butir'

// Satuan dari riwayat dipakai apa adanya, hanya dirapikan. Baris kunjungan lama
// ditulis sebelum berkas ini ada; rekam medis harus tetap terbaca persis
// seperti yang dicatat waktu itu, jadi nilai di luar daftar pun tidak diubah.
export const satuanAtau = (u) => {
  const s = (u || '').trim()
  return s || SATUAN_BAWAAN
}
