-- ENUM DEFINITIONS
CREATE TYPE "role_enum" AS ENUM (
  'MAHASISWA',
  'SARPRAS',
  'SECURITY',
  'ADMIN'
);

CREATE TYPE "jenis_org_enum" AS ENUM (
  'UKM',
  'ORMAWA'
);

CREATE TYPE "peminjaman_status_enum" AS ENUM (
  'PENDING',
  'APPROVED',
  'REJECTED',
  'ONGOING',
  'FINISHED',
  'CANCELLED'
);

CREATE TYPE "mailbox_status_enum" AS ENUM (
  'TERKIRIM',
  'DIBACA'
);

CREATE TYPE "mailbox_jenis_enum" AS ENUM (
  'PENGAJUAN_DIBUAT',
  'STATUS_APPROVED',
  'STATUS_REJECTED',
  'REMINDER_KEHADIRAN'
);

CREATE TYPE kehadiran_status_enum AS ENUM (
  'HADIR',
  'TIDAK_HADIR',
  'BATAL'
);


-- TABLES
CREATE TABLE "organisasi" (
  "kode_organisasi" varchar PRIMARY KEY,
  "nama" varchar,
  "jenis_organisasi" jenis_org_enum,
  "kontak" varchar,
  "created_at" timestamp DEFAULT now()
);

CREATE TABLE "users" (
  "kode_user" varchar PRIMARY KEY,
  "nama" varchar,
  "email" varchar UNIQUE,
  "password_hash" varchar,
  "role" role_enum,
  "organisasi_kode" varchar, 
  "created_at" timestamp DEFAULT now()
);

CREATE TABLE "ruangan" (
  "kode_ruangan" varchar PRIMARY KEY,
  "nama_ruangan" varchar,
  "lokasi" varchar,
  "kapasitas" int,
  "deskripsi" text
);

CREATE TABLE "barang" (
  "kode_barang" varchar PRIMARY KEY,
  "nama_barang" varchar,
  "deskripsi" text,
  "jumlah_total" int,
  "ruangan_kode" varchar
);

CREATE TABLE "kegiatan" (
  "kode_kegiatan" varchar PRIMARY KEY,
  "nama_kegiatan" varchar,
  "deskripsi" text,
  "tanggal_mulai" timestamp,
  "tanggal_selesai" timestamp,
  "organisasi_kode" varchar,
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp
);

CREATE TABLE "peminjaman" (
  "kode_peminjaman" varchar PRIMARY KEY,
  "kode_user" varchar,
  "kode_ruangan" varchar,
  "kode_kegiatan" varchar,
  "tanggal_mulai" timestamp,
  "tanggal_selesai" timestamp,
  "keperluan" text,
  "status" peminjaman_status_enum,
  "path_surat_digital" text,
  "verified_by" varchar,
  "verified_at" timestamp,
  "catatan_verifikasi" text,
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp
);

CREATE TABLE "peminjaman_barang" (
  "kode_peminjaman_barang" varchar PRIMARY KEY,
  "kode_peminjaman" varchar,
  "kode_barang" varchar,
  "jumlah" int
);

CREATE TABLE "mailbox" (
  "kode_mailbox" varchar PRIMARY KEY,
  "kode_user" varchar,
  "kode_peminjaman" varchar,
  "jenis_mailbox" mailbox_jenis_enum,
  "pesan" text,
  "status" mailbox_status_enum,
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp
);

CREATE TABLE "log_aktivitas" (
  "kode_log" varchar PRIMARY KEY,
  "kode_user" varchar,
  "kode_peminjaman" varchar,
  "aksi" varchar,
  "keterangan" text,
  "waktu" timestamp DEFAULT now(),
  "updated_at" timestamp
);

CREATE TABLE kehadiran_peminjam (
  "kode_kehadiran" varchar PRIMARY KEY,
  "kode_peminjaman" varchar REFERENCES peminjaman (kode_peminjaman),
  "status_kehadiran" kehadiran_status_enum,
  "waktu_konfirmasi" timestamp DEFAULT now(),
  "keterangan" text,
  "diverifikasi_oleh" varchar REFERENCES users (kode_user),
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp DEFAULT now()
);

-- FOREIGN KEYS
ALTER TABLE "users" 
  ADD FOREIGN KEY ("organisasi_kode") REFERENCES "organisasi" ("kode_organisasi");

ALTER TABLE "barang" 
  ADD FOREIGN KEY ("ruangan_kode") REFERENCES "ruangan" ("kode_ruangan");

ALTER TABLE "kegiatan" 
  ADD FOREIGN KEY ("organisasi_kode") REFERENCES "organisasi" ("kode_organisasi");

ALTER TABLE "peminjaman" 
  ADD FOREIGN KEY ("kode_user") REFERENCES "users" ("kode_user");

ALTER TABLE "peminjaman" 
  ADD FOREIGN KEY ("kode_ruangan") REFERENCES "ruangan" ("kode_ruangan");

ALTER TABLE "peminjaman" 
  ADD FOREIGN KEY ("kode_kegiatan") REFERENCES "kegiatan" ("kode_kegiatan");

ALTER TABLE "peminjaman" 
  ADD FOREIGN KEY ("verified_by") REFERENCES "users" ("kode_user");

ALTER TABLE "peminjaman_barang" 
  ADD FOREIGN KEY ("kode_peminjaman") REFERENCES "peminjaman" ("kode_peminjaman");

ALTER TABLE "peminjaman_barang" 
  ADD FOREIGN KEY ("kode_barang") REFERENCES "barang" ("kode_barang");

ALTER TABLE "mailbox" 
  ADD FOREIGN KEY ("kode_user") REFERENCES "users" ("kode_user");

ALTER TABLE "mailbox" 
  ADD FOREIGN KEY ("kode_peminjaman") REFERENCES "peminjaman" ("kode_peminjaman");

ALTER TABLE "log_aktivitas" 
  ADD FOREIGN KEY ("kode_user") REFERENCES "users" ("kode_user");

ALTER TABLE "log_aktivitas" 
  ADD FOREIGN KEY ("kode_peminjaman") REFERENCES "peminjaman" ("kode_peminjaman");
