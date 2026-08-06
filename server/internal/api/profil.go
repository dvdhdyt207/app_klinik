package api

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"klinikbidanpit/internal/models"
)

// Batas panjang mengikuti lebar kolom di schema.sql. Diperiksa dalam satuan
// rune, bukan byte: nama dan alamat berbahasa Indonesia bisa memuat karakter
// multi-byte, dan menolak berdasarkan len() akan menolak teks yang sebenarnya
// muat.
const (
	klinikMaks = 120
	alamatMaks = 255
)

type profilReq struct {
	Clinic   string `json:"clinic"`
	Alamat   string `json:"alamat"`
	WhatsApp string `json:"whatsapp"`
}

// normalWA menyeragamkan penulisan nomor menjadi bentuk yang dipakai tautan
// wa.me: hanya angka, diawali kode negara 62.
//
// Dinormalkan, bukan ditolak, karena yang mengetik adalah bidan — bukan
// programmer — dan `0812…`, `+62 812…`, `0812-3456-789` semuanya cara wajar
// menuliskan nomor yang sama. Yang tidak boleh terjadi justru kebalikannya:
// nomor berawalan 0 diterima apa adanya lalu menghasilkan tautan wa.me yang
// terbuka normal tapi tidak menemukan nomornya — gagal tanpa pesan apa pun,
// dan baru ketahuan saat ada pasien yang tidak bisa menghubungi klinik.
//
// String kosong sah dan berarti "jangan tampilkan tombol WhatsApp".
func normalWA(raw string) (string, bool) {
	var b strings.Builder
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	n := b.String()
	if n == "" {
		return "", true
	}
	switch {
	case strings.HasPrefix(n, "62"):
		// sudah benar
	case strings.HasPrefix(n, "0"):
		n = "62" + strings.TrimLeft(n, "0")
	case strings.HasPrefix(n, "8"):
		n = "62" + n
	default:
		return "", false
	}
	// 62 + 9..13 angka. Nomor seluler Indonesia terpendek 10 digit lokal,
	// terpanjang 13; di luar rentang itu hampir pasti salah ketik.
	if len(n) < 11 || len(n) > 15 {
		return "", false
	}
	return n, true
}

// PUT /api/profil — mengubah nama, alamat, dan nomor WhatsApp klinik.
//
// Sebelum ini ketiganya tertulis mati di Patient.vue, sehingga mengubah nomor
// telepon klinik menuntut commit, build ulang, dan deploy — sesuatu yang tidak
// mungkin dilakukan bidannya sendiri.
func (s *Server) putProfil(w http.ResponseWriter, r *http.Request) {
	var req profilReq
	if err := readJSON(r, &req); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}

	clinic := strings.TrimSpace(req.Clinic)
	if clinic == "" {
		badRequest(w, "nama klinik wajib diisi")
		return
	}
	if utf8.RuneCountInString(clinic) > klinikMaks {
		badRequest(w, "nama klinik terlalu panjang")
		return
	}
	alamat := strings.TrimSpace(req.Alamat)
	if utf8.RuneCountInString(alamat) > alamatMaks {
		badRequest(w, "alamat terlalu panjang")
		return
	}
	wa, ok := normalWA(req.WhatsApp)
	if !ok {
		badRequest(w, "nomor WhatsApp tidak dikenali — tulis seperti 081234567890")
		return
	}

	// Satu transaksi karena nilainya tersebar di dua tabel. Tanpa itu, nama
	// klinik bisa tersimpan sementara alamatnya gagal, dan halaman pasien
	// menampilkan gabungan yang tidak pernah diisi siapa pun.
	tx, err := s.DB.BeginTx(r.Context(), nil)
	if err != nil {
		serverError(w, err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	now := nowMs()
	// ON DUPLICATE KEY hanya menyentuh kolom clinic: baris ini juga memegang
	// status hadir/keluar bidan, dan menyimpan profil tidak boleh diam-diam
	// mengembalikan statusnya ke "hadir".
	if _, err := tx.ExecContext(r.Context(),
		`INSERT INTO clinic_status (id, bidan_hadir, away_note, away_until, updated_ts, clinic)
		 VALUES (1,1,'',NULL,?,?)
		 ON DUPLICATE KEY UPDATE clinic=VALUES(clinic)`,
		now, clinic); err != nil {
		serverError(w, err)
		return
	}
	if _, err := tx.ExecContext(r.Context(),
		`INSERT INTO clinic_profile (id, alamat, whatsapp, updated_ts)
		 VALUES (1,?,?,?)
		 ON DUPLICATE KEY UPDATE alamat=VALUES(alamat), whatsapp=VALUES(whatsapp),
		   updated_ts=VALUES(updated_ts)`,
		alamat, wa, now); err != nil {
		serverError(w, err)
		return
	}
	if err := tx.Commit(); err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, models.Profil{Clinic: clinic, Alamat: alamat, WhatsApp: wa})
}
