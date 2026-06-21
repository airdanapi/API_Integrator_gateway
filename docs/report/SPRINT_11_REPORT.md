# Sprint 11 Report: Gateway Routing & API Integration

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 11 - Gateway Routing & API Integration |
| Baseline | `8537d3e feat: implement sprint 10 chat system with backend service, repository, frontend ChatDrawer, RBAC, TDD, and Docker verification` |
| Tanggal laporan | 21 Juni 2026 |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint 11 mengimplementasikan gateway routing untuk menerima request antar aplikasi, memvalidasi payload, meneruskan request ke upstream bila URL tersedia, dan mencatat request/response ke tabel `request_logs`.

Scope yang selesai:

- Backend package `internal/gateway` untuk model payload, validasi, HTTP forwarder, retry transient error, service routing, RBAC, dan request logging.
- Endpoint protected `POST /gateway/payment`, `/gateway/smartbank`, `/gateway/marketplace`, `/gateway/logistics`, dan `/gateway/supplier`.
- Wiring config upstream dari environment variable `GATEWAY_*_URL`, dependency injection, route registration, dan startup server.
- Frontend service module `src/services/gateway.js` untuk memanggil semua endpoint gateway.
- Test TDD backend/frontend, lint, vet, build, Docker checklist, dan smoke test gateway.

## Backend

### Package `internal/gateway`

Package ini berisi:

- `models.go`: request contract untuk payment, SmartBank, Marketplace, LogistiKita, SupplierHub, response standar gateway, dan sentinel error.
- `forwarder.go`: HTTP forwarder dengan `POST`, `Content-Type: application/json`, timeout 10 detik, retry 2 kali untuk error transient, dan durasi request.
- `service.go`: role check, payload validation, upstream selection, response assembly, dan insert ke `request_logs`.

Role access:

- `admin_gateway`: boleh memanggil endpoint gateway.
- `app_user`: boleh memanggil endpoint gateway.
- `monitoring_user`: ditolak dengan `403 forbidden`.

Jika upstream URL kosong, request tetap diterima dengan `forwarded=false` dan `upstream=not_configured`. Ini menjaga setup development tetap bisa berjalan tanpa layanan eksternal.

### Environment

Variabel upstream baru:

```text
GATEWAY_SMARTBANK_URL=
GATEWAY_MARKETPLACE_URL=
GATEWAY_LOGISTICS_URL=
GATEWAY_SUPPLIERHUB_URL=
```

Nilai kosong berarti upstream belum dikonfigurasi. Semua nilai dibaca dari `backend/config` dan di-trim.

### API Contract

Semua endpoint membutuhkan header:

```http
Authorization: Bearer <token>
Content-Type: application/json
```

#### `POST /gateway/payment`

```json
{
  "from_app": "Marketplace",
  "from_user": "user1",
  "to_user": "user2",
  "amount": 10000,
  "metadata": { "order_id": "ORD-001" },
  "service_type": "payment"
}
```

#### `POST /gateway/smartbank`

```json
{ "action": "check_balance", "payload": { "account": "123456" } }
```

#### `POST /gateway/marketplace`

```json
{ "action": "get_order", "payload": { "order_id": "ORD-001" } }
```

#### `POST /gateway/logistics`

```json
{
  "order_id": "ORD-001",
  "address": "Jl. Merdeka No. 1",
  "distance": 7.5,
  "shipping_type": "express"
}
```

#### `POST /gateway/supplier`

```json
{
  "supplier_id": "SUP-001",
  "material": "Beras",
  "qty": 12,
  "total_cost": 240000
}
```

Response sukses standar:

```json
{
  "status": "success",
  "data": {
    "status": "success",
    "transaction_id": "gw-payment-...",
    "message": "request accepted, upstream not configured",
    "forwarded": false,
    "upstream": "not_configured"
  }
}
```

Error behavior:

- `400 invalid_request`: body bukan JSON valid.
- `400 invalid_payload`: field wajib kosong atau nilai numeric tidak valid.
- `401 unauthorized`: token tidak ada atau invalid.
- `403 forbidden`: role tidak boleh mengakses gateway.
- `502 upstream_failure`: upstream gagal dihubungi setelah retry.
- `503 upstream_not_configured`: reserved di handler; service saat ini memperlakukan upstream kosong sebagai accepted development mode.

## Frontend

File `frontend/src/services/gateway.js` menambahkan helper:

- `sendGatewayPayment(apiClient, payload)`
- `sendGatewaySmartBank(apiClient, payload)`
- `sendGatewayMarketplace(apiClient, payload)`
- `sendGatewayLogistics(apiClient, payload)`
- `sendGatewaySupplier(apiClient, payload)`

Setiap helper mengembalikan `response.data.data` agar konsisten dengan service frontend lain.

## TDD dan Test Coverage

### Backend

Test yang ditambahkan:

- Gateway model validation untuk semua payload dan invalid field.
- HTTP forwarder: URL kosong, success, upstream 500 preserved, retry transient error, retries exhausted, duration, context cancelled.
- Gateway service: success path, admin/app_user access, monitoring forbidden, invalid payload, upstream kosong, upstream failure, log content, nil log store.
- Server handler: auth required, invalid JSON, forbidden, success semua endpoint, invalid payload, upstream failure, nil service.
- Config loader: upstream gateway URL dibaca dan di-trim dari environment.

### Frontend

Test yang ditambahkan:

- Gateway client service mengirim payload ke endpoint yang benar untuk payment, SmartBank, Marketplace, LogistiKita, dan SupplierHub.

## Hasil Verifikasi

### Backend

```text
go test ./...
PASS

go vet ./...
PASS

go build -o $env:TEMP\api-integrator-server-sprint11.exe ./cmd/server
PASS
```

### Frontend

```text
npm run lint
PASS

npm test -- --run
Test Files  13 passed (13)
Tests       72 passed (72)

npm run build
PASS, dist generated by Vite
```

Catatan Windows lokal: `npm test` dan `npm run build` perlu dijalankan di luar sandbox karena Vite/Tailwind native binding ditolak dengan `spawn EPERM` saat berada di sandbox.

## Docker Checklist

| Pemeriksaan | Hasil |
| --- | --- |
| `docker --version` | Docker version 29.5.2, build 79eb04c |
| `docker compose version` | Docker Compose version v5.1.4 |
| `docker compose config --quiet` | PASS |
| `docker compose build --pull` | PASS |
| `docker compose up --detach` | PASS |
| `docker compose ps` | MySQL healthy, Backend healthy, Frontend healthy |

### Smoke Test Gateway

Dengan upstream URL kosong di `.env`, request valid diterima dan tidak diteruskan keluar.

| Skenario | Hasil |
| --- | --- |
| Login app user seed | 200, token diterbitkan |
| `POST /gateway/payment` app_user | 200, `forwarded=false`, `upstream=not_configured` |
| `POST /gateway/smartbank` app_user | 200, `forwarded=false`, `upstream=not_configured` |
| `POST /gateway/marketplace` app_user | 200, `forwarded=false`, `upstream=not_configured` |
| `POST /gateway/logistics` app_user | 200, `forwarded=false`, `upstream=not_configured` |
| `POST /gateway/supplier` app_user | 200, `forwarded=false`, `upstream=not_configured` |
| Tanpa token `POST /gateway/payment` | 401 unauthorized |
| Payload invalid `POST /gateway/payment` | 400 invalid_payload |
| Monitoring user `POST /gateway/payment` | 403 forbidden |
| Dashboard admin audit logs | 5 gateway logs terbaca, status 200, source_app Marketplace |

### Docker Audit

| Komponen | Versi / Image | Keputusan |
| --- | --- | --- |
| Docker Engine | 29.5.2 | Tidak diubah, dikelola Docker Desktop. |
| Docker Compose | v5.1.4 | Tidak diubah. |
| Frontend base image | `node:24-alpine3.24@sha256:156b55f92e98ccd5ef49578a8cea0df4679826564bad1c9d4ef04462b9f0ded6` | Dipertahankan, build `--pull` sukses. |
| Backend build image | `golang:1.26.4-alpine3.24@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648` | Dipertahankan, build `--pull` sukses. |
| Backend runtime image | `alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b` | Dipertahankan, build `--pull` sukses. |
| Database image | `mysql:9.7@sha256:e370cd5f64599d46985b7729b452f2153825246f88d82753ec595c5dfc6fef6a` | Dipertahankan, tidak ada perubahan schema/volume. |

Tidak ada backup volume karena Sprint 11 tidak mengubah schema database, image database, atau named volume.

## File Utama yang Dibuat / Diubah

| File | Status | Keterangan |
| --- | --- | --- |
| `backend/internal/gateway/models.go` | BARU | Payload contract, validation, response, sentinel error. |
| `backend/internal/gateway/models_test.go` | BARU | Unit test validasi payload gateway. |
| `backend/internal/gateway/forwarder.go` | BARU | HTTP forwarder, timeout, retry transient error. |
| `backend/internal/gateway/forwarder_test.go` | BARU | Unit test forwarder dan retry behavior. |
| `backend/internal/gateway/service.go` | BARU | Gateway routing service, RBAC, logging. |
| `backend/internal/gateway/service_test.go` | BARU | Unit test service routing dan logging. |
| `backend/internal/server/gateway.go` | BARU | Fiber handler untuk endpoint gateway. |
| `backend/internal/server/gateway_test.go` | BARU | Handler test gateway endpoints. |
| `backend/config/config.go` | DIUBAH | Tambah config upstream gateway. |
| `backend/config/config_test.go` | DIUBAH | Tambah test load upstream gateway URL. |
| `backend/internal/server/auth.go` | DIUBAH | Tambah `GatewayService` ke Dependencies. |
| `backend/internal/server/app.go` | DIUBAH | Register route gateway. |
| `backend/cmd/server/main.go` | DIUBAH | Wire gateway forwarder dan service. |
| `.env.example` | DIUBAH | Tambah env upstream gateway. |
| `frontend/src/services/gateway.js` | BARU | Client service untuk gateway endpoints. |
| `frontend/src/services/gateway.test.js` | BARU | Test client service gateway. |

## Risiko Tersisa dan Handoff ke Sprint 12

- Upstream eksternal belum tersedia di environment lokal; smoke test memakai mode `not_configured`.
- Mapping response spesifik dari upstream nyata perlu divalidasi saat URL SmartBank, Marketplace, LogistiKita, dan SupplierHub tersedia.
- Belum ada load test gateway; ini masuk area Sprint 12 testing dan deployment readiness.