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
// multiStatements=true agar skrip schema/seed bisa mengirim banyak perintah.
func Open() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true",
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
