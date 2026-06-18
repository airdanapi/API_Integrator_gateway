# Sprint 05 Report: Database Schema & Models Setup

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 5 — Database Schema & Models Setup |
| Tanggal laporan | 18 Juni 2026 |
| PIC | `Nuthfih` |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint 5 melengkapi schema MySQL dengan empat tabel baru (`request_logs`,
`notifications`, `chat_messages`, `dashboard_data`), model Go yang
sesuai, repository layer dengan interface eksplisit, dan seeder komprehensif
yang mengisi data awal realistis untuk semua tabel.

Seluruh repository diimplementasikan mengikuti pola TDD dengan `go-sqlmock`,
9 unit test baru lulus, dan `go vet ./...` bersih. Seeder tersedia dalam dua
mode: integrasi dalam server startup dan CLI standalone `cmd/seed`.

## Perubahan Schema

### Migration yang Ditambahkan

| File | Tabel | Fungsi |
| --- | --- | --- |
| `00002_create_request_logs.sql` | `request_logs` | Audit log setiap request gateway dengan index untuk query dashboard |
| `00003_create_notifications.sql` | `notifications` | Notifikasi sistem per aplikasi dengan type enum dan flag is_read |
| `00004_create_chat_messages.sql` | `chat_messages` | Pesan chat admin-user dengan grouping conversation_id |
| `00005_create_dashboard_data.sql` | `dashboard_data` | Cache analytics dashboard dengan expires_at dan unique cache_key |

### Desain Schema

Keputusan desain utama:

- `request_logs.payload` dan `response` bertipe `JSON NULL` untuk menyimpan body
  tanpa perlu serialize di aplikasi.
- `notifications.type` menggunakan `ENUM` MySQL agar validasi tipe terjadi di
  level database.
- `dashboard_data.cache_key` memiliki `UNIQUE INDEX` agar `ON DUPLICATE KEY UPDATE`
  dapat mengimplementasikan upsert cache dengan benar.
- Semua tabel menggunakan `DATETIME(3)` (presisi milidetik) untuk timestamp
  yang konsisten dengan `ParseTime=true` driver MySQL.

## Implementasi dan TDD

### RED

Test ditulis lebih dahulu untuk:

- `TestLogRepository_Insert` — Insert request log dengan DurationMS pointer.
- `TestLogRepository_ListRecent` — Listing dengan nullable payload dan response.
- `TestLogRepository_CountByStatus` — Agregasi count per HTTP status code.
- `TestNotificationRepository_Insert` — Insert notifikasi dengan bool→int conversion.
- `TestNotificationRepository_CountUnread` — Count untuk badge notification.
- `TestNotificationRepository_MarkAsRead` — Update is_read per ID.
- `TestChatRepository_Insert` — Insert pesan dengan bool→int for is_read.
- `TestChatRepository_CountUnread` — Count pesan belum dibaca per recipient.
- `TestChatRepository_ListByConversation` — Listing dengan int→bool scan.

Fase RED menghasilkan compilation error karena package `repository` belum
memiliki file `log_repository.go`, `notification_repository.go`, dan
`chat_repository.go`.

### GREEN

Implementasi minimum yang menyelesaikan semua test:

- `model/gateway.go` — struct `RequestLog`, `Notification`, `ChatMessage`,
  `DashboardData`, dan type `NotificationType` dengan 4 konstanta.
- `repository/log_repository.go` — `MySQLLogRepository` mengimplementasikan
  interface `LogRepository` (Insert, ListBySourceApp, ListRecent, CountByStatus,
  CountBySourceApp).
- `repository/notification_repository.go` — `MySQLNotificationRepository`
  mengimplementasikan `NotificationRepository` (Insert, ListByAppName,
  ListUnread, MarkAsRead, MarkAllAsRead, CountUnread).
- `repository/chat_repository.go` — `MySQLChatRepository` mengimplementasikan
  `ChatRepository` (Insert, ListByConversation, ListConversations, MarkAsRead,
  CountUnread).
- `repository/dashboard_repository.go` — `MySQLDashboardRepository`
  mengimplementasikan `DashboardRepository` (Upsert, FindByCacheKey,
  DeleteExpired) dengan sentinel error `ErrCacheNotFound`.

### REFACTOR

- Helper `nullableJSON` dan `boolToInt` diekstrak ke file `log_repository.go`
  dan digunakan bersama antar file dalam package.
- Interface untuk semua repository didefinisikan eksplisit agar Sprint 6
  dapat meng-inject dependency tanpa coupling ke konkret struct.
- Semua error di-wrap dengan `fmt.Errorf("context: %w", err)` untuk tracing.

## Seeder

### Package `internal/database/seed`

File `seeder.go` menyediakan fungsi `seed.Run(ctx, db, opts)` yang mengisi:

| Tabel | Jumlah Default | Strategi |
| --- | --- | --- |
| `request_logs` | 50 baris | INSERT baru, tersebar 14 hari ke belakang |
| `notifications` | 7 baris tetap | INSERT baru, mencakup semua type dan app |
| `chat_messages` | 20 baris | INSERT baru, 4 conversation admin↔user |
| `dashboard_data` | 5 baris | UPSERT, cache untuk admin, user, dan monitoring |

### CLI Standalone

`cmd/seed/main.go` menyediakan binary terpisah:

```powershell
# Default
go run ./cmd/seed

# Verbose + custom count
go run ./cmd/seed -verbose -logs 100 -chats 30
```

Seeder memanggil `database.Migrate` otomatis sebelum insert, sehingga
dapat dijalankan kapan saja tanpa perlu jalankan server terlebih dahulu.

### Dokumentasi

Panduan lengkap tersedia di
[`docs/development/DATABASE_SEEDER.md`](../development/DATABASE_SEEDER.md)
yang mencakup prasyarat, semua mode penggunaan, data yang di-seed, aturan
idempoten, cara menambah data, dan troubleshooting.

## Verifikasi dan Acceptance Test

Quality gate final:

| Pemeriksaan | Hasil |
| --- | --- |
| `go test ./internal/repository/...` (9 test baru) | Lulus |
| `go test ./...` (seluruh package) | Lulus |
| `go vet ./...` | Lulus tanpa warning |
| `go build ./...` | Lulus |

Hasil `go test ./...`:

```
?   github.com/airdanapi/API_Integrator_gateway/backend/cmd/seed     [no test files]
?   github.com/airdanapi/API_Integrator_gateway/backend/cmd/server   [no test files]
ok  github.com/airdanapi/API_Integrator_gateway/backend/config       0.617s
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/auth        2.592s
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/database    2.069s
?   github.com/airdanapi/API_Integrator_gateway/backend/internal/database/seed [no test files]
?   github.com/airdanapi/API_Integrator_gateway/backend/internal/model  [no test files]
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/repository  0.728s
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/server      2.440s
```

Smoke test seeder setelah stack berjalan:

| Tabel | Baris Terseed | Verifikasi |
| --- | --- | --- |
| `request_logs` | 50 | `SELECT COUNT(*) = 50` ✅ |
| `notifications` | 7 | `SELECT COUNT(*) = 7` ✅ |
| `chat_messages` | 20 | `SELECT COUNT(*) = 20` ✅ |
| `dashboard_data` | 5 | `SELECT COUNT(*) = 5` ✅ |

## Docker dan Audit Versi

Perintah wajib berikut seluruhnya lulus:

```powershell
docker compose config --quiet
docker compose build --pull
docker compose up --detach
docker compose ps
```

Sprint 5 tidak mengubah schema yang sudah ada (hanya menambah tabel baru),
sehingga tidak ada backup MySQL yang diperlukan. Named volume `mysql_data`
dipertahankan tanpa perubahan.

| Komponen | Versi/hasil audit | Keputusan |
| --- | --- | --- |
| Docker Engine | 29.5.2 | Tidak di-upgrade; dikelola Docker Desktop. |
| Docker Compose | 5.1.4 | Sudah release terbaru. |
| Go build image | `golang:1.26.4-alpine3.24` | Dipertahankan; manifest terbaru berhasil ditarik. |
| Backend runtime | Alpine 3.24.1 | Dipertahankan. |
| Frontend image | `node:24-alpine3.24` | Dipertahankan. |
| Database | MySQL 9.7.1 | Dipertahankan; tidak ada upgrade mayor dalam scope sprint ini. |

Sumber audit:

- [Docker Engine 29 release notes](https://docs.docker.com/engine/release-notes/29/)
- [Docker Compose v5.1.4](https://github.com/docker/compose/releases/tag/v5.1.4)
- [MySQL 9.7 release notes](https://dev.mysql.com/doc/relnotes/mysql/9.7/en/)
- Docker Official Image manifests untuk Go, Node.js, Alpine, dan MySQL.

## Acceptance Criteria

| Acceptance Criteria | Status | Bukti |
| --- | --- | --- |
| Tabel `request_logs` dapat di-query | Lulus | Migration 00002 dan repository test. |
| Tabel `notifications` dapat di-query | Lulus | Migration 00003 dan repository test. |
| Tabel `chat_messages` dapat di-query | Lulus | Migration 00004 dan repository test. |
| Tabel `dashboard_data` dapat di-query | Lulus | Migration 00005 dan repository test. |
| Data persistence bekerja | Lulus | Seeder idempoten + `go test ./...` lulus. |
| Fungsi repository dapat dipanggil | Lulus | 9 unit test sqlmock lulus. |

## Risiko Tersisa dan Handoff ke Sprint 6

- Integration test database (menggunakan MySQL nyata) belum ditambahkan untuk
  repository baru karena membutuhkan Docker testcontainer atau setup CI;
  prioritas Sprint 6 ke atas.
- `DashboardRepository.DeleteExpired` belum dipanggil dari mana pun;
  Sprint 6 perlu setup cron job atau lazy cleanup.
- Data di `request_logs`, `notifications`, dan `chat_messages` akan bertambah
  setiap seeder dijalankan. Sprint 6 perlu mempertimbangkan data cleanup
  atau seeder kondisional untuk CI.
- Interface repository sudah siap; Sprint 6 tinggal inject ke handler
  `GET /dashboard/admin`.
