# Diagram Metodologi Scrum

## 1. Alur Proses Scrum

```mermaid
flowchart LR
    subgraph Input
        PB["📋 Product Backlog"]
    end
    
    subgraph Sprint
        direction TB
        SP["📝 Sprint Planning"]
        SB["📌 Sprint Backlog"]
        DEV["⚙️ Development"]
        DS["🗣️ Daily Standup"]
        
        SP --> SB
        SB --> DEV
        DEV <--> DS
    end
    
    subgraph Output
        INC["📦 Product Increment"]
        SR["🔍 Sprint Review"]
        RET["🔄 Retrospective"]
    end
    
    PB --> SP
    DEV --> INC
    INC --> SR
    SR --> RET
    RET -.->|Feedback| PB
    SR -.->|Next Sprint| SP
```

---

## 2. Siklus Sprint

```mermaid
flowchart TB
    A["1. SPRINT PLANNING"]
    B["2. DAILY STANDUP"]
    C["3. DEVELOPMENT"]
    D["4. SPRINT REVIEW"]
    E["5. RETROSPECTIVE"]
    
    A --> B
    B --> C
    C --> B
    C --> D
    D --> E
    E --> A
```

---

## 3. Struktur Tim Scrum

```mermaid
flowchart TB
    subgraph Team
        PO["🎯 Product Owner"]
        SM["🛡️ Scrum Master"]
        DT["💻 Development Team"]
    end
    
    PO <--> SM
    SM <--> DT
    DT <--> PO
```

---

## 4. Artefak Scrum

```mermaid
flowchart LR
    PB["📋 Product Backlog"]
    SB["📌 Sprint Backlog"]
    PI["📦 Product Increment"]
    
    PB --> SB --> PI
```

---

## 5. Timeline Sprint Proyek

```mermaid
gantt
    title Timeline Sprint
    dateFormat  YYYY-MM-DD
    
    section Sprint 1
    Fondasi dan Autentikasi    :s1, 2024-01-01, 14d
    
    section Sprint 2
    Core Peminjaman          :s2, after s1, 21d
    
    section Sprint 3
    Administrasi dan Kehadiran :s3, after s2, 14d
    
    section Sprint 4
    Notifikasi dan Finalisasi  :s4, after s3, 14d
```

---

## 6. Alur Detail Scrum

```mermaid
flowchart TD
    subgraph Backlog
        PB["Product Backlog<br/>56 Fitur"]
    end
    
    subgraph Sprint1
        S1["Sprint 1<br/>11 Fitur"]
    end
    
    subgraph Sprint2
        S2["Sprint 2<br/>16 Fitur"]
    end
    
    subgraph Sprint3
        S3["Sprint 3<br/>14 Fitur"]
    end
    
    subgraph Sprint4
        S4["Sprint 4<br/>15 Fitur"]
    end
    
    subgraph Result
        PROD["Production Ready<br/>System"]
    end
    
    PB --> S1
    S1 --> S2
    S2 --> S3
    S3 --> S4
    S4 --> PROD
```

---

## 7. Event Scrum dengan Durasi

```mermaid
flowchart LR
    E1["Sprint<br/>2-4 minggu"]
    E2["Planning<br/>4-8 jam"]
    E3["Daily<br/>15 menit"]
    E4["Review<br/>2-4 jam"]
    E5["Retro<br/>1-3 jam"]
    
    E1 --- E2 --- E3 --- E4 --- E5
```

---

## Catatan

Untuk melihat diagram:
1. VS Code dengan extension Markdown Preview Mermaid
2. Paste ke mermaid.live
3. Export sebagai gambar untuk laporan
