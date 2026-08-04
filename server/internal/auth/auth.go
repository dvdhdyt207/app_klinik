// Package auth memuat bagian yang berdiri sendiri dari autentikasi: membuat &
// memeriksa token sesi, mengelola kata sandi, dan membatasi laju percobaan.
//
// Sesi memakai token acak buram (opaque), bukan JWT. Untuk satu akun bidan,
// JWT hanya menambah bagian yang bisa salah: ada rahasia penanda tangan yang
// harus dikelola, ada masa berlaku yang harus dipahami, dan token yang sudah
// dikeluarkan tidak bisa ditarik kembali. Token acak yang dicatat di tabel
// sessions lebih sederhana dan justru lebih mudah dicabut.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// TokenBytes = 32 byte acak (256 bit). Cukup jauh di atas batas praktis
// tebakan, dan tetap pendek sebagai nilai cookie.
const TokenBytes = 32

// TokenBaru mengembalikan token sesi acak dalam bentuk aman-URL.
//
// Memakai crypto/rand, bukan math/rand: math/rand dapat ditebak seluruhnya
// bila sebagian keluarannya diketahui, dan itu sudah cukup untuk memalsukan
// sesi orang lain.
func TokenBaru() (string, error) {
	b := make([]byte, TokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("gagal membangkitkan token sesi: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken mengembalikan SHA-256 heksadesimal dari token.
//
// Yang disimpan di database hanya hash ini. SHA-256 tanpa salt memang cukup di
// sini — berbeda dengan kata sandi, tokennya sudah 256 bit acak, jadi tidak ada
// yang bisa dicari lewat kamus atau tabel pelangi.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// SamaAman membandingkan dua string tanpa membocorkan lewat lamanya waktu.
func SamaAman(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// HashSandi membuat hash bcrypt dari kata sandi.
func HashSandi(sandi string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(sandi), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// CocokSandi memeriksa kata sandi terhadap hash bcrypt.
func CocokSandi(hash, sandi string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(sandi)) == nil
}
