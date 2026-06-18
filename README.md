# API Integrator Gateway

Menyediakan spesifikasi produk untuk modul API Integrator yang menjadi pintu masuk semua komunikasi antar aplikasi dalam ekosistem UMKM. API Integrator memastikan routing, keamanan, validasi, logging, dan standarisasi semua request sebelum diteruskan ke SmartBank atau layanan lain.

## Kebutuhan

- Docker Desktop dengan Docker Compose v2; atau
- Node.js 22.12 atau lebih baru dan npm 10; serta
- Go 1.26 atau lebih baru; dan
- MySQL 8.4 untuk setup tanpa Docker.

## Setup tercepat dengan Docker

Salin konfigurasi contoh:

```powershell
Copy-Item .env.example .env
```

Jalankan seluruh layanan:

```powershell
docker compose up --build
```

Layanan yang tersedia:

- Frontend: <http://localhost:5173>
- Backend health check: <http://localhost:8080/health>
- Backend landing content: <http://localhost:8080/landing>
- MySQL: `localhost:3306` secara default

Hentikan layanan dengan `docker compose down`. Data MySQL dipertahankan pada
named volume `mysql_data`. Gunakan `docker compose down --volumes` hanya jika
data lokal memang ingin dihapus.

Nilai pada `.env.example` hanya untuk development. Ubah password dan jangan
commit file `.env`. Jika port MySQL host sudah digunakan, ubah
`MYSQL_HOST_PORT` pada `.env`; port internal `DB_PORT` tetap `3306`.

## Setup lokal tanpa Docker

### Frontend

```powershell
Set-Location frontend
npm.cmd ci
npm.cmd run dev
```

Gunakan `npm.cmd` bila PowerShell menolak menjalankan `npm.ps1` karena execution
policy.

### Backend

Pastikan MySQL sudah aktif, lalu siapkan environment:

```powershell
$env:APP_ENV = "development"
$env:BACKEND_PORT = "8080"
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_NAME = "api_integrator"
$env:DB_USER = "gateway"
$env:DB_PASSWORD = "your-local-password"

Set-Location backend
go run ./cmd/server
```

Sprint 1 memvalidasi konfigurasi database, tetapi belum membuat koneksi atau
schema. Implementasi persistence dijadwalkan pada sprint berikutnya.

## Landing page Sprint 2

Landing page dapat diakses tanpa login dan menyediakan:

- Penjelasan layanan dan manfaat API Integrator.
- Alur request dari aplikasi menuju gateway dan layanan tujuan.
- Peran SmartBank, Marketplace, POS, SupplierHub, LogistiKita, UMKM Insight,
  dan API Gateway.
- Use case integrasi, FAQ, serta tautan ke repositori.
- Navigasi responsif dengan dukungan keyboard dan mobile menu.

CTA login sengaja berstatus `Segera hadir`. Autentikasi dan halaman login
berada di luar scope Sprint 2.

## Test dan build

Frontend:

```powershell
Set-Location frontend
npm.cmd test
npm.cmd run lint
npm.cmd run build
```

Backend:

```powershell
Set-Location backend
go test ./...
go vet ./...
go build ./cmd/server
```

Validasi Docker:

```powershell
Copy-Item .env.example .env
docker compose config
docker compose build
```

## Struktur proyek

```text
.
|-- backend/
|   |-- cmd/server/          # Entrypoint HTTP server
|   |-- config/              # Environment configuration
|   `-- internal/server/     # Fiber app factory dan routes
|-- frontend/
|   `-- src/
|       |-- components/      # Komponen landing page
|       |-- data/            # Konten statis landing page
|       `-- test/            # Setup pengujian frontend
|-- docs/
|   |-- architecture/       # Diagram dan dokumentasi arsitektur
|   |-- development/        # Panduan teknis pengembangan
|   |-- planning/           # Roadmap implementasi sprint
|   |-- project-management/ # Backlog dan administrasi proyek
|   |-- requirements/       # Kebutuhan dan kontrak produk
|   |-- report/             # Laporan pelaksanaan sprint
|   `-- misc/source-data/   # Data sumber tugas besar
|-- .github/workflows/       # Continuous integration
`-- compose.yaml             # Frontend, backend, dan MySQL
```

## Kontrak health check

`GET /health` menghasilkan HTTP `200`:

```json
{
  "status": "success",
  "data": {
    "service": "api-integrator-gateway",
    "environment": "development"
  }
}
```

## Kontrak landing page

`GET /landing` bersifat publik dan menghasilkan HTTP `200` dengan struktur:

```json
{
  "status": "success",
  "data": {
    "service_overview": {},
    "application_roles": [],
    "integration_flow": [],
    "contact_info": {
      "repository_url": "https://github.com/airdanapi/API_Integrator_gateway",
      "login_status": "coming_soon"
    }
  }
}
```
