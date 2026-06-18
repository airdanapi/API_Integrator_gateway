# Sprint Implementation Plan - API Integrator Ekosistem UMKM
## 12 Sprint Development Roadmap

---

## Pendahuluan
Dokumen ini merupakan rencana implementasi untuk pengembangan sistem API Integrator dalam 12 sprint (durasi ~ 3 bulan, asumsi 1 sprint = 1 minggu). Setiap sprint memiliki goal, task, deliverable, dan acceptance criteria yang jelas.

---

## Sprint 1: Project Setup & Infrastructure
**Duration:** Week 1  
**Goal:** Menyiapkan infrastruktur dasar, setup repository, dan konfigurasi awal.

### Tasks
1. Setup repository Git dan struktur folder `frontend/` dan `backend/`
2. Inisialisasi proyek React dengan Vite
3. Inisialisasi proyek Go dengan Go Modules
4. Setup konfigurasi Tailwind CSS
5. Setup konfigurasi environment variables untuk MySQL
6. Dokumentasi setup lokal untuk tim developer
7. Setup GitHub CI/CD workflow (opsional untuk awal)

### Deliverables
- Repository dengan folder structure yang rapi
- Frontend dapat di-run dengan `npm run dev`
- Backend dapat di-run dengan `go run ./cmd/server`
- Dokumentasi setup lokal

### Acceptance Criteria
- Repo dapat di-clone dan di-setup dalam < 5 menit
- Frontend Vite berjalan di `localhost:5173`
- Backend Fiber siap di `localhost:8080`

---

## Sprint 2: Landing Page & Static Content
**Duration:** Week 2  
**Goal:** Membuat landing page publik yang menjelaskan layanan API Integrator.

### Tasks (Frontend)
1. Buat komponen Layout (Header, Footer, Navigation)
2. Buat hero section dengan penjelasan API Integrator
3. Buat section features/benefits sistem
4. Buat section integrasi flow antar aplikasi
5. Buat section testimoni atau use cases
6. Buat section FAQ
7. Buat section contact/CTA ke login
8. Styling dengan Tailwind CSS

### Tasks (Backend)
1. Buat endpoint `GET /landing` untuk serve data landing page (opsional)
2. Setup CORS middleware

### Deliverables
- Landing page fully responsive
- Assets (images, icons) siap
- Mobile-friendly design

### Acceptance Criteria
- Landing page dapat diakses tanpa login
- Responsif di desktop, tablet, mobile
- Loading time < 2 detik

---

## Sprint 3: Authentication System - Backend
**Duration:** Week 3  
**Goal:** Implementasi sistem login dan JWT authentication di backend.

### Tasks
1. Setup MySQL database dan create table `users`
2. Buat model User dengan fields (ID, Username, PasswordHash, Role, AppName)
3. Buat JWT service untuk generate dan validate token
4. Buat endpoint `POST /auth/login` dengan validasi username/password
5. Buat password hashing dengan bcrypt
6. Implementasikan JWT middleware untuk protected routes
7. Buat endpoint `POST /auth/logout` (opsional)
8. Buat seed data untuk test users (admin, app_user, monitoring_user)

### Deliverables
- JWT token yang valid dapat dihasilkan
- MySQL schema untuk users
- Middleware protection working

### Acceptance Criteria
- Login endpoint respond dengan token dalam < 100ms
- Token berisi claim `role` dan `app_name`
- Protected endpoint menolak request tanpa token
- Seed user data tersedia

---

## Sprint 4: Authentication System - Frontend
**Duration:** Week 4  
**Goal:** Implementasi login page dan session management di frontend.

### Tasks
1. Buat LoginPage component dengan form username/password
2. Setup React Router untuk routing publik vs protected
3. Setup context atau state management (Redux/Zustand) untuk auth state
4. Implementasikan login API call ke backend
5. Store token di localStorage
6. Buat interceptor axios untuk attach token ke semua request
7. Implementasikan logout functionality
8. Setup protected route wrapper component
9. Styling login page dengan Tailwind

### Deliverables
- Login page fully functional
- Token management setup
- Protected routes working

### Acceptance Criteria
- User dapat login dengan valid credentials
- Invalid credentials menampilkan error message
- Token tersimpan di localStorage
- Refresh page tetap terlogin jika token valid
- Logout menghapus token

---

## Sprint 5: Database Schema & Models Setup
**Duration:** Week 5  
**Goal:** Lengkapi schema MySQL dan model untuk seluruh sistem.

### Tasks (Backend)
1. Create table `request_logs` dengan fields (ID, Timestamp, SourceApp, Endpoint, Payload, Status, Response)
2. Create table `notifications` dengan fields (ID, CreatedAt, AppName, Type, Message, IsRead)
3. Create table `chat_messages` dengan fields (ID, ConversationID, FromUser, ToUser, Message, Timestamp)
4. Create table `dashboard_data` untuk cache analytics
5. Buat repository layer untuk akses data (UserRepo, LogRepo, NotificationRepo, ChatRepo)
6. Setup database connection pooling di config
7. Buat migration script atau seed data

### Deliverables
- Semua table MySQL tersedia
- Repository layer siap
- Database connection working

### Acceptance Criteria
- Semua table dapat query dengan benar
- Data persistence working
- Repository functions callable

---

## Sprint 6: Dashboard Admin - Backend
**Duration:** Week 6  
**Goal:** Implementasi endpoint backend untuk dashboard admin dengan data analitik.

### Tasks
1. Buat endpoint `GET /dashboard/admin` untuk fetch:
   - Traffic summary (total request, success rate)
   - Active sessions
   - Error logs dari request_logs
   - Service indicator (mana API yang tidak aktif 1 minggu)
   - Analytics data (aggregate dari logs)
2. Implementasikan logic untuk detect inactive API (1 minggu tanpa request)
3. Buat query aggregation untuk traffic graphs
4. Implementasikan pagination untuk audit logs
5. Buat caching untuk dashboard data (opsional, untuk performance)

### Deliverables
- Endpoint `/dashboard/admin` fully functional
- Data aggregation working
- Inactive API detection logic

### Acceptance Criteria
- Endpoint respond < 500ms
- Data akurat dan real-time
- Pagination working

---

## Sprint 7: Dashboard Admin - Frontend
**Duration:** Week 7  
**Goal:** Implementasi AdminDashboard UI dengan grafik analitik.

### Tasks
1. Buat AdminDashboard.jsx component
2. Setup chart library (Chart.js atau Recharts) untuk visualisasi
3. Buat cards untuk traffic summary, error rate, active sessions
4. Buat gauge/progress chart untuk service health
5. Buat tables untuk audit logs dengan filtering dan sorting
6. Buat service indicator visualization (status setiap API)
7. Implementasikan data fetching dari `/dashboard/admin` endpoint
8. Styling dengan Tailwind

### Deliverables
- AdminDashboard fully functional
- Charts dan graphs rendering correctly
- Real-time data update

### Acceptance Criteria
- Charts responsive dan readable
- Data refresh setiap 30 detik (polling)
- Filtering/sorting working
- Loading states handled

---

## Sprint 8: Dashboard User & Monitoring - Backend & Frontend
**Duration:** Week 8  
**Goal:** Implementasi dashboard untuk app users dan monitoring-only users.

### Tasks (Backend)
1. Buat endpoint `GET /dashboard/user` untuk app users (SmartBank, Marketplace, POS, SupplierHub, LogistiKita):
   - Service status (API yang user gunakan)
   - Request history per user
   - Performance graphs
2. Buat endpoint `GET /dashboard/monitoring` untuk UMKM Insight (read-only):
   - Summary data
   - Analytics links
3. Implementasikan role-based filtering data

### Tasks (Frontend)
1. Buat UserDashboard.jsx dengan:
   - Service status cards
   - Request history table
   - Performance graphs
2. Buat MonitoringDashboard.jsx dengan:
   - Read-only analytics view
   - Summary statistics
3. Styling dengan Tailwind

### Deliverables
- Dashboard user & monitoring fully functional
- Role-based data filtering
- Charts dan tables working

### Acceptance Criteria
- User hanya melihat data aplikasi mereka
- Monitoring user hanya read-only access
- Charts update real-time

---

## Sprint 9: Notifications System
**Duration:** Week 9  
**Goal:** Implementasi sistem notifikasi untuk API inactive dan alerts lainnya.

### Tasks (Backend)
1. Buat endpoint `GET /notifications` untuk fetch notifikasi user
2. Implementasikan background job untuk detect inactive API (1 minggu tanpa request)
3. Buat logic untuk generate notifikasi:
   - API inactive alert
   - Error rate alert
   - Response time alert
4. Buat endpoint `POST /notifications/read` untuk mark notifikasi sebagai read
5. Setup scheduler (cron job) untuk periodic checks
6. Implementasikan notification persistence di table `notifications`

### Tasks (Frontend)
1. Buat Notification component untuk display notifikasi
2. Buat notification bell icon di header dengan badge count
3. Buat notification dropdown/modal
4. Implementasikan real-time notification fetch (polling atau WebSocket)
5. Buat dismiss/read functionality
6. Styling dengan Tailwind

### Deliverables
- Notification system fully functional
- Notifikasi inactive API working
- Real-time notification display

### Acceptance Criteria
- Notifikasi muncul dalam 5 menit setelah trigger
- Mark as read working
- Notification badge count accurate

---

## Sprint 10: Chat System
**Duration:** Week 10  
**Goal:** Implementasi fitur chat untuk komunikasi admin-user.

### Tasks (Backend)
1. Buat endpoint `POST /chat/message` untuk send pesan:
   - Input: from_user, to_user, message, timestamp
   - Persist ke table chat_messages
2. Buat endpoint `GET /chat/history` untuk fetch riwayat chat:
   - Pagination support
   - Filter by conversation
3. Buat endpoint `GET /chat/conversations` untuk list aktif conversations
4. Implementasikan role-based access (admin bisa chat dengan siapa saja, user hanya dengan admin)

### Tasks (Frontend)
1. Buat ChatWindow component untuk display chat history
2. Buat ChatInput component untuk send pesan
3. Buat ChatList component untuk list conversations
4. Implementasikan real-time chat fetch (polling 2 detik)
5. Buat chat modal/drawer di dashboard
6. Styling dengan Tailwind
7. Notifikasi incoming message (bell + badge)

### Deliverables
- Chat system fully functional
- Message persistence
- Real-time chat display

### Acceptance Criteria
- Pesan terkirim dan tersimpan
- Chat history loadable
- Real-time message appear dalam 2 detik
- Admin & user bisa communicate

---

## Sprint 11: Gateway Routing & API Integration
**Duration:** Week 11  
**Goal:** Implementasi gateway routing untuk forward request ke SmartBank dan aplikasi lain.

### Tasks (Backend)
1. Buat endpoint `POST /gateway/payment` untuk forward payment request ke SmartBank
2. Buat endpoint `POST /gateway/smartbank` untuk general SmartBank operations
3. Buat endpoint `POST /gateway/marketplace` untuk Marketplace routing
4. Buat endpoint `POST /gateway/logistics` untuk LogistiKita routing
5. Buat endpoint `POST /gateway/supplier` untuk SupplierHub routing
6. Implementasikan request logging di table request_logs setiap forward terjadi
7. Buat error handling & retry logic
8. Implementasikan request validation sesuai kontrak di PRD
9. Setup gateway service untuk handle routing logic

### Deliverables
- Gateway endpoints fully functional
- Request/response logging working
- Routing logic implemented

### Acceptance Criteria
- Payment request forwarded ke SmartBank dengan benar
- Log setiap request/response
- Error handling working
- Validation reject invalid payload

---

## Sprint 12: Testing, Documentation & Deployment
**Duration:** Week 12  
**Goal:** Testing menyeluruh, finalisasi dokumentasi, dan preparation untuk deployment.

### Tasks
1. **Backend Testing:**
   - Unit test untuk services (auth, gateway, notification)
   - Integration test untuk endpoints
   - Load testing untuk gateway performance
   - Test dengan semua role (admin, app_user, monitoring_user)

2. **Frontend Testing:**
   - Component testing untuk pages
   - Integration test untuk auth flow
   - E2E test untuk critical user journey (login → dashboard → navigate)
   - Responsive design testing

3. **Security Testing:**
   - JWT validation test
   - SQL injection prevention
   - CORS policy test
   - Rate limiting test

4. **Performance Optimization:**
   - Database query optimization
   - Frontend bundle size optimization
   - Dashboard loading time < 1 detik

5. **Documentation:**
   - Update API documentation
   - Create deployment guide
   - Create user manual
   - Create troubleshooting guide

6. **Bug Fixes:**
   - Resolve semua bugs yang ditemukan di testing
   - Performance tuning

7. **Deployment Preparation:**
   - Setup production database
   - Configure environment variables
   - Setup Docker (opsional)
   - Create deployment checklist

### Deliverables
- All tests passing
- Final documentation
- Deployment-ready code
- Release notes

### Acceptance Criteria
- > 80% code coverage
- All critical paths tested
- Zero critical bugs
- Load test success (> 100 req/sec)
- All documentation complete

---

## Development Timeline Summary

| Sprint | Week | Focus Area | Status |
|--------|------|-----------|--------|
| 1 | 1 | Project Setup | Selesai |
| 2 | 2 | Landing Page | Selesai |
| 3 | 3 | Auth Backend | Selesai |
| 4 | 4 | Auth Frontend | Planning |
| 5 | 5 | Database Schema | Planning |
| 6 | 6 | Dashboard Admin Backend | Planning |
| 7 | 7 | Dashboard Admin Frontend | Planning |
| 8 | 8 | Dashboard User & Monitoring | Planning |
| 9 | 9 | Notifications | Planning |
| 10 | 10 | Chat System | Planning |
| 11 | 11 | Gateway & Integration | Planning |
| 12 | 12 | Testing & Deployment | Planning |

---

## Team Requirements
- **Frontend Developer:** 1-2 orang (React, Tailwind, Chart library)
- **Backend Developer:** 1-2 orang (Go, Fiber, MySQL)
- **DevOps/Infrastructure:** 0-1 orang (database setup, CI/CD)
- **QA/Tester:** 1 orang (testing, documentation)

Total: ~4-5 developer untuk optimal progress.

---

## Key Milestones
- **End of Sprint 4:** Auth system complete (backend + frontend)
- **End of Sprint 8:** All dashboards functional
- **End of Sprint 11:** Gateway integration complete
- **End of Sprint 12:** System ready for UAT/production

---

## Risk & Mitigation
| Risk | Mitigation |
|------|-----------|
| Database performance | Implement indexing, caching di Sprint 6 |
| API integration complexity | Early integration test di Sprint 11 |
| Frontend responsiveness | Test early & often throughout project |
| Team availability | Overlap sprints jika diperlukan acceleration |

---

## Notes
- Durasi sprint bisa disesuaikan (1 sprint = 1-2 minggu tergantung kapasitas tim)
- Sprint planning bisa dilakukan setiap Senin
- Daily standup recommended untuk tracking progress
- Retrospective setiap akhir sprint untuk improvement

---

Dokumen ini dibuat berdasarkan PRD dan Dokumentasi Pengembangan Aplikasi API Integrator UMKM.
