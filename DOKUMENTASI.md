# Dokumentasi Proyek — Klinik Bidan Pit

Aplikasi manajemen klinik untuk **Bidan Pit** (bidan praktik mandiri). Terdiri dari **dua bagian yang tersinkron**:

1. **Bidan App** — aplikasi Android (privat) untuk bidan: catat kunjungan, kelola stok obat, rekap, dan atur status kehadiran + jadwal.
2. **Halaman Web Pasien** — halaman publik read-only: pasien bisa lihat apakah bidan sedang di klinik, perkiraan kembali, dan agenda mendatang.

Dibangun dari design handoff di [`README.md`](README.md) (spesifikasi desain hi-fi, Bahasa Indonesia, tokens warna/tipografi).

> Catatan: `README.md` = spesifikasi desain. Dokumen ini (`DOKUMENTASI.md`) = dokumentasi sistem yang dibangun.

---

## 1. Arsitektur

```
┌──────────────────────────┐     HTTP/REST     ┌─────────────────────────┐     ┌──────────┐
│  Bidan App (Expo RN)      │ ───────────────►  │  Backend Node/Express   │ ──► │  MySQL   │
│  di HP Android (Expo Go)  │ ◄───────────────  │  (port 4000)            │     │ app_klinik│
└──────────────────────────┘                   │                         │     └──────────┘
┌──────────────────────────┐   HTTP + polling   │  juga men-serve         │
│  Halaman Web Pasien       │ ◄───────────────  │  halaman web pasien     │
│  (browser)                │                   └─────────────────────────┘
└──────────────────────────┘
```

- Aplikasi native **tidak** konek MySQL langsung — semua lewat API backend.
- Sinkronisasi live memakai **DB + polling** (mengganti localStorage/BroadcastChannel pada prototipe desain). Bentuk payload status tetap sesuai README.

## 2. Stack Teknologi

| Lapisan | Teknologi |
|---|---|
| Aplikasi mobile | **Expo (React Native)** SDK 57, React 19 |
| Backend/API | **Node.js** v22 + **Express** + **mysql2** |
| Database | **MySQL 8** (`app_klinik` @ 127.0.0.1:3306) |
| Halaman pasien | HTML/CSS/JS statis (di-serve backend) |

## 3. Struktur Folder

```
app_klinik/
├── README.md                    # Spesifikasi desain (design handoff)
├── DOKUMENTASI.md               # File ini
├── Klinik Bidan Pit App.dc.html # Prototipe desain Bidan App (referensi)
├── Status Bidan Web.dc.html     # Prototipe desain halaman pasien (referensi)
├── .claude/launch.json          # Konfigurasi preview (backend :4000, mobile-web :8081)
│
├── backend/                     # API + halaman pasien
│   ├── .env                     # Kredensial DB & PORT (JANGAN commit)
│   ├── db/
│   │   ├── schema.sql           # Skema 5 tabel
│   │   └── seed.js              # Data contoh (npm run seed)
│   ├── src/
│   │   ├── db.js               # Pool koneksi MySQL
│   │   ├── catalog.js          # Aturan konversi satuan + katalog obat
│   │   └── server.js           # Express: semua endpoint API
│   └── public/index.html        # Halaman web pasien
│
└── mobile/                      # Aplikasi Bidan (Expo)
    ├── App.js                   # Semua layar + modal + tab bar
    └── src/
        ├── api.js              # API client (auto-deteksi host)
        ├── theme.js            # Design tokens (warna, shadow)
        ├── catalog.js          # Aturan satuan (mirror backend)
        ├── format.js           # Format tanggal/waktu Indonesia
        ├── typography.js       # Text kustom (font Plus Jakarta Sans)
        ├── ui.js               # Komponen: Toggle, Stepper, BottomSheet, dll.
        ├── DateTimeField.js    # Date/time picker lintas-platform
        └── store.js            # useKlinik: state + actions (mirror prototipe)
```

## 4. Database — skema `app_klinik`

| Tabel | Isi |
|---|---|
| `clinic_status` | Baris tunggal (id=1): status hadir, keterangan keluar, perkiraan kembali, waktu update |
| `events` | Agenda/jadwal: judul, seharian?, mulai, selesai (epoch ms) |
| `medicines` | Stok obat: nama, kategori (Tablet/Sirup/Sachet), qty (dalam base unit) |
| `visits` | Kunjungan: nama, umur, gejala, waktu |
| `visit_items` | Obat yang diberikan tiap kunjungan (relasi ke `visits`) |

**Aturan konversi satuan** (penting): qty selalu dalam *base unit*.
- **Tablet** → `butir`. Tambah via Box(100)/Strip(10)/Butir. Ambang menipis ≤ 20.
- **Sirup** → `botol`. Tambah via Botol(1). Ambang ≤ 5.
- **Sachet** → `sachet`. Tambah via Box(100)/Sachet(1). Ambang ≤ 20.
- **Danger** (merah) bila qty ≤ ambang×0.5, selain itu **warning** (kuning).

## 5. API Backend (base `http://<host>:4000`)

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/api/health` | Cek server hidup |
| GET | `/api/state` | Snapshot penuh (status, events, medicines, visits) |
| GET | `/api/catalog` | Katalog obat + aturan satuan |
| PUT | `/api/status` | Ubah status bidan `{bidanHadir, awayNote, awayUntil}` |
| POST | `/api/events` | Tambah agenda |
| PUT | `/api/events/:id` | Ubah agenda |
| DELETE | `/api/events/:id` | Hapus agenda |
| POST | `/api/medicines/stock` | Tambah stok / buat obat baru `{id?, name, cat, amount}` |
| POST | `/api/visits` | Catat kunjungan + kurangi stok (transaksi) |
| GET | `/api/public/status` | Payload untuk halaman pasien |
| GET | `/` | Halaman web pasien |

## 6. Konfigurasi

**Backend** — [`backend/.env`](backend/.env):
```
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=<PASSWORD_ANDA>
DB_NAME=app_klinik
PORT=4000
```

**Mobile** — [`mobile/src/api.js`](mobile/src/api.js): host backend **dideteksi otomatis** dari server Expo (`Constants.expoConfig.hostUri`), jadi tak perlu ganti IP manual saat jaringan berubah. Untuk memaksa manual (mis. saat build APK standalone), isi `MANUAL_HOST`.

## 7. Cara Menjalankan (development)

**Persiapan awal (sekali saja — sudah dilakukan):**
```powershell
cd backend; npm install; # buat DB:
Get-Content db/schema.sql | & "C:\mysql\bin\mysql.exe" -u root -p<PASSWORD_ANDA>
npm run seed
cd ../mobile; npm install
```

**Menjalankan:**
```powershell
# Terminal 1 — backend
cd "C:\Web\PROJECT APP\app_klinik\backend"
npm start                      # http://localhost:4000

# Terminal 2 — aplikasi bidan
cd "C:\Web\PROJECT APP\app_klinik\mobile"
npx expo start                 # tekan 'w' untuk web, atau scan QR dg Expo Go
```

**Halaman pasien:** buka `http://localhost:4000` (PC) atau `http://<IP-PC>:4000` (HP).

**Menjalankan di HP Android (Expo Go):**
1. HP & PC harus satu jaringan. Jika PC pakai LAN/Ethernet, nyalakan **Mobile Hotspot** di PC (share dari Ethernet) dan sambungkan HP ke situ.
2. Izinkan firewall (PowerShell admin, sekali):
   `New-NetFirewallRule -DisplayName "Klinik Bidan Pit" -Direction Inbound -Protocol TCP -LocalPort 4000,8081 -Action Allow`
3. Jika pakai hotspot, paksa IP-nya saat start Expo:
   `$env:REACT_NATIVE_PACKAGER_HOSTNAME="192.168.137.1"; npx expo start`
4. Scan QR dengan Expo Go.

**Segarkan data contoh** (agar "hari ini" akurat): `cd backend; npm run seed`.

## 8. Status & Masalah Diketahui

- ✅ Backend + DB + halaman pasien: **selesai & terverifikasi**.
- ✅ Aplikasi Bidan (semua layar/modal): **selesai & terverifikasi** via Expo web.
- ⚠️ **Expo Go di HP incompatible dengan SDK 57** (versi Expo Go terlalu lama untuk perangkat). Menjalankan lewat Expo Go terhambat → perlu **build APK** (lihat langkah berikut).
- ℹ️ Web menampilkan warning kosmetik `"shadow*" deprecated, use "boxShadow"` — hanya di web, aman di native.
- ℹ️ Setelah `expo install`, **restart Metro** agar modul baru terdeteksi (kalau tidak, muncul error "Unable to resolve").

## 9. Langkah Lanjut (belum dikerjakan)

- **Build APK standalone via EAS Build** (direkomendasikan): pakai ulang seluruh kode, hasil app native yang bisa di-install permanen tanpa Expo Go. Perlu akun Expo gratis; build di cloud (~10–20 mnt). Saat build, set `MANUAL_HOST` di `api.js` ke IP backend (mis. IP hotspot).
  - Alternatif build: Docker (SDK di dalam container, kode tetap lokal) atau install Android SDK lokal.
- **Hosting backend** bila ingin app jalan tanpa PC menyala.
- Penyempurnaan opsional: autentikasi/PIN untuk Bidan App, ikon aplikasi & splash screen, WhatsApp deep link asli di halaman pasien.
