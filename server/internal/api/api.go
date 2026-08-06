// Package api memuat HTTP handler & routing untuk API Klinik Bidan Pit.
package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"klinikbidanpit/internal/auth"

	"golang.org/x/time/rate"
)

// Config — setelan yang memengaruhi perilaku keamanan handler.
type Config struct {
	// Produksi menyalakan cookie Secure & HSTS. Sengaja mati saat
	// pengembangan karena localhost diakses lewat http biasa.
	Produksi bool
}

// Server memegang dependensi handler.
type Server struct {
	DB  *sql.DB
	Cfg Config

	// loginFail membatasi kegagalan login per nama pengguna.
	loginFail *auth.KeyedLimiter

	// dummyHash dibandingkan saat nama pengguna tidak ada, supaya lamanya
	// respons tidak membocorkan akun mana yang terdaftar.
	dummyHash string
}

// New menyiapkan Server beserta pembatas laju dan hash pembanding.
func New(db *sql.DB, cfg Config) (*Server, error) {
	// 5 kegagalan beruntun, lalu hanya 1 percobaan tiap 20 detik.
	limiter := auth.NewKeyedLimiter(rate.Every(20*time.Second), 5, 15*time.Minute, 5000)

	// Kalau ini gagal, lebih baik berhenti daripada diam-diam kehilangan
	// perlindungan terhadap kebocoran lewat selisih waktu.
	hash, err := auth.HashSandi("penyeimbang-waktu")
	if err != nil {
		limiter.Stop()
		return nil, err
	}
	return &Server{DB: db, Cfg: cfg, loginFail: limiter, dummyHash: hash}, nil
}

// Close melepas sumber daya latar milik Server.
func (s *Server) Close() {
	if s.loginFail != nil {
		s.loginFail.Stop()
	}
}

// Routes mendaftarkan semua endpoint API (Go 1.22+ method+pattern routing)
// dan mengembalikan *http.ServeMux agar pemanggil bisa menambah handler statis "/".
//
// Hanya tiga jalur yang boleh terbuka tanpa sesi:
//   - health   : dipanggil pemantau/CI, tidak memuat data apa pun
//   - login    : pintu masuknya sendiri
//   - public/* : halaman pasien — status bidan & agenda, tanpa data pasien
//
// Sisanya WAJIB lewat wajibLogin. `/api/state` khususnya memuat nama, umur,
// gejala, dan obat tiap pasien; itu tidak boleh bisa diambil tanpa masuk.
func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// --- terbuka ---
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("GET /api/public/status", s.publicStatus)
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.HandleFunc("POST /api/auth/logout", s.logout)

	// --- wajib masuk ---
	mux.HandleFunc("GET /api/auth/me", s.wajibLogin(s.me))
	mux.HandleFunc("POST /api/auth/password", s.wajibLogin(s.gantiSandi))
	mux.HandleFunc("GET /api/catalog", s.wajibLogin(s.catalog))
	mux.HandleFunc("GET /api/state", s.wajibLogin(s.state))
	mux.HandleFunc("PUT /api/status", s.wajibLogin(s.putStatus))
	mux.HandleFunc("PUT /api/profil", s.wajibLogin(s.putProfil))
	mux.HandleFunc("POST /api/events", s.wajibLogin(s.createEvent))
	mux.HandleFunc("PUT /api/events/{id}", s.wajibLogin(s.updateEvent))
	mux.HandleFunc("DELETE /api/events/{id}", s.wajibLogin(s.deleteEvent))
	mux.HandleFunc("POST /api/medicines/stock", s.wajibLogin(s.addStock))
	mux.HandleFunc("PUT /api/medicines/{id}", s.wajibLogin(s.updateMedicine))
	mux.HandleFunc("DELETE /api/medicines/{id}", s.wajibLogin(s.deleteMedicine))
	mux.HandleFunc("POST /api/visits", s.wajibLogin(s.createVisit))

	return mux
}

// nowMs = epoch milidetik saat ini (sepadan dengan Date.now() di Node).
func nowMs() int64 { return time.Now().UnixMilli() }

const base36 = "0123456789abcdefghijklmnopqrstuvwxyz"

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = base36[rand.IntN(36)]
	}
	return string(b)
}

// uid meniru pola id Node: prefix + base36(waktu) + 4 char acak.
//
// Ini hanya untuk id baris, BUKAN untuk apa pun yang bersifat rahasia —
// math/rand dapat ditebak. Token sesi memakai crypto/rand; lihat paket auth.
func uid(prefix string) string {
	return prefix + strconv.FormatInt(nowMs(), 36) + randStr(4)
}

// ---------- helper respons JSON ----------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// writeErr membalas galat yang memang layak dibaca pemakai (nama bentrok,
// baris tidak ditemukan) — beda dari serverError, yang menyembunyikan sebabnya.
func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func badRequest(w http.ResponseWriter, msg string) {
	writeErr(w, http.StatusBadRequest, msg)
}

// serverError mencatat sebabnya di log server dan membalas pesan umum.
//
// Sebelumnya err.Error() dikirim apa adanya ke klien, sehingga nama tabel,
// bentuk query, dan pesan driver MySQL ikut keluar — bahan gratis bagi siapa
// pun yang sedang meraba-raba isi sistem.
func serverError(w http.ResponseWriter, err error) {
	log.Printf("ERROR: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "terjadi gangguan di server"})
}

// logErr mencatat kegagalan pada operasi sampingan yang tidak boleh
// menggagalkan request utama.
func (s *Server) logErr(apa string, err error) {
	if err != nil {
		log.Printf("PERINGATAN: %s: %v", apa, err)
	}
}

func readJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
