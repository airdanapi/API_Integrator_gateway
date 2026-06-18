# Dokumentasi API Integrator Gateway

Folder ini menjadi pusat dokumentasi produk, perencanaan, pengembangan, dan
pelaporan proyek API Integrator Gateway.

## Indeks

| Kategori | Isi |
| --- | --- |
| [Requirements](requirements/PRD_API_Integrator.md) | Product Requirements Document dan kontrak utama sistem. |
| [Planning](planning/SPRINT_IMPLEMENTATION_PLAN.md) | Roadmap implementasi untuk Sprint 1 sampai Sprint 12. |
| [Development](development/DOKUMENTASI_PENGEMBANGAN_APLIKASI.md) | Panduan arsitektur, teknologi, struktur proyek, API, dan pengembangan. |
| [Architecture](architecture/ALUR_KERJA.png) | Diagram alur kerja ekosistem aplikasi. |
| [Project Management](project-management/GITHUB_PROJECT_ISSUES.md) | Backlog issue, pembagian sprint, dan assignee GitHub Project `RPL 2`. |
| [Sprint 1 Report](report/SPRINT_01_REPORT.md) | Laporan hasil Project Setup & Infrastructure pada 18 Juni 2026. |
| [Sprint 2 Report](report/SPRINT_02_REPORT.md) | Laporan landing page, endpoint publik, TDD, dan verifikasi responsif pada 18 Juni 2026. |
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
- Data sumber ditempatkan di `misc/source-data/`; kesimpulan atau hasil
  pengolahannya ditempatkan pada kategori dokumen yang relevan.
