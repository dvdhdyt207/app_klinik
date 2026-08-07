---
name: klinik-ui
description: "Aturan UI wajib untuk app_klinik (Klinik Bidan Pit) — peruntukan aplikasi, palet teal yang sudah disetujui, larangan hex di luar tokens.css, dan cara membuktikan tampilan dengan potret Playwright. Pakai SETIAP KALI menyentuh berkas di web/src: .vue, tokens.css, atau apa pun yang mengubah tampilan halaman pasien, Bidan App, atau halaman masuk. Juga pakai saat memakai skill ui-ux-pro-max di repo ini, karena saran umumnya perlu disaring lebih dulu."
---

# Aturan UI app_klinik

## Untuk apa aplikasi ini — baca ini dulu

Dua pemakai, dua pekerjaan yang sama sekali berbeda. Rancangan yang tidak membedakan keduanya akan salah.

**Halaman pasien** (`klinik.catering-app.com`, `web/src/views/Patient.vue`) — dibuka pasien dan keluarganya untuk **satu pertanyaan**: bidannya ada atau tidak, dan kalau tidak, kapan kembali. Sering dibuka terburu-buru, di HP murah, kadang tengah malam, kadang saat ada yang mau melahirkan. Jawabannya harus terbaca dalam sekali lihat tanpa menggulir. Halaman ini bukan brosur klinik: jangan tambahkan bagian pemasaran, testimoni, atau ajakan bertindak.

**Bidan App** (`bidan.catering-app.com`, `web/src/views/BidanApp.vue` + `components/screens/`) — rekam medis elektronik untuk **satu bidan praktik mandiri**. Dipakai sambil bekerja, sering satu tangan, sering berdiri, kadang jam dua pagi. Isinya rekam medis pasien orang lain. Yang menentukan mutu: **seberapa cepat satu kunjungan bisa dicatat** — bukan seberapa indah layarnya.

Aplikasi ini **tidak mengelola stok obat**, dan itu keputusan user (7 Agustus 2026): layar Stok, rekap obat, status "stok menipis", dan kategori obat dilepas semuanya. Jangan menghidupkannya kembali, dan jangan menambahkan angka persediaan ke rancangan mana pun. Daftar obat dan daftar pasien **dibentuk dari riwayat kunjungan**: nama yang sekali dicatat muncul lagi sebagai saran berikutnya.

**Halaman masuk** (`web/src/views/Login.vue`) — pintu bidan. Tidak boleh ada apa pun yang bisa menghalangi ia masuk ke rekam medisnya sendiri; hiasan yang gagal dimuat harus gagal diam-diam.

Sasaran perangkat, urut kepentingan: **HP (390px) → laptop (1180px) → monitor lebar (1440px+)**. Bidan memakai HP; monitor lebar tetap harus terisi masuk akal, bukan kolom HP yang ditaruh di tengah.

## Palet — sudah diputuskan, jangan ditawar ulang

Ditetapkan user 6 Agustus 2026 setelah menolak biru korporat lama. **Jangan mengusulkan palet lain**, termasuk bila skill lain menyarankannya.

- teal `#12806a` = identitas klinik **dan** keadaan "bidan hadir"
- amber `#b8791c` = keadaan "bidan tidak ada" — perlu perhatian, bukan bahaya
- merah = kerusakan dan penghapusan **saja**; jangan dipakai untuk "tutup" atau "habis"

Sumber tunggalnya `web/src/styles/tokens.css`.

## Larangan yang paling sering dilanggar

**Jangan menulis nilai warna di luar `tokens.css` — termasuk di berkas `.js`.** Nilai yang berakhir di `:style` inline pun ditulis `var(--token)`; var() ikut diselesaikan di atribut style. Pelanggaran ini pernah membuat penggantian palet tidak berpengaruh apa-apa pada titik agenda, angka stok, dan tint kategori obat, karena hex-nya hidup di `stores/klinik.js` dan `lib/catalog.js` (berkas itu kini `lib/obat.js`, dan tint kategorinya sudah hilang bersama stok). Kalau ada warna yang "tidak mau ikut berubah", **cari hex di berkas `.js` dulu, bukan di `.css`**.

**Jangan menulis mati nama klinik.** Bidan bisa mengubahnya lewat Setelan → Profil klinik. Ambil dari `k.profil.clinic` / `k.status.clinic`; monogram lewat `monogram()` di `web/src/lib/format.js`.

## Cara membuktikan tampilan — jangan menebak

Wajib dilihat, bukan dibayangkan. Potret dengan Playwright MCP di **1440 / 1180 / 390**, baca PNG-nya, perbaiki, potret ulang.

Menyalakan tanpa MySQL: buat `web/vite.mock.config.js` sementara (salinan `vite.config.js` + satu plugin `configureServer` yang menjawab `/api/public/status`, `/api/auth/me`, `/api/state` dari berkas JSON di scratchpad), jalankan `npx vite --config vite.mock.config.js`, **hapus berkasnya setelah selesai**. Confignya wajib di dalam `web/` — Node me-resolve `vite` dari letak berkas config, bukan dari `root`.

**Keadaan yang wajib dipotret, bukan hanya yang bagus:**

- bidan hadir **dan** bidan tidak ada
- **profil kosong** (alamat & WhatsApp belum diisi) — ini keadaan produksi sekarang
- daftar kosong: belum ada kunjungan, belum ada agenda, belum ada obat yang pernah dicatat (layar Cari Obat pada klinik baru)
- galat login
- modal/dialog terbuka

**Batas yang sudah diketahui:** lapisan gelap di belakang dialog (`rgba(20,33,28,.45)`) **tidak tertangkap** potret Chromium headless, padahal terbukti tampil pada alpha 1,0. Jangan menaikkan alpha supaya "terlihat di potret" — itu memperbaiki gambar, bukan aplikasi. Minta user melihatnya di browser sungguhan.

## Memakai skill ui-ux-pro-max di repo ini

Berguna untuk aksesibilitas, ukuran target sentuh, kontras, keadaan fokus, dan pedoman Vue. **Saring dulu hal-hal ini:**

- Ia sering mengembalikan pola **halaman pemasaran** (hero, CTA tunggal, "Minimal Single Column"). Bidan App dan halaman pasien **bukan landing page** — abaikan.
- Ia menyarankan palet sendiri (mis. biru `#3B82F6`). **Abaikan**, palet sudah diputuskan di atas.
- Saran gaya seperti glassmorphism, dark-mode OLED, dan data-dense analytics datang dari kategori dashboard analitik. Aplikasi ini pencatatan, bukan analitik.

Yang layak diambil apa adanya: kontras teks minimal 4.5:1, target sentuh 44×44px, cincin fokus terlihat, `prefers-reduced-motion`, ikon SVG bukan emoji, transisi 150–300ms, label form yang benar, dan elemen semantik (`<button>`, `<nav>`, `<main>`) alih-alih `div` untuk segalanya.

## Jebakan perkakas

**PowerShell 5.1 merusak UTF-8** pada penggantian teks massal: `Get-Content -Raw` + `Set-Content -Encoding utf8` mengubah `·` jadi `Â·`. Pakai `[System.IO.File]::ReadAllText/WriteAllText` dengan `New-Object System.Text.UTF8Encoding($false)`.

## Kejujuran soal hasil

Kalau yang dikerjakan ternyata hanya mengganti warna dan merapikan, **sebut begitu sejak awal** — jangan menyebutnya perombakan atau rekonstruksi. Sebutkan eksplisit apa yang berubah secara struktur dan apa yang tidak. User pernah menagih ini, dan ia benar.
