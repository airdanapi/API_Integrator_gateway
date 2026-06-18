# Sprint 02 Report: Landing Page & Static Content

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Sprint | Sprint 2 — Landing Page & Static Content |
| Tanggal laporan | 18 Juni 2026 |
| PIC | `venalism` |
| Commit implementasi | `93dd9e8` |
| Status | Selesai |

## Ringkasan Eksekutif

Sprint 2 menghasilkan landing page publik berbahasa Indonesia dengan gaya
light corporate, layout responsif, navigasi desktop dan mobile, konten manfaat,
alur integrasi, peran aplikasi, use case, FAQ, CTA, serta footer. Seluruh visual
menggunakan Tailwind CSS, CSS lokal, dan inline SVG tanpa font atau aset jaringan
eksternal.

Backend menyediakan endpoint publik `GET /landing` sesuai kontrak PRD. Frontend
tetap merender konten statis agar landing page tidak bergantung pada
ketersediaan backend. CTA login menampilkan status `Segera hadir` tanpa
menghasilkan dead route.

Implementasi dilakukan dengan Test-Driven Development. Fase RED dibuktikan oleh
7 test frontend yang gagal dan response `404` dari test endpoint `/landing`.
Setelah GREEN dan REFACTOR, 9 test frontend dan 5 test backend lulus bersama
lint, vet, production build, dan Docker image build.

## Tujuan dan Ruang Lingkup

Tujuan Sprint 2 adalah menyediakan portal publik yang menjelaskan fungsi API
Integrator dan hubungan layanan di dalam ekosistem UMKM.

Ruang lingkup yang diselesaikan:

- Header, footer, desktop navigation, dan mobile navigation.
- Hero, benefits, integration flow, application roles, use cases, FAQ, dan CTA.
- Desain light corporate yang responsif pada mobile, tablet, dan desktop.
- Semantic HTML, skip link, keyboard focus, ARIA state, dan native disclosure.
- Inline SVG dan CSS tanpa dependency aset eksternal.
- Endpoint publik `GET /landing`.
- Verifikasi CORS untuk konsumsi lintas origin.
- Test komponen dan kontrak HTTP.

Di luar ruang lingkup:

- Halaman login, React Router, autentikasi, dan session management.
- Database, migration, dan persistence konten.
- Dashboard dan gateway forwarding.
- Pengambilan konten landing oleh frontend saat runtime.

## Hasil Implementasi

### Frontend

| Bagian | Implementasi |
| --- | --- |
| Layout | Header sticky, main content, CTA, dan footer |
| Hero | Penjelasan layanan, statistik, dan diagram aplikasi |
| Benefits | Keamanan terpusat, routing konsisten, operasional terpantau |
| Integration flow | Empat tahap request hingga standardized response |
| Application roles | Tujuh node aplikasi sesuai PRD |
| Use cases | Marketplace checkout, POS payment, supplier order |
| FAQ | Native `details/summary` yang keyboard-accessible |
| CTA | Repositori resmi dan login berstatus `Segera hadir` |
| Assets | Inline SVG dan CSS lokal, tanpa request eksternal |

Konten frontend dipisahkan dari komponen sehingga dapat dipelihara tanpa
mengubah struktur presentasi. Landing page tidak melakukan fetch ke backend;
keputusan ini menjaga halaman tetap tersedia bila API sedang tidak aktif.

### Backend

Endpoint baru:

```http
GET /landing
```

Endpoint tidak membutuhkan header `Authorization` dan menghasilkan HTTP `200`.
Handler dipisahkan dari app factory agar kontrak landing tidak bercampur dengan
registrasi middleware dan route lain.

Struktur response:

```json
{
  "status": "success",
  "data": {
    "service_overview": {
      "name": "API Integrator Gateway",
      "tagline": "Satu pintu aman untuk setiap komunikasi antar aplikasi.",
      "description": "...",
      "benefits": [
        {
          "title": "Keamanan terpusat",
          "description": "..."
        }
      ]
    },
    "application_roles": [
      {
        "application": "SmartBank",
        "role": "Core keuangan",
        "interaction": "..."
      }
    ],
    "integration_flow": [
      {
        "step": 1,
        "title": "Aplikasi mengirim request",
        "description": "..."
      }
    ],
    "contact_info": {
      "repository_url": "https://github.com/airdanapi/API_Integrator_gateway",
      "login_status": "coming_soon"
    }
  }
}
```

`application_roles` berisi SmartBank, Marketplace, POS, SupplierHub,
LogistiKita, UMKM Insight, dan API Gateway.

## Pelaksanaan TDD

### RED

Test ditulis sebelum implementasi untuk:

- Seluruh section landing page.
- Anchor navigation.
- State dan atribut ARIA mobile menu.
- CTA login tanpa dead route.
- FAQ dan tujuh nama aplikasi.
- Link repositori resmi.
- Kontrak response `GET /landing`.
- Akses endpoint tanpa token.
- Header CORS.

Hasil awal:

- Frontend: 7 test baru gagal karena shell Sprint 1 belum memiliki section dan
  interaction Sprint 2.
- Backend: `GET /landing` menghasilkan HTTP `404`.

### GREEN

Implementasi minimum ditambahkan sampai:

- 2 file test frontend dengan 9 test case lulus.
- 2 package test backend dengan total 5 test case lulus.
- Endpoint runtime menghasilkan HTTP `200`, CORS `*`, tujuh application roles,
  dan `login_status: "coming_soon"`.

### REFACTOR

Setelah test lulus:

- Komponen dipisahkan berdasarkan section landing page.
- Konten statis dipindahkan ke module data terpusat.
- Handler dan response type landing dipisahkan ke file backend khusus.
- Cleanup Testing Library dibuat eksplisit untuk isolasi setiap test.
- Breakpoint navigasi diubah menjadi `1024px` setelah visual QA menunjukkan
  header tablet terlalu padat.
- Rule ESLint `react/prop-types` dinonaktifkan karena proyek tidak menggunakan
  runtime PropTypes dan kontrak perilaku dijaga oleh component test.

## Responsivitas dan Performa

Production preview diuji menggunakan Microsoft Edge headless `149.0.4022.69`.
Browser in-app tidak tersedia pada sesi implementasi, sehingga verifikasi
dijalankan melalui Chrome DevTools Protocol pada browser lokal.

| Viewport | Overflow horizontal | Navigasi | Navigation Timing |
| --- | --- | --- | --- |
| `375 × 812` | Tidak ada | Mobile menu, open state lulus | `67.6 ms` |
| `768 × 1024` | Tidak ada | Compact mobile navigation | `59.1 ms` |
| `1440 × 900` | Tidak ada | Desktop navigation | `50.6 ms` |

Temuan tambahan:

- Elemen pertama saat menekan `Tab` adalah `Lewati ke konten utama`.
- Mobile menu berubah ke `aria-expanded="true"` dan menampilkan 5 link.
- Bahasa dokumen adalah `id`.
- Tidak ada resource yang dimuat dari origin eksternal.
- Seluruh hasil berada jauh di bawah acceptance criteria 2 detik.

### Screenshot Verifikasi

| Mobile | Tablet | Desktop |
| --- | --- | --- |
| [375 × 812](screenshots/sprint2-mobile-375x812.png) | [768 × 1024](screenshots/sprint2-tablet-768x1024.png) | [1440 × 900](screenshots/sprint2-desktop-1440x900.png) |

## Ukuran Production Bundle

| Artifact | Ukuran | Gzip |
| --- | --- | --- |
| `dist/index.html` | `0.66 kB` | `0.39 kB` |
| CSS | `29.79 kB` | `6.42 kB` |
| JavaScript | `213.57 kB` | `65.99 kB` |

Tidak ada image, font, atau library UI eksternal yang perlu diunduh saat
landing page dimuat.

## Deliverables

| Deliverable | Status | Catatan |
| --- | --- | --- |
| Landing page publik | Selesai | Dapat diakses tanpa autentikasi. |
| Header, footer, dan navigation | Selesai | Desktop dan mobile navigation tersedia. |
| Hero dan benefits | Selesai | Menjelaskan peran dan manfaat gateway. |
| Integration flow | Selesai | Empat tahap integrasi dan tujuh application roles. |
| Use cases dan FAQ | Selesai | Menggunakan skenario faktual dari dokumen proyek. |
| CTA login | Selesai | Nonaktif dengan status `Segera hadir`. |
| Responsive design | Selesai | Diverifikasi pada tiga viewport. |
| Endpoint `GET /landing` | Selesai | Public JSON contract dan CORS tersedia. |
| Test otomatis | Selesai | 9 frontend dan 5 backend test lulus. |
| Screenshot verifikasi | Selesai | Tersedia dalam folder report. |

## Acceptance Criteria

| Acceptance Criteria | Status | Bukti |
| --- | --- | --- |
| Landing page dapat diakses tanpa login | Lulus | Frontend tidak memiliki auth guard; endpoint `/landing` menerima request tanpa token. |
| Responsif pada desktop, tablet, dan mobile | Lulus | Tidak ada horizontal overflow pada ketiga viewport. |
| Mobile-friendly design | Lulus | Mobile menu, touch-sized control, dan single-column content tersedia. |
| Loading time kurang dari 2 detik | Lulus | Navigation Timing maksimum `67.6 ms` pada audit lokal. |
| Assets siap | Lulus | Seluruh ikon menggunakan inline SVG; tidak ada resource eksternal. |
| CORS tersedia | Lulus | Runtime dan test menghasilkan `Access-Control-Allow-Origin: *`. |

## Bukti Verifikasi

Perintah quality gate:

```powershell
# Frontend
npm.cmd test
npm.cmd run lint
npm.cmd run build

# Backend
go test ./...
go vet ./...
go build ./cmd/server

# Docker
docker compose config --quiet
docker compose build
```

Hasil:

- Frontend: 2 test files, 9 test cases, seluruhnya lulus.
- Backend: 5 test cases, seluruhnya lulus.
- ESLint lulus tanpa error atau warning.
- Vite production build berhasil.
- Go vet dan build berhasil.
- Docker Compose config valid.
- Image frontend dan backend berhasil dibangun.
- Runtime `GET /landing`: HTTP `200`, CORS `*`, tujuh application roles.

## Kendala dan Penyelesaian

### Browser In-App Tidak Tersedia

Browser in-app tidak terhubung pada sesi implementasi. Verifikasi dipindahkan ke
Microsoft Edge headless lokal melalui Chrome DevTools Protocol. Audit tetap
menggunakan production build dan viewport yang ditetapkan.

### Isolasi Sandbox Preview

Browser headless dan preview awal berjalan pada konteks proses berbeda sehingga
localhost sempat ditolak. Production preview kemudian dijalankan pada konteks
yang sama dengan browser audit. Tidak ada perubahan konfigurasi aplikasi yang
diperlukan.

### Header Tablet Terlalu Padat

Visual QA pertama menunjukkan desktop navigation terlalu padat pada lebar
`768px`. Breakpoint desktop navigation dinaikkan dari `768px` ke `1024px`,
kemudian seluruh test, build, dan audit viewport diulang.

## Risiko Tersisa

- Nilai Navigation Timing berasal dari mesin lokal dan perlu diukur ulang pada
  environment deployment dengan latency jaringan nyata.
- Konten frontend dan endpoint backend dikelola terpisah; perubahan copy perlu
  dijaga konsisten sampai tersedia content management atau backend-driven UI.
- CTA login masih nonaktif sampai Sprint 4 menyelesaikan authentication
  frontend.

## Handoff ke Sprint 3

Sprint 3 dapat melanjutkan authentication backend dari baseline ini:

1. Pertahankan `/health` dan `/landing` sebagai public routes.
2. Tambahkan schema user, password hashing, JWT service, dan `POST /auth/login`.
3. Terapkan auth middleware hanya pada endpoint terlindungi.
4. Pertahankan response JSON yang konsisten dengan endpoint saat ini.
5. Tambahkan test RED untuk credential validation, token claims, dan protected route sebelum implementasi.