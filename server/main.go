// Klinik Bidan Pit — server Go.
// Menyajikan REST API (/api/*) + halaman web hasil build Vue (SPA fallback).
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"klinikbidanpit/internal/api"
	"klinikbidanpit/internal/db"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// daftar memecah nilai env berisi banyak entri dipisah koma.
func daftar(v string) []string {
	var out []string
	for _, s := range strings.Split(v, ",") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// spaHandler menyajikan file statis dari dir; bila path tak ada, fallback ke
// index.html (agar vue-router mode history bekerja). Bila build belum ada,
// tampilkan pesan ramah.
func spaHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Jalur /api yang tidak terdaftar dijawab 404, bukan halaman SPA.
		//
		// Tanpa ini setiap /api/apa-saja jatuh ke sini dan dibalas 200 berisi
		// index.html — termasuk POST. Endpoint yang baru saja dihapus (tambah
		// stok, ubah/hapus obat) jadi tampak "berhasil": aplikasi versi lama
		// yang masih terbuka di browser mengirim permintaannya, menerima 200,
		// dan baru gagal saat mencoba membaca JSON dari HTML — dengan pesan
		// yang tidak menyinggung sebab sesungguhnya sama sekali.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"endpoint tidak dikenal"}`))
			return
		}
		clean := filepath.Clean(r.URL.Path)
		full := filepath.Join(dir, clean)
		if info, err := os.Stat(full); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		index := filepath.Join(dir, "index.html")
		if _, err := os.Stat(index); err != nil {
			http.Error(w,
				"Frontend Vue belum di-build.\n\nSaat development, buka Vite dev server (mis. http://localhost:5173).\nUntuk produksi: jalankan `npm run build` di folder web/ lalu restart server ini.",
				http.StatusServiceUnavailable)
			return
		}
		http.ServeFile(w, r, index)
	})
}

func main() {
	// Muat .env dari folder kerja (server/). Diamkan bila tak ada.
	_ = godotenv.Load()

	produksi := strings.EqualFold(env("APP_ENV", "development"), "production")

	database := db.MustOpen()
	defer database.Close()

	srv, err := api.New(database, api.Config{Produksi: produksi})
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	// Akun disiapkan sebelum server menerima request. Kalau ini gagal, lebih
	// baik tidak hidup sama sekali daripada hidup tanpa penjaga.
	ctxAwal, batalAwal := context.WithTimeout(context.Background(), 30*time.Second)
	if err := api.PastikanAkun(ctxAwal, database, env("ADMIN_USERNAME", "bidan"), os.Getenv("ADMIN_PASSWORD")); err != nil {
		batalAwal()
		log.Fatal(err)
	}
	batalAwal()

	mux := srv.Routes()
	mux.Handle("/", spaHandler(env("WEB_DIR", "../web/dist")))

	// Urutan pembungkus: header keamanan terluar supaya ikut terpasang pada
	// balasan 429 dan 413 juga; batas body terdalam supaya berlaku tepat
	// sebelum handler membaca.
	var handler http.Handler = mux
	handler = api.BatasBody(1<<20, handler) // 1 MB
	handler = api.CORS(daftar(os.Getenv("CORS_ORIGINS")), handler)
	handler = api.BatasLaju(20, 60, handler)
	handler = api.HeaderKeamanan(produksi, handler)

	// Di produksi server ini berada di belakang Caddy, jadi cukup mendengar di
	// loopback. Kalau ufw suatu saat salah setel, port ini tetap tidak terlihat
	// dari internet karena memang tidak terikat ke antarmuka publik.
	alamat := env("BIND_ADDR", "127.0.0.1") + ":" + env("PORT", "4000")

	httpSrv := &http.Server{
		Addr:    alamat,
		Handler: handler,
		// Tanpa batas waktu, satu sambungan yang menggantung menahan slotnya
		// selamanya — cukup beberapa untuk membuat server berhenti melayani.
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// Sapu sesi mati secara berkala supaya tabelnya tidak tumbuh selamanya.
	ctxSapu, hentikanSapu := context.WithCancel(context.Background())
	defer hentikanSapu()
	go func() {
		t := time.NewTicker(6 * time.Hour)
		defer t.Stop()
		srv.SapuSesiKedaluwarsa(ctxSapu)
		for {
			select {
			case <-ctxSapu.Done():
				return
			case <-t.C:
				srv.SapuSesiKedaluwarsa(ctxSapu)
			}
		}
	}()

	// Berhenti dengan rapi: request yang sedang berjalan diberi waktu selesai
	// dulu, jadi restart saat deploy tidak memutus bidan di tengah menyimpan.
	selesai := make(chan struct{})
	go func() {
		sinyal := make(chan os.Signal, 1)
		signal.Notify(sinyal, os.Interrupt, syscall.SIGTERM)
		<-sinyal
		log.Println("menghentikan server…")
		ctx, batal := context.WithTimeout(context.Background(), 15*time.Second)
		defer batal()
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Printf("shutdown: %v", err)
		}
		close(selesai)
	}()

	log.Printf("Klinik Bidan Pit (Go) mendengar di %s (produksi=%v)", alamat, produksi)
	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	<-selesai
	log.Println("server berhenti")
}
