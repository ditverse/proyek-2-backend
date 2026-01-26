# Sequence Diagram - Sistem Peminjaman Sarana Prasarana Kampus

## Sequence Diagram Keseluruhan Sistem

```mermaid
sequenceDiagram
    autonumber
    
    box Aktor
        participant M as Mahasiswa
        participant SP as Sarpras
        participant SC as Security
    end
    
    box Frontend
        participant FE as Frontend
    end
    
    box Backend
        participant API as REST API
        participant SVC as Services
        participant REPO as Repositories
    end
    
    box External
        participant DB as Database
        participant SB as Supabase Storage
        participant GM as Gmail API
    end

    Note over M,GM: FASE 1: PENGAJUAN PEMINJAMAN

    M->>FE: Login ke sistem
    FE->>API: POST /api/auth/login
    API->>SVC: AuthService.Login()
    SVC->>REPO: UserRepository.GetByEmail()
    REPO->>DB: SELECT user
    DB-->>REPO: User data
    REPO-->>SVC: User object
    SVC-->>API: JWT Token
    API-->>FE: Token + User info
    FE-->>M: Dashboard Mahasiswa

    M->>FE: Ajukan peminjaman + upload surat
    FE->>SB: Upload surat digital
    SB-->>FE: File path
    FE->>API: POST /api/peminjaman
    API->>SVC: PeminjamanService.Create()
    SVC->>REPO: PeminjamanRepository.Create()
    REPO->>DB: INSERT peminjaman
    DB-->>REPO: Created
    SVC->>REPO: MailboxRepository.Create()
    REPO->>DB: INSERT notifikasi untuk Sarpras
    API-->>FE: Peminjaman created
    FE-->>M: Pengajuan berhasil

    Note over M,GM: FASE 2: VERIFIKASI OLEH SARPRAS

    SP->>FE: Login ke sistem
    FE->>API: POST /api/auth/login
    API-->>FE: Token + User info
    FE-->>SP: Dashboard Sarpras

    SP->>FE: Lihat pengajuan pending
    FE->>API: GET /api/peminjaman/pending
    API->>REPO: GetPendingPeminjaman()
    REPO->>DB: SELECT WHERE status=PENDING
    DB-->>REPO: List peminjaman
    API-->>FE: Daftar pending
    FE-->>SP: Tampilkan daftar

    SP->>FE: Verifikasi peminjaman (Approve)
    FE->>API: POST /api/peminjaman/{id}/verifikasi
    API->>SVC: PeminjamanService.Verifikasi()
    SVC->>REPO: UpdateStatus(APPROVED)
    REPO->>DB: UPDATE status
    SVC->>REPO: Create notifikasi untuk Mahasiswa
    REPO->>DB: INSERT notifikasi
    SVC->>GM: Send email notification
    GM-->>SVC: Email sent
    API-->>FE: Verifikasi berhasil
    FE-->>SP: Tampilkan sukses

    Note over M,GM: FASE 3: VERIFIKASI KEHADIRAN OLEH SECURITY

    SC->>FE: Login ke sistem
    FE->>API: POST /api/auth/login
    API-->>FE: Token + User info
    FE-->>SC: Dashboard Security

    SC->>FE: Lihat jadwal aktif hari ini
    FE->>API: GET /api/jadwal-aktif-belum-verifikasi
    API->>REPO: GetJadwalAktif()
    REPO->>DB: SELECT WHERE date=today
    DB-->>REPO: List jadwal
    API-->>FE: Jadwal aktif
    FE-->>SC: Tampilkan jadwal

    SC->>FE: Verifikasi kehadiran (HADIR)
    FE->>API: POST /api/kehadiran
    API->>SVC: KehadiranService.Catat()
    SVC->>REPO: KehadiranRepository.Create()
    REPO->>DB: INSERT kehadiran
    SVC->>REPO: UpdateStatus(ONGOING)
    REPO->>DB: UPDATE peminjaman
    SVC->>REPO: Create notifikasi untuk Sarpras
    REPO->>DB: INSERT notifikasi
    API-->>FE: Kehadiran tercatat
    FE-->>SC: Tampilkan sukses

    Note over M,GM: FASE 4: SELESAI

    SC->>FE: Kegiatan selesai
    FE->>API: Update status FINISHED
    API->>REPO: UpdateStatus(FINISHED)
    REPO->>DB: UPDATE peminjaman
    API-->>FE: Status updated
    FE-->>SC: Peminjaman selesai
```

## Sequence Diagram: Proses Login

```mermaid
sequenceDiagram
    autonumber
    participant U as User
    participant FE as Frontend
    participant H as AuthHandler
    participant S as AuthService
    participant R as UserRepository
    participant DB as Database

    U->>FE: Masukkan email dan password
    FE->>H: POST /api/auth/login
    H->>S: Login(email, password)
    S->>R: GetByEmail(email)
    R->>DB: SELECT * FROM users WHERE email = ?
    DB-->>R: User data
    R-->>S: User object
    
    alt Password Valid
        S->>S: bcrypt.CompareHashAndPassword()
        S->>S: GenerateJWT(user)
        S-->>H: JWT Token + User
        H-->>FE: 200 OK {token, user}
        FE->>FE: Simpan token di localStorage
        FE-->>U: Redirect ke Dashboard
    else Password Invalid
        S-->>H: Error: Invalid credentials
        H-->>FE: 401 Unauthorized
        FE-->>U: Tampilkan error message
    end
```

## Sequence Diagram: Proses Pengajuan Peminjaman

```mermaid
sequenceDiagram
    autonumber
    participant M as Mahasiswa
    participant FE as Frontend
    participant PH as PeminjamanHandler
    participant PS as PeminjamanService
    participant PR as PeminjamanRepository
    participant MR as MailboxRepository
    participant SS as StorageService
    participant DB as Database
    participant SB as Supabase Storage

    M->>FE: Isi form peminjaman
    M->>FE: Upload surat digital PDF
    FE->>SS: Upload file ke Supabase
    SS->>SB: PUT /storage/v1/object/surat-digital
    SB-->>SS: File path
    SS-->>FE: path_surat_digital
    
    FE->>PH: POST /api/peminjaman
    Note over FE,PH: Header: Authorization: Bearer token
    
    PH->>PH: Validasi JWT Token
    PH->>PS: CreatePeminjaman(request)
    
    PS->>PS: Validasi tanggal tidak bentrok
    PS->>PR: CheckRuanganAvailability()
    PR->>DB: SELECT * FROM peminjaman WHERE ruangan AND tanggal
    DB-->>PR: Existing bookings
    PR-->>PS: Availability status
    
    alt Ruangan Tersedia
        PS->>PS: GenerateKodePeminjaman()
        PS->>PR: Create(peminjaman)
        PR->>DB: INSERT INTO peminjaman
        DB-->>PR: Created record
        PR-->>PS: Peminjaman object
        
        PS->>MR: CreateNotifikasi(to: SARPRAS)
        MR->>DB: INSERT INTO mailbox
        DB-->>MR: Notifikasi created
        
        PS-->>PH: Peminjaman berhasil
        PH-->>FE: 201 Created {peminjaman}
        FE-->>M: Tampilkan success message
    else Ruangan Tidak Tersedia
        PS-->>PH: Error: Jadwal bentrok
        PH-->>FE: 400 Bad Request
        FE-->>M: Tampilkan error: Jadwal bentrok
    end
```

## Sequence Diagram: Proses Verifikasi Peminjaman oleh Sarpras

```mermaid
sequenceDiagram
    autonumber
    participant SP as Sarpras
    participant FE as Frontend
    participant PH as PeminjamanHandler
    participant PS as PeminjamanService
    participant PR as PeminjamanRepository
    participant MR as MailboxRepository
    participant ES as EmailService
    participant DB as Database
    participant GM as Gmail API

    SP->>FE: Buka daftar pengajuan pending
    FE->>PH: GET /api/peminjaman/pending
    PH->>PR: GetPendingPeminjaman()
    PR->>DB: SELECT * FROM peminjaman WHERE status = 'PENDING'
    DB-->>PR: List peminjaman
    PR-->>PH: Peminjaman list
    PH-->>FE: 200 OK {peminjaman[]}
    FE-->>SP: Tampilkan daftar pending

    SP->>FE: Pilih peminjaman untuk review
    FE->>PH: GET /api/peminjaman/{id}
    PH->>PR: GetByID(id)
    PR->>DB: SELECT * FROM peminjaman WHERE kode = ?
    DB-->>PR: Peminjaman detail
    PR-->>PH: Peminjaman object
    PH-->>FE: 200 OK {peminjaman}
    FE-->>SP: Tampilkan detail peminjaman

    SP->>FE: Approve/Reject peminjaman
    FE->>PH: POST /api/peminjaman/{id}/verifikasi
    Note over FE,PH: Body: {status: "APPROVED/REJECTED", catatan: "..."}
    
    PH->>PS: VerifikasiPeminjaman(id, status, catatan, verified_by)
    PS->>PR: UpdateStatus(id, status, verified_by, catatan)
    PR->>DB: UPDATE peminjaman SET status = ?, verified_by = ?
    DB-->>PR: Updated
    PR-->>PS: Success

    PS->>MR: CreateNotifikasi(to: Mahasiswa)
    MR->>DB: INSERT INTO mailbox
    DB-->>MR: Notifikasi created

    PS->>ES: SendEmailAsync(to: Mahasiswa)
    ES->>GM: Send email notification
    GM-->>ES: Email sent

    PS-->>PH: Verifikasi berhasil
    PH-->>FE: 200 OK {message: "Berhasil diverifikasi"}
    FE-->>SP: Tampilkan success message
```

## Sequence Diagram: Proses Verifikasi Kehadiran oleh Security

```mermaid
sequenceDiagram
    autonumber
    participant SC as Security
    participant FE as Frontend
    participant KH as KehadiranHandler
    participant KS as KehadiranService
    participant KR as KehadiranRepository
    participant PR as PeminjamanRepository
    participant MR as MailboxRepository
    participant DB as Database

    SC->>FE: Buka jadwal aktif hari ini
    FE->>KH: GET /api/jadwal-aktif-belum-verifikasi
    KH->>PR: GetJadwalAktifBelumVerifikasi(today)
    PR->>DB: SELECT * FROM peminjaman WHERE status = 'APPROVED' AND date = today
    DB-->>PR: List jadwal aktif
    PR-->>KH: Jadwal list
    KH-->>FE: 200 OK {jadwal[]}
    FE-->>SC: Tampilkan jadwal aktif

    SC->>FE: Verifikasi kehadiran peminjam
    FE->>KH: POST /api/kehadiran
    Note over FE,KH: Body: {kode_peminjaman, status_kehadiran: "HADIR/TIDAK_HADIR/BATAL"}

    KH->>KS: CatatKehadiran(kode_peminjaman, status, verified_by)
    KS->>KR: Create(kehadiran)
    KR->>DB: INSERT INTO kehadiran
    DB-->>KR: Kehadiran created
    KR-->>KS: Kehadiran object

    alt Status HADIR
        KS->>PR: UpdateStatus(kode_peminjaman, "ONGOING")
        PR->>DB: UPDATE peminjaman SET status = 'ONGOING'
        DB-->>PR: Updated
    else Status TIDAK_HADIR atau BATAL
        KS->>PR: UpdateStatus(kode_peminjaman, status)
        PR->>DB: UPDATE peminjaman SET status = ?
        DB-->>PR: Updated
    end

    KS->>MR: CreateNotifikasi(to: Sarpras)
    MR->>DB: INSERT INTO mailbox
    DB-->>MR: Notifikasi created

    KS-->>KH: Kehadiran tercatat
    KH-->>FE: 200 OK {message: "Kehadiran berhasil dicatat"}
    FE-->>SC: Tampilkan success message
```

## Sequence Diagram: Proses Upload Surat Digital

```mermaid
sequenceDiagram
    autonumber
    participant M as Mahasiswa
    participant FE as Frontend
    participant PH as PeminjamanHandler
    participant SS as StorageService
    participant PR as PeminjamanRepository
    participant SB as Supabase Storage
    participant DB as Database

    M->>FE: Pilih file PDF surat
    FE->>FE: Validasi file (max 2MB, PDF only)
    
    alt File Valid
        FE->>PH: POST /api/peminjaman/{id}/upload-surat
        Note over FE,PH: Content-Type: multipart/form-data
        
        PH->>SS: UploadFile(file, bucket)
        SS->>SB: PUT /storage/v1/object/surat-digital/{filename}
        SB-->>SS: Public URL / Path
        SS-->>PH: file_path
        
        PH->>PR: UpdatePathSurat(id, file_path)
        PR->>DB: UPDATE peminjaman SET path_surat_digital = ?
        DB-->>PR: Updated
        PR-->>PH: Success
        
        PH-->>FE: 200 OK {path: file_path}
        FE-->>M: Tampilkan preview surat
    else File Invalid
        FE-->>M: Error: File harus PDF dan maksimal 2MB
    end
```

## Sequence Diagram: Proses Export Laporan PDF/Excel

```mermaid
sequenceDiagram
    autonumber
    participant SP as Sarpras
    participant FE as Frontend
    participant EH as ExportHandler
    participant ES as ExportService
    participant PR as PeminjamanRepository
    participant DB as Database

    SP->>FE: Buka menu laporan peminjaman
    SP->>FE: Set filter (tanggal, status)
    SP->>FE: Klik tombol Export PDF/Excel
    
    alt Export PDF
        FE->>EH: GET /api/export/peminjaman/pdf?start=...&end=...&status=...
        EH->>ES: ExportPeminjamanPDF(filters)
        ES->>PR: GetPeminjamanByFilters(start, end, status)
        PR->>DB: SELECT * FROM peminjaman WHERE date BETWEEN ? AND ?
        DB-->>PR: Peminjaman list
        PR-->>ES: Data peminjaman
        
        ES->>ES: GeneratePDF(data)
        ES-->>EH: PDF byte stream
        EH-->>FE: Content-Type: application/pdf
        FE-->>SP: Download file PDF
    else Export Excel
        FE->>EH: GET /api/export/peminjaman/excel?start=...&end=...&status=...
        EH->>ES: ExportPeminjamanExcel(filters)
        ES->>PR: GetPeminjamanByFilters(start, end, status)
        PR->>DB: SELECT * FROM peminjaman WHERE date BETWEEN ? AND ?
        DB-->>PR: Peminjaman list
        PR-->>ES: Data peminjaman
        
        ES->>ES: GenerateExcel(data)
        ES-->>EH: Excel byte stream
        EH-->>FE: Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
        FE-->>SP: Download file Excel
    end
```

## Sequence Diagram: Proses Melihat dan Menandai Notifikasi

```mermaid
sequenceDiagram
    autonumber
    participant U as User
    participant FE as Frontend
    participant MH as MailboxHandler
    participant MR as MailboxRepository
    participant DB as Database

    U->>FE: Klik icon notifikasi
    FE->>MH: GET /api/notifikasi/count
    MH->>MR: GetUnreadCount(user_id)
    MR->>DB: SELECT COUNT(*) FROM mailbox WHERE user_id = ? AND is_read = false
    DB-->>MR: Count
    MR-->>MH: Unread count
    MH-->>FE: 200 OK {count: N}
    FE->>FE: Update badge notifikasi

    FE->>MH: GET /api/notifikasi/me
    MH->>MR: GetByUserID(user_id)
    MR->>DB: SELECT * FROM mailbox WHERE user_id = ? ORDER BY created_at DESC
    DB-->>MR: Notifikasi list
    MR-->>MH: Notifikasi list
    MH-->>FE: 200 OK {notifikasi[]}
    FE-->>U: Tampilkan daftar notifikasi

    U->>FE: Klik notifikasi untuk dibaca
    FE->>MH: PATCH /api/notifikasi/{id}/dibaca
    MH->>MR: MarkAsRead(id)
    MR->>DB: UPDATE mailbox SET is_read = true WHERE id = ?
    DB-->>MR: Updated
    MR-->>MH: Success
    MH-->>FE: 200 OK
    FE->>FE: Update tampilan notifikasi
    FE-->>U: Notifikasi ditandai sudah dibaca
```
