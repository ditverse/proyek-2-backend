-- Index untuk foreign key yang sering di-JOIN
CREATE INDEX IF NOT EXISTS idx_peminjaman_kode_user ON peminjaman(kode_user);
CREATE INDEX IF NOT EXISTS idx_peminjaman_kode_ruangan ON peminjaman(kode_ruangan);
CREATE INDEX IF NOT EXISTS idx_peminjaman_kode_kegiatan ON peminjaman(kode_kegiatan);

-- Index untuk peminjaman_barang
CREATE INDEX IF NOT EXISTS idx_peminjaman_barang_kode_peminjaman ON peminjaman_barang(kode_peminjaman);
CREATE INDEX IF NOT EXISTS idx_peminjaman_barang_kode_barang ON peminjaman_barang(kode_barang);

-- Index untuk created_at (ordering)
CREATE INDEX IF NOT EXISTS idx_peminjaman_created_at ON peminjaman(created_at DESC);
