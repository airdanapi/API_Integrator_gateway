# Instruksi Agent API Integrator Gateway

Instruksi ini berlaku untuk seluruh repository dan wajib dibaca sebelum
mengerjakan sprint berikutnya.

## Metode Pengembangan Wajib

- Gunakan Test-Driven Development: tulis test, buktikan fase RED, implementasikan
  perubahan minimum sampai GREEN, lalu lakukan REFACTOR.
- Jangan menandai sprint selesai sebelum test, lint/vet, build, dan acceptance
  test yang relevan lulus.
- Pertahankan perubahan lokal pengguna yang tidak terkait dan jangan menghapus
  data, volume, atau artifact pengguna tanpa persetujuan eksplisit.

## Checklist Docker Setelah Setiap Sprint

Setelah satu sprint selesai, agent wajib:

1. Audit versi Docker Engine, Docker Compose, dan seluruh base image dari sumber
   resmi.
2. Upgrade image bila tersedia versi stabil/LTS yang direkomendasikan dan
   kompatibel. Jika tidak di-upgrade, tulis alasan teknis dalam laporan sprint.
3. Buat backup sebelum upgrade database atau perubahan volume. Jangan pernah
   menjalankan `docker compose down --volumes`, menghapus named volume, atau
   mengganti data persisten tanpa persetujuan eksplisit.
4. Jalankan:

   ```powershell
   docker compose config --quiet
   docker compose build --pull
   docker compose up --detach
   docker compose ps
   ```

5. Pastikan seluruh service healthy dan jalankan smoke test endpoint/alur utama
   yang berubah pada sprint tersebut.
6. Catat versi image, nama backup bila ada, hasil build, healthcheck, smoke test,
   dan kendala Docker di `docs/report/SPRINT_NN_REPORT.md`.

## Dokumentasi Sprint

- Tambahkan laporan baru menggunakan format `SPRINT_NN_REPORT.md`.
- Tautkan laporan dari `docs/README.md`.
- Perbarui status roadmap dan README bila setup, environment, API contract, atau
  proses operasional berubah.
