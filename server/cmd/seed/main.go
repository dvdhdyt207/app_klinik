// Seed data awal (meniru prototipe design handoff). Idempotent: mengosongkan tabel dulu.
// Jalankan dari folder server/:  go run ./cmd/seed
package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"klinikbidanpit/internal/db"
)

const dayMs int64 = 86400000

// atMs = epoch ms untuk jam:menit tertentu pada tanggal `base` (waktu lokal).
func atMs(base time.Time, h, m int) int64 {
	d := time.Date(base.Year(), base.Month(), base.Day(), h, m, 0, 0, base.Location())
	return d.UnixMilli()
}

type event struct {
	id      string
	title   string
	allDay  int
	startTs int64
	endTs   int64
}
type visitItem struct {
	name string
	qty  int
	unit string
}
type visit struct {
	id     string
	name   string
	age    int
	ts     int64
	gejala string
	items  []visitItem
}

func main() {
	_ = godotenv.Load()

	// Perkakas ini meng-TRUNCATE tabel kunjungan dan agenda. Dijalankan
	// tanpa sadar di server produksi, yang terhapus adalah data klinik
	// sungguhan dan tidak ada jalan kembali. Menolak jalan di produksi jauh
	// lebih murah daripada mengandalkan ingatan orang yang sedang buru-buru.
	if strings.EqualFold(os.Getenv("APP_ENV"), "production") && os.Getenv("SEED_PAKSA") != "ya" {
		log.Fatal("APP_ENV=production: seed menolak jalan karena akan MENGHAPUS data klinik. " +
			"Bila memang disengaja, jalankan dengan SEED_PAKSA=ya")
	}

	database := db.MustOpen()
	defer database.Close()

	now := time.Now()
	nowMs := now.UnixMilli()
	sot := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sotMs := sot.UnixMilli()

	events := []event{
		{"e1", "Posyandu Desa Sukamaju", 0, atMs(sot, 13, 0), atMs(sot, 15, 0)},
		{"e2", "Keluar kota — acara keluarga", 1, sotMs + dayMs, sotMs + 3*dayMs},
	}
	// Tidak ada daftar obat yang di-seed: daftar obat bukan tabel tersendiri
	// lagi. Nama obat pada kunjungan di bawah inilah yang nanti muncul sebagai
	// saran di layar Cari Obat.
	visits := []visit{
		{"v1", "Siti Aminah", 27, nowMs - 2*3600000, "Demam 2 hari, batuk kering",
			[]visitItem{{"Paracetamol sirup", 1, "botol"}}},
		{"v2", "Budi Santoso", 45, nowMs - 5*3600000, "Nyeri lambung, mual",
			[]visitItem{{"Amoxicillin 500mg", 10, "butir"}, {"Vitamin B Complex", 10, "butir"}}},
		{"v3", "Rina Wati", 5, nowMs - dayMs - 3600000, "Panas tinggi, rewel",
			[]visitItem{{"Paracetamol 500mg", 10, "butir"}}},
		{"v4", "Joko Priyono", 60, nowMs - 2*dayMs, "",
			[]visitItem{{"Antasida", 20, "butir"}, {"Vitamin C 500mg", 10, "butir"}}},
	}

	conn, err := database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Rollback()

	stmts := []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		"TRUNCATE TABLE visit_items",
		"TRUNCATE TABLE visits",
		"TRUNCATE TABLE events",
		"TRUNCATE TABLE clinic_status",
		"SET FOREIGN_KEY_CHECKS = 1",
	}
	for _, q := range stmts {
		if _, err := conn.Exec(q); err != nil {
			log.Fatalf("%s: %v", q, err)
		}
	}

	if _, err := conn.Exec(
		"INSERT INTO clinic_status (id,bidan_hadir,away_note,away_until,updated_ts,clinic) VALUES (1,1,'',NULL,?,?)",
		nowMs, "Klinik Bidan Pit"); err != nil {
		log.Fatal(err)
	}
	for _, e := range events {
		if _, err := conn.Exec("INSERT INTO events (id,title,all_day,start_ts,end_ts) VALUES (?,?,?,?,?)",
			e.id, e.title, e.allDay, e.startTs, e.endTs); err != nil {
			log.Fatal(err)
		}
	}
	for _, v := range visits {
		if _, err := conn.Exec("INSERT INTO visits (id,name,age,gejala,ts) VALUES (?,?,?,?,?)",
			v.id, v.name, v.age, v.gejala, v.ts); err != nil {
			log.Fatal(err)
		}
		for _, it := range v.items {
			if _, err := conn.Exec("INSERT INTO visit_items (visit_id,name,qty,unit) VALUES (?,?,?,?)",
				v.id, it.name, it.qty, it.unit); err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := conn.Commit(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Seed selesai: %d agenda, %d kunjungan.", len(events), len(visits))
}
