// Package models memuat struct data yang dipertukarkan lewat JSON API.
// Bentuk JSON dijaga sama persis dengan backend Node lama (camelCase).
package models

// Status bidan (baris tunggal clinic_status).
type Status struct {
	BidanHadir bool   `json:"bidanHadir"`
	AwayNote   string `json:"awayNote"`
	AwayUntil  *int64 `json:"awayUntil"` // epoch ms; null = belum dipastikan
	UpdatedTs  int64  `json:"updatedTs"`
	Clinic     string `json:"clinic"`
}

// Event = agenda/jadwal bidan.
type Event struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	AllDay  bool   `json:"allDay"`
	StartTs int64  `json:"startTs"`
	EndTs   int64  `json:"endTs"`
}

// Medicine = stok obat. Qty selalu dalam base unit.
type Medicine struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Cat      string `json:"cat"`
	Qty      int    `json:"qty"`
	BaseUnit string `json:"baseUnit"`
}

// VisitItem = obat yang diberikan pada satu kunjungan.
type VisitItem struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
	Unit string `json:"unit"`
}

// Visit = catatan kunjungan pasien beserta obatnya.
type Visit struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Age    int         `json:"age"`
	Gejala string      `json:"gejala"`
	Ts     int64       `json:"ts"`
	Items  []VisitItem `json:"items"`
}

// State = snapshot penuh untuk bootstrap Bidan App (GET /api/state).
type State struct {
	Status    Status     `json:"status"`
	Events    []Event    `json:"events"`
	Medicines []Medicine `json:"medicines"`
	Visits    []Visit    `json:"visits"`
}

// PublicStatus = payload halaman web pasien (bentuk "klinik_bidan_status").
type PublicStatus struct {
	Hadir     bool    `json:"hadir"`
	AwayNote  string  `json:"awayNote"`
	AwayUntil *int64  `json:"awayUntil"`
	Ts        int64   `json:"ts"`
	Clinic    string  `json:"clinic"`
	Events    []Event `json:"events"`
}
