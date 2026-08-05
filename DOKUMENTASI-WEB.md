# Dokumentasi — Klinik Bidan Pit (versi Web: Go + Vue)

Recreate dari versi lama (Node/Express + Expo React Native) menjadi **website penuh**:
**backend Go** + **frontend Vue 3**, dengan **fungsionalitas sama persis**. Menyelesaikan
masalah Expo Go SDK 57 (tidak perlu HP/native lagi — cukup browser).

> Stack lama (Node/Express + Expo React Native) sudah **dihapus** — digantikan
> `server/` (Go) & `web/` (Vue).

## Stack

| Lapisan | Teknologi |
|---|---|
| Backend/API | **Go** (`net/http` stdlib, `database/sql` + `go-sql-driver/mysql`, `godotenv`) |
| Frontend | **Vue 3 + Vite** (`<script setup>`, vue-router, Pinia) |
| Database | **MySQL 8** (`app_klinik` @ 127.0.0.1:3306) — skema tidak berubah |
| Font | Plus Jakarta Sans di-bundle via `@fontsource` (offline-friendly) |

## Arsitektur

```
Browser (Vue SPA)
  ├─ '/'     Halaman Pasien   (publik, read-only, polling 10 dtk)
  └─ '/app'  Bidan App        (privat: Beranda/Kunjungan/Stok/Rekap + modal)
        │  fetch /api/... (URL relatif, satu origin)
        ▼
  Server Go (:4000)  ── database/sql ──►  MySQL app_klinik
   • REST API /api/*
   • men-serve build Vue (web/dist) + SPA fallback
```

Dev: Vite dev server (`:5173`, HMR) mem-**proxy** `/api` ke Go (`:4000`).
Produksi: `go build` → satu `server.exe` yang men-serve API **dan** halaman web di origin sama.

## Struktur folder

```
server/                      # Backend Go
├── main.go                  # wiring: API + serve SPA (web/dist) + CORS
├── go.mod
├── .env                     # kredensial DB & PORT (JANGAN commit)
├── migrations/schema.sql    # skema 5 tabel (idempotent)
├── cmd/seed/main.go         # `go run ./cmd/seed` — isi data contoh
└── internal/
    ├── catalog/catalog.go   # aturan satuan + katalog obat (mirror web)
    ├── models/models.go     # struct JSON (camelCase, sama dgn versi Node)
    ├── db/db.go             # pool koneksi MySQL dari .env
    └── api/                 # api.go (routing+helper), store.go (query), handlers.go (endpoint)

web/                         # Frontend Vue 3
├── vite.config.js           # proxy /api -> :4000
└── src/
    ├── main.js, App.vue, router/index.js
    ├── styles/tokens.css     # design tokens (CSS variables) + primitive
    ├── lib/{catalog,format}.js   # logika bisnis JS murni (dipakai ulang dari mobile)
    ├── api/client.js         # fetch client (URL relatif /api)
    ├── stores/klinik.js      # Pinia — port dari mobile/src/store.js (useKlinik)
    ├── views/{Patient,BidanApp}.vue
    └── components/
        ├── ui/{Toggle,Stepper,BottomSheet,ModalHeader}.vue
        ├── screens/{Beranda,Kunjungan,Stok,Rekap}.vue
        ├── modals/{VisitModal,PickModal,AddStockSheet,AwaySheet,JadwalModal,EventFormModal}.vue
        └── TabBar.vue
```

## API (base `http://<host>:4000`)

Kolom **Sesi** menandai endpoint yang menuntut cookie sesi sah. Yang terbuka
hanya tiga: health, halaman pasien, dan pintu masuknya sendiri.

| Method | Endpoint | Sesi | Fungsi |
|---|---|:--:|---|
| GET | `/api/health` | — | Cek server hidup |
| GET | `/api/public/status` | — | Payload halaman pasien (status bidan + agenda; **tanpa data pasien**) |
| POST | `/api/auth/login` | — | Masuk `{username, password}` → memasang cookie `klinik_sesi` |
| POST | `/api/auth/logout` | — | Cabut sesi & hapus cookie |
| GET | `/api/auth/me` | ✓ | Siapa yang sedang masuk |
| GET | `/api/state` | ✓ | Snapshot penuh (status, events, medicines, **visits**) |
| GET | `/api/catalog` | ✓ | Katalog obat + aturan satuan |
| PUT | `/api/status` | ✓ | Ubah status `{bidanHadir, awayNote, awayUntil}` |
| POST/PUT/DELETE | `/api/events[/:id]` | ✓ | CRUD agenda |
| POST | `/api/medicines/stock` | ✓ | Tambah stok / buat obat baru `{id?, name, cat, amount}` |
| POST | `/api/visits` | ✓ | Catat kunjungan + kurangi stok (transaksi) |
| GET | `/` , `/masuk`, `/app`, `/assets/*` | — | Halaman web (SPA) |

Halaman SPA-nya sendiri memang tersaji tanpa sesi — yang dijaga adalah
**datanya**. Membuka `/app` tanpa masuk hanya menghasilkan 401 dari API, lalu
penjaga rute mengantar ke `/masuk`.

## Autentikasi

Satu akun bidan. Sesi memakai **token acak 256-bit** (bukan JWT) di cookie
`klinik_sesi`: `HttpOnly`, `SameSite=Lax`, dan `Secure` otomatis saat
`APP_ENV=production`.

- Database hanya menyimpan **SHA-256 dari token**, bukan tokennya — isi database
  yang bocor tidak bisa dipakai untuk masuk.
- Tabel `sessions` membuat sesi **bisa dicabut**; itu yang tidak bisa dilakukan
  token yang berdiri sendiri seperti JWT polos.
- Umur sesi 14 hari, diperpanjang selama dipakai (paling sering sekali per jam);
  sesi mati disapu berkala.
- Kata sandi disimpan sebagai hash **bcrypt**.
- Login dibatasi **per nama pengguna** (5 gagal beruntun, lalu 1 percobaan tiap
  20 detik). Pembatas per IP saja tidak cukup: serangan dari banyak IP tetap
  lolos selama tiap IP pelan. Hanya kegagalan yang memotong jatah, jadi bidan
  yang mengetik benar tidak pernah terkena.
- Nama pengguna salah dan sandi salah memberi pesan **yang sama** dan memakan
  waktu yang sama — membedakannya akan memberitahu penebak bahwa separuh
  tebakannya sudah benar.

**Akun pertama** dibuat dari `ADMIN_USERNAME` + `ADMIN_PASSWORD` saat tabel
`users` masih kosong. Setelah akun ada, nilai itu **diabaikan sepenuhnya** —
sandi yang diganti lewat aplikasi tidak akan tertimpa isi `.env`. Bila belum ada
akun dan `ADMIN_PASSWORD` kosong (atau kurang dari 10 karakter), server
**menolak jalan**: lebih baik mati daripada hidup tanpa penjaga.

## Pengaturan lingkungan (`server/.env`)

Contoh lengkap ada di `server/.env.example`. Yang berpengaruh pada keamanan:

| Variabel | Bawaan | Catatan |
|---|---|---|
| `APP_ENV` | `development` | `production` menyalakan cookie `Secure` + HSTS. **Hanya** bila situs sudah lewat https — di http, peramban menolak cookie `Secure` dan login gagal diam-diam |
| `BIND_ADDR` | `127.0.0.1` | Di produksi biarkan loopback; hanya reverse proxy yang menghadap internet |
| `CORS_ORIGINS` | kosong | Daftar tertutup dipisah koma. Di produksi kosong (satu origin). Saat dev dengan Vite: `http://localhost:5173` |
| `ADMIN_USERNAME` / `ADMIN_PASSWORD` | `bidan` / — | Hanya untuk akun pertama; sebaiknya dikosongkan setelah terbentuk |

Pengerasan lain yang sudah terpasang: header keamanan, batas laju per IP, batas
badan request 1 MB, timeout baca/tulis/idle, shutdown yang rapi, pesan galat
yang tidak lagi membocorkan isi database, dan `multiStatements` yang **mati** di
koneksi aplikasi (hanya perkakas schema/seed yang memakainya).

`cmd/seed` **menolak jalan saat `APP_ENV=production`** — ia meng-`TRUNCATE`
tabel kunjungan, obat, dan agenda.

## Cara menjalankan

Prasyarat: **Go** (`go version`; jika "not recognized", tambahkan `C:\Program Files\Go\bin` ke PATH), **Node**, **MySQL** hidup.

**Sekali saja — siapkan DB & akun:**
```powershell
cd server
& "C:\mysql\bin\mysql.exe" -u root -p<PASSWORD_ANDA> -e "CREATE DATABASE IF NOT EXISTS app_klinik CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
Get-Content migrations/schema.sql | & "C:\mysql\bin\mysql.exe" -u root -p<PASSWORD_ANDA> app_klinik
go run ./cmd/seed        # isi data contoh (idempotent; TIDAK menyentuh users/sessions)
```

Lalu salin `.env.example` ke `.env` dan isi minimal `DB_*`, `ADMIN_PASSWORD`
(≥10 karakter), serta `CORS_ORIGINS=http://localhost:5173` bila memakai Vite.
Akun bidan terbentuk otomatis saat server pertama kali hidup.

**Mode development (rekomendasi saat ngoding — ada HMR):**
```powershell
# Terminal 1 — backend Go
cd server ; go run .                 # http://localhost:4000

# Terminal 2 — frontend Vue (HMR)
cd web ; npm run dev                 # buka http://localhost:5173
```
Buka `http://localhost:5173/` (pasien) atau `http://localhost:5173/app` (Bidan App).

**Mode produksi (satu origin / satu binary):**
```powershell
cd web ; npm run build               # hasil ke web/dist
cd ../server ; go run .              # atau: go build -o server.exe . ; ./server.exe
```
Buka `http://localhost:4000/` (pasien) & `http://localhost:4000/app` (Bidan App).

**Segarkan data contoh** (agar "hari ini" akurat): `cd server ; go run ./cmd/seed`.

## Catatan
- Bentuk JSON semua endpoint dijaga sama persis dengan versi Node (camelCase), jadi paritas penuh.
- Aturan satuan ada di dua tempat (mirror): `server/internal/catalog` (validasi backend) & `web/src/lib/catalog.js` (UI) — sama seperti desain lama (mobile me-mirror backend).
- Untuk install di HP sebagai PWA nanti: perlu HTTPS (mis. Cloudflare Tunnel) — belum dikerjakan, saat ini web-only lokal.
- Ganti nomor WhatsApp asli di `web/src/views/Patient.vue` (`wa.me/62...`).
```
