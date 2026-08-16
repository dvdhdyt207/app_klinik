# Klinik Bidan Pit — instruksi tetap

## Ini proyek apa

Aplikasi rekam medis untuk **satu bidan praktik mandiri**. Stack: **Go** (`net/http`, `database/sql`) + **Vue 3 + Vite** + **MySQL 8**. Satu biner Go menyajikan API sekaligus hasil build Vue.

**Ini BUKAN proyek Laravel/PHP.** Beberapa skill `best-setup` yang terpasang berbicara soal `app/Controller`, Eloquent, Deptrac, dan Pest — semuanya untuk stack lain. **Jangan menerapkannya di sini.** Yang berlaku: `klinik-ui` (wajib, lihat bawah) dan penilaian biasa.

Dua permukaan dengan dua tingkat kepercayaan:

- **Halaman Pasien** (`/`, `web/src/views/Patient.vue`) — publik, tanpa login. Menjawab **satu pertanyaan**: bidannya ada atau tidak, dan kalau tidak, kapan kembali. Bukan brosur klinik.
- **Bidan App** (`/app`) — wajib login, memuat **rekam medis pasien orang lain**. Yang menentukan mutunya: seberapa cepat satu kunjungan bisa dicatat.

Penjaga data ada di **server**. Penjaga rute di `web/src/router/index.js` hanya kenyamanan — jangan pernah memindahkan keputusan "boleh lihat atau tidak" ke frontend.

## Menjalankan (WSL / Linux)

README di repo ini masih menulis perintah Windows/PowerShell. Di WSL pakai ini:

```bash
# sekali — siapkan DB
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS app_klinik CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p app_klinik < server/migrations/schema.sql
cd server && cp .env.example .env   # isi DB_PASSWORD dan ADMIN_PASSWORD (min. 10 karakter)
go run ./cmd/seed

# jalankan — dua terminal
cd server && go run .     # API  :4000
cd web    && npm run dev  # UI   :5173  ← buka ini
```

Versi yang dipakai produksi: **Go 1.26**, **Node 22**, **MySQL 8**. Tanpa `ADMIN_PASSWORD` server sengaja menolak jalan.

**Test:** belum ada satu pun (`*_test.go` maupun `*.test.js`). Jalankan test lewat `/uji` — tidak ada hook yang menjalankannya otomatis.

## Yang tidak boleh ditawar

- **Palet teal `#12806a` (hadir/identitas) & amber `#b8791c` (bidan pergi).** Merah hanya untuk kerusakan & penghapusan. Diputuskan user 6 Agustus 2026 — jangan usulkan palet lain.
- **Tidak ada nilai warna di luar `web/src/styles/tokens.css`**, termasuk di berkas `.js`/`.ts`. Kalau ada warna yang "tidak mau ikut berubah", cari hex di berkas JS dulu, bukan di CSS.
- **Aplikasi ini tidak mengelola stok obat.** Dilepas atas permintaan bidan 7 Agustus 2026. Daftar obat & pasien dibentuk dari riwayat kunjungan. Jangan menghidupkannya kembali.
- **Nama klinik tidak pernah ditulis mati** — ambil dari `k.profil.clinic` / `k.status.clinic`.
- **Jangan pernah membangun di VPS.** Artefak dibangun di CI; server hanya menariknya.
- **Jangan deploy.** Produksi tidak disentuh sampai pemilik repo menyatakan siap.

## Skill

**`klinik-ui` wajib dipakai setiap kali menyentuh `web/src`** — apa pun yang mengubah tampilan halaman pasien, Bidan App, atau halaman masuk. Ia memuat aturan palet, peruntukan tiap halaman, dan cara membuktikan tampilan dengan potret Playwright di 1440/1180/390.

Tampilan **wajib dilihat, bukan dibayangkan.** Termasuk keadaan yang mudah terlupa: bidan pergi, profil kosong, daftar kosong, galat login, modal terbuka.

## Penamaan

- **Identifier Go & TypeScript: bahasa Inggris.** (`Visit`, `createVisit`, `mysqlRepo`)
- **Komentar dan pesan commit: bahasa Indonesia**, menjelaskan *kenapa*, bukan *apa*.
- Bentuk JSON di kabel: camelCase Inggris — sudah publik, jangan diubah tanpa alasan.

Kode lama masih memakai identifier Indonesia (`buatSesi`, `wajibLogin`, `BatasLaju`). Itu diterjemahkan saat berkasnya memang sedang ditulis ulang — bukan disapu sekaligus dalam commit tersendiri.

## Rewrite yang sedang berjalan

Rencana rewrite ke modular monolith ada di **GitHub issue #1**, dan itu sumber kebenarannya — bukan berkas di repo. Baca lebih dulu sebelum mengerjakan tahap mana pun; ia memuat nomor temuan (P1-P7 performa, K1-K5 keamanan, N1-N4 proxy) yang dirujuk tiap tahapnya, batas yang tidak boleh ditawar, dan aturan otonomi per tahap.

`main` tetap milik versi produksi yang sekarang berjalan. Rewrite tinggal di cabangnya sendiri.
