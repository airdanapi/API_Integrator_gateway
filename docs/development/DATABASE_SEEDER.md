# Panduan Database Seeder — Sprint 5

Dokumen ini menjelaskan cara menggunakan database seeder yang dibuat pada Sprint 5.  
Seeder mengisi data awal (sample/development data) ke empat tabel baru:
`request_logs`, `notifications`, `chat_messages`, dan `dashboard_data`.

---

## Daftar Isi

1. [Prasyarat](#1-prasyarat)
2. [Cara Menjalankan Seeder](#2-cara-menjalankan-seeder)
   - [2a. Via Server Startup (SEED_USERS_ENABLED)](#2a-via-server-startup)
   - [2b. Via CLI Standalone (`cmd/seed`)](#2b-via-cli-standalone)
   - [2c. Via Docker Compose](#2c-via-docker-compose)
3. [Konfigurasi Seeder](#3-konfigurasi-seeder)
4. [Data yang Di-seed](#4-data-yang-di-seed)
5. [Idempoten & Keamanan](#5-idempoten--keamanan)
6. [Cara Menambah Data Seed Baru](#6-cara-menambah-data-seed-baru)
7. [Troubleshooting](#7-troubleshooting)

---

## 1. Prasyarat

| Prasyarat | Keterangan |
|-----------|-----------|
| Go 1.21+ | Untuk menjalankan CLI seed runner |
| MySQL 8.0+ | Database target yang sudah berjalan |
| File `.env` | Environment variables database sudah dikonfigurasi |
| Migrasi sudah berjalan | Tabel sudah dibuat (seeder memanggil migrate otomatis) |

Pastikan file `.env` sudah ada di root project. Lihat [`.env.example`](../../.env.example) sebagai referensi.

---

## 2. Cara Menjalankan Seeder

### 2a. Via Server Startup

Seeder **users** (Sprint 3) berjalan otomatis saat server start jika `SEED_USERS_ENABLED=true`.

Untuk seeder Sprint 5 (logs, notifikasi, chat, dashboard cache), gunakan CLI standalone di bawah.

---

### 2b. Via CLI Standalone

CLI seed runner tersedia di `backend/cmd/seed/`.

#### Menjalankan dengan Default

```powershell
# Dari root project
cd backend
go run ./cmd/seed
```

Output yang diharapkan:

```
[seed] menghubungkan ke database...
[seed] menjalankan migrasi...
[seed] mulai seeding Sprint 5 (logs=50, chats=20, verbose=false)...
[seed] selesai seeding Sprint 5 dengan sukses.
```

#### Menjalankan dengan Verbose

Flag `-verbose` menampilkan detail setiap record yang di-seed:

```powershell
go run ./cmd/seed -verbose
```

Output yang diharapkan:

```
[seed] menghubungkan ke database...
[seed] menjalankan migrasi...
[seed] mulai seeding Sprint 5 (logs=50, chats=20, verbose=true)...
[seed] mulai seeding tabel: request_logs
  [request_logs] #1: Marketplace → /gateway/payment (200)
  [request_logs] #2: POS → /gateway/smartbank (200)
  ...
[seed] selesai seeding tabel: request_logs
[seed] mulai seeding tabel: notifications
  [notifications] #1: [api_inactive] Marketplace: API Marketplace tidak aktif...
  ...
[seed] selesai seeding tabel: notifications
[seed] selesai seeding Sprint 5 dengan sukses.
```

#### Mengatur Jumlah Record

```powershell
# Seed 100 request logs dan 30 pesan chat
go run ./cmd/seed -logs 100 -chats 30

# Seed hanya sedikit data untuk testing cepat
go run ./cmd/seed -logs 10 -chats 5 -verbose
```

#### Semua Flag

| Flag | Default | Keterangan |
|------|---------|-----------|
| `-verbose` | `false` | Tampilkan log detail setiap record |
| `-logs <N>` | `50` | Jumlah baris `request_logs` yang di-seed |
| `-chats <N>` | `20` | Jumlah baris `chat_messages` yang di-seed |

---

### 2c. Via Docker Compose

Jalankan seeder di dalam container Docker yang sudah berjalan:

```powershell
# Pastikan stack sudah berjalan
docker compose up --detach

# Jalankan seeder di container backend
docker compose exec backend go run ./cmd/seed -verbose

# Atau dengan jumlah record tertentu
docker compose exec backend go run ./cmd/seed -logs 100 -chats 30
```

> **Catatan:** Perintah ini membutuhkan Go toolchain di dalam image backend.
> Jika image production tidak menyertakan Go, build binary terlebih dahulu:
>
> ```powershell
> go build -o seed.exe ./cmd/seed
> # Copy binary ke container jika perlu
> ```

---

## 3. Konfigurasi Seeder

Seeder membaca konfigurasi database dari environment variables yang sama dengan server utama:

| Variable | Contoh | Keterangan |
|----------|--------|-----------|
| `DB_HOST` | `localhost` atau `mysql` | Host MySQL |
| `DB_PORT` | `3306` | Port MySQL |
| `DB_NAME` | `api_integrator` | Nama database |
| `DB_USER` | `gateway` | Username MySQL |
| `DB_PASSWORD` | `change-me` | Password MySQL |

Tidak ada environment variable tambahan yang diperlukan untuk seeder Sprint 5.

---

## 4. Data yang Di-seed

### `request_logs` (50 baris default)

Mensimulasikan request dari semua aplikasi gateway selama 14 hari terakhir.

| Kolom | Nilai Simulasi |
|-------|---------------|
| `source_app` | Bergilir: Marketplace, POS, SupplierHub, LogistiKita, SmartBank |
| `endpoint` | Bergilir: `/gateway/payment`, `/gateway/smartbank`, dll |
| `method` | `POST` |
| `status` | Mayoritas 200, ada 400 dan 500 untuk simulasi error |
| `duration_ms` | 50–200ms simulasi |
| `timestamp` | Tersebar 0–13 hari ke belakang |

### `notifications` (7 baris tetap)

Notifikasi contoh untuk semua tipe dan aplikasi:

| App | Tipe | Is Read |
|-----|------|---------|
| Marketplace | `api_inactive` | ❌ belum dibaca |
| POS | `error_rate` | ❌ belum dibaca |
| SmartBank | `response_time` | ❌ belum dibaca |
| SupplierHub | `response_time` | ✅ sudah dibaca |
| LogistiKita | `api_inactive` | ✅ sudah dibaca |
| API Gateway | `system` | ✅ sudah dibaca |
| UMKM Insight | `system` | ✅ sudah dibaca |

### `chat_messages` (20 baris default)

Percakapan antara `admin` dan user aplikasi di 4 conversation:

| Conversation ID | Pihak |
|----------------|-------|
| `conv-admin-marketplace` | admin ↔ marketplace |
| `conv-admin-pos` | admin ↔ pos |
| `conv-admin-supplierhub` | admin ↔ supplierhub |
| `conv-admin-insight` | admin ↔ insight |

### `dashboard_data` (5 baris tetap)

Cache analytics siap pakai untuk dashboard Sprint 6–8:

| Cache Key | Keterangan |
|-----------|-----------|
| `admin:traffic_summary` | Ringkasan total request, success rate, avg duration |
| `admin:service_indicators` | Status aktif/inaktif tiap aplikasi |
| `user:Marketplace:summary` | Ringkasan per-app untuk Marketplace |
| `user:POS:summary` | Ringkasan per-app untuk POS |
| `monitoring:summary` | Summary untuk UMKM Insight monitoring |

---

## 5. Idempoten & Keamanan

> [!IMPORTANT]
> Seeder dirancang **aman dijalankan berulang kali** (idempoten) untuk tabel yang mendukungnya.

| Tabel | Strategi Idempoten |
|-------|-------------------|
| `users` | `ON DUPLICATE KEY UPDATE` berdasarkan `(username, app_name)` |
| `dashboard_data` | `ON DUPLICATE KEY UPDATE` berdasarkan `cache_key` |
| `request_logs` | **Selalu INSERT baru** — setiap run menambah data baru |
| `notifications` | **Selalu INSERT baru** — setiap run menambah notifikasi baru |
| `chat_messages` | **Selalu INSERT baru** — setiap run menambah pesan baru |

> [!WARNING]
> Menjalankan seeder berulang kali pada `request_logs`, `notifications`, dan `chat_messages`
> akan menambah data (bukan replace). Ini disengaja untuk mensimulasikan pertumbuhan data.
> Jika ingin reset, gunakan:
> ```sql
> TRUNCATE TABLE request_logs;
> TRUNCATE TABLE notifications;
> TRUNCATE TABLE chat_messages;
> ```

> [!CAUTION]
> `SEED_USERS_ENABLED` tidak dapat diset `true` di environment `production`.
> Config server akan menolak startup jika terdeteksi. Seeder Sprint 5 CLI tidak memiliki
> pembatasan ini, namun **jangan jalankan di database production**.

---

## 6. Cara Menambah Data Seed Baru

### Menambah Request Log Seed Custom

Edit file [`backend/internal/database/seed/seeder.go`](../../backend/internal/database/seed/seeder.go),
bagian variabel `sourceApps`, `endpoints`, atau `httpStatuses`:

```go
var sourceApps = []string{
    "Marketplace", "POS", "SupplierHub", "LogistiKita", "SmartBank",
    "AppBaru", // ← tambahkan di sini
}
```

### Menambah Notifikasi Seed

Edit slice `notifSeeds` di file yang sama:

```go
var notifSeeds = []notifSeed{
    // ... entri yang sudah ada ...
    {
        AppName: "AppBaru",
        Type:    model.NotificationTypeSystem,
        Message: "AppBaru berhasil terhubung ke gateway.",
        IsRead:  false,
        DaysAgo: 0,
    },
}
```

### Menambah Percakapan Chat

Edit fungsi `buildChatSeeds` atau tambah entri baru di awal slice `seeds`:

```go
{"conv-admin-appbaru", "admin", "appbaru", "Selamat datang di gateway!", 10, true},
```

### Menambah Cache Dashboard

Edit slice `caches` di fungsi `seedDashboardData`:

```go
{
    key:     "user:AppBaru:summary",
    appName: "AppBaru",
    data: map[string]any{
        "total_requests":   0,
        "success_rate_pct": 100.0,
    },
},
```

---

## 7. Troubleshooting

### Error: `invalid configuration`

```
load config: invalid configuration: DB_HOST is required; DB_NAME is required; ...
```

**Solusi:** Pastikan file `.env` ada dan sudah dikonfigurasi. Jalankan:

```powershell
cp .env.example .env
# Edit nilai sesuai environment lokal
```

### Error: `ping MySQL`

```
open database: ping MySQL: dial tcp: connection refused
```

**Solusi:** Pastikan MySQL atau Docker Compose sudah berjalan:

```powershell
docker compose up --detach mysql
# Tunggu beberapa detik lalu coba lagi
```

### Error: `apply database migrations`

Jika migration gagal karena tabel sudah ada (di luar goose), jalankan:

```sql
-- Di MySQL client
INSERT INTO goose_db_version (version_id, is_applied) VALUES (2, 1), (3, 1), (4, 1), (5, 1)
ON DUPLICATE KEY UPDATE is_applied = 1;
```

Atau reset migration sepenuhnya (hati-hati di data yang ada):

```powershell
# Hanya untuk development
go run ./cmd/seed  # Akan retry migration otomatis
```

### Seeder Berjalan Tapi Data Tidak Muncul

Periksa apakah database yang aktif sudah benar:

```sql
SELECT DATABASE();  -- harus menampilkan 'api_integrator'
SELECT COUNT(*) FROM request_logs;
SELECT COUNT(*) FROM notifications;
SELECT COUNT(*) FROM chat_messages;
SELECT COUNT(*) FROM dashboard_data;
```

---

Dokumen ini dibuat untuk Sprint 5 API Integrator Gateway oleh Nuthfih.
