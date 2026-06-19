# Sprint 08 Report: Dashboard User & Monitoring

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 8 — Dashboard App User & Monitoring |
| Tanggal laporan | 19 Juni 2026 |
| PIC | `Nuthfih` |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint 8 mengimplementasikan dua dashboard role baru yang melengkapi ekosistem
dashboard yang sudah dibangun di Sprint 6–7:

- **Dashboard User** (`/dashboard/user`, role `app_user`) — menampilkan traffic
  summary, status layanan, dan riwayat request yang difilter per aplikasi milik
  user yang login.
- **Dashboard Monitoring** (`/dashboard/monitoring`, role `monitoring_user`) —
  tampilan read-only lintas semua aplikasi: traffic global, status per layanan,
  dan breakdown request per aplikasi dalam bentuk bar chart.

Implementasi mencakup layer backend (repository, service, handler, routing),
frontend (service functions, page components, routing update), unit testing TDD,
dan verifikasi Docker build + smoke test end-to-end.

## Fitur yang Diimplementasikan

### Backend

#### Repository (`repository/log_repository.go`)

- **`CountByStatusForApp(ctx, appName, since)`** — query baru yang menghitung
  jumlah log per HTTP status code untuk satu aplikasi saja:
  ```sql
  SELECT status, COUNT(*) FROM request_logs
  WHERE source_app = ? AND timestamp >= ? GROUP BY status
  ```

#### Service (`dashboard/service.go`)

Dua method baru ditambahkan ke `DashboardService`:

**`GetUserDashboard(ctx, appName, page, limit)`**
- Memanggil `CountByStatusForApp` untuk mendapatkan traffic summary yang
  difilter per aplikasi
- Menghitung `service_status`: `active` jika ada request dalam 7 hari terakhir,
  `inactive` jika tidak
- Memanggil `ListBySourceApp` untuk mendapatkan riwayat log terpaginasi
- Return struct `UserDashboard`:
  ```go
  type UserDashboard struct {
    MyApp         string
    TrafficSummary TrafficSummary
    ServiceStatus  string  // "active" | "inactive"
    RecentLogs     []AuditLogEntry
    TotalLogs      int64
    Page, Limit    int
  }
  ```

**`GetMonitoringDashboard(ctx)`**
- Memanggil `CountByStatus` untuk traffic summary global (semua aplikasi)
- Memanggil `CountBySourceApp` untuk mendapatkan total request per aplikasi
- Membangun `ServiceIndicators` (5 app) dan `AppBreakdown` (5 app)
- Return struct `MonitoringDashboard`:
  ```go
  type MonitoringDashboard struct {
    TrafficSummary    TrafficSummary
    ServiceIndicators []ServiceIndicator
    AppBreakdown      []AppStat
  }
  ```

#### Handler & Routing

**`userDashboardHandler`** (`server/dashboard.go`):
- Membaca `app_name` dari JWT claims (`auth.Claims`)
- Pagination param `page` dan `limit` dari query string
- Route: `GET /dashboard/user` → `requireRole(app_user)`

**`monitoringDashboardHandler`** (`server/dashboard.go`):
- Tidak memerlukan parameter tambahan, men-delegate ke service
- Route: `GET /dashboard/monitoring` → `requireRole(monitoring_user)`

### Frontend

#### Service Functions (`services/dashboard.js`)

```js
// Difilter berdasarkan app_name dari JWT (backend)
export async function fetchUserDashboard(apiClient, params = {}) {
  const { page = 1, limit = 20 } = params
  const response = await apiClient.get('/dashboard/user', { params: { page, limit } })
  return response.data.data
}

// Global, tidak memerlukan parameter
export async function fetchMonitoringDashboard(apiClient) {
  const response = await apiClient.get('/dashboard/monitoring')
  return response.data.data
}
```

#### UserDashboardPage (`pages/UserDashboardPage.jsx`)

- **4 Traffic Stats Cards** — Total Request, Sukses, Error, Success Rate (per app)
- **ServiceStatusBadge** — badge "Aktif" (hijau) / "Tidak Aktif" (merah)
- **Riwayat Request Table** — endpoint, method, status badge, timestamp
- **Pagination** — Sebelumnya/Berikutnya berdasarkan `total_logs`
- **Auto-polling** — `setInterval(load, 30_000)` dengan cleanup
- **Loading/Error state** — `role="status"` dan `role="alert"`
- **Refresh manual** — tombol dengan animasi spinner saat refresh aktif

#### MonitoringDashboardPage (`pages/MonitoringDashboardPage.jsx`)

- **4 Summary Cards** — traffic global semua aplikasi
- **ServiceIndicatorGrid** — daftar status 5 aplikasi dengan dot berwarna
- **AppBreakdownTable** — bar chart CSS horizontal per aplikasi (lebar proporsional
  terhadap `max(total_requests)`)
- Read-only, tidak ada fungsi modifikasi
- **Auto-polling** dan **Refresh manual** sama seperti UserDashboardPage

#### Routing Update (`App.jsx`)

```jsx
// Sebelumnya: placeholder DashboardPage
<Route path="/dashboard/user"       element={<UserDashboardPage />} />
<Route path="/dashboard/monitoring" element={<MonitoringDashboardPage />} />
```

## TDD Phases

### RED → GREEN

| Test File | Jumlah Test | Status |
| --- | --- | --- |
| `repository/log_repository_test.go` | +1 (`CountByStatusForApp`) | ✅ PASS |
| `dashboard/service_test.go` | +3 (`GetUserDashboard` × 2, `GetMonitoringDashboard` × 1) | ✅ PASS |
| `server/dashboard_test.go` | +6 (user: 401, 403, 200; monitoring: 401, 403, 200) | ✅ PASS |
| `services/dashboard.test.js` | +5 (fetchUserDashboard × 3, fetchMonitoringDashboard × 2) | ✅ PASS |
| `pages/UserDashboardPage.test.jsx` | 7 | ✅ PASS |
| `pages/MonitoringDashboardPage.test.jsx` | 7 | ✅ PASS |

### Pola Mock (Frontend)

Test menggunakan `vi.mock('../auth/auth-context', ...)` dan `renderPage(fetchData)`
(injectable prop), konsisten dengan pola `AdminDashboardPage.test.jsx` Sprint 7.

## Hasil Test

### Backend (`go test ./...`)

```
ok  github.com/airdanapi/API_Integrator_gateway/backend/config         (cached)
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/auth  (cached)
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/dashboard  0.803s
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/database  (cached)
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/repository  0.849s
ok  github.com/airdanapi/API_Integrator_gateway/backend/internal/server  1.915s
```

### Frontend (`npx vitest run`)

```
Test Files  8 passed (8)
Tests       48 passed (48)
```

File baru Sprint 8:
```
✓ src/services/dashboard.test.js         (10 test total, +5 baru)
✓ src/pages/UserDashboardPage.test.jsx   (7 test)
✓ src/pages/MonitoringDashboardPage.test.jsx (7 test)
✓ src/App.test.jsx                       (13 test, heading diperbarui)
```

## Smoke Test Endpoint

Semua 10 check berikut **LULUS**:

| # | Endpoint / Skenario | Expected | Hasil |
| --- | --- | --- | --- |
| 1 | `GET /health` | `status=success` | ✅ |
| 2 | Login admin | token + `role=admin_gateway` | ✅ |
| 3 | `GET /dashboard/admin` (admin token) | 200, `traffic_summary`, `service_indicators[5]` | ✅ `total_requests=28` |
| 4 | Login Marketplace | token + `role=app_user` | ✅ |
| 5 | `GET /dashboard/user` (marketplace token) | 200, `my_app=Marketplace`, `service_status=active` | ✅ `total_requests=6` |
| 6 | Login UMKM Insight | token + `role=monitoring_user` | ✅ |
| 7 | `GET /dashboard/monitoring` (monitor token) | 200, `service_indicators[5]`, `app_breakdown[5]` | ✅ `total_requests=28` |
| 8 | `app_user` akses `/dashboard/admin` | 403 Forbidden | ✅ |
| 9 | `monitoring_user` akses `/dashboard/user` | 403 Forbidden | ✅ |
| 10 | Tanpa token akses `/dashboard/monitoring` | 401 Unauthorized | ✅ |

## Docker dan Audit Versi

| Pemeriksaan | Hasil |
| --- | --- |
| `docker compose config --quiet` | ✅ Lulus |
| `docker compose build --pull` | ✅ Image Built (backend + frontend) |
| `docker compose up --detach` | ✅ Semua container Up |
| `docker compose ps` | ✅ MySQL Healthy, Backend Healthy, Frontend Up |

| Komponen | Versi | Keputusan |
| --- | --- | --- |
| Docker Engine | 29.5.2 | Tidak di-upgrade (dikelola Docker Desktop) |
| Docker Compose | 5.1.4 | Sudah terbaru |
| Frontend image | `node:24-alpine3.24` | Dipertahankan |
| Backend image | `golang:1.26.4-alpine3.24` | Dipertahankan |
| Database | MySQL 9.7 | Tidak ada perubahan schema |
| vitest | 4.1.9 | Versi project package.json |

## File yang Dibuat / Diubah

| File | Status | Keterangan |
| --- | --- | --- |
| `backend/internal/repository/log_repository.go` | DIUBAH | Tambah `CountByStatusForApp` |
| `backend/internal/repository/log_repository_test.go` | DIUBAH | Test `CountByStatusForApp` |
| `backend/internal/dashboard/service.go` | DIUBAH | Tambah `UserDashboard`, `MonitoringDashboard`, 2 method |
| `backend/internal/dashboard/service_test.go` | DIUBAH | Stub extended + 3 test baru |
| `backend/internal/server/dashboard.go` | DIUBAH | Handler user + monitoring, interface extended |
| `backend/internal/server/dashboard_test.go` | DIUBAH | 6 test baru (user + monitoring) |
| `backend/internal/server/app.go` | DIUBAH | 2 route baru |
| `frontend/src/services/dashboard.js` | DIUBAH | Tambah `fetchUserDashboard`, `fetchMonitoringDashboard` |
| `frontend/src/services/dashboard.test.js` | DIUBAH | +5 test baru |
| `frontend/src/pages/UserDashboardPage.jsx` | BARU | Dashboard app_user |
| `frontend/src/pages/UserDashboardPage.test.jsx` | BARU | 7 test |
| `frontend/src/pages/MonitoringDashboardPage.jsx` | BARU | Dashboard monitoring_user |
| `frontend/src/pages/MonitoringDashboardPage.test.jsx` | BARU | 7 test |
| `frontend/src/App.jsx` | DIUBAH | Route user + monitoring → page baru |
| `frontend/src/App.test.jsx` | DIUBAH | Mock extended, heading assertion diperbarui |

## Risiko Tersisa dan Handoff ke Sprint 9

- **`avg_duration_ms` belum ada di UserDashboard/MonitoringDashboard** —
  perlu query `AVG(duration_ms)` tambahan di repository.
- **Sorting kolom tabel** belum ada di UserDashboard — belum diimplementasikan.
- **WebSocket/SSE** masih polling 30s — cukup untuk kebutuhan saat ini.
- **`ListBySourceApp` menggunakan offset pagination** — pertimbangkan cursor-based
  untuk dataset besar di Sprint selanjutnya.
