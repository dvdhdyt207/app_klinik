package api

import (
	"net/http"
	"strings"

	"klinikbidanpit/internal/models"
)

// GET /api/health
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// GET /api/state — snapshot penuh untuk bootstrap Bidan App.
func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status, err := s.getStatus(ctx)
	if err != nil {
		serverError(w, err)
		return
	}
	events, err := s.getEvents(ctx)
	if err != nil {
		serverError(w, err)
		return
	}
	visits, err := s.getVisits(ctx)
	if err != nil {
		serverError(w, err)
		return
	}
	alamat, wa, err := s.profilRow(ctx)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, models.State{
		Status: status,
		Profil: models.Profil{Clinic: status.Clinic, Alamat: alamat, WhatsApp: wa},
		Events: events, Visits: visits,
	})
}

// PUT /api/status — ubah status bidan (hadir / keluar / extend).
func (s *Server) putStatus(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BidanHadir bool   `json:"bidanHadir"`
		AwayNote   string `json:"awayNote"`
		AwayUntil  *int64 `json:"awayUntil"`
	}
	if err := readJSON(r, &body); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	hadir := 0
	if body.BidanHadir {
		hadir = 1
	}
	_, err := s.DB.ExecContext(r.Context(),
		`INSERT INTO clinic_status (id, bidan_hadir, away_note, away_until, updated_ts)
		 VALUES (1,?,?,?,?)
		 ON DUPLICATE KEY UPDATE bidan_hadir=VALUES(bidan_hadir), away_note=VALUES(away_note),
		   away_until=VALUES(away_until), updated_ts=VALUES(updated_ts)`,
		hadir, body.AwayNote, body.AwayUntil, nowMs())
	if err != nil {
		serverError(w, err)
		return
	}
	st, err := s.getStatus(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}

type eventReq struct {
	Title   string `json:"title"`
	AllDay  bool   `json:"allDay"`
	StartTs int64  `json:"startTs"`
	EndTs   int64  `json:"endTs"`
}

// normalisasi: title dipangkas, end minimal = start.
func (e eventReq) normalized() (title string, allDay int, start, end int64, ok bool) {
	title = strings.TrimSpace(e.Title)
	if title == "" || e.StartTs <= 0 {
		return "", 0, 0, 0, false
	}
	end = e.StartTs
	if e.EndTs >= e.StartTs {
		end = e.EndTs
	}
	allDay = 0
	if e.AllDay {
		allDay = 1
	}
	return title, allDay, e.StartTs, end, true
}

// POST /api/events — tambah agenda.
func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var req eventReq
	if err := readJSON(r, &req); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	title, allDay, start, end, ok := req.normalized()
	if !ok {
		badRequest(w, "title & startTs wajib")
		return
	}
	id := uid("e")
	_, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO events (id,title,all_day,start_ts,end_ts) VALUES (?,?,?,?,?)",
		id, title, allDay, start, end)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, models.Event{
		ID: id, Title: title, AllDay: allDay != 0, StartTs: start, EndTs: end,
	})
}

// PUT /api/events/{id} — ubah agenda.
func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req eventReq
	if err := readJSON(r, &req); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	title, allDay, start, end, ok := req.normalized()
	if !ok {
		badRequest(w, "title & startTs wajib")
		return
	}
	_, err := s.DB.ExecContext(r.Context(),
		"UPDATE events SET title=?, all_day=?, start_ts=?, end_ts=? WHERE id=?",
		title, allDay, start, end, id)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, models.Event{
		ID: id, Title: title, AllDay: allDay != 0, StartTs: start, EndTs: end,
	})
}

// DELETE /api/events/{id} — hapus agenda.
func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	_, err := s.DB.ExecContext(r.Context(), "DELETE FROM events WHERE id=?", r.PathValue("id"))
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// POST /api/visits — catat kunjungan beserta obat yang diberikan (transaksi).
//
// TIDAK ada stok yang dikurangi di sini. Aplikasi ini sengaja tidak lagi
// mengelola persediaan obat: yang dicatat hanyalah apa yang diberikan kepada
// pasien. Nama & satuan obat disimpan sebagai teks pada visit_items, dan dari
// riwayat itulah daftar saran obat di layar Cari Obat dibentuk — jadi obat yang
// sudah pernah dipakai bisa dipakai ulang tanpa perlu didata lebih dulu.
func (s *Server) createVisit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name   string             `json:"name"`
		Age    int                `json:"age"`
		Gejala string             `json:"gejala"`
		Items  []models.VisitItem `json:"items"`
	}
	if err := readJSON(r, &body); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" || len(body.Items) == 0 {
		badRequest(w, "name & minimal 1 obat wajib")
		return
	}

	tx, err := s.DB.BeginTx(r.Context(), nil)
	if err != nil {
		serverError(w, err)
		return
	}
	defer tx.Rollback() // no-op bila sudah commit

	id := uid("v")
	if _, err := tx.ExecContext(r.Context(),
		"INSERT INTO visits (id,name,age,gejala,ts) VALUES (?,?,?,?,?)",
		id, name, body.Age, strings.TrimSpace(body.Gejala), nowMs()); err != nil {
		serverError(w, err)
		return
	}
	for _, it := range body.Items {
		if _, err := tx.ExecContext(r.Context(),
			"INSERT INTO visit_items (visit_id,name,qty,unit) VALUES (?,?,?,?)",
			id, strings.TrimSpace(it.Name), it.Qty, it.Unit); err != nil {
			serverError(w, err)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// GET /api/public/status — payload halaman pasien.
func (s *Server) publicStatus(w http.ResponseWriter, r *http.Request) {
	payload, err := s.statusPayload(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}
