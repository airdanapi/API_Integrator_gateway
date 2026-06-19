# Sprint 06 Report: Dashboard Admin — Backend

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 6 — Dashboard Admin (Backend) |
| Tanggal laporan | 19 Juni 2026 |
| PIC | `Nuthfih` |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint 6 mengimplementasikan endpoint backend `GET /dashboard/admin` yang
menyediakan data analitik gateway untuk dashboard admin. Endpoint dilindungi
double middleware: autentikasi JWT dan role-based access control (hanya
`admin_gateway`). Response mencakup tiga blok data: traffic summary, service
indicator aktif/inaktif, dan audit log berpage.

Layer baru yang ditambahkan:

- **`internal/dashboard/service.go`** — business logic untuk traffic summary,
  deteksi inactive API (> 7 hari tanpa request), dan paginated audit log.
- **`internal/server/dashboard.go`** — handler HTTP + `requireRole` middleware.
- **`cmd/server/main.go`** — wiring `DashboardService` ke server dependencies.

## API Contract

### `GET /dashboard/admin`

**Headers required:** `Authorization: Bearer <jwt_token>` (role `admin_gateway`)

**Query params:**
| Param | Default | Keterangan |
|-------|---------|-----------|
| `page` | `1` | Halaman audit log |
| `limit` | `20` | Jumlah item per halaman (max 100) |

**Response 200:**
```json
{
  "status": "success",
  "data": {
    "traffic_summary": {
      "total_requests": 50,
      "success_count": 42,
      "error_count": 8,
      "success_rate_pct": 84.0,
      "avg_duration_ms": 0
    },
    "service_indicators": [
      { "app_name": "Marketplace", "status": "inactive", "last_request": "..." },
      { "app_name": "POS", "status": "active", "last_request": "..." }
    ],
    "audit_logs": {
      "items": [ { "id": 1, "source_app": "POS", "endpoint": "...", "status": 200, ... } ],
      "total": 50,
      "page": 1,
      "limit": 20
    }
  }
}
```

**Response 401:** Token tidak ada atau tidak valid.

**Response 403:** Token valid tapi role bukan `admin_gateway`.

## Implementasi dan TDD

### RED

Test ditulis lebih dahulu untuk:

- `TestDashboardAdminRequiresAuthentication` — tanpa token → 401.
- `TestDashboardAdminForbidsNonAdminRole` — role `app_user` → 403.
- `TestDashboardAdminForbidsMonitoringRole` — role `monitoring_user` → 403.
- `TestDashboardAdminReturnsContract` — admin dapat data lengkap → 200 + contract JSON.
- `TestDashboardAdminDefaultsPaginationParams` — tanpa `page`/`limit` → default 1/20.
- `TestGetTrafficSummary_CalculatesRateCorrectly` — 80+10+10 = 100 request, 80% rate.
- `TestGetTrafficSummary_EmptyReturnsZeroRate` — tanpa data → rate 0%.
- `TestGetServiceIndicators_ActiveAndInactive` — Marketplace/POS aktif, 3 lainnya inactive.
- `TestGetServiceIndicators_AllInactive` — semua inactive jika tidak ada request dalam 7 hari.
- `TestGetAuditLogs_PaginatesCorrectly` — limit/offset, total, DurationMS nil handling.

Fase RED menghasilkan compilation error karena package `dashboard` dan handler
`dashboard.go` belum ada, serta `requireRole` middleware belum diimplementasikan.

### GREEN

Implementasi minimum:

- `dashboard.Service` dengan `LogQuerier` interface untuk testability.
- `GetTrafficSummary` — aggregate `CountByStatus` → total/success/error/rate.
- `GetServiceIndicators` — `CountBySourceApp` 7 hari terakhir → active jika count > 0.
- `GetAuditLogs` — `ListRecent` dengan pagination + `CountByStatus` untuk total.
- `requireRole` middleware — cek `auth_claims` locals → 403 jika role tidak cocok.
- `adminDashboardHandler` — compose ketiga service call → satu response JSON.
- `parseQueryInt` helper — safe parsing query string ke int dengan fallback.

### REFACTOR

- `DashboardService` interface didefinisikan di package `server` agar handler
  tidak bergantung langsung pada konkret struct `dashboard.Service`.
- `parseQueryInt` diekstrak sebagai pure function untuk reuse di handler lain.
- Import `dashboard` di `main.go` untuk wiring tanpa perlu mengubah interface
  server (open/closed principle).

## Verifikasi dan Acceptance Test

Quality gate final:

| Pemeriksaan | Hasil |
| --- | --- |
| `go test ./...` (10 test baru + semua lama) | Lulus |
| `go vet ./...` | Lulus tanpa warning |
| `go build ./...` | Lulus |

Ringkasan hasil `go test ./...`:

```
ok  github.com/airdanapi/API_Integrator_gateway/backend/config
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/auth
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/dashboard    (5 test baru)
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/database
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/repository
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/server       (5 test baru)
```

## Smoke Test Endpoint

Setelah Docker build dan `docker compose up --detach`:

```powershell
# Login untuk mendapatkan token admin
$login = Invoke-WebRequest -Uri "http://localhost:8080/auth/login" `
  -Method POST -ContentType "application/json" `
  -Body '{"username":"admin","password":"admin-development-password","app_name":"API Gateway"}' `
  -UseBasicParsing | ConvertFrom-Json
$token = $login.data.token

# Hit dashboard admin
Invoke-WebRequest -Uri "http://localhost:8080/dashboard/admin" `
  -Headers @{ Authorization = "Bearer $token" } `
  -UseBasicParsing | ConvertFrom-Json
```

| Test | Ekspektasi | Hasil |
| --- | --- | --- |
| GET /dashboard/admin tanpa token | 401 | ✅ |
| GET /dashboard/admin dengan token app_user | 403 | ✅ |
| GET /dashboard/admin dengan token admin_gateway | 200 + data JSON | ✅ |

## Docker dan Audit Versi

Perintah wajib berikut seluruhnya lulus:

```powershell
docker compose config --quiet
docker compose build --pull
docker compose up --detach
docker compose ps
```

Sprint 6 menambahkan package `internal/dashboard` di backend. Image di-build
ulang dengan `--pull` untuk memastikan base image terbaru. Tidak ada perubahan
schema database sehingga tidak diperlukan backup volume.

| Komponen | Versi/hasil audit | Keputusan |
| --- | --- | --- |
| Docker Engine | 29.5.2 | Tidak di-upgrade; dikelola Docker Desktop. |
| Docker Compose | 5.1.4 | Sudah release terbaru. |
| Go build image | `golang:1.26.4-alpine3.24` | Dipertahankan; manifest terbaru ditarik. |
| Backend runtime | Alpine 3.24.1 | Dipertahankan. |
| Frontend image | `node:24-alpine3.24` | Dipertahankan. |
| Database | MySQL 9.7.1 | Tidak ada perubahan schema. |

## Acceptance Criteria

| Acceptance Criteria | Status | Bukti |
| --- | --- | --- |
| `GET /dashboard/admin` memerlukan token valid | Lulus | Test 401 tanpa token. |
| Hanya `admin_gateway` dapat akses | Lulus | Test 403 untuk app_user dan monitoring_user. |
| Response berisi traffic summary | Lulus | Contract test field `traffic_summary`. |
| Response berisi service indicator aktif/inaktif | Lulus | Contract test 5 indicator. |
| Response berisi paginated audit logs | Lulus | Contract test total, page, limit, items. |
| Pagination default page=1, limit=20 | Lulus | Default pagination test. |
| Inactive API terdeteksi jika > 7 hari tanpa request | Lulus | Service unit test inactive detection. |

## Risiko Tersisa dan Handoff ke Sprint 7

- `avg_duration_ms` di traffic summary selalu 0 — query rata-rata duration belum
  diimplementasikan (ditunda ke Sprint 8 setelah semua data cukup).
- `last_request` di service indicator adalah estimasi, bukan waktu request
  terakhir yang akurat dari database. Query `MAX(timestamp)` per app bisa
  ditambahkan di Sprint 8.
- Route `/dashboard/user` dan `/dashboard/monitoring` belum ada — Sprint 7 dan
  Sprint 8 akan mengerjakan ini.
- Data yang ditampilkan di `/dashboard/admin` masih menggunakan data seeder;
  real-time logging dari request aktual baru akan ada setelah Sprint 9.
