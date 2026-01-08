-- ============================================================================
-- DATABASE SEEDER: Initial Data for Development/Testing
-- ============================================================================
-- Jalankan script ini di Supabase SQL Editor setelah cleanup dummy data
-- 
-- Data yang di-seed:
-- - 3 Organisasi (HIMATIF, BEM, UKM Badminton)
-- - 3 Users (Mahasiswa, Security, Sarpras)
-- ============================================================================

-- =====================
-- STEP 1: INSERT ORGANISASI (Parent table - harus diinsert dulu)
-- =====================
INSERT INTO organisasi (kode_organisasi, nama, jenis_organisasi, kontak) VALUES
    ('HIMATIF', 'Himpunan Mahasiswa Teknik Informatika', 'ORMAWA', NULL),
    ('BEM', 'Badan Eksekutif Mahasiswa', 'ORMAWA', NULL),
    ('UKM-BADMINTON', 'UKM Badminton', 'UKM', NULL);

-- =====================
-- STEP 2: INSERT USERS
-- =====================
-- Password hashes generated using bcrypt (cost 10):
-- - "mahasiswa" -> $2a$10$hYlNynmV0XwX9Os5K4en2O0UozPwXCl4vVT/y2OruhZqmtF5VQQvm
-- - "hidupmahasiswa" -> $2a$10$Xe7vct6Iyjmm.RLKNzPXVe1hVfiSYYTgMmq10wr0mGuV5DTmDRP5y

INSERT INTO users (kode_user, nama, email, password_hash, role, organisasi_kode) VALUES
    -- Mahasiswa (password: mahasiswa)
    ('USR-MHS-001', 'Mahasiswa', 'mahasiswaulbi54@gmail.com', 
     '$2a$10$hYlNynmV0XwX9Os5K4en2O0UozPwXCl4vVT/y2OruhZqmtF5VQQvm', 
     'MAHASISWA', 'HIMATIF'),
    
    -- Security (password: hidupmahasiswa)
    ('USR-SEC-001', 'Security', 'securityulbi54@gmail.com', 
     '$2a$10$Xe7vct6Iyjmm.RLKNzPXVe1hVfiSYYTgMmq10wr0mGuV5DTmDRP5y', 
     'SECURITY', NULL),
    
    -- Sarpras (password: hidupmahasiswa)
    ('USR-SAR-001', 'Sarpras', 'sarprasulbi54@gmail.com', 
     '$2a$10$Xe7vct6Iyjmm.RLKNzPXVe1hVfiSYYTgMmq10wr0mGuV5DTmDRP5y', 
     'SARPRAS', NULL);

-- ============================================================================
-- VERIFIKASI: Pastikan data sudah ter-insert
-- ============================================================================
SELECT 'organisasi' as tabel, COUNT(*) as jumlah FROM organisasi
UNION ALL SELECT 'users', COUNT(*) FROM users;

-- ============================================================================
-- DETAIL DATA YANG DI-INSERT
-- ============================================================================
-- 
-- ORGANISASI:
-- | kode_organisasi | nama                                    | jenis     |
-- |-----------------|-----------------------------------------|-----------|
-- | HIMATIF         | Himpunan Mahasiswa Teknik Informatika   | ORMAWA    |
-- | BEM             | Badan Eksekutif Mahasiswa               | ORMAWA    |
-- | UKM-BADMINTON   | UKM Badminton                           | UKM       |
--
-- USERS:
-- | email                       | password       | role      | organisasi |
-- |-----------------------------|----------------|-----------|------------|
-- | mahasiswaulbi54@gmail.com   | mahasiswa      | MAHASISWA | HIMATIF    |
-- | securityulbi54@gmail.com    | hidupmahasiswa | SECURITY  | -          |
-- | sarprasulbi54@gmail.com     | hidupmahasiswa | SARPRAS   | -          |
--
-- ============================================================================
