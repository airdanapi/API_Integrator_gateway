# GitHub Project: RPL 2

Project organisasi: `airdanapi`  
Repository: `airdanapi/API_Integrator_gateway`  
Status awal seluruh item: `Todo`

## Pembagian

| Sprint | Assignee |
| --- | --- |
| 1-4 | `venalism` |
| 5-8 | `Nuthfih` |
| 9-12 | `saripudin14` |

## [Sprint 1] Project Setup

**Assignee:** `venalism`

### Goal

Menyiapkan infrastruktur dasar, setup repository, dan konfigurasi awal.

### Tasks

- [ ] Setup struktur folder `frontend/` dan `backend/`
- [ ] Inisialisasi React dengan Vite
- [ ] Inisialisasi Go Modules dan Fiber
- [ ] Setup Tailwind CSS dan environment variables MySQL
- [ ] Dokumentasikan setup lokal
- [ ] Setup CI dan Docker Compose

### Deliverables

- Struktur proyek frontend/backend
- Frontend dan backend dapat dijalankan
- Dokumentasi setup lokal
- Docker Compose full-stack

### Acceptance Criteria

- [ ] Setup operasional kurang dari 5 menit, di luar download awal
- [ ] Frontend berjalan di `localhost:5173`
- [ ] Backend berjalan di `localhost:8080`
- [ ] Test, lint, vet, build, dan healthcheck lulus

## [Sprint 2] Landing Page

**Assignee:** `venalism`

### Goal

Membuat landing page publik yang menjelaskan layanan API Integrator.

### Tasks

- [ ] Buat layout, header, footer, dan navigasi
- [ ] Buat hero, features, integration flow, use cases, FAQ, dan CTA
- [ ] Styling responsif dengan Tailwind CSS
- [ ] Tambahkan `GET /landing` dan CORS bila data disediakan backend

### Deliverables

- Landing page responsif
- Assets dan tampilan mobile-ready

### Acceptance Criteria

- [ ] Landing page dapat diakses tanpa login
- [ ] Responsif pada desktop, tablet, dan mobile
- [ ] Loading time kurang dari 2 detik

## [Sprint 3] Auth Backend

**Assignee:** `venalism`

### Goal

Mengimplementasikan login dan autentikasi JWT pada backend.

### Tasks

- [ ] Buat schema dan model user
- [ ] Buat JWT service dan password hashing
- [ ] Implementasikan `POST /auth/login`
- [ ] Implementasikan middleware JWT
- [ ] Sediakan seed user untuk seluruh role

### Deliverables

- JWT valid dengan klaim role dan app
- Schema user, middleware, dan seed data

### Acceptance Criteria

- [ ] Login menghasilkan token dalam kurang dari 100 ms
- [ ] Token memuat `role` dan `app_name`
- [ ] Route terproteksi menolak request tanpa token

## [Sprint 4] Auth Frontend

**Assignee:** `venalism`

### Goal

Mengimplementasikan login page dan session management frontend.

### Tasks

- [ ] Buat login form dan validasi error
- [ ] Setup React Router dan protected routes
- [ ] Implementasikan auth state dan login API
- [ ] Simpan token dan pasang request interceptor
- [ ] Implementasikan logout

### Deliverables

- Login page dan token management
- Protected route berdasarkan sesi

### Acceptance Criteria

- [ ] Kredensial valid dapat login
- [ ] Kredensial invalid menampilkan error
- [ ] Sesi bertahan setelah refresh bila token valid
- [ ] Logout menghapus token

## [Sprint 5] Database

**Assignee:** `Nuthfih`

### Goal

Melengkapi schema MySQL dan model inti sistem.

### Tasks

- [ ] Buat tabel request logs, notifications, chat messages, dan dashboard data
- [ ] Buat repository layer
- [ ] Setup connection pooling
- [ ] Buat migration dan seed data

### Deliverables

- Seluruh tabel dan repository tersedia
- Koneksi database siap digunakan

### Acceptance Criteria

- [ ] Seluruh tabel dapat di-query
- [ ] Persistence berjalan
- [ ] Fungsi repository dapat dipanggil

## [Sprint 6] Admin Backend

**Assignee:** `Nuthfih`

### Goal

Mengimplementasikan endpoint analitik dashboard admin.

### Tasks

- [ ] Buat `GET /dashboard/admin`
- [ ] Agregasikan traffic, sessions, errors, dan service indicators
- [ ] Deteksi API tidak aktif selama satu minggu
- [ ] Tambahkan pagination audit log

### Deliverables

- Endpoint dashboard admin
- Agregasi dan inactive API detection

### Acceptance Criteria

- [ ] Response kurang dari 500 ms
- [ ] Data akurat
- [ ] Pagination berjalan

## [Sprint 7] Admin Frontend

**Assignee:** `Nuthfih`

### Goal

Mengimplementasikan UI dashboard admin dengan visualisasi analitik.

### Tasks

- [ ] Buat Admin Dashboard
- [ ] Tambahkan charts dan summary cards
- [ ] Tambahkan audit table dengan filter dan sort
- [ ] Tampilkan status setiap service
- [ ] Poll data setiap 30 detik

### Deliverables

- Dashboard admin dengan charts dan tables
- Refresh data berkala

### Acceptance Criteria

- [ ] Chart responsif dan terbaca
- [ ] Polling 30 detik berjalan
- [ ] Filter, sort, loading, dan error state ditangani

## [Sprint 8] User/Monitoring Dashboard

**Assignee:** `Nuthfih`

### Goal

Mengimplementasikan dashboard app user dan monitoring user.

### Tasks

- [ ] Buat endpoint dashboard user dan monitoring
- [ ] Terapkan role-based data filtering
- [ ] Buat User Dashboard dan Monitoring Dashboard
- [ ] Tambahkan status, history, statistics, dan charts

### Deliverables

- Dashboard user dan monitoring
- Filtering berbasis role

### Acceptance Criteria

- [ ] User hanya melihat data aplikasinya
- [ ] Monitoring user memiliki akses read-only
- [ ] Chart diperbarui secara berkala

## [Sprint 9] Notifications

**Assignee:** `saripudin14`

### Goal

Mengimplementasikan notifikasi untuk inactive API dan alert operasional.

### Tasks

- [ ] Buat endpoint list dan mark-as-read notification
- [ ] Buat background scheduler
- [ ] Deteksi inactive API, error rate, dan response time
- [ ] Buat notification UI, badge, dan polling

### Deliverables

- Notification persistence dan scheduler
- Notification UI

### Acceptance Criteria

- [ ] Notifikasi muncul maksimal 5 menit setelah trigger
- [ ] Mark-as-read berjalan
- [ ] Badge count akurat

## [Sprint 10] Chat

**Assignee:** `saripudin14`

### Goal

Mengimplementasikan komunikasi chat antara admin dan app user.

### Tasks

- [ ] Buat endpoint send message, history, dan conversations
- [ ] Terapkan role-based chat access
- [ ] Buat Chat Window, Input, dan Conversation List
- [ ] Poll incoming message setiap 2 detik

### Deliverables

- Chat persistence dan UI
- Riwayat dan pembaruan pesan

### Acceptance Criteria

- [ ] Pesan terkirim dan tersimpan
- [ ] Riwayat dapat dimuat
- [ ] Pesan baru tampil maksimal 2 detik
- [ ] Admin dan user dapat berkomunikasi sesuai role

## [Sprint 11] Gateway Integration

**Assignee:** `saripudin14`

### Goal

Mengimplementasikan routing gateway ke SmartBank dan aplikasi lain.

### Tasks

- [ ] Buat endpoint payment, SmartBank, Marketplace, logistics, dan supplier
- [ ] Implementasikan validation, forwarding, retry, dan error handling
- [ ] Catat setiap request dan response
- [ ] Ikuti kontrak payload PRD

### Deliverables

- Gateway endpoints dan routing service
- Request/response logging

### Acceptance Criteria

- [ ] Payment diteruskan ke SmartBank
- [ ] Setiap forwarding dicatat
- [ ] Error handling berjalan
- [ ] Payload invalid ditolak

## [Sprint 12] Testing/Deployment

**Assignee:** `saripudin14`

### Goal

Memfinalisasi testing, dokumentasi, optimasi, dan kesiapan deployment.

### Tasks

- [ ] Lengkapi unit, integration, E2E, security, dan load test
- [ ] Optimalkan query dan frontend bundle
- [ ] Finalisasi API docs, deployment guide, user manual, dan troubleshooting
- [ ] Perbaiki bug dan siapkan deployment checklist

### Deliverables

- Seluruh test lulus
- Dokumentasi final
- Release yang siap deployment

### Acceptance Criteria

- [ ] Coverage lebih dari 80%
- [ ] Critical path teruji
- [ ] Tidak ada critical bug
- [ ] Load test lebih dari 100 request/detik
- [ ] Dokumentasi lengkap
