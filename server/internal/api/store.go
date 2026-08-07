package api

import (
	"context"
	"database/sql"

	"klinikbidanpit/internal/models"
)

// getStatus membaca baris tunggal clinic_status (id=1); membuatnya bila belum ada.
func (s *Server) getStatus(ctx context.Context) (models.Status, error) {
	var (
		hadir     int
		note      string
		awayUntil sql.NullInt64
		updated   int64
		clinic    string
	)
	err := s.DB.QueryRowContext(ctx,
		"SELECT bidan_hadir, away_note, away_until, updated_ts, clinic FROM clinic_status WHERE id = 1").
		Scan(&hadir, &note, &awayUntil, &updated, &clinic)

	if err == sql.ErrNoRows {
		ts := nowMs()
		_, err = s.DB.ExecContext(ctx,
			"INSERT INTO clinic_status (id, bidan_hadir, away_note, away_until, updated_ts, clinic) VALUES (1,1,'',NULL,?,?)",
			ts, "Klinik Bidan Pit")
		if err != nil {
			return models.Status{}, err
		}
		return models.Status{BidanHadir: true, AwayNote: "", AwayUntil: nil, UpdatedTs: ts, Clinic: "Klinik Bidan Pit"}, nil
	}
	if err != nil {
		return models.Status{}, err
	}

	var au *int64
	if awayUntil.Valid {
		v := awayUntil.Int64
		au = &v
	}
	return models.Status{
		BidanHadir: hadir != 0,
		AwayNote:   note,
		AwayUntil:  au,
		UpdatedTs:  updated,
		Clinic:     clinic,
	}, nil
}

func (s *Server) getEvents(ctx context.Context) ([]models.Event, error) {
	rows, err := s.DB.QueryContext(ctx,
		"SELECT id, title, all_day, start_ts, end_ts FROM events ORDER BY start_ts ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var e models.Event
		var allDay int
		if err := rows.Scan(&e.ID, &e.Title, &allDay, &e.StartTs, &e.EndTs); err != nil {
			return nil, err
		}
		e.AllDay = allDay != 0
		events = append(events, e)
	}
	return events, rows.Err()
}

func (s *Server) getVisits(ctx context.Context) ([]models.Visit, error) {
	rows, err := s.DB.QueryContext(ctx,
		"SELECT id, name, age, gejala, ts FROM visits ORDER BY ts DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := []models.Visit{}
	order := []string{} // jaga urutan; kumpulkan item per kunjungan setelah ini
	idx := map[string]int{}
	for rows.Next() {
		var v models.Visit
		var gejala sql.NullString
		if err := rows.Scan(&v.ID, &v.Name, &v.Age, &gejala, &v.Ts); err != nil {
			return nil, err
		}
		v.Gejala = gejala.String
		v.Items = []models.VisitItem{}
		idx[v.ID] = len(visits)
		order = append(order, v.ID)
		visits = append(visits, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	_ = order

	// Ambil semua item sekaligus, lalu kelompokkan ke kunjungannya.
	irows, err := s.DB.QueryContext(ctx,
		"SELECT visit_id, name, qty, unit FROM visit_items")
	if err != nil {
		return nil, err
	}
	defer irows.Close()
	for irows.Next() {
		var visitID string
		var it models.VisitItem
		if err := irows.Scan(&visitID, &it.Name, &it.Qty, &it.Unit); err != nil {
			return nil, err
		}
		if i, ok := idx[visitID]; ok {
			visits[i].Items = append(visits[i].Items, it)
		}
	}
	return visits, irows.Err()
}

// statusPayload = bentuk payload untuk halaman web pasien (README: klinik_bidan_status).
func (s *Server) statusPayload(ctx context.Context) (models.PublicStatus, error) {
	st, err := s.getStatus(ctx)
	if err != nil {
		return models.PublicStatus{}, err
	}
	events, err := s.getEvents(ctx)
	if err != nil {
		return models.PublicStatus{}, err
	}
	alamat, wa, err := s.profilRow(ctx)
	if err != nil {
		return models.PublicStatus{}, err
	}
	return models.PublicStatus{
		Hadir:     st.BidanHadir,
		AwayNote:  st.AwayNote,
		AwayUntil: st.AwayUntil,
		Ts:        st.UpdatedTs,
		Clinic:    st.Clinic,
		Alamat:    alamat,
		WhatsApp:  wa,
		Events:    events,
	}, nil
}

// profilRow membaca baris clinic_profile, membuatnya bila belum ada.
//
// Dibuat saat dibaca, bukan lewat INSERT di schema.sql, karena database
// produksi sudah berjalan sebelum tabel ini ada — dan tabel yang masih kosong
// tidak boleh membuat halaman pasien gagal dimuat. Baris kosong berarti "belum
// diisi", dan halaman pasien memang menyembunyikan bagian yang kosong.
func (s *Server) profilRow(ctx context.Context) (alamat, whatsapp string, err error) {
	err = s.DB.QueryRowContext(ctx,
		"SELECT alamat, whatsapp FROM clinic_profile WHERE id = 1").Scan(&alamat, &whatsapp)
	if err == sql.ErrNoRows {
		if _, err = s.DB.ExecContext(ctx,
			"INSERT INTO clinic_profile (id, alamat, whatsapp, updated_ts) VALUES (1,'','',?)",
			nowMs()); err != nil {
			return "", "", err
		}
		return "", "", nil
	}
	if err != nil {
		return "", "", err
	}
	return alamat, whatsapp, nil
}
