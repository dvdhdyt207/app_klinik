# Dokumentasi — Klinik Bidan Pit (versi Web: Go + Vue)

Recreate dari versi lama (Node/Express + Expo React Native) menjadi **website penuh**:
**backend Go** + **frontend Vue 3**, dengan **fungsionalitas sama persis**. Menyelesaikan
masalah Expo Go SDK 57 (tidak perlu HP/native lagi — cukup browser).

> Versi lama tetap ada sebagai referensi: `backend/` (Node) & `mobile/` (Expo).
> Versi web baru: `server/` (Go) & `web/` (Vue). Lihat juga `DOKUMENTASI.md` (versi lama).

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

## API (base `http://<host>:4000`) — identik dengan versi Node

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/health` | Cek server hidup |
| GET | `/api/state` | Snapshot penuh (status, events, medicines, visits) |
| GET | `/api/catalog` | Katalog obat + aturan satuan |
| PUT | `/api/status` | Ubah status `{bidanHadir, awayNote, awayUntil}` |
| POST/PUT/DELETE | `/api/events[/:id]` | CRUD agenda |
| POST | `/api/medicines/stock` | Tambah stok / buat obat baru `{id?, name, cat, amount}` |
| POST | `/api/visits` | Catat kunjungan + kurangi stok (transaksi) |
| GET | `/api/public/status` | Payload halaman pasien |
| GET | `/` , `/app`, `/assets/*` | Halaman web (SPA) |

## Cara menjalankan

Prasyarat: **Go** (`go version`; jika "not recognized", tambahkan `C:\Program Files\Go\bin` ke PATH), **Node**, **MySQL** hidup.

**Sekali saja — siapkan DB:**
```powershell
cd server
Get-Content migrations/schema.sql | & "C:\mysql\bin\mysql.exe" -u root -p<PASSWORD_ANDA>
go run ./cmd/seed        # isi data contoh (idempotent)
```

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
