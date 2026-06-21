# Pre-Sprint 9 Hardening: Baseline Frontend dan Verifikasi

## Metadata

| Atribut | Nilai |
| --- | --- |
| Proyek | API Integrator Gateway |
| Checkpoint | Pre-Sprint 9 Baseline Hardening |
| Tanggal laporan | 21 Juni 2026 |
| Status | Selesai |

## Ringkasan

Checkpoint ini membersihkan baseline sebelum masuk Sprint 9 Notifications.
Roadmap tidak digeser: Sprint 9 tetap berfokus pada sistem notifikasi.

Tidak ada perubahan fitur, API backend, route, schema database, image Docker,
atau volume data. Perubahan hanya menjaga lint frontend, stabilitas Vitest, dan
dokumentasi verifikasi.

## Baseline RED

Sebelum hardening, `docker compose exec -T frontend npm run lint` gagal dengan
4 error:

- `react-hooks/set-state-in-effect` pada initial load dashboard admin.
- `react-hooks/set-state-in-effect` pada initial load dashboard user.
- `react-hooks/set-state-in-effect` pada initial load dashboard monitoring.
- `no-unused-vars` untuk import `beforeEach` yang tidak dipakai.

Full frontend test pernah gagal timeout pada test landing page dengan default
5 detik. Saat checkpoint ini direproduksi, suite sempat lulus, tetapi test
landing tetap menjadi jalur paling berat sehingga timeout global dinaikkan ke
10 detik untuk mengurangi flakiness tanpa mengubah script test.

## Perubahan

- Initial load dashboard admin, user, dan monitoring dijadwalkan melalui
  `setTimeout(..., 0)` serta dibersihkan pada cleanup effect.
- Polling dashboard tetap berjalan setiap 30 detik dengan `setInterval`.
- Import `beforeEach` yang tidak dipakai dihapus dari test dashboard admin.
- `testTimeout: 10000` ditambahkan pada konfigurasi Vitest.
- Dokumentasi checkpoint ditautkan dari indeks `docs/README.md`.

## Hasil Verifikasi

Semua command di bawah lulus setelah perubahan. Karena service frontend tidak
menggunakan bind mount source host, image frontend direbuild sebelum checks
berbasis `docker compose exec`.

```powershell
docker compose build frontend
docker compose up --detach frontend
docker compose exec -T frontend npm run lint
docker compose exec -T frontend npm test
docker compose exec -T frontend npm run build
docker run --rm --workdir /src --volume "<repo>/backend:/src" golang:1.26.4-alpine3.24 go test ./...
docker run --rm --workdir /src --volume "<repo>/backend:/src" golang:1.26.4-alpine3.24 go vet ./...
docker run --rm --workdir /src --volume "<repo>/backend:/src" golang:1.26.4-alpine3.24 go build -o /tmp/api-integrator-server ./cmd/server
docker compose config --quiet
docker compose up --detach
docker compose ps
```

Smoke test host:

- `GET http://localhost:8080/health`
- `GET http://localhost:8080/landing`
- Login seed admin: `admin` / `admin-development-password` / `API Gateway`

## Catatan Docker

Tidak ada upgrade image dan tidak ada backup database pada checkpoint ini karena
tidak ada perubahan Dockerfile, image tag, migration, schema, atau volume.
Environment tetap menggunakan container karena versi Node dan Go lokal belum
sesuai kebutuhan README.