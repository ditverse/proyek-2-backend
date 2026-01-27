# Architecture Overview - Sistem Peminjaman Sarana Prasarana

Dokumen ini menjelaskan arsitektur sistem secara keseluruhan.

---

## Arsitektur Sistem

```
+------------------+          +------------------+          +------------------+
|                  |   HTTP   |                  |   SQL    |                  |
|    Frontend      |<-------->|    Backend API   |<-------->|    Database      |
|    (React/Next)  |   REST   |    (Go Native)   |   pgx    |   (PostgreSQL)   |
|                  |          |                  |          |    Supabase      |
+------------------+          +------------------+          +------------------+
                                      |
                                      | OAuth2
                                      v
                              +------------------+
                              |                  |
                              |    Gmail API     |
                              |   (Notification) |
                              |                  |
                              +------------------+
                                      
                              +------------------+
                              |                  |
                              | Supabase Storage |
                              |   (File Upload)  |
                              |                  |
                              +------------------+
```

---

## Diagram Arsitektur Detail

```mermaid
flowchart TB
    subgraph Client["Client Layer"]
        WEB[Web Browser]
        MOBILE[Mobile App]
    end

    subgraph API["API Layer"]
        LB[Load Balancer / Reverse Proxy]
        
        subgraph Backend["Go Backend"]
            ROUTER[Router]
            MW[Middleware]
            
            subgraph Handlers["Handlers"]
                AH[Auth Handler]
                PH[Peminjaman Handler]
                RH[Ruangan Handler]
                BH[Barang Handler]
                KH[Kehadiran Handler]
                EH[Export Handler]
            end
            
            subgraph Services["Services"]
                AS[Auth Service]
                PS[Peminjaman Service]
                KS[Kehadiran Service]
                ES[Email Service]
                XS[Export Service]
            end
            
            subgraph Repos["Repositories"]
                UR[User Repo]
                PR[Peminjaman Repo]
                RR[Ruangan Repo]
                BR[Barang Repo]
                KR[Kehadiran Repo]
                MR[Mailbox Repo]
                LR[Log Repo]
            end
        end
    end

    subgraph External["External Services"]
        GMAIL[Gmail API]
        STORAGE[Supabase Storage]
    end

    subgraph Data["Data Layer"]
        DB[(PostgreSQL / Supabase)]
    end

    WEB --> LB
    MOBILE --> LB
    LB --> ROUTER
    ROUTER --> MW
    MW --> Handlers
    Handlers --> Services
    Services --> Repos
    Repos --> DB
    
    ES --> GMAIL
    PS --> STORAGE
```

---

## Layered Architecture

Sistem menggunakan arsitektur berlapis (Layered Architecture):

```
+-----------------------------------------------------------+
|                      PRESENTATION LAYER                    |
|  (Handlers - HTTP Request/Response)                        |
+-----------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------+
|                      BUSINESS LOGIC LAYER                  |
|  (Services - Core Business Rules)                          |
+-----------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------+
|                      DATA ACCESS LAYER                     |
|  (Repositories - Database Operations)                      |
+-----------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------+
|                      DATABASE LAYER                        |
|  (PostgreSQL via Supabase)                                 |
+-----------------------------------------------------------+
```

### Penjelasan Setiap Layer

| Layer | Komponen | Tanggung Jawab |
|-------|----------|----------------|
| **Presentation** | Handlers | Menerima HTTP request, validasi input, mengembalikan response |
| **Business Logic** | Services | Implementasi aturan bisnis, koordinasi antar repository |
| **Data Access** | Repositories | Operasi CRUD ke database, query building |
| **Database** | PostgreSQL | Penyimpanan data persisten |

---

## Komponen Utama

### 1. Internal Config

```
internal/
├── config/
│   ├── config.go          # Environment variables
│   ├── gmail.go           # Gmail OAuth2 config
│   └── supabase.go        # Supabase connection
├── db/
│   └── db.go              # Database connection pool
├── router/
│   └── router.go          # Route definitions
└── services/
    ├── email_service.go   # Gmail API integration
    ├── email_templates.go # HTML email templates
    ├── status_scheduler.go# Auto status update
    └── storage_service.go # Supabase storage
```

### 2. Business Domain

```
models/          # Data structures
repositories/    # Database operations
services/        # Business logic
handlers/        # HTTP controllers
middleware/      # Auth, CORS
```

---

## Request Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router
    participant M as Middleware
    participant H as Handler
    participant S as Service
    participant Repo as Repository
    participant DB as Database

    C->>R: HTTP Request
    R->>M: Route Match
    M->>M: CORS Check
    M->>M: JWT Validation
    M->>M: Role Authorization
    M->>H: Pass Request
    H->>H: Parse Request Body
    H->>S: Call Service Method
    S->>S: Business Logic
    S->>Repo: Database Operation
    Repo->>DB: SQL Query
    DB-->>Repo: Result Set
    Repo-->>S: Domain Object
    S-->>H: Result
    H-->>C: HTTP Response (JSON)
```

---

## Security Architecture

### Authentication Flow

```mermaid
flowchart LR
    subgraph Login
        A[User Submit Credentials]
        B[Validate Email/Password]
        C[Generate JWT Token]
        D[Return Token + User Data]
    end

    subgraph Protected Request
        E[Request with Bearer Token]
        F[Extract Token from Header]
        G[Validate JWT Signature]
        H[Check Token Expiry]
        I[Extract User Claims]
        J[Check Role Permission]
        K[Process Request]
    end

    A --> B --> C --> D
    E --> F --> G --> H --> I --> J --> K
```

### Security Measures

| Aspek | Implementasi |
|-------|--------------|
| **Authentication** | JWT Bearer Token |
| **Password** | bcrypt hashing (cost factor 10) |
| **Authorization** | Role-based Access Control (RBAC) |
| **CORS** | Configurable allowed origins |
| **SQL Injection** | Parameterized queries (pgx) |
| **Token Expiry** | 24 hours default |

---

## Database Design

### Schema Overview

```mermaid
flowchart TB
    subgraph Master["Master Data"]
        USERS[users]
        ORGANISASI[organisasi]
        RUANGAN[ruangan]
        BARANG[barang]
    end

    subgraph Transaction["Transaction Data"]
        PEMINJAMAN[peminjaman]
        PEMINJAMAN_BARANG[peminjaman_barang]
        KEGIATAN[kegiatan]
        KEHADIRAN[kehadiran_peminjam]
    end

    subgraph System["System Data"]
        MAILBOX[mailbox]
        LOG[log_aktivitas]
    end

    USERS --> PEMINJAMAN
    ORGANISASI --> USERS
    ORGANISASI --> KEGIATAN
    RUANGAN --> PEMINJAMAN
    RUANGAN --> BARANG
    BARANG --> PEMINJAMAN_BARANG
    PEMINJAMAN --> PEMINJAMAN_BARANG
    KEGIATAN --> PEMINJAMAN
    PEMINJAMAN --> KEHADIRAN
    PEMINJAMAN --> MAILBOX
    PEMINJAMAN --> LOG
    USERS --> LOG
    USERS --> MAILBOX
    USERS --> KEHADIRAN
```

### Database Features

- **Auto-generated Codes**: Triggers untuk generate kode unik
- **Timestamps**: created_at dan updated_at otomatis
- **Foreign Keys**: Referential integrity
- **Enums**: Type-safe status values

---

## External Integrations

### 1. Supabase Storage

```
Purpose: File storage untuk surat digital
Bucket: surat-digital
File Path: peminjaman/{kode_peminjaman}/surat.pdf
Max Size: 2MB
Format: PDF only
```

### 2. Gmail API

```
Purpose: Notifikasi email otomatis
OAuth2: Service account atau user credentials
Templates: HTML email templates
Types:
  - New Submission (ke Sarpras)
  - Approved/Rejected (ke Mahasiswa)
  - Security Notification (ke Security)
  - Cancellation (ke semua pihak)
```

---

## Deployment Architecture

```
+------------------+     +------------------+     +------------------+
|   Development    |     |     Staging      |     |    Production    |
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  localhost:8000  |     |  staging.xyz.com |     |  api.xyz.com     |
|                  |     |                  |     |                  |
|  Air hot-reload  |     |  Docker          |     |  Docker/K8s      |
|  Local Postgres  |     |  Supabase Dev    |     |  Supabase Prod   |
|                  |     |                  |     |  Cloudflare CDN  |
|                  |     |                  |     |  SSL/TLS         |
+------------------+     +------------------+     +------------------+
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `PORT` | Server port | No (default: 8000) |
| `JWT_SECRET` | Secret for JWT signing | Yes |
| `SUPABASE_URL` | Supabase project URL | Yes |
| `SUPABASE_SERVICE_KEY` | Supabase service role key | Yes |
| `SUPABASE_BUCKET_NAME` | Storage bucket name | Yes |
| `GMAIL_CREDENTIALS_FILE` | Path to Gmail credentials | No |
| `GMAIL_TOKEN_FILE` | Path to Gmail OAuth token | No |
| `CORS_ALLOWED_ORIGIN` | Allowed CORS origins | No (default: *) |

---

## Performance Considerations

### Optimizations Implemented

1. **Connection Pooling**: pgx connection pool for database
2. **Batch Queries**: GetPeminjamanBarangByIDs for bulk fetch
3. **Async Email**: Goroutines for non-blocking email sending
4. **Minimal JSON**: Password hash excluded from responses
5. **Index**: Primary keys dan foreign keys indexed

### Recommendations

1. Add Redis for session caching
2. Implement rate limiting
3. Add request logging middleware
4. Consider GraphQL for flexible queries
5. Add health check endpoints for monitoring

---

## Error Handling

```go
// Standard error response format
{
  "error": "Error message description"
}

// HTTP Status Codes used:
// 200 - Success
// 201 - Created
// 400 - Bad Request (validation error)
// 401 - Unauthorized (no/invalid token)
// 403 - Forbidden (insufficient role)
// 404 - Not Found
// 500 - Internal Server Error
```
