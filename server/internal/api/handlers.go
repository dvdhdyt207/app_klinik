package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"klinikbidanpit/internal/catalog"
	"klinikbidanpit/internal/models"
)

// GET /api/health
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// GET /api/catalog — katalog + aturan satuan (layar Cari Obat / Tambah Stok).
func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"catalog": catalog.CATALOG,
		"cat":     catalog.CAT,
	})
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
	meds, err := s.getMedicines(ctx)
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
		Events: events, Medicines: meds, Visits: visits,
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

// POST /api/medicines/stock — tambah stok (buat obat baru bila belum ada).
// amount = jumlah dalam base unit.
func (s *Server) addStock(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Cat    string `json:"cat"`
		Amount int    `json:"amount"`
	}
	if err := readJSON(r, &body); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" || body.Amount <= 0 {
		badRequest(w, "name & amount(>0) wajib")
		return
	}
	cat := body.Cat
	if cat == "" {
		cat = "Tablet"
	}

	// cari berdasarkan id atau nama (case-insensitive)
	var foundID string
	err := s.DB.QueryRowContext(r.Context(),
		"SELECT id FROM medicines WHERE id=? OR LOWER(name)=LOWER(?) LIMIT 1",
		body.ID, name).Scan(&foundID)
	if err == nil {
		if _, err := s.DB.ExecContext(r.Context(),
			"UPDATE medicines SET qty = qty + ? WHERE id=?", body.Amount, foundID); err != nil {
			serverError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"id": foundID})
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		serverError(w, err)
		return
	}

	newID := body.ID
	if newID == "" {
		newID = uid("m")
	}
	if _, err := s.DB.ExecContext(r.Context(),
		"INSERT INTO medicines (id,name,cat,qty) VALUES (?,?,?,?)",
		newID, name, cat, body.Amount); err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": newID})
}

// PUT /api/medicines/{id} — ubah nama, kategori, & koreksi jumlah.
//
// Jumlah ikut bisa diubah di sini karena tanpanya stok hanya bisa naik (tambah
// stok) atau turun lewat kunjungan — tidak ada cara memperbaiki angka yang
// telanjur meleset dari hitungan fisik, dan angka yang tidak bisa dikoreksi
// pelan-pelan berhenti dipercaya.
func (s *Server) updateMedicine(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Name string `json:"name"`
		Cat  string `json:"cat"`
		Qty  int    `json:"qty"`
	}
	if err := readJSON(r, &body); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" {
		badRequest(w, "name wajib")
		return
	}
	if body.Qty < 0 {
		badRequest(w, "qty tidak boleh negatif")
		return
	}
	cat := body.Cat
	if cat == "" {
		cat = "Tablet"
	}
	if !catalog.ValidCat(cat) {
		badRequest(w, "kategori tidak dikenal")
		return
	}

	// Nama harus tetap unik. Pengurangan stok saat mencatat kunjungan
	// mencocokkan LOWER(name) (lihat createVisit), jadi dua obat bernama sama
	// membuat kunjungan mengurangi baris yang salah — diam-diam, tanpa error.
	var lain string
	err := s.DB.QueryRowContext(r.Context(),
		"SELECT id FROM medicines WHERE LOWER(name)=LOWER(?) AND id<>? LIMIT 1",
		name, id).Scan(&lain)
	if err == nil {
		writeErr(w, http.StatusConflict, "sudah ada obat lain dengan nama itu")
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		serverError(w, err)
		return
	}

	res, err := s.DB.ExecContext(r.Context(),
		"UPDATE medicines SET name=?, cat=?, qty=? WHERE id=?", name, cat, body.Qty, id)
	if err != nil {
		serverError(w, err)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		// Bedakan dari "berhasil tapi nilainya sama": baris yang tidak ada
		// sama sekali berarti layar pemanggilnya memegang data basi.
		var ada string
		if err := s.DB.QueryRowContext(r.Context(),
			"SELECT id FROM medicines WHERE id=?", id).Scan(&ada); errors.Is(err, sql.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "obat tidak ditemukan")
			return
		}
	}
	writeJSON(w, http.StatusOK, models.Medicine{
		ID: id, Name: name, Cat: cat, Qty: body.Qty, BaseUnit: catalog.BaseUnit(cat),
	})
}

// DELETE /api/medicines/{id} — hapus obat dari daftar stok.
//
// Riwayat kunjungan TIDAK ikut terhapus: visit_items menyimpan nama obatnya
// sendiri sebagai teks, tanpa foreign key ke tabel ini. Itu disengaja — rekam
// medis harus menyimpan apa yang waktu itu diberikan, bukan apa yang hari ini
// masih ada di lemari.
func (s *Server) deleteMedicine(w http.ResponseWriter, r *http.Request) {
	_, err := s.DB.ExecContext(r.Context(), "DELETE FROM medicines WHERE id=?", r.PathValue("id"))
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// POST /api/visits — catat kunjungan + kurangi stok (transaksi).
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
			id, it.Name, it.Qty, it.Unit); err != nil {
			serverError(w, err)
			return
		}
		// kurangi stok berdasarkan nama, clamp di 0
		if _, err := tx.ExecContext(r.Context(),
			"UPDATE medicines SET qty = GREATEST(0, qty - ?) WHERE LOWER(name)=LOWER(?)",
			it.Qty, it.Name); err != nil {
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
