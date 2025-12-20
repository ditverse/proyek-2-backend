# Test Export Laporan Peminjaman to Excel

## Prerequisites
1. Server harus sudah running (`air` atau `go run`)
2. Sudah punya JWT token dengan role **SARPRAS** atau **ADMIN**

---

## Test Cases

### 1. Export All Data (Default - 1 bulan terakhir)

**Windows PowerShell:**
```powershell
$token = "YOUR_JWT_TOKEN_HERE"
Invoke-WebRequest -Uri "http://localhost:8080/api/laporan/peminjaman/export" `
    -Headers @{"Authorization" = "Bearer $token"} `
    -OutFile "Laporan_Test_Default.xlsx"
```

**curl (Git Bash / WSL):**
```bash
curl -X GET "http://localhost:8080/api/laporan/peminjaman/export" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  --output Laporan_Test_Default.xlsx
```

---

### 2. Export with Date Range (Januari 2025)

**PowerShell:**
```powershell
$token = "YOUR_JWT_TOKEN_HERE"
$startDate = "2025-01-01T00:00:00Z"
$endDate = "2025-01-31T23:59:59Z"

Invoke-WebRequest -Uri "http://localhost:8080/api/laporan/peminjaman/export?start=$startDate&end=$endDate" `
    -Headers @{"Authorization" = "Bearer $token"} `
    -OutFile "Laporan_Januari_2025.xlsx"
```

**curl:**
```bash
curl -X GET "http://localhost:8080/api/laporan/peminjaman/export?start=2025-01-01T00:00:00Z&end=2025-01-31T23:59:59Z" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  --output Laporan_Januari_2025.xlsx
```

---

### 3. Export Only APPROVED Status

**PowerShell:**
```powershell
$token = "YOUR_JWT_TOKEN_HERE"

Invoke-WebRequest -Uri "http://localhost:8080/api/laporan/peminjaman/export?status=APPROVED" `
    -Headers @{"Authorization" = "Bearer $token"} `
    -OutFile "Laporan_Approved.xlsx"
```

**curl:**
```bash
curl -X GET "http://localhost:8080/api/laporan/peminjaman/export?status=APPROVED" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  --output Laporan_Approved.xlsx
```

---

### 4. Export REJECTED in December 2025

**PowerShell:**
```powershell
$token = "YOUR_JWT_TOKEN_HERE"
$startDate = "2025-12-01T00:00:00Z"
$endDate = "2025-12-31T23:59:59Z"

Invoke-WebRequest -Uri "http://localhost:8080/api/laporan/peminjaman/export?start=$startDate&end=$endDate&status=REJECTED" `
    -Headers @{"Authorization" = "Bearer $token"} `
    -OutFile "Laporan_Rejected_Dec2025.xlsx"
```

**curl:**
```bash
curl -X GET "http://localhost:8080/api/laporan/peminjaman/export?start=2025-12-01T00:00:00Z&end=2025-12-31T23:59:59Z&status=REJECTED" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  --output Laporan_Rejected_Dec2025.xlsx
```

---

## Expected Results

### Success (200 OK)
- File `.xlsx` ter-download
- File bisa dibuka di Microsoft Excel / LibreOffice Calc / Google Sheets
- Header berwarna orange
- Status badge berwarna (green/red/yellow)
- Data lengkap dengan relasi

### Errors

#### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```
**Solution**: Login ulang untuk mendapatkan token baru

#### 403 Forbidden
```json
{
  "error": "Forbidden"
}
```
**Solution**: Pastikan user adalah SARPRAS atau ADMIN

#### 400 Bad Request
```json
{
  "error": "Invalid start date"
}
```
**Solution**: Pastikan format tanggal menggunakan RFC3339

---

## Verification Checklist

After download, verify:

- [ ] File dapat dibuka tanpa error
- [ ] Header section complete (title, filter info, generated timestamp)
- [ ] Column headers berwarna orange dengan text putih
- [ ] Data rows ada border
- [ ] Status badge berwarna sesuai:
  - APPROVED = hijau
  - REJECTED = merah
  - PENDING = kuning
- [ ] Freeze panes bekerja (header tetap saat scroll)
- [ ] Data relasi complete:
  - [ ] Nama Peminjam
  - [ ] Organisasi
  - [ ] Nama Kegiatan
  - [ ] Ruangan (atau "-")
  - [ ] Barang dengan jumlah (atau "-")
  - [ ] Verifikator (atau "-")
- [ ] Date format: DD/MM/YYYY HH:MM
- [ ] Multi-line text (barang & catatan) wrap dengan benar

---

## Quick Login to Get Token

**Request:**
```bash
curl -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "sarpras@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "kode_user": "USR-001",
    "nama": "Petugas Sarpras",
    "email": "sarpras@example.com",
    "role": "SARPRAS"
  }
}
```

Copy the `token` value and use it in the export requests above.

---

## Notes

1. **File size**: Tergantung jumlah data. ~50KB untuk 100 records.
2. **Performance**: Export 1000 records ≈ 2-3 detik
3. **Browser download**: Otomatis trigger download di browser
4. **Filename**: Auto-generated dengan timestamp current date

---

*Happy Testing!* 🚀
