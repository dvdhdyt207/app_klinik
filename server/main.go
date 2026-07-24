// Klinik Bidan Pit — server Go.
// Menyajikan REST API (/api/*) + halaman web hasil build Vue (SPA fallback).
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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

// withCORS mengizinkan Vite dev server (origin lain) memanggil API saat development.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// spaHandler menyajikan file statis dari dir; bila path tak ada, fallback ke
// index.html (agar vue-router mode history bekerja). Bila build belum ada,
// tampilkan pesan ramah.
func spaHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	database := db.MustOpen()
	defer database.Close()

	srv := &api.Server{DB: database}
	mux := srv.Routes()
	mux.Handle("/", spaHandler(env("WEB_DIR", "../web/dist")))

	port := env("PORT", "4000")
	log.Printf("Klinik Bidan Pit (Go) berjalan di http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
