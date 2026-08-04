package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"klinikbidanpit/internal/auth"
)

// SandiMinimal — panjang minimum kata sandi bidan.
//
// Halaman login terbuka ke internet, jadi sandi pendek adalah satu-satunya
// hal yang memisahkan rekam medis pasien dari siapa pun yang mau menebak.
const SandiMinimal = 10

// PastikanAkun membuat akun bidan pertama bila tabel users masih kosong.
//
// Akun yang SUDAH ada tidak pernah diubah — termasuk sandinya. Kalau
// ADMIN_PASSWORD ikut mengganti sandi setiap kali server hidup, mengganti
// sandi lewat aplikasi jadi sia-sia dan nilai lama di .env diam-diam menang.
//
// Bila belum ada akun dan sandinya tidak disediakan, server sengaja MENOLAK
// jalan. Alternatifnya adalah hidup dengan sandi bawaan yang bisa ditebak,
// dan itu sama saja dengan tidak punya login.
func PastikanAkun(ctx context.Context, db *sql.DB, username, sandi string) error {
	username = strings.ToLower(strings.TrimSpace(username))
	if username == "" {
		username = "bidan"
	}

	var ada int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&ada)
	if err != nil {
		return fmt.Errorf("memeriksa akun: %w", err)
	}
	if ada > 0 {
		log.Printf("akun sudah ada (%d), ADMIN_PASSWORD diabaikan", ada)
		return nil
	}

	if sandi == "" {
		return errors.New("belum ada akun dan ADMIN_PASSWORD kosong — " +
			"isi ADMIN_PASSWORD di .env (minimal " + fmt.Sprint(SandiMinimal) + " karakter) lalu jalankan lagi")
	}
	if utf8.RuneCountInString(sandi) < SandiMinimal {
		return fmt.Errorf("ADMIN_PASSWORD terlalu pendek, minimal %d karakter", SandiMinimal)
	}

	hash, err := auth.HashSandi(sandi)
	if err != nil {
		return fmt.Errorf("membuat hash sandi: %w", err)
	}
	now := nowMs()
	if _, err := db.ExecContext(ctx,
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES (?,?,?,?)",
		username, hash, now, now); err != nil {
		return fmt.Errorf("membuat akun: %w", err)
	}
	log.Printf("akun '%s' dibuat dari ADMIN_PASSWORD — sebaiknya hapus nilainya dari .env setelah ini", username)
	return nil
}
