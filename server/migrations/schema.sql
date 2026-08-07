-- Skema database Klinik Bidan Pit
-- MySQL 8.x. Bentuk data mengikuti "State Management" pada README design handoff.
--
-- Jalankan dengan menyebut databasenya di baris perintah:
--
--   mysql -u root -p app_klinik < schema.sql
--
-- Berkas ini SENGAJA tidak memuat CREATE DATABASE maupun USE. Sebelumnya ia
-- memuat keduanya, dan itu berarti dijalankan di mana pun ia selalu menyentuh
-- database bernama `app_klinik` — termasuk saat yang kamu maksud database uji,
-- atau saat di server sudah ada `app_klinik` berisi data pasien sungguhan.
-- Nama database ditentukan pemanggilnya; berkas ini hanya menggambarkan bentuk
-- tabelnya.
--
-- Databasenya dibuat terpisah, dengan collation yang eksplisit:
--
--   CREATE DATABASE app_klinik CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
--
-- Collation itu penting disebut, bukan dibiarkan ikut bawaan server: ia
-- menentukan kapan dua string dianggap sama. Perbandingan nama obat di aplikasi
-- ini bergantung padanya.

-- Akun untuk masuk ke Bidan App. Praktiknya hanya satu (bidan), tapi tabel
-- tetap dipakai supaya menambah akun kedua tidak perlu mengubah skema.
-- Kata sandi disimpan sebagai hash bcrypt — tidak pernah dalam bentuk asli.
CREATE TABLE IF NOT EXISTS users (
  id            INT AUTO_INCREMENT PRIMARY KEY,
  username      VARCHAR(100) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  created_at    BIGINT       NOT NULL,                   -- epoch ms
  updated_at    BIGINT       NOT NULL
) ENGINE=InnoDB;

-- Satu baris = satu sesi login di satu perangkat.
--
-- Keberadaan tabel inilah yang membuat sesi bisa DICABUT. Dengan token yang
-- berdiri sendiri (mis. JWT tanpa catatan di server), token yang bocor tetap
-- sah sampai kedaluwarsa dan tidak ada cara menghentikannya.
--
-- Yang disimpan hanya SHA-256 dari tokennya, bukan tokennya. Kalau isi
-- database bocor, isinya tidak bisa dipakai untuk masuk.
CREATE TABLE IF NOT EXISTS sessions (
  id           INT AUTO_INCREMENT PRIMARY KEY,
  user_id      INT          NOT NULL,
  token_hash   CHAR(64)     NOT NULL UNIQUE,             -- hex SHA-256
  user_agent   VARCHAR(255) NOT NULL DEFAULT '',
  ip           VARCHAR(64)  NOT NULL DEFAULT '',
  created_at   BIGINT       NOT NULL,                    -- epoch ms
  last_seen_at BIGINT       NOT NULL,
  expires_at   BIGINT       NOT NULL,
  revoked_at   BIGINT       NULL,
  CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_sessions_expires (expires_at)
) ENGINE=InnoDB;

-- Status bidan (baris tunggal, id selalu = 1)
CREATE TABLE IF NOT EXISTS clinic_status (
  id          TINYINT      NOT NULL PRIMARY KEY DEFAULT 1,
  bidan_hadir TINYINT(1)   NOT NULL DEFAULT 1,
  away_note   VARCHAR(255) NOT NULL DEFAULT '',
  away_until  BIGINT       NULL,                       -- epoch ms, null = belum dipastikan
  updated_ts  BIGINT       NOT NULL,                   -- epoch ms terakhir diubah
  clinic      VARCHAR(120) NOT NULL DEFAULT 'Klinik Bidan Pit',
  CONSTRAINT chk_status_singleton CHECK (id = 1)
) ENGINE=InnoDB;

-- Profil klinik yang tampil di halaman pasien (baris tunggal, id selalu = 1).
--
-- Tabel TERPISAH, bukan kolom tambahan di clinic_status. Alasannya praktis:
-- proyek ini tidak punya migration runner — schema.sql dijalankan ulang apa
-- adanya — dan MySQL 8 tidak mengenal ALTER TABLE ... ADD COLUMN IF NOT EXISTS.
-- Tabel baru dengan CREATE TABLE IF NOT EXISTS aman dijalankan di database
-- kosong maupun di database produksi yang sudah berisi data.
--
-- Nama kliniknya sendiri TIDAK di sini: ia sudah tinggal di clinic_status.clinic
-- sejak awal, dan menyalinnya ke sini akan membuat dua sumber kebenaran untuk
-- satu nama yang sama.
--
-- Alamat & nomor boleh kosong, dan itu bukan kelalaian: halaman pasien
-- menyembunyikan bagiannya kalau kosong. Tidak menampilkan apa-apa lebih baik
-- daripada menampilkan alamat karangan, yang justru keadaan sebelum tabel ini
-- ada — `Jl. Melati No. 12` dan `wa.me/6281234567890` tertulis mati di
-- Patient.vue dan sudah terlihat publik.
CREATE TABLE IF NOT EXISTS clinic_profile (
  id         TINYINT      NOT NULL PRIMARY KEY DEFAULT 1,
  alamat     VARCHAR(255) NOT NULL DEFAULT '',
  whatsapp   VARCHAR(20)  NOT NULL DEFAULT '',   -- hanya angka, diawali 62
  updated_ts BIGINT       NOT NULL,
  CONSTRAINT chk_profile_singleton CHECK (id = 1)
) ENGINE=InnoDB;

-- Agenda / jadwal bidan
CREATE TABLE IF NOT EXISTS events (
  id       VARCHAR(40)  NOT NULL PRIMARY KEY,
  title    VARCHAR(255) NOT NULL,
  all_day  TINYINT(1)   NOT NULL DEFAULT 0,
  start_ts BIGINT       NOT NULL,                       -- epoch ms
  end_ts   BIGINT       NOT NULL,                       -- epoch ms
  INDEX idx_events_start (start_ts)
) ENGINE=InnoDB;

-- Tabel `medicines` (stok obat) SUDAH TIDAK ADA di berkas ini, dan itu
-- disengaja: pengelolaan persediaan dilepas atas permintaan bidan — yang
-- dicatat cukup obat apa yang diberikan ke pasien, bukan berapa sisa di lemari.
--
-- Database yang sudah berjalan tetap memiliki tabelnya; tidak ada satu pun
-- kode yang menyentuhnya lagi, jadi ia hanya menganggur. Menghapusnya dilakukan
-- sendiri bila diinginkan, SETELAH versi ini jalan dengan baik:
--
--   DROP TABLE medicines;
--
-- Riwayat kunjungan tidak terpengaruh: visit_items menyimpan nama & satuan
-- obatnya sendiri sebagai teks, tanpa foreign key ke tabel itu.

-- Kunjungan pasien
CREATE TABLE IF NOT EXISTS visits (
  id     VARCHAR(40)  NOT NULL PRIMARY KEY,
  name   VARCHAR(255) NOT NULL,
  age    INT          NOT NULL DEFAULT 0,
  gejala TEXT         NULL,
  ts     BIGINT       NOT NULL,                          -- epoch ms
  INDEX idx_visits_ts (ts)
) ENGINE=InnoDB;

-- Obat yang diberikan pada tiap kunjungan (unit = base unit saat pencatatan)
CREATE TABLE IF NOT EXISTS visit_items (
  id       INT AUTO_INCREMENT PRIMARY KEY,
  visit_id VARCHAR(40)  NOT NULL,
  name     VARCHAR(255) NOT NULL,
  qty      INT          NOT NULL,
  unit     VARCHAR(40)  NOT NULL,
  CONSTRAINT fk_visit_items_visit
    FOREIGN KEY (visit_id) REFERENCES visits(id) ON DELETE CASCADE,
  INDEX idx_visit_items_visit (visit_id)
) ENGINE=InnoDB;
