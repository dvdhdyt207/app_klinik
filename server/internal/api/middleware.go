package api

import (
	"net/http"
	"time"

	"klinikbidanpit/internal/auth"

	"golang.org/x/time/rate"
)

// HeaderKeamanan menambahkan header keamanan dasar pada setiap response.
//
// HSTS hanya di produksi: di localhost aksesnya lewat http, dan HSTS akan
// membuat peramban mengingat https untuk localhost sehingga menyulitkan
// proyek lain di laptop yang sama.
func HeaderKeamanan(produksi bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		// Filter XSS bawaan peramban lama justru pernah menimbulkan celah
		// sendiri dan sudah dihapus dari peramban modern; 0 mematikannya.
		h.Set("X-XSS-Protection", "0")
		if produksi {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

// BatasLaju membatasi jumlah request per alamat IP.
//
// Alamat diambil lewat clientIP, yang hanya memercayai X-Forwarded-For dari
// sambungan loopback (reverse proxy kita sendiri).
func BatasLaju(rps float64, burst int, next http.Handler) http.Handler {
	limiter := auth.NewKeyedLimiter(rate.Limit(rps), burst, 10*time.Minute, 10000)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow(clientIP(r)) {
			writeJSON(w, http.StatusTooManyRequests,
				map[string]string{"error": "terlalu banyak permintaan, coba lagi nanti"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// BatasBody menolak badan request yang melebihi maksBytes, supaya kiriman
// raksasa tidak menghabiskan memori server sebelum sempat divalidasi.
func BatasBody(maksBytes int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maksBytes)
		}
		next.ServeHTTP(w, r)
	})
}

// CORS mengizinkan origin lain memanggil API — hanya yang terdaftar.
//
// Dulu header ini selalu "*". Digabung dengan cookie sesi itu berbahaya:
// situs mana pun bisa memanggil API atas nama bidan yang sedang login.
// Peramban sendiri menolak "*" begitu Allow-Credentials menyala, jadi origin
// harus disebut satu per satu. Di produksi frontend disajikan dari origin yang
// sama, sehingga daftar ini normalnya kosong dan CORS tidak dipakai sama sekali.
func CORS(diizinkan []string, next http.Handler) http.Handler {
	set := make(map[string]bool, len(diizinkan))
	for _, o := range diizinkan {
		if o != "" {
			set[o] = true
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && set[origin] {
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Credentials", "true")
			h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Content-Type")
			h.Set("Access-Control-Max-Age", "600")
			// Balasan berbeda-beda menurut Origin; tanpa Vary, cache bisa
			// menyajikan izin milik origin lain.
			h.Add("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
