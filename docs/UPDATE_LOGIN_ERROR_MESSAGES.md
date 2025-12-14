# Update: Login Error Messages untuk Frontend

## Tanggal Update
14 Desember 2025

## Latar Belakang
Frontend login sudah disesuaikan untuk menampilkan pesan error yang spesifik berdasarkan response dari backend. Sebelumnya, backend hanya mengembalikan pesan generik "invalid credentials" untuk semua jenis kesalahan login, sehingga frontend tidak dapat membedakan apakah error disebabkan oleh email yang tidak terdaftar atau password yang salah.

## Masalah Sebelumnya
Backend mengembalikan pesan error yang sama untuk dua kondisi berbeda:

```go
// Ketika email tidak ditemukan
return nil, errors.New("invalid credentials")

// Ketika password salah
return nil, errors.New("invalid credentials")
```

**Dampak**: Frontend tidak bisa menampilkan pesan error di field yang tepat (email atau password).

## Solusi yang Diterapkan

### File yang Diubah
**File**: `services/auth_service.go`

### Perubahan Kode

```diff
func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
    user, err := s.UserRepo.GetByEmail(email)
    if err != nil {
        return nil, err
    }
    if user == nil {
-       return nil, errors.New("invalid credentials")
+       return nil, errors.New("Email tidak ditemukan")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
-       return nil, errors.New("invalid credentials")
+       return nil, errors.New("Password salah")
    }
    // ... rest of the function
}
```

## Alur Logic Login

```
┌─────────────────────────────────────┐
│         POST /auth/login            │
│   {"email": "...", "password": ...} │
└───────────────┬─────────────────────┘
                │
                ▼
        ┌───────────────┐
        │ Cek email di  │
        │   database    │
        └───────┬───────┘
                │
        ┌───────┴───────┐
        │  Email ada?   │
        └───────┬───────┘
                │
       NO ◄─────┴─────► YES
        │               │
        ▼               ▼
   ┌──────────┐  ┌─────────────┐
   │   401    │  │ Cek password│
   │  "Email  │  └──────┬──────┘
   │  tidak   │         │
   │ditemukan"│  ┌──────┴──────┐
   └──────────┘  │Password OK? │
                 └──────┬──────┘
                        │
               NO ◄─────┴─────► YES
                │               │
                ▼               ▼
           ┌─────────┐   ┌───────────┐
           │   401   │   │    200    │
           │"Password│   │  {token,  │
           │  salah" │   │   user}   │
           └─────────┘   └───────────┘
```

## API Response Format

### 1. Login Sukses (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "kode_user": "USR-251214-0001",
    "nama": "User Name",
    "email": "user@example.com",
    "role": "mahasiswa",
    "organisasi_kode": null,
    "created_at": "2025-12-14T20:00:00Z"
  }
}
```

### 2. Email Tidak Ditemukan (401 Unauthorized)
```json
{
  "error": "Email tidak ditemukan"
}
```
**Kata kunci**: `email`, `tidak ditemukan`

### 3. Password Salah (401 Unauthorized)
```json
{
  "error": "Password salah"
}
```
**Kata kunci**: `password`, `salah`

## Integrasi dengan Frontend

Frontend mendeteksi jenis error berdasarkan kata kunci dalam pesan error:

| Kata Kunci dalam Error | Field Error di Frontend |
|------------------------|-------------------------|
| `email`, `tidak ditemukan`, `not found`, `user` | Input Email |
| `password`, `salah`, `wrong`, `invalid` | Input Password |

### Contoh Implementasi Frontend
```javascript
const errorMessage = response.error.toLowerCase();

if (errorMessage.includes('email') || 
    errorMessage.includes('tidak ditemukan') ||
    errorMessage.includes('not found') ||
    errorMessage.includes('user')) {
    // Tampilkan error di field email
    showEmailError(response.error);
} else if (errorMessage.includes('password') ||
           errorMessage.includes('salah') ||
           errorMessage.includes('wrong') ||
           errorMessage.includes('invalid')) {
    // Tampilkan error di field password
    showPasswordError(response.error);
}
```

## Role User untuk Redirect

Field `role` dalam response user digunakan untuk redirect ke dashboard yang sesuai:

| Role       | Dashboard Redirect    |
|------------|----------------------|
| mahasiswa  | /dashboard/mahasiswa |
| sarpras    | /dashboard/sarpras   |
| security   | /dashboard/security  |

## Testing

### Test Case 1: Email Tidak Terdaftar
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "notexist@example.com", "password": "anypassword"}'

# Expected Response (401):
# {"error": "Email tidak ditemukan"}
```

### Test Case 2: Password Salah
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "existing@example.com", "password": "wrongpassword"}'

# Expected Response (401):
# {"error": "Password salah"}
```

### Test Case 3: Login Berhasil
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "existing@example.com", "password": "correctpassword"}'

# Expected Response (200):
# {"token": "eyJ...", "user": {...}}
```

## Status

✅ **Implemented** - Pesan error sudah diperbarui di `auth_service.go`
✅ **Compatible** - Mendukung deteksi kata kunci frontend
✅ **Documented** - Dokumentasi tersedia di `docs/UPDATE_LOGIN_ERROR_MESSAGES.md`
