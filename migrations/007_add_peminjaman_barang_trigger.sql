-- =====================================================
-- MIGRATION: Add Trigger for kode_peminjaman_barang
-- Tanggal: 2026-01-29
-- Deskripsi: Menambahkan trigger auto-generate kode_peminjaman_barang
--            dengan format PMB-YYYY-MM-DD-XXXX agar selaras dengan trigger lain
-- =====================================================

-- ==========================================
-- PEMINJAMAN BARANG (kode_peminjaman_barang)
-- Format: PMB-YYYY-MM-DD-XXXX
-- ==========================================
CREATE OR REPLACE FUNCTION generate_kode_peminjaman_barang()
RETURNS TRIGGER AS $$
DECLARE
    today_str TEXT;
    next_seq INT;
    new_kode TEXT;
BEGIN
    IF NEW.kode_peminjaman_barang IS NOT NULL AND NEW.kode_peminjaman_barang != '' THEN
        RETURN NEW;
    END IF;

    -- Format: YYYY-MM-DD
    today_str := TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD');

    -- Find max sequence for today
    SELECT COALESCE(MAX(
        CAST(SUBSTRING(kode_peminjaman_barang FROM 'PMB-[0-9]{4}-[0-9]{2}-[0-9]{2}-([0-9]{4})') AS INTEGER)
    ), 0) + 1
    INTO next_seq
    FROM peminjaman_barang
    WHERE kode_peminjaman_barang LIKE 'PMB-' || today_str || '-%';

    new_kode := 'PMB-' || today_str || '-' || LPAD(next_seq::TEXT, 4, '0');
    NEW.kode_peminjaman_barang := new_kode;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_kode_peminjaman_barang ON peminjaman_barang;
CREATE TRIGGER trigger_generate_kode_peminjaman_barang
    BEFORE INSERT ON peminjaman_barang
    FOR EACH ROW
    EXECUTE FUNCTION generate_kode_peminjaman_barang();

-- =====================================================
-- CONTOH PENGGUNAAN:
-- 
-- INSERT INTO peminjaman_barang (kode_peminjaman, kode_barang, jumlah)
-- VALUES ('PMJ-2026-01-29-0001', 'BRG-0001', 5);
-- 
-- Hasil: kode_peminjaman_barang = 'PMB-2026-01-29-0001'
-- =====================================================
