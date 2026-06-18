# Sprint 01 Report: Project Setup & Infrastructure

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 1 — Project Setup & Infrastructure |
| Tanggal laporan | 18 Juni 2026 |
| PIC | `venalism` |
| Commit implementasi | `0d51525` |
| Status | Selesai dengan tindak lanjut eksternal |

## Ringkasan Eksekutif

Sprint 1 menghasilkan fondasi teknis untuk pengembangan API Integrator Gateway:
frontend React, backend Go Fiber, konfigurasi environment MySQL, test otomatis,
Docker Compose full-stack, dokumentasi setup lokal, dan workflow CI.

Implementasi dikerjakan dengan Test-Driven Development. Test kontrak dibuat
lebih dahulu, dijalankan dalam kondisi gagal karena implementasi belum tersedia,
kemudian dilanjutkan dengan implementasi minimum dan refactor sampai seluruh
test, lint, vet, dan build lulus.

Pekerjaan eksternal untuk membuat 12 GitHub Issues dan organization project
`RPL 2` belum dapat diselesaikan karena token GitHub MCP tidak memiliki akses
tulis dan browser GitHub tidak tersedia. Backlog lengkap telah disiapkan pada
[dokumen GitHub Project](../project-management/GITHUB_PROJECT_ISSUES.md).

## Tujuan dan Ruang Lingkup

Tujuan Sprint 1 adalah menyiapkan repository dan infrastruktur dasar agar
pengembangan fitur pada sprint berikutnya dapat dimulai dari baseline yang
konsisten dan dapat diuji.

Ruang lingkup yang diselesaikan:

- Struktur proyek `frontend/` dan `backend/`.
- React 19 dengan Vite 8 dan Tailwind CSS 4.
- Go 1.26 dengan Fiber 3.
- Konfigurasi environment untuk aplikasi dan MySQL.
- Endpoint publik `GET /health`.
- Middleware recovery, CORS, dan request logging.
- Graceful shutdown backend.
- Unit/component test frontend dan backend.
- Docker Compose untuk frontend, backend, dan MySQL.
- GitHub Actions untuk test, lint, vet, build, dan Docker validation.
- Dokumentasi setup lokal dan Docker.

Di luar ruang lingkup Sprint 1:

- Koneksi runtime backend ke MySQL.
- Database migration dan schema.
- Landing page final Sprint 2.
- Authentication, dashboard, gateway routing, dan fitur bisnis.

## Hasil Implementasi

### Frontend

Frontend menyediakan application shell minimum sebagai bukti bahwa fondasi UI
siap digunakan tanpa mengambil scope landing page Sprint 2.

| Komponen | Implementasi |
| --- | --- |
| Framework | React `19.2.7` |
| Build tool | Vite `8.0.16` |
| Styling | Tailwind CSS `4.3.1` melalui plugin Vite |
| Testing | Vitest dan Testing Library |
| Quality | ESLint dengan plugin React dan React Hooks |
| Default URL | `http://localhost:5173` |

Konfigurasi API membaca `VITE_API_BASE_URL` dan menggunakan
`http://localhost:8080` sebagai fallback.

### Backend

Backend menggunakan app factory agar Fiber app dapat diuji tanpa membuka
network listener.

| Komponen | Implementasi |
| --- | --- |
| Bahasa | Go `1.26` |
| Framework | Fiber `3.3.0` |
| Entrypoint | `backend/cmd/server` |
| Middleware | Recover, CORS, dan logger |
| Lifecycle | Graceful shutdown dengan timeout 10 detik |
| Default URL | `http://localhost:8080` |

Konfigurasi backend memberikan default untuk `APP_ENV`, `BACKEND_PORT`, dan
`DB_PORT`, serta menolak startup bila konfigurasi database wajib tidak tersedia.

### Environment Variables

| Variable | Default/contoh | Keterangan |
| --- | --- | --- |
| `APP_ENV` | `development` | Nama environment aplikasi. |
| `BACKEND_PORT` | `8080` | Port HTTP backend. |
| `DB_HOST` | `mysql` | Host MySQL pada network Compose. |
| `DB_PORT` | `3306` | Port internal MySQL. |
| `MYSQL_HOST_PORT` | `3306` | Port MySQL yang diekspos pada host. |
| `DB_NAME` | `api_integrator` | Nama database. |
| `DB_USER` | `gateway` | User database development. |
| `DB_PASSWORD` | Nilai development | Password wajib yang tidak boleh disimpan pada `.env` di Git. |
| `VITE_API_BASE_URL` | `http://localhost:8080` | Base URL API untuk frontend. |

Sumber konfigurasi contoh tersedia pada
[`/.env.example`](../../.env.example).

## Kontrak Health Check

Request:

```http
GET /health
```

Response sukses:

```json
{
  "status": "success",
  "data": {
    "service": "api-integrator-gateway",
    "environment": "development"
  }
}
```

Endpoint menghasilkan HTTP `200` dan dipakai oleh Docker Compose untuk
menentukan kesiapan backend.

## Docker Compose

```text
Developer
   │
   ├── localhost:5173 ──> Frontend
   │                         │ depends_on: healthy
   ├── localhost:8080 ──> Backend
   │                         │ depends_on: healthy
   └── localhost:3306 ──> MySQL 8.4
```

Compose memakai named volume `mysql_data` dan healthcheck pada ketiga service.
Urutan startup adalah MySQL, backend, kemudian frontend.

Pada mesin implementasi, port host `3306` sudah digunakan. Verifikasi runtime
dilakukan menggunakan `MYSQL_HOST_PORT=3307`, sementara port internal MySQL dan
default repository tetap `3306`.

## Pelaksanaan TDD

### RED

Test ditulis sebelum implementasi berikut tersedia:

- React `App` shell.
- Helper konfigurasi API frontend.
- `Config.Load` backend.
- Fiber app factory `NewApp`.
- Kontrak response `GET /health`.

Eksekusi awal gagal karena module implementasi belum tersedia. Kondisi ini
menjadi bukti fase RED.

### GREEN

Implementasi minimum ditambahkan hingga:

- 3 frontend test lulus.
- 3 backend test lulus.
- Frontend lint dan production build lulus.
- Backend vet dan build pada Go 1.26 lulus.

### REFACTOR

Setelah test lulus, struktur dirapikan dengan:

- Pemisahan environment configuration.
- Fiber app factory yang dapat diuji.
- Graceful shutdown pada entrypoint.
- Middleware recovery, CORS, dan logger.
- Konfigurasi ESLint dan script build/test konsisten.
- Multi-stage backend Docker image dengan non-root runtime user.

## Deliverables

| Deliverable | Status | Catatan |
| --- | --- | --- |
| Struktur frontend/backend | Selesai | Source dipisahkan berdasarkan runtime. |
| Frontend dapat dijalankan | Selesai | `npm run dev`, port `5173`. |
| Backend dapat dijalankan | Selesai | `go run ./cmd/server`, port `8080`. |
| Konfigurasi Tailwind | Selesai | Menggunakan plugin Tailwind Vite. |
| Environment MySQL | Selesai | Template dan validasi tersedia. |
| Dokumentasi setup lokal | Selesai | Tersedia pada [README proyek](../../README.md). |
| CI workflow | Selesai | Frontend, backend, dan Docker jobs. |
| Docker Compose | Selesai | Frontend, backend, dan MySQL dengan healthcheck. |
| GitHub Issues dan Project | Tertunda | Terhalang izin GitHub MCP dan browser. |

## Acceptance Criteria

| Acceptance Criteria | Status | Bukti |
| --- | --- | --- |
| Repository dapat disiapkan kurang dari 5 menit | Belum diukur independen | Alur disederhanakan menjadi salin `.env` dan `docker compose up --build`; download awal dikecualikan. |
| Frontend berjalan pada `localhost:5173` | Lulus | HTTP frontend berhasil dan healthcheck Compose healthy. |
| Backend berjalan pada `localhost:8080` | Lulus | `/health` menghasilkan HTTP `200`. |
| Frontend test, lint, dan build lulus | Lulus | 3 test, ESLint, dan Vite production build berhasil. |
| Backend test, vet, dan build lulus | Lulus | 3 test serta Go vet/build berhasil pada toolchain 1.26. |
| Docker image dapat dibangun | Lulus | Image frontend dan backend berhasil dibangun. |
| Seluruh service Compose healthy | Lulus | Frontend, backend, dan MySQL terverifikasi healthy. |

## Bukti Verifikasi

Perintah yang dijalankan selama Sprint 1:

```powershell
# Frontend
npm.cmd test
npm.cmd run lint
npm.cmd run build

# Backend dengan Go 1.26
go test ./...
go vet ./...
go build ./cmd/server

# Docker
docker compose config --quiet
docker compose build
docker compose up --detach

# Runtime
curl.exe --fail http://localhost:8080/health
curl.exe --fail http://localhost:5173
```

Hasil utama:

- Frontend: 2 test files, 3 test cases, seluruhnya lulus.
- Backend: package config dan server lulus, total 3 test cases.
- Vite production build berhasil.
- Backend Go 1.26 test, vet, dan build berhasil.
- Docker Compose config dan image build berhasil.
- Response `/health` sesuai kontrak.
- Frontend, backend, dan MySQL mencapai status healthy.

## Kendala dan Penyelesaian

### Konflik Port MySQL

Port `3306` pada host sudah digunakan. Compose diperbarui agar port host dapat
diatur melalui `MYSQL_HOST_PORT`. Port internal service tetap `3306`, sehingga
konfigurasi antar-container tidak berubah.

### Toolchain Go Lokal

Mesin implementasi memiliki Go 1.23.4, sedangkan Fiber 3 membutuhkan Go yang
lebih baru. Test, vet, build, dan dependency resolution backend dijalankan
menggunakan image Go 1.26 agar sesuai dengan baseline proyek.

### GitHub Project

Pembuatan issue melalui GitHub MCP ditolak dengan HTTP `403 Resource not
accessible by personal access token`. Browser in-app juga tidak tersedia.
Akibatnya, 12 issue dan project organisasi `RPL 2` belum dibuat. Isi issue dan
pembagian assignee telah disiapkan pada
[GITHUB_PROJECT_ISSUES.md](../project-management/GITHUB_PROJECT_ISSUES.md).

### Batasan Database Sprint 1

Sprint ini hanya menyiapkan dan memvalidasi konfigurasi MySQL. Backend belum
membuka koneksi database dan belum menjalankan migration. Pekerjaan persistence
tetap mengikuti sprint database pada roadmap.

## Risiko Tersisa

- Workflow CI baru dapat dibuktikan pada GitHub setelah commit diproses oleh
  GitHub Actions.
- Waktu setup kurang dari lima menit perlu diukur oleh contributor lain pada
  clone bersih.
- Hak tulis GitHub MCP dan akses organization Projects perlu diperbaiki sebelum
  backlog dapat dipublikasikan.
- Kredensial `.env.example` hanya ditujukan untuk development dan wajib diganti
  pada environment bersama atau production.

## Handoff ke Sprint 2

Sprint 2 dapat menggunakan baseline ini untuk membangun landing page publik.
Prioritas handoff:

1. Buat layout, navigation, hero, feature, integration flow, FAQ, dan CTA.
2. Pertahankan shell dan konfigurasi API yang sudah diuji.
3. Tambahkan test responsivitas dan komponen untuk landing page.
4. Putuskan apakah konten landing bersifat statis atau disediakan melalui
   endpoint opsional `GET /landing`.
5. Setelah izin GitHub tersedia, buat issue Sprint 1–12 dan project `RPL 2`
   menggunakan backlog yang sudah disiapkan.
