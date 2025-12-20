# Fitur Export Laporan Peminjaman ke Excel

Dokumentasi implementasi fitur export laporan peminjaman ke format Excel (.xlsx).

---

## 📋 Overview

Fitur ini memungkinkan user dengan role **SARPRAS** dan **ADMIN** untuk mengekspor data laporan peminjaman ke file Excel dengan format yang rapi dan profesional.

### **Endpoint**
```
GET /api/laporan/peminjaman/export
```

### **Authentication**
- ✅ Requires: JWT Authentication
- ✅ Roles: `SARPRAS`, `ADMIN`

### **Query Parameters**
| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `start` | string (RFC3339) | No | Tanggal mulai filter | `2025-01-01T00:00:00Z` |
| `end` | string (RFC3339) | No | Tanggal akhir filter | `2025-12-31T23:59:59Z` |
| `status` | string (enum) | No | Filter berdasarkan status | `APPROVED`, `REJECTED`, `PENDING` |

**Default Values:**
- `start`: 1 bulan yang lalu dari sekarang
- `end`: Waktu sekarang
- `status`: Semua status

---

## 🔧 Technical Implementation

### **1. Library Used**
```go
github.com/xuri/excelize/v2
```
Library Excel yang powerful dan populer untuk Golang.

### **2. File Structure**
```
services/
├── export_service.go       # Logic untuk generate Excel file

handlers/
├── export_handler.go       # HTTP handler untuk endpoint export

internal/router/
├── router.go              # Route registration
```

### **3. Architecture Flow**

```
[Frontend]
    ↓
[GET /api/laporan/peminjaman/export?start=...&end=...&status=...]
    ↓
[Export Handler]
    ↓
[Fetch Data from Repository] → Enrich with Relations
    ↓
[Export Service] → Generate Excel File
    ↓
[Return .xlsx File Download]
```

---

## 📊 Excel File Features

### **Formatting & Styling**

#### **1. Header Section**
- **Title**: LAPORAN PEMINJAMAN SARANA PRASARANA (merged cells, bold, 14pt)
- **Filter Info**: Periode & status filter yang digunakan
- **Generated At**: Timestamp pembuatan file

#### **2. Column Headers**
```
No | Kode Peminjaman | Nama Peminjam | Organisasi | Nama Kegiatan | 
Ruangan | Barang | Tanggal Mulai | Tanggal Selesai | Status | 
Verifikator | Tanggal Verifikasi | Catatan Verifikasi
```

- **Background**: Orange (`#FF6B35`) - sesuai theme aplikasi
- **Font**: White, Bold, Calibri 11pt
- **Alignment**: Center, Middle
- **Border**: Black border semua sisi

#### **3. Data Cells**
- **Border**: Light gray (`#CCCCCC`)
- **Text wrap**: Enabled
- **Vertical align**: Top

#### **4. Status Badge Styling**

| Status | Background Color | Text Color | Style |
|--------|-----------------|------------|-------|
| **APPROVED** | Light green (`#D4EDDA`) | Dark green (`#006400`) | Bold, Center |
| **REJECTED** | Light red (`#F8D7DA`) | Dark red (`#8B0000`) | Bold, Center |
| **PENDING** | Light yellow (`#FFF3CD`) | Dark yellow (`#856404`) | Bold, Center |

#### **5. Special Features**
- ✅ **Freeze panes**: Header row tetap terlihat saat scroll
- ✅ **Auto-width columns**: Ukuran kolom disesuaikan dengan konten
- ✅ **Multi-line support**: List barang dan catatan panjang otomatis wrap

---

## 💾 Data Fields

### **Kolom Excel**

1. **No**: Nomor urut
2. **Kode Peminjaman**: ID unik peminjaman (e.g., `PMJ-20250120-001`)
3. **Nama Peminjam**: Nama lengkap user peminjam
4. **Organisasi**: Nama organisasi peminjam (ORMAWA/UKM)
5. **Nama Kegiatan**: Judul kegiatan
6. **Ruangan**: Nama ruangan yang dipinjam (atau "-" jika hanya barang)
7. **Barang**: List barang yang dipinjam dengan jumlah
   - Format: `Nama Barang (x2)`
   - Multi-line jika lebih dari 1 barang
8. **Tanggal Mulai**: Format `DD/MM/YYYY HH:MM`
9. **Tanggal Selesai**: Format `DD/MM/YYYY HH:MM`
10. **Status**: APPROVED / REJECTED / PENDING (dengan warna)
11. **Verifikator**: Nama petugas yang verifikasi (atau "-")
12. **Tanggal Verifikasi**: Waktu verifikasi (atau "-")
13. **Catatan Verifikasi**: Catatan dari sarpras (atau "-")

### **Data Relations**
File Excel include semua relasi:
- ✅ Peminjam (User)
- ✅ Organisasi
- ✅ Kegiatan
- ✅ Ruangan (jika ada)
- ✅ List Barang dengan detail (jika ada)
- ✅ Verifier (User yang approve/reject)

---

## 🔄 Use Cases

### **Use Case 1: Export All Data (Default)**
```http
GET /api/laporan/peminjaman/export
Authorization: Bearer <token>
```
**Result**: File `Laporan_Peminjaman_2025-12-20.xlsx` dengan data 1 bulan terakhir, semua status.

### **Use Case 2: Export Specific Date Range**
```http
GET /api/laporan/peminjaman/export?start=2025-01-01T00:00:00Z&end=2025-01-31T23:59:59Z
Authorization: Bearer <token>
```
**Result**: Data bulan Januari 2025 saja.

### **Use Case 3: Export Only Approved**
```http
GET /api/laporan/peminjaman/export?status=APPROVED
Authorization: Bearer <token>
```
**Result**: Hanya peminjaman yang sudah disetujui.

### **Use Case 4: Export Rejected in December**
```http
GET /api/laporan/peminjaman/export?start=2025-12-01T00:00:00Z&end=2025-12-31T23:59:59Z&status=REJECTED
Authorization: Bearer <token>
```
**Result**: Peminjaman yang ditolak di bulan Desember.

---

## 📝 Filename Convention

```
Laporan_Peminjaman_YYYY-MM-DD.xlsx
```

**Example:**
- `Laporan_Peminjaman_2025-12-20.xlsx` (generated on Dec 20, 2025)
- `Laporan_Peminjaman_2025-01-15.xlsx` (generated on Jan 15, 2025)

---

## 🧪 Testing Checklist

- [ ] **Auth Test**: Endpoint hanya accessible oleh SARPRAS/ADMIN
- [ ] **Filter Test**: Query params bekerja sesuai ekspektasi
- [ ] **Data Test**: Semua field terisi dengan benar
- [ ] **Relation Test**: Data nested (organisasi, barang, dll) muncul
- [ ] **Styling Test**: Warna status badge sesuai
- [ ] **Edge Case**: Handle peminjaman tanpa ruangan / tanpa barang
- [ ] **Performance**: Test dengan 100+ records
- [ ] **Download Test**: File ter-download dengan nama yang benar

---

## 🔮 Future Enhancements

### **Potential Improvements:**
1. **Multiple Formats**
   - Export ke PDF
   - Export ke CSV

2. **Advanced Filtering**
   - Filter by organisasi
   - Filter by ruangan
   - Filter by peminjam

3. **Custom Columns**
   - Allow user memilih kolom yang ingin di-export

4. **Charts & Analytics**
   - Add chart sheet dengan statistik peminjaman
   - Summary by status, organisasi, dll

5. **Scheduled Export**
   - Auto-generate monthly report
   - Email report to admin

6. **Template Customization**
   - Allow admin customize Excel template
   - Logo kampus di header

---

## 🚨 Error Handling

### **Common Errors**

| Error | HTTP Status | Description | Solution |
|-------|-------------|-------------|----------|
| Invalid date format | 400 | Query param `start`/`end` format salah | Use RFC3339 format |
| Unauthorized | 401 | Token invalid/expired | Login ulang |
| Forbidden | 403 | User bukan SARPRAS/ADMIN | Check user role |
| No data | 200 | File Excel kosong (hanya header) | Valid, no error |
| Server Error | 500 | Gagal generate Excel | Check logs |

---

## 📚 Code Example

### **Frontend (Vanilla JS)**
```javascript
async function exportToExcel() {
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;
    const status = document.getElementById('statusFilter').value;
    
    const params = new URLSearchParams();
    if (startDate) params.append('start', new Date(startDate).toISOString());
    if (endDate) params.append('end', new Date(endDate).toISOString());
    if (status) params.append('status', status);
    
    const response = await fetch(`${API_URL}/laporan/peminjaman/export?${params}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    });
    
    if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `Laporan_Peminjaman_${new Date().toISOString().split('T')[0]}.xlsx`;
        a.click();
        window.URL.revokeObjectURL(url);
    } else {
        alert('Gagal export: ' + response.statusText);
    }
}
```

---

## 🎯 Summary

Fitur export Excel ini memberikan:
- ✅ **Professional formatting** dengan styling yang menarik
- ✅ **Flexible filtering** by date & status
- ✅ **Complete data** dengan semua relasi
- ✅ **Easy to use** - tinggal klik download
- ✅ **Role-based access** untuk keamanan data

**Status**: ✅ **Ready for Production**

---

*Last Updated: 2025-12-20*
*Developed by: Backend Team*
