# Sprint 03 Report: Authentication System - Backend

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 3 — Authentication System - Backend |
| Tanggal laporan | 18 Juni 2026 |
| PIC | `venalism` |
| Commit implementasi | `c5dcabc` |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint 3 menghasilkan autentikasi backend berbasis MySQL, bcrypt, dan JWT.
Backend sekarang menjalankan migration Goose saat startup, menyediakan seed
development idempotent untuk tiga role, menerima login melalui
`POST /auth/login`, dan memvalidasi Bearer token melalui `GET /auth/me`.

Implementasi dilakukan dengan Test-Driven Development. Fase RED dibuktikan oleh
test yang gagal karena model, konfigurasi JWT, migration, repository, service,
dan middleware belum tersedia. Setelah GREEN dan REFACTOR, 22 test backend,
race detector, vet, build, integration test MySQL 9.7, serta seluruh regresi
frontend lulus.

Docker diperbarui ke Go 1.26.4, Alpine 3.24.1, Node 24 LTS, dan MySQL 9.7 LTS.
Volume MySQL 8.4 dibackup sebelum upgrade dan upgrade in-place berhasil.

## Hasil Implementasi

### Database dan User

- Koneksi MySQL menggunakan connection pool dan validasi startup.
- Migration Goose embedded membuat tabel `users`.
- Kombinasi `(username, app_name)` memiliki unique constraint.
- Password disimpan sebagai bcrypt cost 10 dan tidak pernah plaintext.
- Seed development bersifat upsert/idempotent untuk:
  - `admin_gateway` pada `API Gateway`
  - `app_user` pada `Marketplace`
  - `monitoring_user` pada `UMKM Insight`
- Seeding hanya berjalan ketika `SEED_USERS_ENABLED=true` dan ditolak pada
  `APP_ENV=production`.

### JWT dan API

JWT menggunakan HS256 dengan klaim:

```text
sub, username, role, app_name, iss, iat, nbf, exp
```

Endpoint publik:

```http
POST /auth/login
```

Input:

```json
{
  "username": "admin",
  "password": "admin-development-password",
  "app_name": "API Gateway"
}
```

Response sukses:

```json
{
  "status": "success",
  "data": {
    "token": "<jwt>",
    "role": "admin_gateway",
    "app_name": "API Gateway",
    "dashboard_url": "/dashboard/admin",
    "expires_in": 3600
  }
}
```

Endpoint terlindungi:

```http
GET /auth/me
Authorization: Bearer <token>
```

Payload tidak lengkap menghasilkan `400`. Kredensial salah dan token invalid
menghasilkan `401` dengan pesan generik. `/health` dan `/landing` tetap publik.

`POST /auth/logout` tidak dibuat karena JWT bersifat stateless. Logout frontend
Sprint 4 cukup menghapus token lokal.

## Pelaksanaan TDD

### RED

Test ditulis lebih dahulu untuk:

- Validasi konfigurasi JWT, TTL, dan seed.
- Bcrypt cost dan verifikasi password.
- Generate/validate JWT serta token expired, salah signature, dan malformed.
- Login seluruh role dan response error generik.
- Login user yang tidak ada tetap menjalankan dummy bcrypt comparison untuk
  mengurangi timing-based username enumeration.
- Repository query dan upsert.
- Migration schema dan idempotensi.
- Seed idempotent dan larangan password plaintext.
- Kontrak HTTP login dan middleware Bearer.

Eksekusi RED gagal karena package `model`, dependency MySQL/JWT, field config,
service, repository, dan migration belum tersedia.

### GREEN

Implementasi minimum ditambahkan sampai seluruh test unit dan kontrak lulus.
Integration test kemudian dijalankan langsung terhadap MySQL 9.7 Compose.

### REFACTOR

- App factory menggunakan dependency injection agar test HTTP tidak memerlukan
  database.
- Kode dipisahkan ke package `auth`, `database`, `model`, `repository`, dan
  `server`.
- Migration disimpan sebagai embedded SQL.
- Docker build memakai BuildKit cache mount dan `.cache` dikeluarkan dari build
  context.

## Docker dan Upgrade Infrastruktur

| Komponen | Sebelum | Sprint 3 |
| --- | --- | --- |
| Build Go | `golang:1.26-alpine` | `golang:1.26.4-alpine3.24` |
| Runtime backend | `alpine:3.23` | `alpine:3.24.1` |
| Frontend | `node:22-alpine` | `node:24-alpine3.24` |
| Database | `mysql:8.4` | `mysql:9.7` |

Sebelum upgrade database, stack dihentikan dan volume disalin ke:

```text
api-integrator-gateway_mysql_data_backup_sprint2_mysql84
```

MySQL melaporkan upgrade server dari `8.4.10` ke `9.7.1` berhasil. Setelah
startup, service MySQL, backend, dan frontend berstatus healthy. Tabel
`goose_db_version` dan `users` tersedia serta tiga seed user berhasil dibuat.

Audit runtime menggunakan Docker Engine `29.5.2`, Docker Compose `5.1.4`,
Node.js `24.16.0`, npm `11.13.0`, dan MySQL `9.7.1`. Image backend berukuran
sekitar `8.59 MB`, frontend `104 MB`, dan MySQL `270 MB`.

## Verifikasi dan Acceptance Criteria

Quality gate:

```powershell
# Backend
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/server

# Frontend regression
npm.cmd test
npm.cmd run lint
npm.cmd run build

# Docker
docker compose config --quiet
docker compose build --pull
docker compose up --detach
docker compose ps
```

Hasil:

- Backend: 22 top-level test lulus.
- Race detector, vet, dan production build lulus.
- Integration test migration MySQL 9.7 lulus dan idempotent.
- Frontend: 9 test, lint, dan production build lulus.
- Login ketiga role dan `GET /auth/me` berhasil.
- `/auth/me` tanpa token menghasilkan `401`.
- Payload login tidak lengkap menghasilkan `400`.
- Kredensial salah menghasilkan `401`.
- Login handler final tercatat `81.02 ms`, `74.21 ms`, dan `83.41 ms`, seluruhnya
  di bawah target 100 ms. Pengukuran PowerShell client mencakup overhead proses
  dan jaringan lokal sehingga tidak dipakai sebagai latency handler.

| Acceptance Criteria | Status | Bukti |
| --- | --- | --- |
| Login menghasilkan token kurang dari 100 ms | Lulus | Server timing maksimum `83.41 ms`. |
| Token memuat `role` dan `app_name` | Lulus | Unit test claims dan smoke test tiga role. |
| Protected route menolak request tanpa token | Lulus | `/auth/me` menghasilkan `401`. |
| Seed user tersedia | Lulus | Admin, Marketplace, dan UMKM Insight tersedia di MySQL. |
| Password tidak disimpan plaintext | Lulus | Bcrypt cost 10 dan query verifikasi hash non-kosong. |
| Migration idempotent | Lulus | Migration dijalankan dua kali dalam integration test. |

## CI

Workflow CI diperbarui untuk:

- Node.js 24.
- Backend unit test, race test, vet, dan build.
- Job integration test dengan service MySQL 9.7.
- Docker build menggunakan `--pull`.

## Risiko Tersisa dan Handoff Sprint 4

- Kredensial `.env.example` hanya untuk development dan wajib diganti pada
  environment bersama.
- JWT belum memiliki revocation list; logout Sprint 4 menghapus token lokal.
- Rate limiting login belum termasuk Sprint 3 dan perlu dipertimbangkan sebelum
  deployment publik.
- Sprint 4 dapat menggunakan `POST /auth/login` dan `GET /auth/me` untuk login,
  pemulihan sesi, role routing, dan protected frontend routes.
