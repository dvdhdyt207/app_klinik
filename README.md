# Klinik Bidan Pit

Aplikasi manajemen klinik untuk bidan praktik mandiri. **Website penuh** (responsif — laptop & HP):

- **Bidan App** (`/app`) — privat: catat kunjungan, kelola stok obat, rekap, atur status kehadiran + jadwal.
- **Halaman Pasien** (`/`) — publik, read-only: status bidan (di klinik / keluar + perkiraan kembali) & agenda mendatang.

## Stack

| Lapisan | Teknologi |
|---|---|
| Backend/API | **Go** (`net/http`, `database/sql` + `go-sql-driver/mysql`) |
| Frontend | **Vue 3 + Vite** (vue-router, Pinia) |
| Database | **MySQL 8** (`app_klinik`) |

## Struktur

```
server/   # Backend Go (API + serve build Vue)
web/      # Frontend Vue 3
```

## Jalankan (development)

Prasyarat: Go, Node, MySQL hidup.

```powershell
# sekali — siapkan DB
cd server
Get-Content migrations/schema.sql | & "C:\mysql\bin\mysql.exe" -u root -p<PASSWORD_ANDA>
go run ./cmd/seed

# jalankan (2 terminal)
cd server ; go run .     # API  :4000
cd web ; npm run dev     # UI   :5173  (buka ini)
```

## Produksi (satu binary)

```powershell
cd web ; npm run build            # -> web/dist
cd ../server ; go build -o server.exe . ; ./server.exe   # serve semua di :4000
```

Detail lengkap: **[DOKUMENTASI-WEB.md](DOKUMENTASI-WEB.md)**.

> Kredensial DB ada di `server/.env` (lihat `server/.env.example`) — tidak di-commit.
> Ganti nomor WhatsApp di `web/src/views/Patient.vue`.
