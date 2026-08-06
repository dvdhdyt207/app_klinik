package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// ---------- ganti sandi ----------

type gantiSandiReq struct {
	SandiLama string `json:"old_password"`
	SandiBaru string `json:"new_password"`
}

// POST /api/auth/password — mengganti sandi akun yang sedang login.
//
// Sebelum endpoint ini ada, sandi bidan adalah nilai yang dituliskan pemilik
// server di ADMIN_PASSWORD dan tidak pernah bisa diubah dari aplikasi:
// menggantinya menuntut UPDATE manual berisi hash bcrypt yang dibuat di luar.
// Untuk aplikasi berisi rekam medis milik pasien orang lain, itu berarti
// pemilik aplikasi selamanya ikut memegang kunci akun kliennya.
func (s *Server) gantiSandi(w http.ResponseWriter, r *http.Request) {
	ss, ok := sesiDari(r.Context())
	if !ok {
		s.tolakBelumMasuk(w)
		return
	}

	var req gantiSandiReq
	if err := readJSON(r, &req); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	if req.SandiLama == "" {
		badRequest(w, "sandi lama wajib diisi")
		return
	}
	if utf8.RuneCountInString(req.SandiBaru) < SandiMinimal {
		badRequest(w, fmt.Sprintf("sandi baru minimal %d karakter", SandiMinimal))
		return
	}
	if req.SandiBaru == req.SandiLama {
		badRequest(w, "sandi baru harus berbeda dari sandi lama")
		return
	}

	// Tebakan sandi lama ikut dibatasi, memakai pembatas yang sama dengan
	// login. Sesi yang dicuri — laptop yang ditinggal terbuka, cookie yang
	// bocor — tidak boleh berubah menjadi tempat menebak sandi asli tanpa
	// henti hanya karena penebaknya sudah berada di dalam.
	kunci := "ganti-sandi:" + strconv.FormatInt(ss.UserID, 10)
	if !s.loginFail.Allow(kunci) {
		writeJSON(w, http.StatusTooManyRequests,
			map[string]string{"error": "terlalu banyak percobaan, coba lagi beberapa saat lagi"})
		return
	}

	var hashLama string
	if err := s.DB.QueryRowContext(r.Context(),
		"SELECT password_hash FROM users WHERE id = ?", ss.UserID).Scan(&hashLama); err != nil {
		serverError(w, err)
		return
	}
	if !auth.CocokSandi(hashLama, req.SandiLama) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "sandi lama salah"})
		return
	}

	hashBaru, err := auth.HashSandi(req.SandiBaru)
	if err != nil {
		serverError(w, err)
		return
	}
	now := nowMs()
	if _, err := s.DB.ExecContext(r.Context(),
		"UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?",
		hashBaru, now, ss.UserID); err != nil {
		serverError(w, err)
		return
	}

	// Sesi LAIN dicabut, sesi yang sedang dipakai dibiarkan hidup.
	//
	// Orang mengganti sandi justru ketika curiga ada yang tahu; membiarkan
	// sesi lama hidup berarti penggantian itu tidak mengusir siapa pun.
	// Sebaliknya, ikut mencabut sesi sendiri akan melempar keluar orang yang
	// baru saja melakukan hal yang benar.
	//
	// Kegagalan di sini dicatat, bukan dibatalkan: sandinya sudah berganti,
	// dan membalas "gagal" akan membuat bidan mengira sandi lamanya masih
	// berlaku.
	if _, err := s.DB.ExecContext(r.Context(),
		"UPDATE sessions SET revoked_at = ? WHERE user_id = ? AND id <> ? AND revoked_at IS NULL",
		now, ss.UserID, ss.ID); err != nil {
		s.logErr("mencabut sesi lain setelah ganti sandi", err)
	}

	s.loginFail.Reset(kunci)
	log.Printf("sandi akun id=%d diganti", ss.UserID)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "message": "sandi berhasil diganti"})
}
