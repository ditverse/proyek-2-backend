-- Migration: Auto-generate kode_user dengan format USR-YYMMDD-0001
-- Trigger untuk tabel users

-- Function untuk generate kode_user
CREATE OR REPLACE FUNCTION generate_kode_user()
RETURNS TRIGGER AS $$
DECLARE
    today_date TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    -- Jika kode_user sudah diisi, skip trigger
    IF NEW.kode_user IS NOT NULL AND NEW.kode_user != '' THEN
        RETURN NEW;
    END IF;

    -- Format tanggal: YYMMDD
    today_date := TO_CHAR(CURRENT_DATE, 'YYMMDD');
    
    -- Hitung urutan berikutnya untuk hari ini
    SELECT COALESCE(MAX(
        CAST(
            SUBSTRING(kode_user FROM 'USR-[0-9]{6}-([0-9]{4})') 
            AS INTEGER
        )
    ), 0) + 1
    INTO next_seq
    FROM users
    WHERE kode_user LIKE 'USR-' || today_date || '-%';
    
    -- Generate kode baru dengan format USR-YYMMDD-0001
    new_kode := 'USR-' || today_date || '-' || LPAD(next_seq::TEXT, 4, '0');
    
    NEW.kode_user := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger yang akan dijalankan sebelum INSERT
DROP TRIGGER IF EXISTS trigger_generate_kode_user ON users;
CREATE TRIGGER trigger_generate_kode_user
    BEFORE INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_user();

-- Trigger untuk tabel lainnya (opsional, sesuaikan dengan kebutuhan)

-- Function untuk generate kode_peminjaman
CREATE OR REPLACE FUNCTION generate_kode_peminjaman()
RETURNS TRIGGER AS $$
DECLARE
    today_date TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_peminjaman IS NOT NULL AND NEW.kode_peminjaman != '' THEN
        RETURN NEW;
    END IF;

    today_date := TO_CHAR(CURRENT_DATE, 'YYMMDD');
    
    SELECT COALESCE(MAX(
        CAST(
            SUBSTRING(kode_peminjaman FROM 'PMJ-[0-9]{6}-([0-9]{4})') 
            AS INTEGER
        )
    ), 0) + 1
    INTO next_seq
    FROM peminjaman
    WHERE kode_peminjaman LIKE 'PMJ-' || today_date || '-%';
    
    new_kode := 'PMJ-' || today_date || '-' || LPAD(next_seq::TEXT, 4, '0');
    
    NEW.kode_peminjaman := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_peminjaman ON peminjaman;
CREATE TRIGGER trigger_generate_kode_peminjaman
    BEFORE INSERT ON peminjaman
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_peminjaman();

-- Function untuk generate kode_ruangan
CREATE OR REPLACE FUNCTION generate_kode_ruangan()
RETURNS TRIGGER AS $$
DECLARE
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_ruangan IS NOT NULL AND NEW.kode_ruangan != '' THEN
        RETURN NEW;
    END IF;
    
    SELECT COALESCE(MAX(
        CAST(
            SUBSTRING(kode_ruangan FROM 'RNG-([0-9]{4})') 
            AS INTEGER
        )
    ), 0) + 1
    INTO next_seq
    FROM ruangan;
    
    new_kode := 'RNG-' || LPAD(next_seq::TEXT, 4, '0');
    
    NEW.kode_ruangan := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_ruangan ON ruangan;
CREATE TRIGGER trigger_generate_kode_ruangan
    BEFORE INSERT ON ruangan
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_ruangan();

-- Function untuk generate kode_barang
CREATE OR REPLACE FUNCTION generate_kode_barang()
RETURNS TRIGGER AS $$
DECLARE
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_barang IS NOT NULL AND NEW.kode_barang != '' THEN
        RETURN NEW;
    END IF;
    
    SELECT COALESCE(MAX(
        CAST(
            SUBSTRING(kode_barang FROM 'BRG-([0-9]{4})') 
            AS INTEGER
        )
    ), 0) + 1
    INTO next_seq
    FROM barang;
    
    new_kode := 'BRG-' || LPAD(next_seq::TEXT, 4, '0');
    
    NEW.kode_barang := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_barang ON barang;
CREATE TRIGGER trigger_generate_kode_barang
    BEFORE INSERT ON barang
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_barang();
