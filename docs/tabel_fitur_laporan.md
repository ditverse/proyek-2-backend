# Tabel Fitur Sistem - Siap Copy ke Word

> **Instruksi:** Copy tabel-tabel di bawah ini ke Microsoft Word. Format akan otomatis terbawa.

---

## Tabel Ringkasan Fitur Sistem

| Kategori | Jumlah Fitur |
|----------|:------------:|
| Autentikasi | 5 |
| Fitur Inti Peminjaman | 21 |
| Fitur Administrasi | 8 |
| Fitur Pendukung | 6 |
| Komponen Teknis | 15 |
| **Total** | **56** |

---

## Tabel 1. Fitur Autentikasi

| No | Nama Fitur | Deskripsi |
|----|------------|-----------|
| 1 | Login Pengguna | Autentikasi pengguna menggunakan email dan password untuk mengakses sistem |
| 2 | Registrasi Pengguna | Pendaftaran akun baru dengan validasi email unik dan enkripsi password |
| 3 | Manajemen Token JWT | Pembuatan dan validasi token akses dengan masa berlaku 24 jam |
| 4 | Enkripsi Password | Pengamanan password menggunakan algoritma bcrypt |
| 5 | Redirect Berdasarkan Role | Pengarahan otomatis ke dashboard sesuai role setelah login berhasil |

---

## Tabel 2. Fitur Inti Peminjaman

| No | Nama Fitur | Deskripsi |
|----|------------|-----------|
| 1 | Form Pengajuan Peminjaman | Formulir untuk mengajukan peminjaman ruangan dan/atau barang |
| 2 | Unggah Surat Digital | Mengunggah dokumen surat peminjaman dalam format PDF |
| 3 | Unduh Surat Digital | Mengunduh atau melihat surat digital yang sudah diunggah |
| 4 | Riwayat Peminjaman | Melihat daftar seluruh peminjaman milik pengguna |
| 5 | Detail Peminjaman | Melihat informasi lengkap satu pengajuan peminjaman |
| 6 | Peminjaman Barang | Menambahkan daftar barang yang ingin dipinjam dalam satu pengajuan |
| 7 | Pembuatan Kegiatan Otomatis | Sistem otomatis membuat record kegiatan baru saat pengajuan |
| 8 | Validasi Konflik Jadwal | Pengecekan ketersediaan ruangan pada rentang waktu yang dipilih |
| 9 | Dashboard Mahasiswa | Halaman utama mahasiswa dengan statistik dan menu peminjaman |
| 10 | Dashboard Admin Sarpras | Halaman admin dengan statistik pengajuan dan menu pengelolaan |
| 11 | Dashboard Security | Halaman petugas keamanan dengan jadwal dan verifikasi kehadiran |
| 12 | Dashboard Tamu | Halaman publik untuk melihat jadwal ruangan tanpa login |
| 13 | Jadwal Ruangan | Tampilan kalender dengan jadwal peminjaman semua ruangan |
| 14 | Jadwal Aktif Harian | Daftar jadwal hari ini dengan status Disetujui atau Berlangsung |
| 15 | Widget Kalender Interaktif | Komponen kalender dengan indikator visual untuk hari yang ada jadwal |
| 16 | Cek Tanggal Terisi | Mendapatkan tanggal-tanggal yang sudah dibooking untuk ruangan tertentu |
| 17 | Verifikasi Kehadiran | Petugas keamanan mencatat kehadiran peminjam |
| 18 | Riwayat Kehadiran | Melihat daftar riwayat kehadiran yang sudah dicatat |
| 19 | Jadwal Belum Diverifikasi | Daftar jadwal hari ini yang belum diverifikasi kehadirannya |

---

## Tabel 3. Fitur Administrasi

| No | Nama Fitur | Deskripsi |
|----|------------|-----------|
| 1 | Verifikasi Peminjaman | Menyetujui atau menolak pengajuan peminjaman dengan catatan |
| 2 | Daftar Pengajuan Pending | Melihat seluruh pengajuan dengan status menunggu verifikasi |
| 3 | Pembatalan Peminjaman | Membatalkan peminjaman yang sudah disetujui dengan alasan |
| 4 | Laporan Peminjaman | Rekap peminjaman dengan filter berdasarkan tanggal dan status |
| 5 | Ekspor ke Excel | Mengunduh laporan peminjaman dalam format file Excel |
| 6 | Pengelolaan Ruangan | Menambah, mengubah, menghapus, dan melihat data ruangan |
| 7 | Pengelolaan Barang | Menambah, mengubah, menghapus, dan melihat data inventaris barang |
| 8 | Log Aktivitas | Melihat catatan seluruh aktivitas penting dalam sistem |

---

## Tabel 4. Fitur Pendukung

| No | Nama Fitur | Deskripsi |
|----|------------|-----------|
| 1 | Notifikasi Email Persetujuan | Email otomatis ke mahasiswa saat peminjaman disetujui |
| 2 | Notifikasi Email Penolakan | Email otomatis ke mahasiswa saat peminjaman ditolak |
| 3 | Notifikasi Email ke Security | Email otomatis ke petugas keamanan tentang jadwal kegiatan |
| 4 | Notifikasi Email Pembatalan | Email otomatis ke mahasiswa saat peminjaman dibatalkan |
| 5 | Pengiriman Email Asinkron | Pengiriman email secara non-blocking agar tidak menghambat respons sistem |
| 6 | Informasi Instansi | Endpoint untuk menampilkan informasi umum instansi |
| 7 | Daftar Organisasi | Menampilkan daftar organisasi untuk dropdown registrasi |
| 8 | Pendaftaran Organisasi | Halaman untuk mendaftarkan organisasi baru |
| 9 | Health Check | Endpoint untuk memantau status ketersediaan server |

---

## Tabel 5. Komponen Teknis (Non-Fungsional)

| No | Nama Komponen | Deskripsi |
|----|---------------|-----------|
| 1 | Middleware Autentikasi JWT | Validasi token JWT pada setiap request ke endpoint terlindungi |
| 2 | Otorisasi Berbasis Role | Pembatasan akses berdasarkan role pengguna (Mahasiswa, Sarpras, Security, Admin) |
| 3 | Middleware CORS | Pengaturan Cross-Origin Resource Sharing untuk akses lintas domain |
| 4 | Penjadwal Status Otomatis | Background job untuk auto-update status APPROVED→ONGOING→FINISHED |
| 5 | Integrasi Supabase Storage | Layanan penyimpanan file PDF surat digital di cloud |
| 6 | Integrasi Gmail API | Layanan pengiriman email menggunakan Google Gmail API |
| 7 | Migrasi Database | Script SQL untuk inisialisasi struktur tabel database |
| 8 | Auto-Generate Kode | Trigger database untuk membuat kode unik otomatis (USR-xxx, PMJ-xxx) |
| 9 | Indexing Database | Optimasi performa query dengan index pada kolom penting |
| 10 | Data Awal (Seed) | Data awal untuk keperluan pengujian sistem |
| 11 | Pencatatan Log Aktivitas | Audit trail untuk mencatat seluruh aksi penting dalam sistem |
| 12 | Modul API Wrapper | Modul JavaScript untuk komunikasi dengan backend API |
| 13 | Manajemen Konfigurasi | Pengelolaan variabel environment (JWT, Database, Supabase, Gmail) |
| 14 | Pemindahan File Storage | Fungsi untuk memindahkan file di storage saat status berubah |
| 15 | Hot Reload Development | Konfigurasi untuk auto-reload saat pengembangan |

---

## Catatan untuk Laporan

> **Fungsional vs Non-Fungsional:**
> - Tabel 1-4 berisi **kebutuhan fungsional** (fitur yang langsung digunakan pengguna)
> - Tabel 5 berisi **kebutuhan non-fungsional** (komponen teknis pendukung sistem)

> **Referensi Silang:**
> - Detail implementasi setiap fitur dapat dilihat pada Bab Implementasi
> - Pengujian setiap fitur dibahas pada Bab Pengujian
