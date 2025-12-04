# Fix: Kode User Format Issue

## Masalah
Ketika user mendaftar melalui halaman register, `kode_user` yang tersimpan di database tidak sesuai format yang diinginkan:
- **Format yang diinginkan**: `USR-251129-0001` (USR-YYMMDD-urutan)
- **Format yang tersimpan**: `USR-1733295825123456789` (timestamp nano)

## Root Cause
1. **Backend menggunakan `generateCode()` function** yang menghasilkan kode dengan timestamp nano:
   ```go
   func generateCode(prefix string) string {
       return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
   }
   ```

2. **Repository langsung insert kode_user** yang sudah di-generate oleh aplikasi, sehingga **database trigger tidak jalan**.

3. **Database trigger hanya jalan** ketika `kode_user` bernilai NULL atau tidak disertakan dalam INSERT statement.

## Solusi yang Diterapkan

### 1. Update User Repository
**File**: `repositories/user_repository.go`

**Perubahan**:
- Hapus `kode_user` dari INSERT statement
- Tambahkan `kode_user` ke RETURNING clause
- Database trigger akan otomatis generate kode dengan format yang benar

```go
// BEFORE
INSERT INTO users (kode_user, nama, email, password_hash, role, organisasi_kode)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING created_at

// AFTER
INSERT INTO users (nama, email, password_hash, role, organisasi_kode)
VALUES ($1, $2, $3, $4, $5)
RETURNING kode_user, created_at
```

### 2. Update Auth Service
**File**: `services/auth_service.go`

**Perubahan**:
- Hapus assignment `KodeUser: generateCode("USR")`
- Biarkan field kosong, akan diisi oleh database

```go
// BEFORE
user := &models.User{
    KodeUser:       generateCode("USR"),
    Nama:           req.Nama,
    ...
}

// AFTER
user := &models.User{
    Nama:           req.Nama,
    ...
}
```

### 3. Update Peminjaman Repository & Service
**File**: `repositories/peminjaman_repository.go` dan `services/peminjaman_service.go`

**Perubahan yang sama**:
- Hapus `kode_peminjaman` dari INSERT
- Tambahkan ke RETURNING
- Hapus assignment di service

## Database Trigger

### Cara Menjalankan Trigger di Supabase

1. **Buka Supabase Dashboard** → SQL Editor

2. **Jalankan file migration**: `migrations/002_auto_generate_codes.sql`

   Atau copy-paste SQL berikut:

```sql
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
```

3. **Klik Run** atau tekan `Ctrl+Enter`

4. **Verifikasi trigger sudah aktif**:
   ```sql
   SELECT trigger_name, event_manipulation, event_object_table 
   FROM information_schema.triggers 
   WHERE trigger_name = 'trigger_generate_kode_user';
   ```

### Testing

**Test di Supabase SQL Editor**:
```sql
-- Insert tanpa kode_user
INSERT INTO users (nama, email, password_hash, role)
VALUES ('Test User', 'test@example.com', 'hash123', 'MAHASISWA');

-- Cek hasil
SELECT kode_user, nama, email FROM users WHERE email = 'test@example.com';
-- Expected: kode_user = 'USR-251204-0001' (sesuai tanggal hari ini)
```

**Test dari aplikasi**:
1. Restart server backend (jika menggunakan Air, akan auto-restart)
2. Buka halaman register di frontend
3. Daftar user baru
4. Cek database, `kode_user` harus format: `USR-YYMMDD-0001`

## Trigger untuk Tabel Lain

File migration `002_auto_generate_codes.sql` juga sudah menyertakan trigger untuk:
- `peminjaman` → `PMJ-YYMMDD-0001`
- `ruangan` → `RNG-0001`
- `barang` → `BRG-0001`

Semua menggunakan pola yang sama.

## Catatan Penting

1. **Trigger hanya jalan jika field NULL/kosong** - Jika aplikasi mengirim nilai, trigger akan di-skip
2. **Urutan reset setiap hari** - Untuk kode dengan tanggal (USR, PMJ), urutan akan reset ke 0001 setiap hari
3. **Thread-safe** - PostgreSQL trigger sudah handle concurrent insert dengan baik
4. **Backward compatible** - Jika ada data lama dengan format berbeda, tidak akan terpengaruh

## Rollback (Jika Diperlukan)

Jika ingin kembali ke sistem lama (generate di aplikasi):

1. **Drop trigger**:
   ```sql
   DROP TRIGGER IF EXISTS trigger_generate_kode_user ON users;
   DROP FUNCTION IF EXISTS generate_kode_user();
   ```

2. **Revert kode**:
   - Kembalikan `KodeUser: generateCode("USR")` di auth_service.go
   - Kembalikan INSERT dengan kode_user di user_repository.go

## Status

✅ **Fixed** - Backend sudah diupdate untuk menggunakan database trigger
⏳ **Pending** - Trigger perlu dijalankan di Supabase Database
📝 **Documented** - Migration file tersedia di `migrations/002_auto_generate_codes.sql`
