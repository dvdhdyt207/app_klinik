#!/usr/bin/env bash
#
# Memasang versi baru Klinik Bidan Pit dari image yang sudah dibangun CI.
#
# Pasang di server (milik root, tidak bisa ditulis deployer):
#   sudo install -m 755 -o root -g root klinik-deploy.sh /usr/local/sbin/klinik-deploy
#
# Izinkan CI memanggilnya. Perhatikan berkas sudoers-nya BUKAN /etc/sudoers.d/deployer
# — berkas itu sudah dipakai catering. Dua baris di dua berkas terpisah supaya
# menghapus salah satu aplikasi tidak diam-diam mencabut hak aplikasi lainnya:
#   echo 'deployer ALL=(root) NOPASSWD: /usr/local/sbin/klinik-deploy' \
#     | sudo tee /etc/sudoers.d/deployer-klinik
#   sudo visudo -cf /etc/sudoers.d/deployer-klinik && sudo chmod 440 /etc/sudoers.d/deployer-klinik
#
# Skrip ini TIDAK PERNAH dikirim lewat CI. Kalau ia ikut terkirim tiap deploy,
# aturan sudoers sempit di atas tidak menahan apa pun: siapa pun yang bisa
# mengubah workflow tinggal mengirim skrip berisi apa saja dan menjalankannya
# sebagai root. Yang boleh melintas lewat jalur CI hanya data, bukan perintah.
set -euo pipefail

APP_DIR=/opt/klinik
PEMILIK=ProjectRunning
SEHAT_URL=http://127.0.0.1:4000/api/health
BATAS_DETIK=60

cd "$APP_DIR"

# --- Konfigurasi ---------------------------------------------------------
# compose.yml ikut versi git, jadi perubahannya harus sampai ke server. Yang
# TIDAK ikut: .env, karena ia di luar git — jadi `reset --hard` tidak
# menyentuhnya.
#
# git dijalankan sebagai pemilik folder, bukan root: deploy key ada di
# ~ProjectRunning/.ssh, dan root punya konfigurasi ssh sendiri yang tidak
# mengenalnya.
#
# `reset --hard`, bukan `pull`: server harus sama persis dengan origin, bukan
# hasil penggabungan dengan apa pun yang mungkin pernah disunting di sana.
sudo -u "$PEMILIK" git fetch --quiet origin main
sudo -u "$PEMILIK" git reset --hard --quiet origin/main
echo "konfigurasi pada $(git rev-parse --short HEAD)"

# --- Image ---------------------------------------------------------------
# Ditarik, bukan dibangun. Yang mendarat di produksi adalah image yang benar-
# benar dibangun CI dari commit ini, bukan hasil build ulang di server.
docker compose pull app
docker compose up -d

# --- Pembuktian ----------------------------------------------------------
for ((i = 1; i <= BATAS_DETIK / 2; i++)); do
  if curl -sf -o /dev/null "$SEHAT_URL"; then
    echo "sehat setelah ~$((i * 2)) detik"
    exit 0
  fi
  sleep 2
done

# Gagal: cetak yang dibutuhkan untuk mendiagnosis, karena deployer tidak boleh
# menjalankan docker sendiri dan tidak bisa mengambilnya belakangan.
echo "TIDAK SEHAT setelah ${BATAS_DETIK} detik" >&2
docker compose ps
docker compose logs --tail 50 app
exit 1
