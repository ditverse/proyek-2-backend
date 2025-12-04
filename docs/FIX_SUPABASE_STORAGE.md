# Fix: Supabase Storage Configuration Error

## Error
```
Gagal upload file ke storage: supabase config incomplete
```

## Root Cause

File `.env` Anda **tidak lengkap**. Backend membutuhkan 3 environment variables untuk Supabase Storage:

### Yang Dibutuhkan:
1. ✅ `SUPABASE_URL` - Sudah ada
2. ❌ `SUPABASE_SERVICE_KEY` - **MISSING!**
3. ⚠️ `SUPABASE_BUCKET_NAME` - Nama variable salah (Anda pakai `SUPABASE_BUCKET`)

---

## Cara Mendapatkan Credentials

### 1. SUPABASE_SERVICE_KEY

**Langkah**:
1. Buka [Supabase Dashboard](https://supabase.com/dashboard)
2. Pilih project Anda
3. Klik **Settings** (⚙️) di sidebar kiri
4. Klik **API**
5. Scroll ke bagian **Project API keys**
6. Copy **`service_role` key** (bukan anon key!)

⚠️ **PENTING**: 
- Gunakan `service_role` key, BUKAN `anon` key
- `service_role` key bypass RLS (Row Level Security)
- Jangan expose key ini di frontend!

### 2. SUPABASE_BUCKET_NAME

**Langkah**:
1. Buka Supabase Dashboard
2. Klik **Storage** di sidebar kiri
3. Lihat nama bucket yang sudah Anda buat
4. Jika belum ada, buat bucket baru:
   - Klik **New bucket**
   - Nama: `surat-digital` (atau nama lain)
   - Public: **No** (private bucket)
   - Klik **Create bucket**

---

## Update File .env

Buka file `.env` dan tambahkan/update:

```env
# Supabase Storage Configuration
SUPABASE_URL=https://hdbproscdii...supabase.co
SUPABASE_SERVICE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...  # ← TAMBAHKAN INI!
SUPABASE_BUCKET_NAME=surat-digital  # ← GANTI dari SUPABASE_BUCKET

# Optional (sudah ada default)
STORAGE_SIGNED_URL_EXPIRES=600
MAX_UPLOAD_SIZE_MB=2
```

### Contoh Lengkap .env:

```env
# Environment Configuration
APP_ENV=development
PORT=8000

# Database Configuration
DATABASE_URL=postgresql://postgres:password@db.xxx.supabase.co:5432/postgres

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production

# Supabase Storage Configuration
SUPABASE_URL=https://xxx.supabase.co
SUPABASE_SERVICE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6...
SUPABASE_BUCKET_NAME=surat-digital

# Storage Settings (optional)
STORAGE_SIGNED_URL_EXPIRES=600
MAX_UPLOAD_SIZE_MB=2

# CORS
CORS_ALLOWED_ORIGIN=*
```

---

## Verifikasi

### 1. Cek Environment Variables Loaded

Tambahkan log sementara di `internal/config/supabase.go`:

```go
func GetSupabaseConfig() SupabaseConfig {
    cfg := SupabaseConfig{
        URL:        os.Getenv("SUPABASE_URL"),
        ServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
        Bucket:     os.Getenv("SUPABASE_BUCKET_NAME"),
        // ...
    }
    
    // Debug log (hapus setelah fix)
    log.Printf("Supabase Config: URL=%s, Bucket=%s, ServiceKey=%s", 
        cfg.URL, cfg.Bucket, 
        if cfg.ServiceKey != "" { "SET" } else { "EMPTY" })
    
    return cfg
}
```

### 2. Restart Server

Setelah update `.env`:

```bash
# Stop server (Ctrl+C)
# Start lagi
air
```

### 3. Test Upload

Coba upload surat lagi dari frontend.

**Expected**: ✅ Upload berhasil, no error

---

## Troubleshooting

### Error: "supabase config incomplete"

**Cek**:
```bash
# Di terminal, jalankan:
echo $env:SUPABASE_URL
echo $env:SUPABASE_SERVICE_KEY
echo $env:SUPABASE_BUCKET_NAME
```

Jika kosong, berarti `.env` tidak ter-load atau nama variable salah.

### Error: "upload gagal: 401 Unauthorized"

**Penyebab**: Service key salah atau expired

**Solusi**: 
1. Pastikan menggunakan `service_role` key, bukan `anon` key
2. Copy ulang key dari Supabase Dashboard
3. Pastikan tidak ada spasi di awal/akhir key

### Error: "upload gagal: 404 Not Found"

**Penyebab**: Bucket tidak ada atau nama salah

**Solusi**:
1. Cek nama bucket di Supabase Dashboard → Storage
2. Pastikan `SUPABASE_BUCKET_NAME` sama persis dengan nama bucket
3. Bucket harus sudah dibuat sebelumnya

### Error: "upload gagal: 403 Forbidden"

**Penyebab**: Bucket policy tidak allow upload

**Solusi**:
1. Buka Supabase Dashboard → Storage
2. Klik bucket Anda
3. Tab **Policies**
4. Pastikan ada policy untuk INSERT/UPDATE
5. Atau gunakan `service_role` key yang bypass RLS

---

## Setup Bucket Policy (Optional)

Jika ingin menggunakan RLS (Row Level Security):

```sql
-- Policy untuk authenticated users bisa upload
CREATE POLICY "Authenticated users can upload surat"
ON storage.objects FOR INSERT
TO authenticated
WITH CHECK (
  bucket_id = 'surat-digital' AND
  (storage.foldername(name))[1] = 'peminjaman'
);

-- Policy untuk authenticated users bisa read own files
CREATE POLICY "Users can read own surat"
ON storage.objects FOR SELECT
TO authenticated
USING (
  bucket_id = 'surat-digital'
);
```

Tapi dengan `service_role` key, policy ini tidak diperlukan karena bypass RLS.

---

## Quick Fix Checklist

- [ ] Copy `service_role` key dari Supabase Dashboard → Settings → API
- [ ] Buka file `.env`
- [ ] Tambahkan `SUPABASE_SERVICE_KEY=...`
- [ ] Ganti `SUPABASE_BUCKET` menjadi `SUPABASE_BUCKET_NAME`
- [ ] Pastikan bucket `surat-digital` sudah dibuat di Supabase Storage
- [ ] Restart server (`air`)
- [ ] Test upload surat lagi

---

**Status**: Waiting for user to add credentials to .env file
