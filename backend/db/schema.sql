-- Skema database Klinik Bidan Pit
-- MySQL 8.x. Jalankan: mysql -u root -p < schema.sql
-- Bentuk data mengikuti "State Management" pada README design handoff.

CREATE DATABASE IF NOT EXISTS app_klinik
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE app_klinik;

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

-- Agenda / jadwal bidan
CREATE TABLE IF NOT EXISTS events (
  id       VARCHAR(40)  NOT NULL PRIMARY KEY,
  title    VARCHAR(255) NOT NULL,
  all_day  TINYINT(1)   NOT NULL DEFAULT 0,
  start_ts BIGINT       NOT NULL,                       -- epoch ms
  end_ts   BIGINT       NOT NULL,                       -- epoch ms
  INDEX idx_events_start (start_ts)
) ENGINE=InnoDB;

-- Stok obat. qty selalu dalam BASE UNIT (butir/botol/sachet).
CREATE TABLE IF NOT EXISTS medicines (
  id   VARCHAR(40)                         NOT NULL PRIMARY KEY,
  name VARCHAR(255)                        NOT NULL,
  cat  ENUM('Tablet','Sirup','Sachet')     NOT NULL DEFAULT 'Tablet',
  qty  INT                                 NOT NULL DEFAULT 0,
  INDEX idx_medicines_name (name)
) ENGINE=InnoDB;

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
