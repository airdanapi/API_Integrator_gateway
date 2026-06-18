# Product Requirements Document (PRD)
## API Integrator / API Gateway untuk Ekosistem UMKM

### 1. Tujuan
Menyediakan spesifikasi produk untuk modul API Integrator yang menjadi pintu masuk semua komunikasi antar aplikasi dalam ekosistem UMKM. API Integrator memastikan routing, keamanan, validasi, logging, dan standarisasi semua request sebelum diteruskan ke SmartBank atau layanan lain.

### 2. Latar Belakang
Berdasarkan dokumen tugas besar, sistem terdiri dari beberapa aplikasi utama: SmartBank, Marketplace, POS, SupplierHub, LogistiKita, UMKM Insight, dan API Gateway/Integrator. API Integrator wajib menjadi jalur tunggal (single entry point) untuk semua permintaan antar-aplikasi.

### 3. Ruang Lingkup
Termasuk:
- Routing request antar aplikasi
- Validasi JWT
- Logging request/response
- Standarisasi format komunikasi JSON
- Penyaringan dan penerusan request ke SmartBank, Marketplace, POS, SupplierHub, dan LogistiKita
- Portal landing page yang menjelaskan layanan API Integrator
- Sistem login role-based untuk admin dan pengguna aplikasi
- Dashboard terpisah untuk admin dan user
- Pemisahan business logic dari orkestrasi

Tidak termasuk:
- Logika bisnis pembayaran atau manajemen saldo
- Penghitungan fee transaksi
- Penyimpanan data transaksi selain log audit

### 4. Stakeholder
- Dosen pembimbing dan evaluasi tugas
- Tim pengembang aplikasi UMKM (mahasiswa)
- SmartBank sebagai sistem pembayaran pusat
- Marketplace, POS, SupplierHub, LogistiKita sebagai konsumen API
- UMKM Insight sebagai pembaca data read-only

### 4.1 User Persona
- `Admin Gateway`:
  - Peran: Pengelola dan pemantau infrastruktur API.
  - Tujuan: Memastikan semua aplikasi terhubung dengan benar, melihat log, dan menangani masalah.
  - Kebutuhkan: Dashboard admin, akses audit logs, status sistem, kontrol akses.
- `User SmartBank`:
  - Peran: Pengguna sistem pembayaran pusat.
  - Tujuan: Mengelola transaksi keuangan dan melihat status integrasi.
  - Kebutuhkan: Dashboard user, status request, notifikasi, keamanan token.
- `User Marketplace`:
  - Peran: Penjual/buyer di marketplace yang menggunakan integrator untuk pembayaran.
  - Tujuan: Menyelesaikan transaksi penjualan/pembelian secara aman.
  - Kebutuhkan: Dashboard user, notifikasi status order, akses API payment.
- `User POS`:
  - Peran: Kasir atau operator toko fisik.
  - Tujuan: Memproses transaksi offline melalui gateway dan memastikan payment request diteruskan.
  - Kebutuhkan: Dashboard user, histori transaksi, status pembayaran.
- `User SupplierHub`:
  - Peran: Manajer bahan baku dan supplier.
  - Tujuan: Mengirim order bahan dan memastikan pembayaran supplier diproses.
  - Kebutuhkan: Dashboard user, status pembayaran supplier, tracking order bahan.
- `User LogistiKita`:
  - Peran: Pengelola pengiriman barang.
  - Tujuan: Menerima request pengiriman dan meminta pembayaran ongkir melalui SmartBank.
  - Kebutuhkan: Dashboard user, monitoring pengiriman, status request.
- `UMKM Insight (monitoring_user)`:
  - Peran: Analitik dan reporting read-only.
  - Tujuan: Memantau performa transaksi tanpa mengubah data.
  - Kebutuhkan: Dashboard monitoring, akses read-only, ringkasan analytics.

### 5. Syarat Fungsional Utama
1. Semua komunikasi antar aplikasi harus melalui API Integrator.
2. API Integrator harus memvalidasi JWT untuk setiap request.
3. API Integrator harus melakukan logging setiap request dan response.
4. API Integrator harus meneruskan request yang valid ke tujuan yang tepat.
5. API Integrator harus menolak request yang tidak memiliki token atau tidak memenuhi skema.
6. Endpoint transaksi keuangan wajib diteruskan ke SmartBank.
7. API Integrator dapat menerapkan rate limiting minimal untuk memastikan kestabilan.
8. API Integrator harus menjaga konsistensi payload JSON dan header standar.
9. Sistem harus menyediakan landing page publik yang menjelaskan layanan integrasi dan peran API Gateway.
10. Harus tersedia login role-based untuk:
   - Admin API Gateway sebagai pengelola sistem
   - User SmartBank, Marketplace, POS, SupplierHub, LogistiKita, UMKM Insight hanya sebagai monitoring/read-only
11. Harus ada dashboard terpisah:
   - Dashboard admin untuk manajemen gateway, monitoring traffic, audit logs, indikator layanan, dan grafik analitik
   - Dashboard user untuk melihat status layanan, histori request, notifikasi, dan grafik performa aplikasi
12. Harus ada fitur notifikasi untuk mendeteksi jika suatu API tidak aktif digunakan selama 1 minggu.
13. Harus ada fitur chat untuk komunikasi antara admin dan user aplikasi.

### 6. Syarat Non-Fungsional
- Respon dalam 200ms untuk request gateway normal.
- Standar keamanan: JWT + HTTPS (simulasi, jika belum tersedia sertifikat TLS maka hanya konseptual).
- Logging minimal: timestamp, source application, endpoint tujuan, status response.
- Struktur code: MVC / Clean Code seperti diinstruksikan dokumen.
- Format komunikasi JSON.

### 7. Business Rules dan Kebijakan
- SmartBank adalah pusat kontrol untuk semua transaksi keuangan.
- Tidak ada aplikasi lain yang boleh mengubah saldo langsung.
- Semua output transaksi harus dalam bentuk payment request yang diarahkan ke SmartBank.
- Gateway hanya melakukan orkestrasi, bukan logika bisnis transaksi.
- Analytics hanya membaca data; tidak boleh melalui API Integrator untuk operasi menulis.
- Tiap endpoint dianggap kontrak sistem dan harus jelas serta konsisten.

### 8. Endpoint API Gateway
Berikut adalah kontrak API utama untuk integrator:

#### 8.1 Autentikasi dan Middleware
- `Authorization: Bearer <token>`
- Validasi JWT wajib untuk semua endpoint.
- Jika token tidak valid: response `401 Unauthorized`.
- Jika format request salah: response `400 Bad Request`.
- Login role-based harus mendukung minimal 3 tipe pengguna:
  - `admin_gateway` untuk admin API Gateway
  - `app_user` untuk user SmartBank, Marketplace, POS, SupplierHub, dan LogistiKita
  - `monitoring_user` untuk UMKM Insight yang hanya melihat data read-only
- Session atau token user harus menyertakan klaim `role` dan `app_name`.

#### 8.2 Landing Page dan Dashboard
- `GET /landing`
  - Deskripsi: Menampilkan halaman informasi layanan API Integrator.
  - Output: `service_overview`, `application_roles`, `integration_flow`, `contact_info`
- `POST /auth/login`
  - Deskripsi: Autentikasi login untuk admin dan pengguna aplikasi.
  - Input: `{ username, password, app_name }`
  - Output: `status`, `token`, `role`, `dashboard_url`
- `GET /dashboard/admin`
  - Deskripsi: Dashboard khusus admin API Gateway.
  - Output: `traffic_summary`, `active_sessions`, `audit_logs`, `error_rate`, `service_indicators`, `analytics_graphs`
- `GET /dashboard/user`
  - Deskripsi: Dashboard untuk pengguna aplikasi (SmartBank, Marketplace, POS, SupplierHub, LogistiKita).
  - Output: `service_status`, `request_history`, `notifications`, `performance_graphs`
- `GET /dashboard/monitoring`
  - Deskripsi: Dashboard read-only untuk UMKM Insight.
  - Output: `summary_data`, `analytics_links`, `read_only_access`
- `GET /notifications`
  - Deskripsi: Menampilkan notifikasi sistem, termasuk API yang tidak aktif selama 1 minggu.
  - Output: `notifications`
- `POST /chat/message`
  - Deskripsi: Mengirim pesan antara admin dan pengguna aplikasi.
  - Input: `{ from_user, to_user, message, timestamp }`
  - Output: `status`, `message_id`
- `GET /chat/history`
  - Deskripsi: Mengambil riwayat chat antara admin dan pengguna.
  - Input: `conversation_id`
  - Output: `messages`

#### 8.3 Routing dan Forwarding
- `POST /gateway/payment`
  - Deskripsi: Menerima semua payment request dari Marketplace, POS, SupplierHub, atau LogistiKita.
  - Input: `{ from_app, from_user, to_user, amount, metadata, service_type }`
  - Proses: validasi token → validasi payload → log request → forward ke SmartBank.
  - Output: `status`, `transaction_id`, `message`.

- `POST /gateway/smartbank`
  - Deskripsi: Jalur umum untuk request ke SmartBank jika aplikasi lain memerlukan operasi khusus.
  - Input: `{ action, payload }`
  - Proses: validasi token → log → forward ke endpoint SmartBank sesuai `action`.
  - Output: `status`, `data`.

- `POST /gateway/marketplace`
  - Deskripsi: Meneruskan request Marketplace selain payment, misal katalog atau order status.
  - Input: `{ action, payload }`
  - Proses: validasi → log → forward ke Marketplace.
  - Output: `status`, `data`.

- `POST /gateway/logistics`
  - Deskripsi: Menerima permintaan pengiriman dari Marketplace/SupplierHub.
  - Input: `{ order_id, address, distance, shipping_type }`
  - Proses: validasi → log → forward ke LogistiKita.
  - Output: `status`, `delivery_id`, `message`.

- `POST /gateway/supplier`
  - Deskripsi: Meneruskan request SupplierHub untuk order bahan atau pembayaran supplier.
  - Input: `{ supplier_id, material, qty, total_cost }`
  - Proses: validasi → log → forward ke SupplierHub / SmartBank.
  - Output: `status`, `order_id`, `message`.

#### 8.4 Response Standar
- `200 OK` jika berhasil dengan body:
  - `status: "success"`
  - `data: {...}`
- `4xx` untuk error validasi
- `5xx` untuk error internal gateway

### 9. Alur Integrasi dengan SmartBank
1. Aplikasi seperti Marketplace/POS/SupplierHub/LogistiKita membuat request transaksi.
2. Request dikirim ke API Gateway dengan header Authorization.
3. Gateway memvalidasi JWT dan payload.
4. Gateway men-logging request, kemudian meneruskan ke SmartBank atau tujuan lain.
5. SmartBank memproses transaksi, mencatat ledger, memotong fee, dan mengembalikan hasil.
6. Gateway mengembalikan respon yang dibersihkan ke aplikasi pemanggil.

### 10. Kebutuhan Data / Payload
- `from_app`: nama aplikasi pemanggil
- `from_user`: ID pengguna pemanggil
- `to_user`: ID penerima atau layanan
- `amount`: nominal transaksi
- `metadata`: informasi order, produk, atau layanan
- `service_type`: kategori, misal `payment`, `loan`, `logistics`

### 11. Acceptance Criteria
- Semua request antar aplikasi melewati `API Gateway`.
- Gateway menolak request tanpa token.
- Gateway mem-forward request pembayaran ke SmartBank.
- Gateway melakukan logging minimal untuk setiap request.
- Response ke klien konsisten dalam format JSON.
- Landing page tersedia dan menjelaskan layanan API Integrator.
- Login role-based tersedia untuk admin gateway, user aplikasi, dan monitoring UMKM.
- Dashboard admin dan dashboard user tersedia secara terpisah.
- Dashboard admin dan user memiliki grafik analitik.
- UMKM Insight hanya dapat mengakses mode monitoring/read-only.
- Sistem mengirim notifikasi jika suatu API tidak aktif selama 1 minggu.
- Admin dan user dapat berkomunikasi melalui fitur chat.

### 12. Risiko dan Asumsi
- Asumsi: SmartBank dan aplikasi lain akan mematuhi kontrak endpoint.
- Risiko: jika aplikasi langsung mengakses SmartBank tanpa lewat gateway, maka integritas sistem terganggu.
- Risiko: validasi JWT harus diterapkan dengan benar untuk mencegah bypass.

### 13. Rekomendasi Pengembangan
- Gunakan framework web yang mendukung middleware JWT (mis. Express + middleware di Node.js atau Flask + decorator di Python).
- Dokumentasikan endpoint gateway dengan jelas.
- Pisahkan logika validasi dan forwarding agar mudah diuji.
- Jalankan pengujian integrasi untuk skenario Marketplace → Gateway → SmartBank.

---

Dokumen ini dibuat berdasarkan data di folder `dokumen_tugas` dan menitikberatkan pada peran API Integrator sebagai pengatur lalu lintas request antar sistem dan pelindung integrasi SmartBank.