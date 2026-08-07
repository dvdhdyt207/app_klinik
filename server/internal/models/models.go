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

// Profil = keterangan klinik yang tampil di halaman pasien dan bisa diubah
// bidan sendiri lewat Bidan App.
//
// Clinic dibaca dari clinic_status.clinic, sisanya dari clinic_profile —
// digabung di sini supaya frontend tidak perlu tahu bahwa asalnya dua tabel.
type Profil struct {
	Clinic   string `json:"clinic"`
	Alamat   string `json:"alamat"`
	WhatsApp string `json:"whatsapp"` // hanya angka, diawali 62; kosong = sembunyikan
}

// Event = agenda/jadwal bidan.
type Event struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	AllDay  bool   `json:"allDay"`
	StartTs int64  `json:"startTs"`
	EndTs   int64  `json:"endTs"`
}

// VisitItem = obat yang diberikan pada satu kunjungan.
//
// Nama dan satuannya disimpan sebagai teks, berdiri sendiri: tidak ada tabel
// daftar obat yang ia tunjuk. Sejak pengelolaan stok dilepas, riwayat inilah
// satu-satunya sumber daftar obat — nama yang pernah dipakai muncul kembali
// sebagai saran di layar Cari Obat.
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
	Status Status  `json:"status"`
	Profil Profil  `json:"profil"`
	Events []Event `json:"events"`
	Visits []Visit `json:"visits"`
}

// PublicStatus = payload halaman web pasien (bentuk "klinik_bidan_status").
type PublicStatus struct {
	Hadir     bool    `json:"hadir"`
	AwayNote  string  `json:"awayNote"`
	AwayUntil *int64  `json:"awayUntil"`
	Ts        int64   `json:"ts"`
	Clinic    string  `json:"clinic"`
	Alamat    string  `json:"alamat"`
	WhatsApp  string  `json:"whatsapp"`
	Events    []Event `json:"events"`
}
