# Dokumentasi Pengembangan Aplikasi

## 1. Ringkasan
Dokumentasi ini menjelaskan arsitektur, tumpukan teknologi, struktur proyek, dan langkah pengembangan untuk aplikasi ekosistem UMKM.

- Frontend: React
- Styling: Tailwind CSS
- Backend: Go (Golang)
- Framework backend: Fiber
- Tujuan: API Gateway/Integrator dengan landing page, login role-based, dashboard admin/user, serta integrasi antar aplikasi.

## 2. Tujuan Pengembangan
Aplikasi dibangun untuk mendukung skema tugas besar: API integrator yang menjadi pintu masuk semua request antar aplikasi, di samping portal frontend untuk landing page, login, dan dashboard.

## 3. Teknologi Utama

### Frontend
- React
- Vite (opsional untuk scaffolding dan kecepatan dev)
- React Router
- Tailwind CSS
- Axios atau Fetch API untuk komunikasi dengan backend

### Backend
- Go
- Fiber
- JWT untuk otentikasi
- Middleware untuk logging, CORS, dan validasi
- Database: MySQL

### Integrasi
- JSON REST API
- JWT authentication
- Role-based access control: `admin_gateway`, `app_user`, `monitoring_user`

## 3.1 Sumber Data API
- API dan endpoint yang akan dihubungkan didasarkan pada file:
  - `dokumen_tugas/TugasBesar_Plan.xlsx - 2. Fungsional.csv`
- Dokumen tersebut memuat daftar fitur, endpoint, input, proses, output, dan kontrak API untuk SmartBank, Marketplace, POS, SupplierHub, dan LogistiKita.
- List sumber data API:
  - `SmartBank`: endpoint terkait registrasi/login, manajemen saldo, transfer, pembayaran transaksi, pinjaman, pajak/biaya, ledger transaksi, biaya layanan bank.
  - `Marketplace`: endpoint terkait manajemen produk, browse produk, checkout, integrasi pembayaran, status order, biaya layanan marketplace.
  - `POS`: endpoint pembayaran kasir dan pengiriman request pembayaran ke SmartBank.
  - `SupplierHub`: endpoint order bahan, konfirmasi supplier, dan request pembayaran supplier.
  - `LogistiKita`: endpoint request pengiriman, perhitungan ongkir, pembayaran logistik, dan update status pengiriman.
  - `UMKM Insight`: sumber data read-only dari SmartBank untuk dashboard analitik.

## 4. Arsitektur Sistem

1. User mengakses landing page React
2. Login melalui backend GoFiber
3. Backend menerbitkan JWT dengan klaim `role` dan `app_name`
4. React menyimpan token di `localStorage` atau `sessionStorage`
5. User mengakses dashboard sesuai role
6. Dashboard memanggil API gateway untuk melihat status, history, dan log

## 5. Struktur Proyek yang Direkomendasikan

```
project-root/
  frontend/
    public/
    src/
      assets/
      components/
      pages/
        LandingPage.jsx
        LoginPage.jsx
        AdminDashboard.jsx
        UserDashboard.jsx
        MonitoringDashboard.jsx
      routes/
      services/
        api.js
      App.jsx
      main.jsx
    tailwind.config.js
    postcss.config.js
    package.json
  backend/
    cmd/
      server/
        main.go
    config/
      config.go
    controllers/
      auth.go
      dashboard.go
      gateway.go
      landing.go
    middlewares/
      auth.go
      logger.go
      cors.go
    models/
      user.go
      request_log.go
      dashboard_data.go
    repositories/
      user_repo.go
      log_repo.go
    services/
      auth_service.go
      jwt_service.go
      gateway_service.go
    utils/
      response.go
      validation.go
    go.mod
    go.sum
  .gitignore
  README.md
```

## 6. Setup Frontend

### 6.1 Inisialisasi Proyek React

```bash
cd project-root
cd frontend
npm create vite@latest . -- --template react
npm install
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

### 6.2 Konfigurasi Tailwind

`tailwind.config.js`
```js
module.exports = {
  content: [
    './index.html',
    './src/**/*.{js,jsx,ts,tsx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

`src/index.css`
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

### 6.3 Struktur Halaman Utama
- `LandingPage.jsx`: menampilkan layanan API Integrator, benefits, dan alur integrasi.
- `LoginPage.jsx`: formulir login untuk admin dan user aplikasi.
- `AdminDashboard.jsx`: ringkasan traffic gateway, audit logs, status session, indikator tiap API, dan grafik analitik.
- `UserDashboard.jsx`: status request, histori permintaan, notifikasi, dan grafik performa aplikasi.
- `MonitoringDashboard.jsx`: tampilan read-only untuk UMKM Insight.

### 6.4 Contoh Service API
`src/services/api.js`
```js
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

export default api;
```

## 7. Setup Backend

### 7.1 Inisialisasi Proyek Go

```bash
cd project-root
cd backend
go mod init github.com/yourusername/project-name
go get github.com/gofiber/fiber/v2
go get github.com/gofiber/jwt/v3
go get github.com/gofiber/cors/v2
go get github.com/gofiber/logger/v2
go get github.com/golang-jwt/jwt/v5
go get github.com/go-sql-driver/mysql
```

### 7.2 Contoh Struktur Folder

- `cmd/server/main.go`: entrypoint aplikasi.
- `config/`: konfigurasi database dan env.
- `controllers/`: handler endpoint.
- `middlewares/`: auth, logging, CORS.
- `services/`: logika JWT, login, routing.
- `repositories/`: akses data MySQL.
- `models/`: definisi entitas.
- `utils/`: response helper.

### 7.3 Contoh `main.go`

```go
package main

import (
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/logger/v2"
  "github.com/gofiber/cors/v2"
  "github.com/yourusername/project-name/backend/controllers"
  "github.com/yourusername/project-name/backend/middlewares"
)

func main() {
  app := fiber.New()

  app.Use(cors.New())
  app.Use(logger.New())

  app.Get("/landing", controllers.GetLanding)
  app.Post("/auth/login", controllers.Login)

  app.Use(middlewares.JWTProtected)
  app.Get("/dashboard/admin", controllers.AdminDashboard)
  app.Get("/dashboard/user", controllers.UserDashboard)
  app.Get("/dashboard/monitoring", controllers.MonitoringDashboard)
  app.Get("/notifications", controllers.GetNotifications)
  app.Get("/chat/history", controllers.GetChatHistory)
  app.Post("/chat/message", controllers.SendChatMessage)
  app.Post("/gateway/payment", controllers.GatewayPayment)
  app.Post("/gateway/smartbank", controllers.GatewaySmartBank)
  app.Post("/gateway/marketplace", controllers.GatewayMarketplace)
  app.Post("/gateway/logistics", controllers.GatewayLogistics)
  app.Post("/gateway/supplier", controllers.GatewaySupplier)

  app.Listen(":8080")
}
```

### 7.4 Middleware JWT

`middlewares/auth.go`
```go
package middlewares

import (
  "github.com/gofiber/fiber/v2"
)

func JWTProtected(c *fiber.Ctx) error {
  token := c.Get("Authorization")
  if token == "" {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
  }
  // validasi token di sini
  return c.Next()
}
```

### 7.5 Endpoints Utama

#### Autentikasi & Header
- Semua endpoint selain `/landing` dan `/auth/login` harus menggunakan header:
  - `Authorization: Bearer <token>`
- Token JWT harus memuat klaim `role` dan `app_name`.

#### Landing
- `GET /landing`
  - Output: `{ service_overview, application_roles, integration_flow, contact_info }`

#### Auth
- `POST /auth/login`
  - Input: `{ username, password, app_name }`
  - Output: `{ status, token, role, dashboard_url }`

#### Dashboard
- `GET /dashboard/admin`
  - Output: `{ traffic_summary, active_sessions, audit_logs, error_rate, service_indicators, analytics_graphs }`
- `GET /dashboard/user`
  - Output: `{ service_status, request_history, notifications, performance_graphs }`
- `GET /dashboard/monitoring`
  - Output: `{ summary_data, analytics_links, read_only_access }`

#### Notifications
- `GET /notifications`
  - Output: `{ notifications }`
  - Contoh: notifikasi API yang tidak aktif selama 1 minggu.

#### Chat
- `POST /chat/message`
  - Input: `{ from_user, to_user, message, timestamp }`
  - Output: `{ status, message_id }`
- `GET /chat/history`
  - Input: `conversation_id`
  - Output: `{ messages }`

#### Gateway
- `POST /gateway/payment`
  - Input: `{ from_app, from_user, to_user, amount, metadata, service_type }`
  - Output: `{ status, transaction_id, message }`
- `POST /gateway/smartbank`
  - Input: `{ action, payload }`
  - Output: `{ status, data }`
- `POST /gateway/marketplace`
  - Input: `{ action, payload }`
  - Output: `{ status, data }`
- `POST /gateway/logistics`
  - Input: `{ order_id, address, distance, shipping_type }`
  - Output: `{ status, delivery_id, message }`
- `POST /gateway/supplier`
  - Input: `{ supplier_id, material, qty, total_cost }`
  - Output: `{ status, order_id, message }`

### 7.6 API Contract
- Semua response sukses harus menggunakan format JSON dengan properti `status` dan `data`.
- `400 Bad Request` untuk payload atau validasi yang salah.
- `401 Unauthorized` untuk token tidak valid atau tidak ada.
- `500 Internal Server Error` untuk kegagalan internal.
- Admin memiliki akses ke `/dashboard/admin` dan semua log.
- `app_user` memiliki akses ke `/dashboard/user` dan endpoint gateway sesuai peran aplikasi.
- `monitoring_user` hanya memiliki akses read-only ke `/dashboard/monitoring` dan tidak dapat memanggil endpoint transaksi.

## 8. Database dan Model Data

### Model Pengguna
Aplikasi backend harus menyediakan model pengguna dengan atribut:
- `ID`
- `Username`
- `PasswordHash`
- `Role` (`admin_gateway`, `app_user`, `monitoring_user`)
- `AppName` (`SmartBank`, `Marketplace`, `POS`, `SupplierHub`, `LogistiKita`, `UMKM Insight`)

### Model Log Request
- `ID`
- `Timestamp`
- `SourceApp`
- `Endpoint`
- `Payload`
- `Status`
- `Response`

### Model Notifications
- `ID`
- `CreatedAt`
- `AppName`
- `Type`
- `Message`
- `IsRead`

### Model ChatMessage
- `ID`
- `ConversationID`
- `FromUser`
- `ToUser`
- `Message`
- `Timestamp`

### MySQL Database
- Gunakan MySQL untuk penyimpanan `users`, `request_logs`, `notifications`, `chat_messages`, dan data dashboard.
- Tabel utama:
  - `users`
  - `request_logs`
  - `notifications`
  - `chat_messages`
  - `dashboard_data`
- Koneksi di backend dapat dikonfigurasi dengan DSN MySQL:
  - `USER:PASSWORD@tcp(HOST:PORT)/DBNAME?charset=utf8mb4&parseTime=True&loc=Local`
- Simpan konfigurasi sensitif di environment variables seperti `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_NAME`.

### Model Dashboard
- `TrafficSummary`
- `ActiveSessions`
- `AuditLogs`
- `ServiceStatus`
- `RequestHistory`

## 9. Alur Pengembangan

### 9.1 Frontend
1. Buat komponen UI untuk landing page.
2. Buat formulir login dan service `api.js` untuk panggilan backend.
3. Implementasikan protected routes berdasarkan role.
4. Buat halaman dashboard dengan pengambilan data dari API.
5. Gunakan Tailwind untuk styling komponen.

### 9.2 Backend
1. Siapkan server Fiber, middleware CORS dan logger.
2. Buat endpoint auth dan keluarkan JWT.
3. Buat middleware untuk memverifikasi token dan role.
4. Buat controller landing page dan dashboard.
5. Buat controller gateway untuk menerima dan meneruskan request.
6. Simpan log akses dalam database.

## 10. Deployment dan Build

### Frontend
```bash
cd frontend
npm run build
```

### Backend
```bash
cd backend
go build ./cmd/server
./server
```

## 11. Testing

- Uji login dengan setiap role user.
- Uji akses dashboard admin, user, monitoring.
- Uji endpoint gateway dan pastikan request diarahkan ke layanan yang tepat.
- Uji landing page dapat diakses tanpa login.
- Uji otorisasi JWT dan pastikan role-based access berfungsi.

## 12. Best Practices

- Pisahkan logic bisnis di service layer.
- Gunakan middleware untuk keamanan dan logging.
- Validasi semua input request.
- Simpan konfigurasi sensitif di environment variables.
- Buat dokumentasi API endpoint.

## 13. Catatan Khusus

Aplikasi ini diposisikan sebagai integrator di antara SmartBank, Marketplace, POS, SupplierHub, dan LogistiKita. Pastikan:
- SmartBank tetap menjadi satu-satunya sumber kebenaran untuk semua transaksi keuangan.
- UMKM Insight hanya memiliki akses read-only.
- Semua aplikasi berkomunikasi melalui gateway.

---

Dokumentasi ini bisa dijadikan acuan pengembangan frontend React dengan Tailwind dan backend GoFiber untuk aplikasi API Integrator UMKM.