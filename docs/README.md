# Dokumentasi API Integrator Gateway

Folder ini menjadi pusat dokumentasi produk, perencanaan, pengembangan, dan
pelaporan proyek API Integrator Gateway.

## Indeks

| Kategori | Isi |
| --- | --- |
| [Requirements](requirements/PRD_API_Integrator.md) | Product Requirements Document dan kontrak utama sistem. |
| [Planning](planning/SPRINT_IMPLEMENTATION_PLAN.md) | Roadmap implementasi untuk Sprint 1 sampai Sprint 12. |
| [Development](development/DOKUMENTASI_PENGEMBANGAN_APLIKASI.md) | Panduan arsitektur, teknologi, struktur proyek, API, dan pengembangan. |
| [Database Seeder](development/DATABASE_SEEDER.md) | Panduan penggunaan seeder untuk tabel Sprint 5: request_logs, notifications, chat_messages, dashboard_data. |
| [Architecture](architecture/ALUR_KERJA.png) | Diagram alur kerja ekosistem aplikasi. |
| [Project Management](project-management/GITHUB_PROJECT_ISSUES.md) | Backlog issue, pembagian sprint, dan assignee GitHub Project `RPL 2`. |
| [Sprint 1 Report](report/SPRINT_01_REPORT.md) | Laporan hasil Project Setup & Infrastructure pada 18 Juni 2026. |
| [Sprint 2 Report](report/SPRINT_02_REPORT.md) | Laporan landing page, endpoint publik, TDD, dan verifikasi responsif pada 18 Juni 2026. |
| [Sprint 3 Report](report/SPRINT_03_REPORT.md) | Laporan autentikasi backend, JWT, migration MySQL, TDD, dan upgrade Docker pada 18 Juni 2026. |
| [Sprint 4 Report](report/SPRINT_04_REPORT.md) | Laporan login frontend, session management, protected route, TDD, Docker, dan smoke test pada 18 Juni 2026. |
| [Sprint 5 Report](report/SPRINT_05_REPORT.md) | Laporan database schema Sprint 5: migrasi, model, repository layer, seeder, dan TDD pada 18 Juni 2026. |
| [Sprint 6 Report](report/SPRINT_06_REPORT.md) | Laporan Dashboard Admin backend: endpoint GET /dashboard/admin, requireRole middleware, service layer, TDD pada 19 Juni 2026. |
| [Sprint 7 Report](report/SPRINT_07_REPORT.md) | Laporan Dashboard Admin frontend: AdminDashboardPage, traffic cards, service indicators, audit table, polling 30s, fix jsdom 28 pada 19 Juni 2026. |
| [Sprint 8 Report](report/SPRINT_08_REPORT.md) | Laporan Dashboard User & Monitoring: endpoint /dashboard/user dan /dashboard/monitoring, UserDashboardPage, MonitoringDashboardPage, RBAC, TDD 48 test pada 19 Juni 2026. |
| [Pre-Sprint 9 Hardening](report/PRE_SPRINT_09_HARDENING.md) | Catatan baseline hardening frontend sebelum Sprint 9 Notifications: lint, stabilitas Vitest, build, backend regression, dan smoke test pada 21 Juni 2026. |
| [Sprint 9 Report](report/SPRINT_09_REPORT.md) | Laporan Notifications System: backend notification API, scheduler alert, frontend bell/dropdown, mark-read, TDD, Docker checklist, dan smoke test pada 21 Juni 2026. |
| [Sprint 10 Report](report/SPRINT_10_REPORT.md) | Laporan Chat System: backend chat API, service, repository, frontend ChatDrawer, polling 2s, RBAC, TDD, dan verifikasi pada 21 Juni 2026. |
| [Sprint 11 Report](report/SPRINT_11_REPORT.md) | Laporan Gateway Routing & API Integration: endpoint gateway, forwarder HTTP, validasi payload, request logging, frontend service, TDD, Docker checklist, dan smoke test pada 21 Juni 2026. |
| [Miscellaneous Source Data](misc/source-data/) | Data CSV sumber tugas besar yang dipertahankan sebagai referensi. |

## Struktur

```text
docs/
├── README.md
├── architecture/       # Diagram dan dokumentasi arsitektur
├── development/        # Panduan teknis pengembangan
├── planning/           # Roadmap dan sprint implementation plan
├── project-management/ # Backlog dan administrasi proyek
├── requirements/       # Kebutuhan dan kontrak produk
├── report/             # Laporan pelaksanaan setiap sprint
└── misc/
    └── source-data/    # Data sumber yang tidak diedit sebagai dokumen utama
```

## Konvensi

- Folder menggunakan nama lowercase dengan pemisah tanda hubung.
- Laporan sprint menggunakan format `SPRINT_NN_REPORT.md`, misalnya
  `SPRINT_02_REPORT.md`.
- Dokumen baru harus ditautkan dari indeks ini.
- Aturan agent dan checklist Docker per sprint tersedia pada
  [`AGENTS.md`](../AGENTS.md).
- Data sumber ditempatkan di `misc/source-data/`; kesimpulan atau hasil
  pengolahannya ditempatkan pada kategori dokumen yang relevan.
