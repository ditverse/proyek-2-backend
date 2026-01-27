-- =====================================================
-- MIGRATION: Redesign Tabel Mailbox untuk Gmail API Notification
-- Tanggal: 2026-01-25
-- Deskripsi: Menyederhanakan tabel mailbox dari 11 atribut menjadi 5 atribut (normalized)
-- =====================================================

-- 1. Backup data lama (jika ada)
CREATE TABLE IF NOT EXISTS mailbox_backup AS SELECT * FROM mailbox;

-- 2. Drop tabel lama beserta constraints
DROP TABLE IF EXISTS mailbox CASCADE;

-- 3. Buat tabel mailbox baru (normalized - 5 atribut)
CREATE TABLE public.mailbox (
    kode_mailbox character varying NOT NULL,
    kode_user character varying NOT NULL,
    kode_peminjaman character varying NOT NULL,
    jenis_pesan character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    
    CONSTRAINT mailbox_pkey PRIMARY KEY (kode_mailbox),
    CONSTRAINT mailbox_kode_user_fkey FOREIGN KEY (kode_user) 
        REFERENCES public.users(kode_user) ON DELETE CASCADE,
    CONSTRAINT mailbox_kode_peminjaman_fkey FOREIGN KEY (kode_peminjaman) 
        REFERENCES public.peminjaman(kode_peminjaman) ON DELETE CASCADE
);

-- 4. Buat sequence untuk urutan harian
CREATE SEQUENCE IF NOT EXISTS mailbox_daily_seq START 1;

-- 5. Function untuk auto-generate kode_mailbox
-- Format: MBX-YYYY-MM-DD-XXXX (contoh: MBX-2026-01-25-0001)
CREATE OR REPLACE FUNCTION generate_kode_mailbox()
RETURNS TRIGGER AS $$
DECLARE
    today_date TEXT;
    seq_num INTEGER;
BEGIN
    -- Format tanggal: YYYY-MM-DD
    today_date := TO_CHAR(NOW(), 'YYYY-MM-DD');
    
    -- Reset sequence jika hari baru (cek dari record terakhir hari ini)
    IF NOT EXISTS (
        SELECT 1 FROM mailbox 
        WHERE kode_mailbox LIKE 'MBX-' || today_date || '-%'
    ) THEN
        -- Reset sequence ke 1 untuk hari baru
        PERFORM setval('mailbox_daily_seq', 1, FALSE);
    END IF;
    
    -- Ambil nomor urut berikutnya
    seq_num := nextval('mailbox_daily_seq');
    
    -- Generate kode: MBX-YYYY-MM-DD-XXXX
    NEW.kode_mailbox := 'MBX-' || today_date || '-' || LPAD(seq_num::TEXT, 4, '0');
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 6. Trigger untuk auto-generate kode_mailbox saat INSERT
CREATE TRIGGER trigger_generate_kode_mailbox
BEFORE INSERT ON mailbox
FOR EACH ROW
WHEN (NEW.kode_mailbox IS NULL OR NEW.kode_mailbox = '')
EXECUTE FUNCTION generate_kode_mailbox();

-- 7. Index untuk optimasi query
CREATE INDEX idx_mailbox_kode_user ON mailbox(kode_user);
CREATE INDEX idx_mailbox_kode_peminjaman ON mailbox(kode_peminjaman);
CREATE INDEX idx_mailbox_jenis_pesan ON mailbox(jenis_pesan);
CREATE INDEX idx_mailbox_created_at ON mailbox(created_at DESC);

-- 8. Comment untuk dokumentasi
COMMENT ON TABLE mailbox IS 'Tabel untuk menyimpan log notifikasi email Gmail API';
COMMENT ON COLUMN mailbox.kode_mailbox IS 'Primary key, auto-generate format MBX-YYYY-MM-DD-XXXX';
COMMENT ON COLUMN mailbox.kode_user IS 'FK ke tabel users - penerima email';
COMMENT ON COLUMN mailbox.kode_peminjaman IS 'FK ke tabel peminjaman - data terkait';
COMMENT ON COLUMN mailbox.jenis_pesan IS 'Tipe email: APPROVED, REJECTED, CANCELLED, SECURITY_NOTIFY';
COMMENT ON COLUMN mailbox.created_at IS 'Timestamp pembuatan record';

-- =====================================================
-- CONTOH PENGGUNAAN:
-- 
-- INSERT INTO mailbox (kode_user, kode_peminjaman, jenis_pesan)
-- VALUES ('USR-MHS-001', 'PMJ-001', 'APPROVED');
-- 
-- Hasil: kode_mailbox = 'MBX-2026-01-25-0001'
-- =====================================================
