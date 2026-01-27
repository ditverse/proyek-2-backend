-- Migration: Update Code Generation Triggers using YYYY-MM-DD format
-- 1. Kegiatan: KGT-YYYY-MM-DD-XXXX
-- 2. Log Aktivitas: LOG-YYYY-MM-DD-XXXX
-- 3. Kehadiran: KHD-YYYY-MM-DD-XXXX
-- 4. Peminjaman: PMJ-YYYY-MM-DD-XXXX
-- 5. User: USR-YYYY-MM-DD-XXXX

-- NOTE: Regex matches the DATE part separately or simplifies matching the suffix.

-- ==========================================
-- 1. KEGIATAN (kode_kegiatan)
-- ==========================================
CREATE OR REPLACE FUNCTION generate_kode_kegiatan()
RETURNS TRIGGER AS $$
DECLARE
    today_str TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_kegiatan IS NOT NULL AND NEW.kode_kegiatan != '' THEN
        RETURN NEW;
    END IF;

    -- Format: YYYY-MM-DD
    today_str := TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD');

    -- Find max sequence for today
    SELECT COALESCE(MAX(
        CAST(SUBSTRING(kode_kegiatan FROM 'KGT-[0-9]{4}-[0-9]{2}-[0-9]{2}-([0-9]{4})') AS INTEGER)
    ), 0) + 1
    INTO next_seq
    FROM kegiatan
    WHERE kode_kegiatan LIKE 'KGT-' || today_str || '-%';

    new_kode := 'KGT-' || today_str || '-' || LPAD(next_seq::TEXT, 4, '0');
    NEW.kode_kegiatan := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_kegiatan ON kegiatan;
CREATE TRIGGER trigger_generate_kode_kegiatan
    BEFORE INSERT ON kegiatan
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_kegiatan();


-- ==========================================
-- 2. LOG AKTIVITAS (kode_log)
-- ==========================================
CREATE OR REPLACE FUNCTION generate_kode_log()
RETURNS TRIGGER AS $$
DECLARE
    today_str TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_log IS NOT NULL AND NEW.kode_log != '' THEN
        RETURN NEW;
    END IF;

    today_str := TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD');

    SELECT COALESCE(MAX(
        CAST(SUBSTRING(kode_log FROM 'LOG-[0-9]{4}-[0-9]{2}-[0-9]{2}-([0-9]{4})') AS INTEGER)
    ), 0) + 1
    INTO next_seq
    FROM log_aktivitas
    WHERE kode_log LIKE 'LOG-' || today_str || '-%';

    new_kode := 'LOG-' || today_str || '-' || LPAD(next_seq::TEXT, 4, '0');
    NEW.kode_log := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_log ON log_aktivitas;
CREATE TRIGGER trigger_generate_kode_log
    BEFORE INSERT ON log_aktivitas
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_log();


-- ==========================================
-- 3. KEHADIRAN PEMINJAM (kode_kehadiran)
-- ==========================================
CREATE OR REPLACE FUNCTION generate_kode_kehadiran()
RETURNS TRIGGER AS $$
DECLARE
    today_str TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_kehadiran IS NOT NULL AND NEW.kode_kehadiran != '' THEN
        RETURN NEW;
    END IF;

    today_str := TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD');

    SELECT COALESCE(MAX(
        CAST(SUBSTRING(kode_kehadiran FROM 'KHD-[0-9]{4}-[0-9]{2}-[0-9]{2}-([0-9]{4})') AS INTEGER)
    ), 0) + 1
    INTO next_seq
    FROM kehadiran_peminjam
    WHERE kode_kehadiran LIKE 'KHD-' || today_str || '-%';

    new_kode := 'KHD-' || today_str || '-' || LPAD(next_seq::TEXT, 4, '0');
    NEW.kode_kehadiran := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_kehadiran ON kehadiran_peminjam;
CREATE TRIGGER trigger_generate_kode_kehadiran
    BEFORE INSERT ON kehadiran_peminjam
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_kehadiran();


-- ==========================================
-- 4. PEMINJAMAN (kode_peminjaman)
-- ==========================================
CREATE OR REPLACE FUNCTION generate_kode_peminjaman()
RETURNS TRIGGER AS $$
DECLARE
    today_str TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_peminjaman IS NOT NULL AND NEW.kode_peminjaman != '' THEN
        RETURN NEW;
    END IF;

    today_str := TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD');

    SELECT COALESCE(MAX(
        CAST(SUBSTRING(kode_peminjaman FROM 'PMJ-[0-9]{4}-[0-9]{2}-[0-9]{2}-([0-9]{4})') AS INTEGER)
    ), 0) + 1
    INTO next_seq
    FROM peminjaman
    WHERE kode_peminjaman LIKE 'PMJ-' || today_str || '-%';

    new_kode := 'PMJ-' || today_str || '-' || LPAD(next_seq::TEXT, 4, '0');
    NEW.kode_peminjaman := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_peminjaman ON peminjaman;
CREATE TRIGGER trigger_generate_kode_peminjaman
    BEFORE INSERT ON peminjaman
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_peminjaman();


-- ==========================================
-- 5. USERS (kode_user)
-- ==========================================
CREATE OR REPLACE FUNCTION generate_kode_user()
RETURNS TRIGGER AS $$
DECLARE
    today_str TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_user IS NOT NULL AND NEW.kode_user != '' THEN
        RETURN NEW;
    END IF;

    today_str := TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD');

    SELECT COALESCE(MAX(
        CAST(SUBSTRING(kode_user FROM 'USR-[0-9]{4}-[0-9]{2}-[0-9]{2}-([0-9]{4})') AS INTEGER)
    ), 0) + 1
    INTO next_seq
    FROM users
    WHERE kode_user LIKE 'USR-' || today_str || '-%';

    new_kode := 'USR-' || today_str || '-' || LPAD(next_seq::TEXT, 4, '0');
    NEW.kode_user := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_user ON users;
CREATE TRIGGER trigger_generate_kode_user
    BEFORE INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_user();
