# API Integrator Gateway

Menyediakan spesifikasi produk untuk modul API Integrator yang menjadi pintu masuk semua komunikasi antar aplikasi dalam ekosistem UMKM. API Integrator memastikan routing, keamanan, validasi, logging, dan standarisasi semua request sebelum diteruskan ke SmartBank atau layanan lain.

Sprint 3 menyediakan autentikasi backend berbasis MySQL, bcrypt, dan JWT.
Laporan implementasi tersedia pada
[Sprint 3 Report](docs/report/SPRINT_03_REPORT.md).

## Kebutuhan

- Docker Desktop dengan Docker Compose v2; atau
- Node.js 24 LTS dan npm 11; serta
- Go 1.26.4 atau lebih baru; dan
- MySQL 9.7 LTS untuk setup tanpa Docker.

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
- Backend login: `POST http://localhost:8080/auth/login`
- Backend current user: `GET http://localhost:8080/auth/me`
- MySQL: `localhost:3306` secara default

Hentikan layanan dengan `docker compose down`. Data MySQL dipertahankan pada
named volume `mysql_data`. Gunakan `docker compose down --volumes` hanya jika
data lokal memang ingin dihapus.

Nilai pada `.env.example` hanya untuk development. Ubah password dan jangan
commit file `.env`. Jika port MySQL host sudah digunakan, ubah
`MYSQL_HOST_PORT` pada `.env`; port internal `DB_PORT` tetap `3306`.
Seed user hanya ditujukan untuk development dan otomatis ditolak ketika
`APP_ENV=production`.

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
$env:JWT_SECRET = "minimum-32-character-secret-change-this"
$env:JWT_TTL = "1h"
$env:JWT_ISSUER = "api-integrator-gateway"
$env:SEED_USERS_ENABLED = "false"

Set-Location backend
go run ./cmd/server
```

Backend membuka koneksi MySQL dan menjalankan migration Goose saat startup.
Jika seed diaktifkan, seluruh username/password seed wajib tersedia melalui
environment variables yang dicontohkan pada `.env.example`.

## Landing page Sprint 2

Landing page dapat diakses tanpa login dan menyediakan:

- Penjelasan layanan dan manfaat API Integrator.
- Alur request dari aplikasi menuju gateway dan layanan tujuan.
- Peran SmartBank, Marketplace, POS, SupplierHub, LogistiKita, UMKM Insight,
  dan API Gateway.
- Use case integrasi, FAQ, serta tautan ke repositori.
- Navigasi responsif dengan dukungan keyboard dan mobile menu.

CTA login sengaja berstatus `Segera hadir`. Autentikasi dan halaman login
frontend berada pada Sprint 4, tetapi backend autentikasi sudah tersedia.

## Authentication backend Sprint 3

Login menggunakan kombinasi `username`, `password`, dan `app_name`:

```http
POST /auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin-development-password",
  "app_name": "API Gateway"
}
```

Response sukses berisi JWT HS256, role, aplikasi, dashboard tujuan, dan waktu
kedaluwarsa. Token memuat `sub`, `username`, `role`, `app_name`, `iss`, `iat`,
`nbf`, dan `exp`.

Validasi token:

```http
GET /auth/me
Authorization: Bearer <token>
```

Role yang tersedia:

- `admin_gateway` → `/dashboard/admin`
- `app_user` → `/dashboard/user`
- `monitoring_user` → `/dashboard/monitoring`

`POST /auth/logout` tidak tersedia karena token bersifat stateless. Sprint 4
melakukan logout dengan menghapus token di sisi frontend.

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
go test -race ./...
go vet ./...
go build ./cmd/server
```

Validasi Docker:

```powershell
Copy-Item .env.example .env
docker compose config --quiet
docker compose build --pull
docker compose up --detach
docker compose ps
```

## Struktur proyek

```text
.
|-- backend/
|   |-- cmd/server/          # Entrypoint HTTP server
|   |-- config/              # Environment configuration
|   `-- internal/
|       |-- auth/            # Bcrypt, JWT, login service, dan seed
|       |-- database/        # Koneksi dan migration Goose
|       |-- model/           # Model dan role user
|       |-- repository/      # Repository MySQL
|       `-- server/          # Fiber app factory, middleware, dan routes
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
|-- AGENTS.md                 # Aturan TDD dan checklist Docker per sprint
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
