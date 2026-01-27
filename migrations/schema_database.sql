-- WARNING: This schema is for context only and is not meant to be run.
-- Table order and constraints may not be valid for execution.

CREATE TABLE public.barang (
  kode_barang character varying NOT NULL,
  nama_barang character varying,
  deskripsi text,
  jumlah_total integer,
  ruangan_kode character varying,
  CONSTRAINT barang_pkey PRIMARY KEY (kode_barang),
  CONSTRAINT barang_ruangan_kode_fkey FOREIGN KEY (ruangan_kode) REFERENCES public.ruangan(kode_ruangan)
);
CREATE TABLE public.kegiatan (
  kode_kegiatan character varying NOT NULL,
  nama_kegiatan character varying,
  deskripsi text,
  tanggal_mulai timestamp without time zone,
  tanggal_selesai timestamp without time zone,
  organisasi_kode character varying,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone,
  CONSTRAINT kegiatan_pkey PRIMARY KEY (kode_kegiatan),
  CONSTRAINT kegiatan_organisasi_kode_fkey FOREIGN KEY (organisasi_kode) REFERENCES public.organisasi(kode_organisasi)
);
CREATE TABLE public.kehadiran_peminjam (
  kode_kehadiran character varying NOT NULL,
  kode_peminjaman character varying,
  status_kehadiran USER-DEFINED,
  waktu_konfirmasi timestamp without time zone DEFAULT now(),
  keterangan text,
  diverifikasi_oleh character varying,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone DEFAULT now(),
  CONSTRAINT kehadiran_peminjam_pkey PRIMARY KEY (kode_kehadiran),
  CONSTRAINT kehadiran_peminjam_kode_peminjaman_fkey FOREIGN KEY (kode_peminjaman) REFERENCES public.peminjaman(kode_peminjaman),
  CONSTRAINT kehadiran_peminjam_diverifikasi_oleh_fkey FOREIGN KEY (diverifikasi_oleh) REFERENCES public.users(kode_user)
);
CREATE TABLE public.log_aktivitas (
  kode_log character varying NOT NULL,
  kode_user character varying,
  kode_peminjaman character varying,
  aksi character varying,
  keterangan text,
  waktu timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone,
  CONSTRAINT log_aktivitas_pkey PRIMARY KEY (kode_log),
  CONSTRAINT log_aktivitas_kode_user_fkey FOREIGN KEY (kode_user) REFERENCES public.users(kode_user),
  CONSTRAINT log_aktivitas_kode_peminjaman_fkey FOREIGN KEY (kode_peminjaman) REFERENCES public.peminjaman(kode_peminjaman)
);
CREATE TABLE public.mailbox (
  kode_mailbox character varying NOT NULL,
  kode_user character varying NOT NULL,
  kode_peminjaman character varying NOT NULL,
  jenis_pesan character varying NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  CONSTRAINT mailbox_pkey PRIMARY KEY (kode_mailbox),
  CONSTRAINT mailbox_kode_user_fkey FOREIGN KEY (kode_user) REFERENCES public.users(kode_user),
  CONSTRAINT mailbox_kode_peminjaman_fkey FOREIGN KEY (kode_peminjaman) REFERENCES public.peminjaman(kode_peminjaman)
);
CREATE TABLE public.organisasi (
  kode_organisasi character varying NOT NULL,
  nama character varying,
  kontak character varying,
  created_at timestamp without time zone DEFAULT now(),
  jenis_organisasi USER-DEFINED,
  CONSTRAINT organisasi_pkey PRIMARY KEY (kode_organisasi)
);
CREATE TABLE public.peminjaman (
  kode_peminjaman character varying NOT NULL,
  kode_user character varying,
  kode_ruangan character varying,
  kode_kegiatan character varying,
  tanggal_mulai timestamp without time zone,
  tanggal_selesai timestamp without time zone,
  status USER-DEFINED,
  path_surat_digital text,
  verified_by character varying,
  verified_at timestamp without time zone,
  catatan_verifikasi text,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone,
  CONSTRAINT peminjaman_pkey PRIMARY KEY (kode_peminjaman),
  CONSTRAINT peminjaman_kode_user_fkey FOREIGN KEY (kode_user) REFERENCES public.users(kode_user),
  CONSTRAINT peminjaman_kode_ruangan_fkey FOREIGN KEY (kode_ruangan) REFERENCES public.ruangan(kode_ruangan),
  CONSTRAINT peminjaman_kode_kegiatan_fkey FOREIGN KEY (kode_kegiatan) REFERENCES public.kegiatan(kode_kegiatan),
  CONSTRAINT peminjaman_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.users(kode_user)
);
CREATE TABLE public.peminjaman_barang (
  kode_peminjaman_barang character varying NOT NULL,
  kode_peminjaman character varying,
  kode_barang character varying,
  jumlah integer,
  CONSTRAINT peminjaman_barang_pkey PRIMARY KEY (kode_peminjaman_barang),
  CONSTRAINT peminjaman_barang_kode_peminjaman_fkey FOREIGN KEY (kode_peminjaman) REFERENCES public.peminjaman(kode_peminjaman),
  CONSTRAINT peminjaman_barang_kode_barang_fkey FOREIGN KEY (kode_barang) REFERENCES public.barang(kode_barang)
);
CREATE TABLE public.ruangan (
  kode_ruangan character varying NOT NULL,
  nama_ruangan character varying,
  lokasi character varying,
  kapasitas integer,
  deskripsi text,
  CONSTRAINT ruangan_pkey PRIMARY KEY (kode_ruangan)
);
CREATE TABLE public.users (
  kode_user character varying NOT NULL,
  nama character varying,
  email character varying UNIQUE,
  password_hash character varying,
  role USER-DEFINED,
  organisasi_kode character varying,
  created_at timestamp without time zone DEFAULT now(),
  no_hp character varying,
  CONSTRAINT users_pkey PRIMARY KEY (kode_user),
  CONSTRAINT users_organisasi_kode_fkey FOREIGN KEY (organisasi_kode) REFERENCES public.organisasi(kode_organisasi)
);