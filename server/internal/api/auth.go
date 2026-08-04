package api

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"klinikbidanpit/internal/auth"
)

// NamaCookieSesi — nama cookie yang membawa token sesi.
const NamaCookieSesi = "klinik_sesi"

// SesiTTL — umur sesi. Diperpanjang selama masih dipakai (lihat segarkanSesi),
// jadi bidan yang memakai aplikasi tiap hari tidak pernah diminta login ulang,
// sementara sesi di perangkat yang ditinggalkan tetap mati sendiri.
const SesiTTL = 14 * 24 * time.Hour

// jedaSegarkan — sesi paling sering diperbarui sekali per jeda ini. Tanpa
// pembatas, setiap request ber-auth menulis ke database tanpa guna.
const jedaSegarkan = time.Hour

// ---------- IP klien ----------

// clientIP menentukan alamat pemanggil.
//
// Di produksi aplikasi ini berada di belakang Caddy, jadi RemoteAddr selalu
// 127.0.0.1 — kalau dipakai apa adanya, pembatas laju per-IP berubah menjadi
// satu jatah bersama untuk seluruh dunia.
//
// X-Forwarded-For hanya dipercaya bila sambungan langsungnya berasal dari
// loopback (artinya reverse proxy kita sendiri), dan yang diambil adalah entri
// TERAKHIR — itulah yang ditambahkan proxy. Entri sebelumnya bisa dikarang
// pengirim, dan memercayainya berarti siapa pun bisa mengganti identitas untuk
// mendapat jatah baru sesuka hati.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return host
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		bagian := strings.Split(xff, ",")
		if akhir := strings.TrimSpace(bagian[len(bagian)-1]); akhir != "" {
			return akhir
		}
	}
	return host
}

// ---------- penyimpanan sesi ----------

type sesi struct {
	ID     int64
	UserID int64
}

// buatSesi mencatat sesi baru dan mengembalikan token mentahnya.
// Token mentah hanya ada di sini dan di cookie — database cuma memegang hash.
func (s *Server) buatSesi(ctx context.Context, userID int64, r *http.Request) (string, error) {
	token, err := auth.TokenBaru()
	if err != nil {
		return "", err
	}
	now := nowMs()
	ua := r.UserAgent()
	if len(ua) > 255 {
		ua = ua[:255]
	}
	_, err = s.DB.ExecContext(ctx,
		`INSERT INTO sessions (user_id, token_hash, user_agent, ip, created_at, last_seen_at, expires_at)
		 VALUES (?,?,?,?,?,?,?)`,
		userID, auth.HashToken(token), ua, clientIP(r), now, now, now+SesiTTL.Milliseconds())
	if err != nil {
		return "", err
	}
	return token, nil
}

// cariSesi mengembalikan sesi yang masih sah untuk sebuah token.
// sql.ErrNoRows berarti token tidak dikenal, sudah dicabut, atau kedaluwarsa —
// ketiganya diperlakukan sama supaya tidak ada yang bisa disimpulkan darinya.
func (s *Server) cariSesi(ctx context.Context, token string) (sesi, error) {
	var out sesi
	var lastSeen int64
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, last_seen_at FROM sessions
		 WHERE token_hash = ? AND revoked_at IS NULL AND expires_at > ?`,
		auth.HashToken(token), nowMs()).Scan(&out.ID, &out.UserID, &lastSeen)
	if err != nil {
		return sesi{}, err
	}
	s.segarkanSesi(ctx, out.ID, lastSeen)
	return out, nil
}

// segarkanSesi memperpanjang masa berlaku sesi yang sedang dipakai.
// Kegagalannya sengaja hanya dicatat: request yang sudah sah tidak boleh gagal
// hanya karena catatan "terakhir dipakai" tidak sempat diperbarui.
func (s *Server) segarkanSesi(ctx context.Context, id int64, lastSeen int64) {
	now := nowMs()
	if now-lastSeen < jedaSegarkan.Milliseconds() {
		return
	}
	if _, err := s.DB.ExecContext(ctx,
		"UPDATE sessions SET last_seen_at = ?, expires_at = ? WHERE id = ?",
		now, now+SesiTTL.Milliseconds(), id); err != nil {
		s.logErr("memperpanjang sesi", err)
	}
}

func (s *Server) cabutSesi(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx,
		"UPDATE sessions SET revoked_at = ? WHERE id = ? AND revoked_at IS NULL", nowMs(), id)
	return err
}

// SapuSesiKedaluwarsa membuang sesi mati yang sudah lama. Dipanggil berkala
// dari main supaya tabelnya tidak tumbuh selamanya.
func (s *Server) SapuSesiKedaluwarsa(ctx context.Context) {
	batas := nowMs() - (30 * 24 * time.Hour).Milliseconds()
	if _, err := s.DB.ExecContext(ctx,
		"DELETE FROM sessions WHERE expires_at < ? OR (revoked_at IS NOT NULL AND revoked_at < ?)",
		batas, batas); err != nil {
		s.logErr("menyapu sesi kedaluwarsa", err)
	}
}

// ---------- cookie ----------

func (s *Server) pasangCookieSesi(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:  NamaCookieSesi,
		Value: token,
		Path:  "/",
		// HttpOnly: JavaScript tidak bisa membacanya, jadi satu celah XSS pun
		// tidak langsung berarti sesi bidan ikut terbawa.
		HttpOnly: true,
		// Lax: cookie tidak ikut terkirim pada permintaan lintas situs yang
		// mengubah data — ini yang menutup CSRF untuk POST/PUT/DELETE.
		SameSite: http.SameSiteLaxMode,
		Secure:   s.Cfg.Produksi,
		MaxAge:   int(SesiTTL.Seconds()),
	})
}

func (s *Server) hapusCookieSesi(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     NamaCookieSesi,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   s.Cfg.Produksi,
		MaxAge:   -1,
	})
}

// ---------- handler ----------

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// POST /api/auth/login
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := readJSON(r, &req); err != nil {
		badRequest(w, "body JSON tidak valid")
		return
	}
	username := strings.ToLower(strings.TrimSpace(req.Username))
	if username == "" || req.Password == "" {
		badRequest(w, "nama pengguna & kata sandi wajib diisi")
		return
	}

	// Jatah dipotong per nama pengguna, dan hanya saat gagal. Pembatas per IP
	// saja tidak cukup: serangan dari banyak IP tetap lolos selama tiap IP
	// pelan. Karena hanya kegagalan yang memotong, bidan yang mengetik benar
	// tidak pernah terkena.
	if !s.loginFail.Allow(username) {
		writeJSON(w, http.StatusTooManyRequests,
			map[string]string{"error": "terlalu banyak percobaan, coba lagi beberapa saat lagi"})
		return
	}

	var userID int64
	var hash string
	err := s.DB.QueryRowContext(r.Context(),
		"SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&userID, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		// Tetap hitung bcrypt terhadap hash boneka. Tanpa ini, balasan untuk
		// nama yang tidak terdaftar datang jauh lebih cepat, dan selisih waktu
		// itu sendiri sudah memberitahu akun mana yang ada.
		auth.CocokSandi(s.dummyHash, req.Password)
		s.tolakLogin(w)
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}
	if !auth.CocokSandi(hash, req.Password) {
		s.tolakLogin(w)
		return
	}

	token, err := s.buatSesi(r.Context(), userID, r)
	if err != nil {
		serverError(w, err)
		return
	}
	s.loginFail.Reset(username)
	s.pasangCookieSesi(w, token)
	writeJSON(w, http.StatusOK, map[string]any{"username": username})
}

// tolakLogin membalas dengan pesan yang sama untuk nama pengguna salah maupun
// kata sandi salah. Membedakannya akan memberitahu penebak bahwa separuh
// tebakannya sudah benar.
func (s *Server) tolakLogin(w http.ResponseWriter) {
	writeJSON(w, http.StatusUnauthorized,
		map[string]string{"error": "nama pengguna atau kata sandi salah"})
}

// POST /api/auth/logout
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if ck, err := r.Cookie(NamaCookieSesi); err == nil && ck.Value != "" {
		if ss, err := s.cariSesi(r.Context(), ck.Value); err == nil {
			if err := s.cabutSesi(r.Context(), ss.ID); err != nil {
				serverError(w, err)
				return
			}
		}
	}
	// Cookie dibersihkan apa pun hasilnya: kalau sesinya memang sudah tidak
	// sah, menyisakan cookie mati hanya membuat pemakai bingung.
	s.hapusCookieSesi(w)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// GET /api/auth/me — dipakai frontend untuk tahu masih login atau belum.
func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	ss, ok := sesiDari(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "belum masuk"})
		return
	}
	var username string
	if err := s.DB.QueryRowContext(r.Context(),
		"SELECT username FROM users WHERE id = ?", ss.UserID).Scan(&username); err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"username": username})
}

// ---------- middleware ----------

type kunciCtx int

const kunciSesi kunciCtx = 0

func sesiDari(ctx context.Context) (sesi, bool) {
	ss, ok := ctx.Value(kunciSesi).(sesi)
	return ss, ok
}

// wajibLogin menolak request tanpa sesi yang sah.
//
// Ini pengaman satu-satunya untuk data pasien. Halaman publik dan health check
// sengaja TIDAK dibungkus — semua sisanya wajib lewat sini.
func (s *Server) wajibLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ck, err := r.Cookie(NamaCookieSesi)
		if err != nil || ck.Value == "" {
			s.tolakBelumMasuk(w)
			return
		}
		ss, err := s.cariSesi(r.Context(), ck.Value)
		if errors.Is(err, sql.ErrNoRows) {
			s.hapusCookieSesi(w)
			s.tolakBelumMasuk(w)
			return
		}
		if err != nil {
			serverError(w, err)
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), kunciSesi, ss)))
	}
}

func (s *Server) tolakBelumMasuk(w http.ResponseWriter) {
	writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "belum masuk"})
}
