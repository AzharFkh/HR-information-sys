# HR Information System (HRIS)

REST API sederhana untuk manajemen pengajuan cuti karyawan menggunakan Go, Gin, GORM, dan PostgreSQL.

## Fitur Utama

- Autentikasi user dengan JWT (`register` dan `login`)
- Role-Based Access Control (RBAC): `admin`, `hr`, dan `employee`
- Pengajuan cuti oleh karyawan serta riwayat pengajuan
- Persetujuan atau penolakan pengajuan cuti oleh HR (`APPROVED` / `REJECTED`)
- Pemotongan kuota cuti otomatis sesuai jenis cuti
- Database auto-migration dan seeder data awal

## Tech Stack

- Go 1.26
- Gin Web Framework
- GORM + PostgreSQL
- JWT (`golang-jwt/jwt/v5`) & bcrypt
- Testify (unit testing)
- Docker & Docker Compose

## Struktur Proyek

```text
.
├── cmd/
│   ├── api/        # Entrypoint API server
│   ├── migrate/    # Database migration & seeding
│   └── seed/       # Data awal (roles, users, leave types)
├── internal/
│   ├── config/     # Konfigurasi env dan koneksi database
│   ├── dto/        # Request dan response struct
│   ├── handler/    # HTTP handler
│   ├── middleware/ # Auth JWT & RBAC middleware
│   ├── mocks/      # Mock repository untuk testing
│   ├── models/     # Model database GORM
│   ├── repository/ # Query layer ke database
│   ├── router/     # Definisi routing API
│   ├── service/    # Business logic & unit test
│   └── utils/      # Helper JWT, context, dan error
├── docker-compose.yml
├── dockerfile
├── go.mod
└── Makefile
```

## Konfigurasi Environment

Buat file `.env` di root proyek dengan variabel berikut:

```env
DB_DSN=host=localhost user=postgres password=postgres dbname=hris port=5432 sslmode=disable
JWT_KEY=your_secret_key
```

## Menjalankan Aplikasi

### Secara Lokal

1. Pastikan PostgreSQL sudah berjalan dan database sudah dibuat.
2. Jalankan migrasi dan seeder:

```bash
make migrate
```

3. Jalankan API server (berjalan di port `7100`):

```bash
make run
```

### Menggunakan Docker

```bash
make docker-up
```

Container akan otomatis menjalankan migrasi sebelum menyalakan server pada `127.0.0.1:7100`. Untuk menghentikan container:

```bash
make docker-down
```

### Perintah Makefile Lainnya

| Perintah | Keterangan |
| --- | --- |
| `make run` | Menjalankan API server |
| `make migrate` | Menjalankan migrasi database dan seeding |
| `make seed` | Menjalankan seeder secara terpisah |
| `make test` | Menjalankan unit test |
| `make fmt` | Memformat kode Go (`gofmt`) |
| `make build` | Build seluruh package |

## Akun Default (Seeder)

Setelah menjalankan migrasi/seeder, akun berikut tersedia dengan password default `1234`:

| Username | Role |
| --- | --- |
| `budi_admin` | `admin` |
| `andi_hr` | `hr` |
| `bagus_employee` | `employee` |

## Daftar Endpoint

Base URL: `http://localhost:7100`

### Public

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| `GET` | `/ping` | Health check |
| `POST` | `/register` | Registrasi akun baru |
| `POST` | `/login` | Login untuk mendapatkan JWT token |

### Protected (Membutuhkan Header `Authorization: Bearer <token>`)

| Method | Endpoint | Role | Keterangan |
| --- | --- | --- | --- |
| `POST` | `/leaves` | Authenticated | Mengajukan permohonan cuti |
| `GET` | `/leaves` | Authenticated | Melihat daftar pengajuan cuti |
| `GET` | `/hr/pending` | `hr` | Melihat daftar pengajuan cuti berstatus pending |
| `PATCH` | `/hr/update-status/:id` | `hr` | Mengubah status cuti (`APPROVED` / `REJECTED`) |
