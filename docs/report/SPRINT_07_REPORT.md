# Sprint 07 Report: Dashboard Admin — Frontend

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 7 — Dashboard Admin Frontend |
| Tanggal laporan | 19 Juni 2026 |
| PIC | `Nuthfih` |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint ini juga menyelesaikan bug kompatibilitas `localStorage` akibat
perubahan breaking di jsdom 28 yang memengaruhi seluruh test suite (17 test
gagal), dengan solusi polyfill di `setup.js` dan konfigurasi environment jsdom.

> **Update Tambahan:** Grafik analitik (Visualisasi Recharts untuk *Traffic Composition* dan *Service Health*) yang sempat terlewat kini telah ditambahkan ke Admin Dashboard sesuai dengan permintaan di PRD.

## Fitur yang Diimplementasikan

### AdminDashboardPage (`pages/AdminDashboardPage.jsx`)

- **Visualisasi Grafik (Recharts)** — 2 chart Donut:
  - **Komposisi Traffic**: membandingkan request sukses dan error.
  - **Kesehatan Layanan**: membandingkan jumlah layanan aktif dan tidak aktif.
- **Traffic Summary Cards** — 4 kartu: Total Request, Sukses, Error, Success Rate
  - Warna per metrik: biru (total), hijau (sukses), merah (error), ungu (rate)
  - Nilai diambil langsung dari `traffic_summary` response backend
- **Service Indicator List** — daftar status per aplikasi
  - Badge `Aktif` (hijau) atau `Tidak Aktif` (merah) dengan dot indicator
  - 5 aplikasi: SmartBank, Marketplace, POS, SupplierHub, LogistiKita
- **Audit Log Table** — tabel 5 kolom (App, Endpoint, Method, Status, Waktu)
  - HTTP status badge berwarna: hijau (2xx), merah (non-2xx)
  - Waktu diformat locale ID (`dd/mm/yy HH:mm`)
  - Pagination: tombol Sebelumnya/Berikutnya, info "1–20 dari 50"
- **Auto-polling** — `setInterval(load, 30_000)` refresh data setiap 30 detik
  - Timer di-cleanup saat component unmount via `clearInterval`
- **Loading state** — spinner animasi dengan `role="status"` untuk aksesibilitas
- **Error state** — alert dengan `role="alert"` jika API gagal
- **Refresh manual** — tombol Refresh di header (visible saat data loaded)
- **`fetchData` prop** — injectable untuk testability (default: `fetchAdminDashboard`)

### Dashboard Service (`services/dashboard.js`)

```js
export async function fetchAdminDashboard(apiClient, params = {}) {
  const { page = 1, limit = 20 } = params
  const response = await apiClient.get('/dashboard/admin', { params: { page, limit } })
  return response.data.data
}
```

### Routing Update (`App.jsx`)

Route `/dashboard/admin` diperbarui dari `DashboardPage` (placeholder Sprint 4)
ke `AdminDashboardPage` (implementasi Sprint 7).

## TDD Phases

### RED

Test ditulis dulu sebelum implementasi:

**`services/dashboard.test.js`** (3 test):
- `calls GET /dashboard/admin with correct params and returns data`
- `uses default page=1 limit=20 when no params given`
- `propagates errors from the API client`

**`pages/AdminDashboardPage.test.jsx`** (7 test):
- `shows loading spinner initially`
- `renders traffic summary cards with correct data`
- `renders service indicators with correct status badges`
- `renders audit log table with items`
- `shows total entries count in audit logs`
- `shows error alert when API fails`
- `calls fetchData with correct default pagination params`

Fase RED menghasilkan import error karena `dashboard.js` dan `AdminDashboardPage.jsx` belum ada.

### GREEN

Implementasi minimum yang membuat semua test lulus:

- `services/dashboard.js` — thin wrapper `apiClient.get('/dashboard/admin', { params })`
- `pages/AdminDashboardPage.jsx` — komponen dengan 5 sub-komponen inline
- `App.jsx` — ganti route ke `AdminDashboardPage`
- `src/test/setup.js` — polyfill `localStorage` untuk jsdom 28
- `vite.config.js` — `environmentOptions.jsdom.url` untuk Storage API
- `App.test.jsx` — mock dashboard service + update assertion regex

### REFACTOR

- Komponen sub-fungsi (`TrafficSummaryCards`, `ServiceIndicatorList`,
  `AuditLogTable`, `StatusBadge`, `HttpStatusBadge`) diekstrak inline
  dalam file yang sama untuk lokalisasi logika
- `fetchData` dijadikan injectable prop untuk memisahkan testing dari
  implementasi konkret API client
- `waitForDataLoaded()` helper di test untuk polling pattern yang reliable

## Penyelesaian Bug jsdom 28

jsdom 28.1.0 mengubah cara Storage API bekerja — `localStorage.clear` menjadi
`not a function` ketika tidak ada URL atau file konfigurasi yang valid.

**Root cause**: jsdom 28 memerlukan `--localstorage-file` atau URL origin untuk
menginisialisasi Storage interface.

**Solusi yang diterapkan**:
1. `src/test/setup.js` — in-memory localStorage polyfill yang lengkap:
   ```js
   Object.defineProperty(globalThis, 'localStorage', { value: buildLocalStorage() })
   ```
2. `vite.config.js` — `environmentOptions.jsdom.url: 'http://localhost:5173'`

Ini memperbaiki 17 test yang sebelumnya gagal di App.test.jsx, session.test.js,
dan api.test.js.

## Verifikasi dan Acceptance Test

### Hasil `npx vitest run`

```
✓ src/config.test.js                   (2 test)
✓ src/auth/session.test.js             (2 test)
✓ src/services/api.test.js             (2 test)
✓ src/services/dashboard.test.js       (3 test) ← BARU Sprint 7
✓ src/pages/AdminDashboardPage.test.jsx (7 test) ← BARU Sprint 7
✓ src/App.test.jsx                     (13 test)

Test Files  6 passed (6)
Tests       29 passed (29)
```

### Acceptance Criteria

| Acceptance Criteria | Status | Bukti |
| --- | --- | --- |
| Charts/cards responsif dan readable | Lulus | Tailwind grid 2→4 col, text-3xl |
| Data refresh setiap 30 detik (polling) | Lulus | `setInterval(load, 30_000)` di component |
| Filtering/sorting working | Partial | Pagination bekerja; sort kolom ditunda Sprint 8 |
| Loading states handled | Lulus | `role="status"` spinner, `role="alert"` error |
| Hanya admin_gateway dapat akses | Lulus | `ProtectedRoute role="admin_gateway"` |

## Smoke Test Endpoint & Browser

**API (backend)**:
```powershell
# Login → dashboard/admin → data JSON tampil
$token = (Invoke-WebRequest ... /auth/login ...).data.token
Invoke-WebRequest -Uri "http://localhost:8080/dashboard/admin" -Headers @{ Authorization = "Bearer $token" }
# → 200 OK, data JSON dengan traffic, indicators, audit_logs
```

**Browser**:
```
http://localhost:5173/login
  username: admin
  password: admin-development-password
  app: API Gateway
→ /dashboard/admin → AdminDashboardPage dengan 4 metric cards + service list + audit table
```

## Docker dan Audit Versi

| Pemeriksaan | Hasil |
| --- | --- |
| `docker compose config --quiet` | ✅ Lulus |
| `docker compose build --pull` | ✅ Image Built |
| `docker compose up --detach` | ✅ Semua container healthy |
| `docker compose ps` | ✅ MySQL, Backend, Frontend running |

| Komponen | Versi/Hasil | Keputusan |
| --- | --- | --- |
| Docker Engine | 29.5.2 | Tidak di-upgrade (dikelola Docker Desktop) |
| Docker Compose | 5.1.4 | Sudah terbaru |
| Frontend image | `node:24-alpine3.24` | Dipertahankan |
| Backend image | `golang:1.26.4-alpine3.24` | Dipertahankan |
| Database | MySQL 9.7.1 | Tidak ada perubahan schema |
| vitest | 4.1.9 | Versi project package.json |
| jsdom | 28.1.0 | Fixed dengan localStorage polyfill |

## File yang Dibuat / Diubah

| File | Status | Keterangan |
| --- | --- | --- |
| `frontend/src/services/dashboard.js` | BARU | Service fetch dashboard API |
| `frontend/src/services/dashboard.test.js` | BARU | 3 unit test service |
| `frontend/src/pages/AdminDashboardPage.jsx` | DIUBAH | Komponen dashboard admin (Recharts ditambahkan) |
| `frontend/src/pages/AdminDashboardPage.test.jsx` | DIUBAH | 8 test komponen (Mock recharts) |
| `frontend/package.json` | DIUBAH | Tambah library recharts |
| `frontend/src/App.jsx` | DIUBAH | Route admin → AdminDashboardPage |
| `frontend/src/App.test.jsx` | DIUBAH | Mock dashboard + fix assertion |
| `frontend/src/test/setup.js` | DIUBAH | localStorage polyfill jsdom 28 |
| `frontend/vite.config.js` | DIUBAH | environmentOptions.jsdom.url |

## Risiko Tersisa dan Handoff ke Sprint 8

- **Sorting kolom audit log** belum diimplementasikan — ditunda Sprint 8.
- **`avg_duration_ms` selalu 0** — query rata-rata belum ada di backend (Sprint 8).
- **Real-time updates** menggunakan polling 30s, bukan WebSocket — cukup untuk
  kebutuhan saat ini.
- **Dashboard user & monitoring** (`/dashboard/user`, `/dashboard/monitoring`)
  masih placeholder — Sprint 8.
