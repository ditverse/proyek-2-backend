# LAPORAN PENGUJIAN BLACK BOX
## Sistem Informasi Peminjaman Sarana dan Prasarana Kampus

**Tanggal Pengujian**: 11 Desember 2025  
**Penguji**: Tim Development  
**Versi Aplikasi**: 1.0

---

## 1. Pendahuluan

### 1.1 Tujuan Pengujian
Pengujian Black Box dilakukan untuk memvalidasi fungsionalitas sistem tanpa melihat struktur internal kode. Fokus pengujian pada input-output dan perilaku sistem sesuai spesifikasi.

### 1.2 Ruang Lingkup
- Modul Authentication (Login & Register)
- Modul Master Data (Ruangan & Barang)
- Modul Peminjaman
- Modul Kehadiran
- Modul Notifikasi

### 1.3 Metode Pengujian
- **Equivalence Partitioning**: Membagi input ke dalam kelas yang valid dan tidak valid
- **Boundary Value Analysis**: Menguji nilai batas input
- **Decision Table Testing**: Menguji kombinasi kondisi dan aksi

---

## 2. Pengujian Modul Authentication

### 2.1 Pengujian Login

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-AUTH-01 | Login dengan email dan password valid | Email: user@test.com, Password: password123 | Login berhasil, mendapat JWT token | Login berhasil, token diterima | ✅ PASS |
| TC-AUTH-02 | Login dengan email tidak terdaftar | Email: notexist@test.com, Password: password123 | Error: "invalid credentials" (401) | Error: invalid credentials | ✅ PASS |
| TC-AUTH-03 | Login dengan password salah | Email: user@test.com, Password: wrongpass | Error: "invalid credentials" (401) | Error: invalid credentials | ✅ PASS |
| TC-AUTH-04 | Login dengan email kosong | Email: "", Password: password123 | Error: "Invalid request body" (400) | Error: Invalid request body | ✅ PASS |
| TC-AUTH-05 | Login dengan password kosong | Email: user@test.com, Password: "" | Error: "invalid credentials" (401) | Error: invalid credentials | ✅ PASS |
| TC-AUTH-06 | Login dengan format email tidak valid | Email: invalidemail, Password: password123 | Error: "invalid credentials" (401) | Error: invalid credentials | ✅ PASS |
| TC-AUTH-07 | Login dengan method selain POST | GET /api/auth/login | Error: "Method not allowed" (405) | Method not allowed | ✅ PASS |

### 2.2 Pengujian Register

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-AUTH-08 | Register dengan data lengkap dan valid | Nama: John Doe, Email: john@test.com, Password: pass123, Role: MAHASISWA | User terdaftar dengan kode_user format USR-YYMMDD-XXXX | User berhasil terdaftar | ✅ PASS |
| TC-AUTH-09 | Register dengan email sudah terdaftar | Email yang sudah ada di database | Error: "email sudah terdaftar" (400) | Error: email sudah terdaftar | ✅ PASS |
| TC-AUTH-10 | Register tanpa nama | Nama: "", Email: new@test.com | Error: "nama, email, dan password wajib diisi" (400) | Error: field wajib diisi | ✅ PASS |
| TC-AUTH-11 | Register tanpa email | Email: "" | Error: "nama, email, dan password wajib diisi" (400) | Error: field wajib diisi | ✅ PASS |
| TC-AUTH-12 | Register tanpa password | Password: "" | Error: "nama, email, dan password wajib diisi" (400) | Error: field wajib diisi | ✅ PASS |
| TC-AUTH-13 | Register tanpa role (default MAHASISWA) | Role: tidak disertakan | User terdaftar dengan role MAHASISWA | Role default MAHASISWA | ✅ PASS |
| TC-AUTH-14 | Register dengan role SARPRAS | Role: SARPRAS | User terdaftar dengan role SARPRAS | Role SARPRAS berhasil | ✅ PASS |
| TC-AUTH-15 | Register dengan role SECURITY | Role: SECURITY | User terdaftar dengan role SECURITY | Role SECURITY berhasil | ✅ PASS |
| TC-AUTH-16 | Register dengan role ADMIN | Role: ADMIN | User terdaftar dengan role ADMIN | Role ADMIN berhasil | ✅ PASS |

---

## 3. Pengujian Modul Master Data Ruangan

### 3.1 Get All Ruangan

| No | Kasus Uji | Kondisi Awal | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|--------------|----------------------|--------------|--------|
| TC-RNG-01 | Get semua ruangan dengan token valid | User terautentikasi | List ruangan ditampilkan (200) | List ruangan berhasil | ✅ PASS |
| TC-RNG-02 | Get semua ruangan tanpa token | Tidak ada header Authorization | Error: "Unauthorized" (401) | Unauthorized | ✅ PASS |
| TC-RNG-03 | Get semua ruangan dengan token expired | Token sudah kadaluarsa | Error: "Unauthorized" (401) | Unauthorized | ✅ PASS |

### 3.2 Get Ruangan by ID

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-RNG-04 | Get ruangan dengan ID valid | kode_ruangan: RNG-0001 | Detail ruangan ditampilkan (200) | Detail berhasil | ✅ PASS |
| TC-RNG-05 | Get ruangan dengan ID tidak ada | kode_ruangan: RNG-9999 | Error: "Not found" (404) | Not found | ✅ PASS |
| TC-RNG-06 | Get ruangan dengan ID format salah | kode_ruangan: invalid | Error: "Invalid ID" (400) | Invalid ID | ✅ PASS |

### 3.3 Create Ruangan (SARPRAS/ADMIN)

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-RNG-07 | Create ruangan dengan data lengkap | nama_ruangan: Lab Komputer, lokasi: Gedung A, kapasitas: 40 | Ruangan berhasil dibuat (201) | Ruangan berhasil | ✅ PASS |
| TC-RNG-08 | Create ruangan oleh role MAHASISWA | Role: MAHASISWA | Error: "Forbidden" (403) | Forbidden | ✅ PASS |
| TC-RNG-09 | Create ruangan dengan data tidak lengkap | nama_ruangan: "" | Error: "Invalid request body" (400) | Error validasi | ✅ PASS |
| TC-RNG-10 | Create ruangan dengan kapasitas negatif | kapasitas: -10 | Error: validasi kapasitas (400) | Error validasi | ✅ PASS |

### 3.4 Update Ruangan (SARPRAS/ADMIN)

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-RNG-11 | Update ruangan dengan data valid | kode_ruangan: RNG-0001, nama baru | Ruangan berhasil diupdate (200) | Update berhasil | ✅ PASS |
| TC-RNG-12 | Update ruangan tidak ada | kode_ruangan: RNG-9999 | Error: "Not found" (404) | Not found | ✅ PASS |
| TC-RNG-13 | Update ruangan oleh role tidak berwenang | Role: MAHASISWA | Error: "Forbidden" (403) | Forbidden | ✅ PASS |

### 3.5 Delete Ruangan (SARPRAS/ADMIN)

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-RNG-14 | Delete ruangan valid | kode_ruangan: RNG-0001 | Ruangan berhasil dihapus (204) | Delete berhasil | ✅ PASS |
| TC-RNG-15 | Delete ruangan tidak ada | kode_ruangan: RNG-9999 | Error (404/500) | Error | ✅ PASS |
| TC-RNG-16 | Delete ruangan oleh role tidak berwenang | Role: SECURITY | Error: "Forbidden" (403) | Forbidden | ✅ PASS |

---

## 4. Pengujian Modul Master Data Barang

### 4.1 CRUD Barang

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-BRG-01 | Get semua barang | Token valid | List barang ditampilkan (200) | List berhasil | ✅ PASS |
| TC-BRG-02 | Get barang by ID valid | kode_barang: BRG-0001 | Detail barang (200) | Detail berhasil | ✅ PASS |
| TC-BRG-03 | Get barang by ID tidak ada | kode_barang: BRG-9999 | Error: "Not found" (404) | Not found | ✅ PASS |
| TC-BRG-04 | Create barang dengan data lengkap | nama_barang, deskripsi, jumlah_total | Barang berhasil dibuat (201) | Barang berhasil | ✅ PASS |
| TC-BRG-05 | Create barang oleh MAHASISWA | Role: MAHASISWA | Error: "Forbidden" (403) | Forbidden | ✅ PASS |
| TC-BRG-06 | Update barang valid | Data update lengkap | Barang berhasil diupdate (200) | Update berhasil | ✅ PASS |
| TC-BRG-07 | Delete barang valid | kode_barang: BRG-0001 | Barang dihapus (204) | Delete berhasil | ✅ PASS |

---

## 5. Pengujian Modul Peminjaman

### 5.1 Create Peminjaman

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-PMJ-01 | Create peminjaman dengan data lengkap | kode_ruangan, tanggal_mulai, tanggal_selesai, keperluan | Peminjaman berhasil dengan status PENDING (201) | Berhasil status PENDING | ✅ PASS |
| TC-PMJ-02 | Create peminjaman tanpa token | Tidak ada Authorization header | Error: "Unauthorized" (401) | Unauthorized | ✅ PASS |
| TC-PMJ-03 | Create peminjaman tanggal_selesai < tanggal_mulai | tanggal_selesai sebelum tanggal_mulai | Error: "tanggal_selesai harus setelah tanggal_mulai" (400) | Error validasi | ✅ PASS |
| TC-PMJ-04 | Create peminjaman format tanggal invalid | tanggal_mulai: "invalid-date" | Error: "format tanggal_mulai tidak valid" (400) | Error format | ✅ PASS |
| TC-PMJ-05 | Create peminjaman dengan barang tidak ada | kode_barang: BRG-9999 | Error: "barang BRG-9999 tidak ditemukan" (400) | Error barang | ✅ PASS |
| TC-PMJ-06 | Create peminjaman dengan multiple barang | Array barang dengan 2-3 item | Peminjaman + relasi barang berhasil (201) | Berhasil | ✅ PASS |

### 5.2 Get Peminjaman

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-PMJ-07 | Get my peminjaman | User terautentikasi | List peminjaman milik user (200) | List berhasil | ✅ PASS |
| TC-PMJ-08 | Get peminjaman by ID valid | kode_peminjaman: PMJ-XXXX | Detail peminjaman + relasi (200) | Detail berhasil | ✅ PASS |
| TC-PMJ-09 | Get peminjaman by ID tidak ada | kode_peminjaman: PMJ-9999 | Error: "Peminjaman not found" (404) | Not found | ✅ PASS |
| TC-PMJ-10 | Get pending peminjaman (SARPRAS) | Role: SARPRAS | List peminjaman PENDING (200) | List PENDING | ✅ PASS |
| TC-PMJ-11 | Get pending peminjaman (MAHASISWA) | Role: MAHASISWA | Error: "Forbidden" (403) | Forbidden | ✅ PASS |

### 5.3 Verifikasi Peminjaman (SARPRAS/ADMIN)

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-PMJ-12 | Verifikasi APPROVED | status: APPROVED | Status berubah APPROVED, notifikasi terkirim | Berhasil | ✅ PASS |
| TC-PMJ-13 | Verifikasi REJECTED dengan catatan | status: REJECTED, catatan_verifikasi | Status berubah REJECTED + catatan | Berhasil | ✅ PASS |
| TC-PMJ-14 | Verifikasi peminjaman sudah diverifikasi | Status bukan PENDING | Error: "peminjaman sudah diverifikasi" (400) | Error | ✅ PASS |
| TC-PMJ-15 | Verifikasi dengan status invalid | status: INVALID_STATUS | Error: "status verifikasi tidak valid" (400) | Error status | ✅ PASS |
| TC-PMJ-16 | Verifikasi oleh MAHASISWA | Role: MAHASISWA | Error: "Forbidden" (403) | Forbidden | ✅ PASS |
| TC-PMJ-17 | Verifikasi peminjaman tidak ada | kode_peminjaman: PMJ-9999 | Error: "peminjaman tidak ditemukan" (400) | Not found | ✅ PASS |

### 5.4 Upload Surat Digital

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-PMJ-18 | Upload file PDF valid | File PDF < 2MB | Upload berhasil, path tersimpan | Berhasil | ✅ PASS |
| TC-PMJ-19 | Upload file non-PDF | File gambar (.jpg) | Error: "File harus berupa PDF" (400) | Error format | ✅ PASS |
| TC-PMJ-20 | Upload file > 2MB | File PDF > 2MB | Error: "Ukuran file maksimal 2MB" (400) | Error size | ✅ PASS |
| TC-PMJ-21 | Upload tanpa file | Form tanpa file | Error: "File surat wajib diupload" (400) | Error required | ✅ PASS |

### 5.5 Jadwal Ruangan

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-PMJ-22 | Get jadwal ruangan tanpa filter | Tanpa query params | Jadwal 1 bulan ke depan (200) | Jadwal berhasil | ✅ PASS |
| TC-PMJ-23 | Get jadwal ruangan dengan filter tanggal | start: 2025-12-01, end: 2025-12-31 | Jadwal sesuai filter (200) | Jadwal filtered | ✅ PASS |
| TC-PMJ-24 | Get jadwal dengan format tanggal invalid | start: invalid-date | Error: "Invalid start date" (400) | Error format | ✅ PASS |

### 5.6 Jadwal Aktif (SECURITY)

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-PMJ-25 | Get jadwal aktif (SECURITY) | Role: SECURITY | List jadwal APPROVED (200) | List berhasil | ✅ PASS |
| TC-PMJ-26 | Get jadwal aktif belum verifikasi | Role: SECURITY | Jadwal belum ada kehadiran (200) | List berhasil | ✅ PASS |
| TC-PMJ-27 | Get jadwal aktif oleh MAHASISWA | Role: MAHASISWA | Error: "Forbidden" (403) | Forbidden | ✅ PASS |

---

## 6. Pengujian Modul Kehadiran

### 6.1 Create Kehadiran (SECURITY/ADMIN)

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-KHD-01 | Catat kehadiran HADIR | kode_peminjaman, status: HADIR | Kehadiran dicatat, status jadi ONGOING | Berhasil | ✅ PASS |
| TC-KHD-02 | Catat kehadiran TIDAK_HADIR | Status: TIDAK_HADIR | Kehadiran dicatat sebagai tidak hadir | Berhasil | ✅ PASS |
| TC-KHD-03 | Catat kehadiran BATAL | Status: BATAL | Kehadiran dicatat sebagai batal | Berhasil | ✅ PASS |
| TC-KHD-04 | Catat kehadiran oleh MAHASISWA | Role: MAHASISWA | Error: "Forbidden" (403) | Forbidden | ✅ PASS |
| TC-KHD-05 | Catat kehadiran tanpa token | Tidak ada Authorization | Error: "Unauthorized" (401) | Unauthorized | ✅ PASS |

### 6.2 Get Kehadiran

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-KHD-06 | Get riwayat kehadiran by security | Role: SECURITY | List kehadiran yang dicatat (200) | List berhasil | ✅ PASS |
| TC-KHD-07 | Get laporan kehadiran | Role: SARPRAS/ADMIN | Laporan kehadiran (200) | Laporan berhasil | ✅ PASS |

---

## 7. Pengujian Modul Notifikasi

### 7.1 Get Notifikasi

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-NTF-01 | Get notifikasi user | User terautentikasi | List notifikasi milik user (200) | List berhasil | ✅ PASS |
| TC-NTF-02 | Get count notifikasi unread | User terautentikasi | Jumlah notifikasi TERKIRIM (200) | Count berhasil | ✅ PASS |
| TC-NTF-03 | Get notifikasi tanpa token | Tidak ada Authorization | Error: "Unauthorized" (401) | Unauthorized | ✅ PASS |

### 7.2 Update Status Notifikasi

| No | Kasus Uji | Data Input | Hasil yang Diharapkan | Hasil Aktual | Status |
|----|-----------|------------|----------------------|--------------|--------|
| TC-NTF-04 | Tandai notifikasi dibaca | kode_notifikasi valid | Status DIBACA (200) | Status berubah | ✅ PASS |
| TC-NTF-05 | Tandai notifikasi milik user lain | kode_notifikasi user lain | Error: "Forbidden" (403) | Forbidden | ✅ PASS |

---

## 8. Pengujian Otorisasi Role-Based Access Control (RBAC)

| No | Endpoint | MAHASISWA | SARPRAS | SECURITY | ADMIN | Status |
|----|----------|-----------|---------|----------|-------|--------|
| TC-RBAC-01 | POST /peminjaman | ✅ | ✅ | ✅ | ✅ | ✅ PASS |
| TC-RBAC-02 | GET /peminjaman/pending | ❌ 403 | ✅ | ❌ 403 | ✅ | ✅ PASS |
| TC-RBAC-03 | POST /peminjaman/{id}/verifikasi | ❌ 403 | ✅ | ❌ 403 | ✅ | ✅ PASS |
| TC-RBAC-04 | POST /ruangan/create | ❌ 403 | ✅ | ❌ 403 | ✅ | ✅ PASS |
| TC-RBAC-05 | POST /barang/create | ❌ 403 | ✅ | ❌ 403 | ✅ | ✅ PASS |
| TC-RBAC-06 | POST /kehadiran | ❌ 403 | ❌ 403 | ✅ | ✅ | ✅ PASS |
| TC-RBAC-07 | GET /jadwal-aktif | ❌ 403 | ❌ 403 | ✅ | ✅ | ✅ PASS |
| TC-RBAC-08 | GET /log-aktivitas | ❌ 403 | ❌ 403 | ❌ 403 | ✅ | ✅ PASS |

---

## 9. Ringkasan Hasil Pengujian

### 9.1 Statistik Pengujian

| Modul | Total Test Case | PASS | FAIL | Persentase |
|-------|-----------------|------|------|------------|
| Authentication | 16 | 16 | 0 | 100% |
| Master Data Ruangan | 16 | 16 | 0 | 100% |
| Master Data Barang | 7 | 7 | 0 | 100% |
| Peminjaman | 27 | 27 | 0 | 100% |
| Kehadiran | 7 | 7 | 0 | 100% |
| Notifikasi | 5 | 5 | 0 | 100% |
| RBAC | 8 | 8 | 0 | 100% |
| **TOTAL** | **86** | **86** | **0** | **100%** |

### 9.2 Kesimpulan
Berdasarkan hasil pengujian Black Box, seluruh fungsionalitas sistem berjalan sesuai dengan spesifikasi yang telah ditentukan. Sistem telah memenuhi kriteria:

1. ✅ Authentication berfungsi dengan baik (login, register, JWT validation)
2. ✅ CRUD Master Data (ruangan, barang) berjalan sesuai spesifikasi
3. ✅ Workflow peminjaman lengkap dari pengajuan hingga verifikasi
4. ✅ Upload surat digital dengan validasi file
5. ✅ Pencatatan kehadiran oleh security
6. ✅ Sistem notifikasi otomatis
7. ✅ Role-Based Access Control (RBAC) konsisten untuk semua role

---

**Dibuat oleh**: Tim Development  
**Tanggal**: 11 Desember 2025
