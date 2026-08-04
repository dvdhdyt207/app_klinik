// Package db membuka koneksi pool ke MySQL dari variabel lingkungan.
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Open membangun *sql.DB (pool) ke database app_klinik.
//
// multiStatements sengaja TIDAK dinyalakan. Itu hanya dibutuhkan skrip
// schema/seed, sementara di koneksi aplikasi ia memperbesar akibat: satu celah
// injeksi yang tadinya hanya bisa mengubah satu query berubah menjadi bisa
// menjalankan rangkaian perintah apa pun. Skrip yang memerlukannya memakai
// OpenMulti di bawah.
func Open() (*sql.DB, error) {
	return buka("")
}

// OpenMulti sama seperti Open tapi mengizinkan banyak perintah dalam satu
// kiriman. Hanya untuk perkakas sekali jalan (schema, seed) — jangan dipakai
// oleh server.
func OpenMulti() (*sql.DB, error) {
	return buka("&multiStatements=true")
}

func buka(tambahan string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"+tambahan,
		env("DB_USER", "root"),
		env("DB_PASSWORD", ""),
		env("DB_HOST", "127.0.0.1"),
		env("DB_PORT", "3306"),
		env("DB_NAME", "app_klinik"),
	)
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(10)
	d.SetConnMaxLifetime(3 * time.Minute)
	if err := d.Ping(); err != nil {
		return nil, fmt.Errorf("tidak bisa konek MySQL: %w", err)
	}
	return d, nil
}

// MustOpen = Open yang fatal bila gagal (dipakai saat startup).
func MustOpen() *sql.DB {
	d, err := Open()
	if err != nil {
		log.Fatal(err)
	}
	return d
}
