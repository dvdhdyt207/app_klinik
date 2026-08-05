# syntax=docker/dockerfile:1
#
# Satu image berisi seluruh aplikasi Klinik Bidan Pit: API Go sekaligus halaman
# Vue-nya. Server Go memang sudah menyajikan build Vue sendiri (spaHandler di
# main.go), jadi tidak ada yang perlu disusun ulang — cukup keduanya masuk ke
# image yang sama.
#
#   docker compose build
#   docker compose up -d
#
# Tiga tahap supaya image akhir hanya memuat yang benar-benar dijalankan. Node
# dan toolchain Go — ratusan MB — berhenti di tahap build.

# --------------------------------------------------------------- 1. frontend
FROM node:22-alpine AS frontend
WORKDIR /src

# Daftar dependensi disalin sendirian dulu. Selama package.json & lock-nya tidak
# berubah, Docker memakai ulang lapisan `npm ci` — di server 1 vCPU itu bedanya
# menit, bukan detik.
COPY web/package.json web/package-lock.json ./
# `npm ci` bukan `npm install`: ia memasang persis versi yang tercatat di
# package-lock.json dan gagal bila keduanya tidak sepakat. Build di server harus
# menghasilkan hal yang sama dengan build di laptop.
RUN npm ci

COPY web/ ./
RUN npm run build

# ---------------------------------------------------------------- 2. backend
FROM golang:1.26-alpine AS backend
WORKDIR /src

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./
# CGO_ENABLED=0 menghasilkan biner statis yang tidak bergantung pada libc,
# sehingga bisa dijalankan di image runtime yang ramping.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /klinik-server .

# ---------------------------------------------------------------- 3. runtime
FROM alpine:3.20

# tzdata wajib. Biner statis tidak membawa basis data zona waktu, jadi tanpa ini
# container berjalan di UTC — dan waktu kunjungan pasien tercatat tujuh jam
# meleset dari kenyataan, tanpa satu pun error yang menandainya.
RUN apk add --no-cache ca-certificates tzdata

# Aplikasi tidak pernah menulis apa pun ke dalam container, jadi ia berjalan
# sebagai pengguna biasa tanpa home.
RUN adduser -D -H -u 10001 klinik

WORKDIR /app
COPY --from=backend  /klinik-server /app/klinik-server
COPY --from=frontend /src/dist      /app/web

# Dibaca main.go lewat env(). Bawaannya ../web/dist yang benar saat dijalankan
# dari folder server/ di laptop, tapi tidak ada artinya di dalam container.
ENV WEB_DIR=/app/web

USER klinik
ENTRYPOINT ["/app/klinik-server"]
