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
)

// Server memegang dependensi handler (koneksi DB).
type Server struct {
	DB *sql.DB
}

// Routes mendaftarkan semua endpoint API (Go 1.22+ method+pattern routing)
// dan mengembalikan *http.ServeMux agar pemanggil bisa menambah handler statis "/".
func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("GET /api/catalog", s.catalog)
	mux.HandleFunc("GET /api/state", s.state)
	mux.HandleFunc("PUT /api/status", s.putStatus)
	mux.HandleFunc("POST /api/events", s.createEvent)
	mux.HandleFunc("PUT /api/events/{id}", s.updateEvent)
	mux.HandleFunc("DELETE /api/events/{id}", s.deleteEvent)
	mux.HandleFunc("POST /api/medicines/stock", s.addStock)
	mux.HandleFunc("POST /api/visits", s.createVisit)
	mux.HandleFunc("GET /api/public/status", s.publicStatus)
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

func badRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("ERROR: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func readJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
