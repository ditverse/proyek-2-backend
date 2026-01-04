# Postman Collection for Time Validation Testing

## Test Collection: Validasi Waktu Operasional

### Base URL
```
http://localhost:8000
```

---

## 1. Valid Requests

### 1.1 Valid Peminjaman - Single Day

**Note**: Run this during weekday (Mon-Fri) 07:00-17:00 WIB

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Rapat Koordinasi Valid",
  "deskripsi": "Test validasi waktu - valid single day",
  "tanggal_mulai": "2026-01-06T09:00:00+07:00",
  "tanggal_selesai": "2026-01-06T16:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 201 Created
```json
{
  "kode_peminjaman": "PMJ-260106-XXXX",
  "kode_user": "USR-...",
  "status": "PENDING",
  ...
}
```

---

### 1.2 Valid Peminjaman - Multi-Day (No Sunday)

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Workshop 3 Hari",
  "deskripsi": "Test multi-day tanpa Minggu",
  "tanggal_mulai": "2026-01-06T08:00:00+07:00",
  "tanggal_selesai": "2026-01-08T17:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 201 Created

---

### 1.3 Valid - Edge Case (Exact 07:00 - 17:00)

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Full Day Event",
  "deskripsi": "Exact boundary test",
  "tanggal_mulai": "2026-01-06T07:00:00+07:00",
  "tanggal_selesai": "2026-01-06T17:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 201 Created

---

### 1.4 Valid - Saturday Rental

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Kegiatan Sabtu",
  "deskripsi": "Test rental di hari Sabtu",
  "tanggal_mulai": "2026-01-10T10:00:00+07:00",
  "tanggal_selesai": "2026-01-10T14:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 201 Created

---

## 2. Invalid Requests - Jam Submit

### 2.1 Invalid - Submit on Weekend

**Precondition**: Run this on Saturday or Sunday

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Test Weekend Submit",
  "tanggal_mulai": "2026-01-06T09:00:00+07:00",
  "tanggal_selesai": "2026-01-06T16:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Pengajuan peminjaman hanya dapat dilakukan pada hari kerja (Senin-Jumat)"
}
```

---

### 2.2 Invalid - Submit Outside Office Hours

**Precondition**: Run this before 07:00 or at/after 17:00

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Test After Hours",
  "tanggal_mulai": "2026-01-07T09:00:00+07:00",
  "tanggal_selesai": "2026-01-07T16:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Pengajuan peminjaman hanya dapat dilakukan pada jam kerja (07:00-17:00 WIB)"
}
```

---

## 3. Invalid Requests - Waktu Peminjaman

### 3.1 Invalid - Rental on Sunday

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Rapat Minggu",
  "deskripsi": "Test rental hari Minggu",
  "tanggal_mulai": "2026-01-11T10:00:00+07:00",
  "tanggal_selesai": "2026-01-11T12:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Peminjaman tidak dapat dilakukan pada hari Minggu"
}
```

---

### 3.2 Invalid - Start Before 07:00

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Early Morning",
  "deskripsi": "Test start terlalu pagi",
  "tanggal_mulai": "2026-01-06T06:30:00+07:00",
  "tanggal_selesai": "2026-01-06T10:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Waktu mulai peminjaman minimal pukul 07:00 WIB"
}
```

---

### 3.3 Invalid - End After 17:00

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Late Event",
  "deskripsi": "Test end terlalu malam",
  "tanggal_mulai": "2026-01-06T15:00:00+07:00",
  "tanggal_selesai": "2026-01-06T18:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Waktu selesai peminjaman maksimal pukul 17:00 WIB"
}
```

---

### 3.4 Invalid - Start At/After 17:00

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Evening Start",
  "deskripsi": "Test start jam 17:00",
  "tanggal_mulai": "2026-01-06T17:00:00+07:00",
  "tanggal_selesai": "2026-01-06T18:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Waktu mulai peminjaman maksimal pukul 16:59 WIB"
}
```

---

### 3.5 Invalid - Cross to Sunday

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Weekend Crossing",
  "deskripsi": "Test Sabtu ke Minggu",
  "tanggal_mulai": "2026-01-10T16:00:00+07:00",
  "tanggal_selesai": "2026-01-11T10:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Peminjaman harus selesai sebelum hari Minggu"
}
```

---

### 3.6 Invalid - Includes Sunday (Multi-day)

```http
POST {{baseUrl}}/api/peminjaman
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "kode_ruangan": "HLP-002",
  "nama_kegiatan": "Long Event with Sunday",
  "deskripsi": "Test Jumat-Senin (melewati Minggu)",
  "tanggal_mulai": "2026-01-09T14:00:00+07:00",
  "tanggal_selesai": "2026-01-12T10:00:00+07:00",
  "path_surat_digital": "temp",
  "barang": []
}
```

**Expected Response**: 400 Bad Request
```json
{
  "error": "Peminjaman tidak dapat melewati hari Minggu"
}
```

---

## 🧪 Testing Checklist

### Jam Submit (Submission Time)
- [ ] ✅ Submit on weekday 10:00 → Success (201)
- [ ] ❌ Submit on Saturday → Error (400)
- [ ] ❌ Submit on Sunday → Error (400)
- [ ] ❌ Submit at 06:30 → Error (400)
- [ ] ❌ Submit at 17:00 → Error (400)
- [ ] ✅ Submit at 16:59 → Success (201)

### Waktu Peminjaman (Rental Period)
- [ ] ✅ Mon 09:00 - Mon 16:00 → Success (201)
- [ ] ✅ Mon 07:00 - Wed 17:00 → Success (201)
- [ ] ✅ Sat 10:00 - Sat 14:00 → Success (201)
- [ ] ❌ Sun 10:00 - Sun 12:00 → Error (400)
- [ ] ❌ Mon 06:00 - Mon 10:00 → Error (400)
- [ ] ❌ Fri 15:00 - Fri 18:00 → Error (400)
- [ ] ❌ Mon 17:00 - Mon 18:00 → Error (400)
- [ ] ❌ Sat 16:00 - Sun 10:00 → Error (400)
- [ ] ❌ Fri 14:00 - Mon 10:00 → Error (400)

---

## 📝 Notes

1. **Timezone**: All timestamps must use `+07:00` (WIB) timezone
2. **Date Format**: Use ISO 8601 format: `YYYY-MM-DDTHH:mm:ss+07:00`
3. **Token**: Get valid JWT token from `/api/auth/login` first
4. **Room Code**: Use existing room code (e.g., `HLP-002`)
5. **Current Date**: Adjust dates to future dates from current time

---

## 📅 Reference Calendar (January 2026)

```
   January 2026
Su Mo Tu We Th Fr Sa
             1  2  3
 4  5  6  7  8  9 10
11 12 13 14 15 16 17
18 19 20 21 22 23 24
25 26 27 28 29 30 31
```

- **Weekdays (Mon-Fri)**: 5, 6, 7, 8, 9, 12, 13, 14, 15, 16...
- **Saturday**: 3, 10, 17, 24, 31
- **Sunday**: 4, 11, 18, 25

Use these dates for testing!
