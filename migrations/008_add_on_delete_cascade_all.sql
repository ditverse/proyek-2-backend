-- =====================================================
-- MIGRATION: Add ON DELETE CASCADE to peminjaman related tables
-- Tanggal: 2026-01-30
-- Deskripsi: Menambahkan ON DELETE CASCADE pada foreign key yang
--            mereferensi tabel peminjaman agar data child terhapus otomatis
-- =====================================================

-- ==========================================
-- 1. KEHADIRAN_PEMINJAM
-- ==========================================
ALTER TABLE public.kehadiran_peminjam
DROP CONSTRAINT IF EXISTS kehadiran_peminjam_kode_peminjaman_fkey;

ALTER TABLE public.kehadiran_peminjam
ADD CONSTRAINT kehadiran_peminjam_kode_peminjaman_fkey 
FOREIGN KEY (kode_peminjaman) 
REFERENCES public.peminjaman(kode_peminjaman) 
ON DELETE CASCADE;

-- ==========================================
-- 2. LOG_AKTIVITAS
-- ==========================================
ALTER TABLE public.log_aktivitas
DROP CONSTRAINT IF EXISTS log_aktivitas_kode_peminjaman_fkey;

ALTER TABLE public.log_aktivitas
ADD CONSTRAINT log_aktivitas_kode_peminjaman_fkey 
FOREIGN KEY (kode_peminjaman) 
REFERENCES public.peminjaman(kode_peminjaman) 
ON DELETE SET NULL;

-- ==========================================
-- 3. PEMINJAMAN_BARANG
-- ==========================================
ALTER TABLE public.peminjaman_barang
DROP CONSTRAINT IF EXISTS peminjaman_barang_kode_peminjaman_fkey;

ALTER TABLE public.peminjaman_barang
ADD CONSTRAINT peminjaman_barang_kode_peminjaman_fkey 
FOREIGN KEY (kode_peminjaman) 
REFERENCES public.peminjaman(kode_peminjaman) 
ON DELETE CASCADE;

-- ==========================================
-- 4. MAILBOX (pastikan sudah CASCADE)
-- ==========================================
ALTER TABLE public.mailbox
DROP CONSTRAINT IF EXISTS mailbox_kode_peminjaman_fkey;

ALTER TABLE public.mailbox
ADD CONSTRAINT mailbox_kode_peminjaman_fkey 
FOREIGN KEY (kode_peminjaman) 
REFERENCES public.peminjaman(kode_peminjaman) 
ON DELETE CASCADE;

-- ==========================================
-- BONUS: FK lainnya yang mungkin bermasalah
-- ==========================================

-- Kehadiran diverifikasi_oleh -> users (SET NULL jika user dihapus)
ALTER TABLE public.kehadiran_peminjam
DROP CONSTRAINT IF EXISTS kehadiran_peminjam_diverifikasi_oleh_fkey;

ALTER TABLE public.kehadiran_peminjam
ADD CONSTRAINT kehadiran_peminjam_diverifikasi_oleh_fkey 
FOREIGN KEY (diverifikasi_oleh) 
REFERENCES public.users(kode_user) 
ON DELETE SET NULL;

-- Log aktivitas kode_user -> users (SET NULL jika user dihapus)
ALTER TABLE public.log_aktivitas
DROP CONSTRAINT IF EXISTS log_aktivitas_kode_user_fkey;

ALTER TABLE public.log_aktivitas
ADD CONSTRAINT log_aktivitas_kode_user_fkey 
FOREIGN KEY (kode_user) 
REFERENCES public.users(kode_user) 
ON DELETE SET NULL;

-- Mailbox kode_user -> users (CASCADE jika user dihapus)
ALTER TABLE public.mailbox
DROP CONSTRAINT IF EXISTS mailbox_kode_user_fkey;

ALTER TABLE public.mailbox
ADD CONSTRAINT mailbox_kode_user_fkey 
FOREIGN KEY (kode_user) 
REFERENCES public.users(kode_user) 
ON DELETE CASCADE;

-- Peminjaman kode_user -> users (SET NULL jika user dihapus)
ALTER TABLE public.peminjaman
DROP CONSTRAINT IF EXISTS peminjaman_kode_user_fkey;

ALTER TABLE public.peminjaman
ADD CONSTRAINT peminjaman_kode_user_fkey 
FOREIGN KEY (kode_user) 
REFERENCES public.users(kode_user) 
ON DELETE SET NULL;

-- Peminjaman verified_by -> users (SET NULL jika verifier dihapus)
ALTER TABLE public.peminjaman
DROP CONSTRAINT IF EXISTS peminjaman_verified_by_fkey;

ALTER TABLE public.peminjaman
ADD CONSTRAINT peminjaman_verified_by_fkey 
FOREIGN KEY (verified_by) 
REFERENCES public.users(kode_user) 
ON DELETE SET NULL;

-- Peminjaman kode_kegiatan -> kegiatan (SET NULL jika kegiatan dihapus)
ALTER TABLE public.peminjaman
DROP CONSTRAINT IF EXISTS peminjaman_kode_kegiatan_fkey;

ALTER TABLE public.peminjaman
ADD CONSTRAINT peminjaman_kode_kegiatan_fkey 
FOREIGN KEY (kode_kegiatan) 
REFERENCES public.kegiatan(kode_kegiatan) 
ON DELETE SET NULL;

-- Peminjaman kode_ruangan -> ruangan (SET NULL jika ruangan dihapus)
ALTER TABLE public.peminjaman
DROP CONSTRAINT IF EXISTS peminjaman_kode_ruangan_fkey;

ALTER TABLE public.peminjaman
ADD CONSTRAINT peminjaman_kode_ruangan_fkey 
FOREIGN KEY (kode_ruangan) 
REFERENCES public.ruangan(kode_ruangan) 
ON DELETE SET NULL;

-- Barang ruangan_kode -> ruangan (SET NULL jika ruangan dihapus)
ALTER TABLE public.barang
DROP CONSTRAINT IF EXISTS barang_ruangan_kode_fkey;

ALTER TABLE public.barang
ADD CONSTRAINT barang_ruangan_kode_fkey 
FOREIGN KEY (ruangan_kode) 
REFERENCES public.ruangan(kode_ruangan) 
ON DELETE SET NULL;

-- Peminjaman_barang kode_barang -> barang (CASCADE jika barang dihapus)
ALTER TABLE public.peminjaman_barang
DROP CONSTRAINT IF EXISTS peminjaman_barang_kode_barang_fkey;

ALTER TABLE public.peminjaman_barang
ADD CONSTRAINT peminjaman_barang_kode_barang_fkey 
FOREIGN KEY (kode_barang) 
REFERENCES public.barang(kode_barang) 
ON DELETE CASCADE;

-- =====================================================
-- RINGKASAN ON DELETE BEHAVIOR:
-- 
-- CASCADE (child ikut terhapus):
--   - kehadiran_peminjam.kode_peminjaman -> peminjaman
--   - peminjaman_barang.kode_peminjaman -> peminjaman
--   - peminjaman_barang.kode_barang -> barang
--   - mailbox.kode_peminjaman -> peminjaman
--   - mailbox.kode_user -> users
--
-- SET NULL (child tetap, referensi jadi NULL):
--   - log_aktivitas.kode_peminjaman -> peminjaman
--   - log_aktivitas.kode_user -> users
--   - kehadiran_peminjam.diverifikasi_oleh -> users
--   - peminjaman.kode_user -> users
--   - peminjaman.verified_by -> users
--   - peminjaman.kode_kegiatan -> kegiatan
--   - peminjaman.kode_ruangan -> ruangan
--   - barang.ruangan_kode -> ruangan
--   - users.organisasi_kode -> organisasi
--   - kegiatan.organisasi_kode -> organisasi
-- =====================================================
