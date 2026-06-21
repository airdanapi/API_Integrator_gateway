# Sprint 10 Report: Chat System

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 10 - Chat System |
| Baseline | `3d5f8ef feat: implement sprint 9 notifications system` |
| Tanggal laporan | 21 Juni 2026 |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint 10 mengimplementasikan fitur chat untuk komunikasi antara admin gateway dan app user. Tabel `chat_messages` dari Sprint 5 dipakai kembali untuk menyimpan pesan.

Scope yang selesai:

- Backend chat service, repository, handler, dan route registration.
- Role-based access: admin bisa chat dengan semua app_user, app_user hanya dengan admin, monitoring_user tidak memiliki akses chat.
- Frontend chat service module dan ChatDrawer component.
- Chat drawer dipasang di AdminDashboardPage dan UserDashboardPage (tidak di MonitoringDashboardPage).
- Polling 2 detik untuk percakapan baru.
- Test TDD backend/frontend, lint, build, dan verifikasi lengkap.

## Backend

### Package `internal/chat`

Package `internal/chat` berisi service utama untuk:

- **ListConversations**: admin melihat semua app_user, app_user melihat admin saja. Monitoring_user mendapat error 403.
- **History**: fetch riwayat pesan per conversation dengan pagination. Visibility check memastikan user hanya bisa melihat conversation mereka sendiri.
- **SendMessage**: admin mengirim ke app_user tertentu (butuh `to_username` + `to_app_name`), app_user mengirim ke admin (otomatis). Pesan max 1000 karakter.
- **MarkRead**: menandai semua pesan dalam conversation sebagai dibaca.

### Conversation ID Format

Conversation ID menggunakan format `{adminUsername}__{appUsername}__{appName}`, contoh: `admin__marketplace__Marketplace`.

### API Contract

#### `GET /chat/conversations`

Protected endpoint. Response sukses:

```json
{
  "status": "success",
  "data": {
    "conversations": [
      {
        "conversation_id": "admin__marketplace__Marketplace",
        "target_username": "marketplace",
        "target_app_name": "Marketplace",
        "unread_count": 1,
        "latest_message": {
          "id": 2,
          "conversation_id": "admin__marketplace__Marketplace",
          "from_user": "marketplace",
          "to_user": "admin",
          "message": "Butuh bantuan integrasi",
          "timestamp": "2026-06-21T10:01:00Z",
          "is_read": false
        }
      }
    ],
    "total_unread": 1
  }
}
```

#### `GET /chat/history`

Protected endpoint. Query: `conversation_id` (required), `page` (default 1), `limit` (default 20, max 50).

Response sukses:

```json
{
  "status": "success",
  "data": {
    "messages": [
      {
        "id": 1,
        "conversation_id": "admin__marketplace__Marketplace",
        "from_user": "admin",
        "to_user": "marketplace",
        "message": "Halo",
        "timestamp": "2026-06-21T10:00:00Z",
        "is_read": true
      }
    ],
    "total_unread": 0,
    "page": 1,
    "limit": 20
  }
}
```

#### `POST /chat/message`

Protected endpoint. Body untuk admin:

```json
{ "to_username": "marketplace", "to_app_name": "Marketplace", "message": "Halo" }
```

Body untuk app_user:

```json
{ "message": "Butuh bantuan" }
```

Response sukses:

```json
{
  "status": "success",
  "data": {
    "message": { "id": 7, "conversation_id": "admin__marketplace__Marketplace", "from_user": "admin", "to_user": "marketplace", "message": "Halo", "timestamp": "...", "is_read": false },
    "total_unread": 0
  }
}
```

#### `POST /chat/read`

Protected endpoint. Body:

```json
{ "conversation_id": "admin__marketplace__Marketplace" }
```

Response sukses:

```json
{
  "status": "success",
  "data": { "total_unread": 0 }
}
```

Error behavior:

- `400` untuk request body invalid atau field required kosong.
- `401` tanpa token atau token invalid.
- `403` untuk role yang tidak memiliki akses chat (monitoring_user).
- `404` jika conversation tidak ada atau tidak visible untuk user.

### Wiring

File `backend/internal/server/auth.go` menambahkan `ChatService` ke struct `Dependencies`. File `backend/internal/server/app.go` mendaftarkan 4 route chat. File `backend/cmd/server/main.go` membuat `chatRepository` dan `chatService`, lalu menginjeksi ke Dependencies.

## Frontend

### Service

File `frontend/src/services/chat.js` menambahkan:

- `fetchChatConversations(apiClient)` - GET `/chat/conversations`
- `fetchChatHistory(apiClient, conversationId, { page, limit })` - GET `/chat/history`
- `sendChatMessage(apiClient, { to_username, to_app_name, message })` - POST `/chat/message`
- `markChatRead(apiClient, conversationId)` - POST `/chat/read`

### Component

File `frontend/src/components/ChatDrawer.jsx` menambahkan drawer dengan:

- Badge unread count (format `99+`).
- Tombol Chat untuk membuka drawer.
- Daftar conversation (admin melihat semua app_user, user melihat admin).
- Klik conversation memuat history dan otomatis mark read.
- Form kirim pesan dengan input "Pesan" dan tombol "Kirim".
- Empty state: "Belum ada percakapan."
- Error state: role="alert" dengan "Gagal memuat chat."
- Polling 2 detik untuk conversations.
- Styling Tailwind CSS responsive.

ChatDrawer dipasang di header:

- `AdminDashboardPage` - ada ChatDrawer
- `UserDashboardPage` - ada ChatDrawer
- `MonitoringDashboardPage` - TIDAK ada ChatDrawer (monitoring user tidak punya akses chat)

## TDD dan Test Coverage

### Backend

Test yang ditambahkan/diupdate:

- Repository: Insert, CountUnread, ListByConversation, LatestByConversation, CountUnreadByConversation, MarkAsRead, ListConversations.
- Service: ListConversations (admin/app_user/monitoring), SendMessage (validasi/admin-to-user/user-to-admin/forbidden/not-found), History (valid/empty/not-visible/forbidden), MarkRead (success/empty/not-found/forbidden), nil dependencies, normalizePagination.
- Server handler: 401 tanpa auth, conversations success, history validation + not found, history success, send validation + forbidden, send + read success.

### Frontend

Test yang ditambahkan/diupdate:

- Service: fetchChatConversations, fetchChatHistory, sendChatMessage, markChatRead.
- Component ChatDrawer: badge, buka drawer, pilih conversation, load history, mark read, kirim pesan, empty state, error state, polling 2 detik.
- Dashboard pages: ChatDrawer trigger di Admin dan User, tidak ada di Monitoring.
- App tests: mock ChatDrawer agar tidak melakukan API call nyata.

## Hasil Verifikasi

### Frontend

```text
docker compose exec -T frontend npm run lint
PASS

docker compose exec -T frontend npm test
Test Files  12 passed (12)
Tests       67 passed (67)

docker compose exec -T frontend npm run build
PASS, dist generated by Vite
```

### Backend

```text
docker run --rm --workdir /src --volume "<repo>/backend:/src" golang:1.26.4-alpine3.24 go test ./...
ok  config           0.039s
ok  internal/auth    2.024s
ok  internal/chat    0.233s
ok  internal/dashboard 0.029s
ok  internal/database  0.045s
ok  internal/notification 0.212s
ok  internal/repository 0.051s
ok  internal/server    0.208s

docker run --rm --workdir /src --volume "<repo>/backend:/src" golang:1.26.4-alpine3.24 go vet ./...
PASS

docker run --rm --workdir /src --volume "<repo>/backend:/src" golang:1.26.4-alpine3.24 go build -o /tmp/api-integrator-server ./cmd/server
PASS
```

### Docker Checklist

| Pemeriksaan | Hasil |
| --- | --- |
| `docker compose config --quiet` | PASS |
| `docker compose build --pull frontend backend` | PASS |
| `docker compose up --detach` | PASS |
| `docker compose ps` | MySQL healthy, Backend healthy, Frontend healthy |

### Smoke Test Chat

| Skenario | Hasil |
| --- | --- |
| Login admin seed | `role=admin_gateway` |
| Login Marketplace seed | `role=app_user` |
| Login UMKM Insight seed | `role=monitoring_user` |
| Admin `GET /chat/conversations` | 200, 1 conversation (marketplace), total_unread=0 |
| Marketplace `GET /chat/conversations` | 200, 1 conversation (admin), total_unread=0 |
| Monitoring `GET /chat/conversations` | 403 forbidden |
| Admin `POST /chat/message` ke marketplace | 200, message id=1, conversation `admin__marketplace__Marketplace` |
| Marketplace `POST /chat/message` ke admin | 200, message id=2, conversation `admin__marketplace__Marketplace` |
| Admin `POST /chat/read` | 200, total_unread=0 |
| Admin `GET /chat/history?conversation_id=admin__marketplace__Marketplace` | 200, 2 messages, page=1, limit=20 |
| Tanpa token `GET /chat/conversations` | 401 unauthorized |

### Docker Audit

| Komponen | Versi / Image | Keputusan |
| --- | --- | --- |
| Docker Engine | 29.5.2 | Tidak diubah, dikelola Docker Desktop. |
| Docker Compose | 5.1.4 | Tidak diubah. |
| Frontend base image | `node:24-alpine3.24` | Dipertahankan, build `--pull` sukses. |
| Backend build image | `golang:1.26.4-alpine3.24` | Dipertahankan, build `--pull` sukses. |
| Backend runtime image | `alpine:3.24.1` | Dipertahankan, build `--pull` sukses. |
| Database image | `mysql:9.7` | Dipertahankan, tidak ada perubahan schema/volume. |

Tidak ada backup volume karena Sprint 10 tidak mengubah schema database, image database, atau named volume.

## File Utama yang Dibuat / Diubah

| File | Status | Keterangan |
| --- | --- | --- |
| `backend/internal/chat/service.go` | BARU (dari sesi sebelumnya) | Chat service, visibility, send, history, mark-read. |
| `backend/internal/chat/service_test.go` | BARU | Unit test service chat. |
| `backend/internal/server/chat.go` | BARU (dari sesi sebelumnya) | Handler chat endpoints. |
| `backend/internal/server/chat_test.go` | BARU (dari sesi sebelumnya) | Handler test chat endpoints. |
| `backend/internal/repository/chat_repository.go` | DIUBAH (dari sesi sebelumnya) | MySQL chat repository. |
| `backend/internal/repository/chat_repository_test.go` | DIUBAH | Tambah MarkAsRead dan ListConversations test. |
| `backend/internal/repository/user_repository.go` | DIUBAH (dari sesi sebelumnya) | Tambah FindFirstByRole, ListByRole. |
| `backend/internal/server/auth.go` | DIUBAH | Tambah ChatService ke Dependencies. |
| `backend/internal/server/app.go` | DIUBAH | Register chat routes. |
| `backend/cmd/server/main.go` | DIUBAH | Wire chat repository dan service. |
| `frontend/src/services/chat.js` | BARU | Client service chat. |
| `frontend/src/services/chat.test.js` | BARU (dari sesi sebelumnya) | Test client service chat. |
| `frontend/src/components/ChatDrawer.jsx` | BARU | Drawer chat UI. |
| `frontend/src/components/ChatDrawer.test.jsx` | BARU (dari sesi sebelumnya) | Test ChatDrawer component. |
| `frontend/src/pages/AdminDashboardPage.jsx` | DIUBAH | Tambah ChatDrawer di header. |
| `frontend/src/pages/UserDashboardPage.jsx` | DIUBAH | Tambah ChatDrawer di header. |
| `frontend/src/pages/MonitoringDashboardPage.jsx` | TIDAK DIUBAH | Tidak ada ChatDrawer (sesuai desain). |

## Risiko Tersisa dan Handoff ke Sprint 11

- Chat menggunakan polling 2 detik, belum WebSocket/SSE.
- Monitoring user tidak memiliki akses chat; jika produk menginginkan monitoring bisa melihat riwayat chat (read-only), perlu perubahan contract.
- Chat drawer hanya tersedia di halaman dashboard, bukan landing page.
