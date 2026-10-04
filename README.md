# ALUSI Backend API (Akses Layanan Unggulan & Sistem Terintegrasi)

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin Framework](https://img.shields.io/badge/Framework-Gin_v1.12-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL_16+-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![Swagger OpenAPI](https://img.shields.io/badge/Docs-Swagger_OpenAPI_2.0-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**Backend RESTful API & Realtime Streaming Engine untuk Portal Aplikasi Terintegrasi Badan Pusat Statistik (BPS) Provinsi Sumatera Utara.**

[Dokumentasi API (Swagger)](#-dokumentasi-api-interaktif-swagger) • [Fitur Utama](#-fitur-utama) • [Arsitektur Sistem](#-arsitektur--struktur-proyek) • [Panduan Instalasi](#-panduan-instalasi--menjalankan-server) • [Skema Database](#-skema-database--migrasi)

</div>

---

## 📖 Tentang ALUSI

**ALUSI** adalah portal terpusat (*Single Portal Gateway*) bagi seluruh aplikasi internal dan layanan publik di lingkungan Badan Pusat Statistik Provinsi Sumatera Utara. Backend ALUSI dibangun dengan performa tinggi menggunakan bahasa **Go (Golang)**, arsitektur **Clean Architecture**, integrasi autentikasi **SSO BPS Sumut**, pemantauan kesehatan layanan secara **Realtime Server-Sent Events (SSE)**, optimasi aset logo otomatis, asisten rekomendasi cerdas berbasis AI, serta pelaporan eksekutif terstandar Satu Data Indonesia.

---

## ✨ Fitur Utama

### 🔐 1. Autentikasi Terpusat SSO BPS Sumut & JWT Session
* Integrasi OAuth 2.0 / OpenID Connect (OIDC) dengan server SSO BPS Sumut (`/authorize`, `/token`, `/userinfo`, `/logout`).
* *Role-Based Access Control* (RBAC) bertingkat: `admin`, `pimpinan`, `user` (Pegawai BPS & Mitra Statistik).
* Token sesi JWT aman dengan *claims* identitas NIP, Satuan Kerja (Satker), dan peran akses.

### 📱 2. Katalog Aplikasi & Personalisasi Pengguna
* Pengelompokan aplikasi berdasarkan Kategori, Target Pengguna, dan Hak Akses Role.
* *Fuzzy Trigram Search* dengan PostgreSQL pg_trgm untuk pencarian aplikasi cepat dan toleran kesalahan ketik.
* Personalisasi: Aplikasi Favorit (*Favorite/Pin*), Riwayat Akses Terakhir (*Recent Apps*), dan Panduan Penggunaan/FAQ per aplikasi.
* *Asynchronous Click Tracking Worker* berbasis Go Channel & Worker Pool untuk mencatat metrik kunjungan tanpa membebani latensi user.

### 📡 3. Realtime Service Health Monitoring (Server-Sent Events)
* *Automated Health Probe Worker* yang memantau ketersediaan (*uptime/downtime/latency*) seluruh aplikasi setiap 5 menit secara non-blocking.
* *Streaming Server-Sent Events* (`GET /api/v1/services/realtime-status`) yang mem-broadcast perubahan status secara instan ke frontend tanpa *polling*.
* Rekapitulasi SLA Uptime dan latensi rata-rata 30 hari.

### 🖼️ 4. Optimasi & Kompresi Logo Otomatis (*Image Pipeline*)
* Endpoint Upload Logo Admin (`POST /api/v1/admin/apps/:id/logo`) dengan batas maksimal 2 MB dan verifikasi *Magic Bytes*.
* **Vektor SVG**: Sanitasi otomatis untuk mencegah risiko *Stored XSS* (<script>, inline event handlers, XXE).
* **Raster (PNG/JPG/JPEG/WebP)**: Resize proporsional otomatis (maksimal 512×512 px) dan kompresi level tinggi (`BestCompression`), menghemat penyimpanan hingga **85%–95%**.
* Static asset delivery dengan *caching header* optimal (`/uploads/logos/*`).

### 🤖 5. Asisten Rekomendasi Berbasis AI
* Pemrosesan bahasa alami (*natural language query*) untuk merekomendasikan aplikasi BPS yang paling sesuai dengan kebutuhan pengguna.
* Pencocokan semantik cerdas terhadap nama, deskripsi, panduan, dan kategori aplikasi disertai alasan rekomendasi kontekstual.

### 📊 6. Pelaporan Eksekutif & Metadata Satu Data Indonesia
* Ringkasan analitik pimpinan: Aplikasi terpopuler, tren klik berkala, dan log gangguan layanan (*disruptions*).
* Ekspor laporan eksekutif format CSV (Katalog & Metrik Analitik).
* Endpoint Open API Metadata Katalog (`GET /api/v1/openapi/apps`) sesuai standar Satu Data Indonesia.

### 🛡️ 7. Audit Logging & Keamanan Berlapis
* Pencatatan log jejak audit administratif pada setiap operasi CRUD aplikasi, kategori, dan panduan.
* *Token Bucket IP Rate Limiter* (100 req/sec dengan burst 50 per IP).
* Security Headers: HSTS, X-Frame-Options (Clickjacking Protection), X-Content-Type-Options, CSP, dan XSS Protection.

---

## 🏛️ Arsitektur & Struktur Proyek

Proyek ini menerapkan prinsip **Clean Architecture & Idiomatic Go Patterns**:

```
Backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entrypoint server, dependency wiring, graceful shutdown
├── config/
│   └── config.go               # Viper configuration loader (env / config.yaml)
├── docs/                       # Swagger / OpenAPI documentation auto-generated
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── delivery/
│   │   └── http/               # Gin HTTP Handlers, Routing, & Middlewares
│   │       ├── middleware/     # Auth, RateLimiter, CORS, SecurityHeaders, Logger
│   │       ├── response/       # Standardized JSON response formatting
│   │       ├── admin_handler.go
│   │       ├── ai_handler.go
│   │       ├── analytics_handler.go
│   │       ├── announcement_handler.go
│   │       ├── audit_handler.go
│   │       ├── auth_handler.go
│   │       ├── catalog_handler.go
│   │       ├── feedback_handler.go
│   │       ├── health_handler.go
│   │       ├── interaction_handler.go
│   │       ├── monitoring_handler.go
│   │       ├── report_handler.go
│   │       └── router.go
│   ├── domain/                 # Core Entities, Structs, & Repository Interfaces
│   ├── repository/
│   │   └── postgres/           # PostgreSQL PGX v5 implementations with transactions
│   └── usecase/                # Business Logic & Application Use Cases
├── migrations/                 # PostgreSQL SQL Migration scripts (001 - 006)
├── pkg/
│   ├── ai/                     # Assistant recommendation semantic engine
│   ├── database/               # PGX pool connector & migration runner
│   ├── exporter/               # CSV / Spreadsheet report generator
│   ├── jwt/                    # JWT token signing & verification service
│   ├── logger/                 # Zerolog structured logging setup
│   ├── media/                  # Image optimizer, compressor, & SVG sanitizer
│   ├── realtime/               # Thread-safe Server-Sent Events (SSE) Hub
│   ├── sso/                    # SSO BPS Sumut OAuth2 / OIDC client
│   └── worker/                 # Click worker & Health probe background workers
├── uploads/                    # Directory for uploaded static assets (logos)
├── go.mod
├── go.sum
└── README.md
```

---

## 🛠️ Persyaratan Sistem & Tech Stack

* **Bahasa**: [Go (Golang)](https://go.dev/) versi **1.26+**
* **Web Framework**: [Gin Gonic](https://gin-gonic.com/) v1.12
* **Database**: [PostgreSQL](https://www.postgresql.org/) versi **16+** dengan ekstensi `pg_trgm` & `uuid-ossp`
* **Driver DB**: [pgx/v5](https://github.com/jackc/pgx) (High-Performance Connection Pool)
* **JWT**: [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
* **Dokumentasi API**: [swaggo/swag](https://github.com/swaggo/swag) v1.16+
* **Logging**: [rs/zerolog](https://github.com/rs/zerolog)

---

## 🚀 Panduan Instalasi & Menjalankan Server

### 1. Clone Repository
```bash
git clone https://github.com/restusatrio11/backend-alusi-go.git
cd backend-alusi-go
```

### 2. Konfigurasi Environment Variable
Salin file `.env.example` menjadi `.env` atau buat file konfigurasi `config.yaml`:

```env
# Application Configuration
APP_NAME=ALUSI Backend API
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=db_alusi
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_MAX_IDLE_TIME=15m
DB_MAX_LIFETIME=1h

# JWT Authentication
JWT_SECRET=super-secret-key-must-be-at-least-32-characters-long
JWT_EXPIRATION_HOURS=24

# SSO BPS Sumut (OAuth 2.0 / OIDC)
SSO_CLIENT_ID=app_1791057531231_943f9ac1
SSO_CLIENT_SECRET=YOUR_SSO_CLIENT_SECRET
SSO_ISSUER_URL=https://otp-dev.bps.web.id
SSO_REDIRECT_URI=https://aron.bps.web.id/callback

# CORS Origins
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://alusi.sumut.bps.go.id
```

### 3. Setup Database PostgreSQL
Pastikan database telah dibuat. Migrasi SQL di folder `migrations/` akan **otomatis dijalankan** saat backend pertama kali dinyalakan:

```sql
CREATE DATABASE db_alusi;
```

### 4. Menjalankan Server
```bash
# Mengunduh dependensi
go mod download

# Menjalankan server dalam mode development
go run cmd/server/main.go
```

Server akan aktif dan siap melayani permintaan di `http://localhost:8080`.

---

## 📚 Dokumentasi API Interaktif (Swagger)

Backend ALUSI dilengkapi dengan dokumentasi OpenAPI 2.0 / Swagger UI interaktif:

* **Swagger Web UI**: [http://localhost:8080/docs](http://localhost:8080/docs) atau [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
* **Raw Swagger JSON**: [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json)

Untuk men-generate ulang Swagger docs setelah mengubah annotasi handler:
```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g main.go -d cmd/server,internal/delivery/http,internal/domain,internal/usecase -o docs --parseDependency --parseInternal
```

---

## 📋 Ringkasan Endpoint API

| Method | Endpoint | Keterangan | Autentikasi |
| :--- | :--- | :--- | :--- |
| `GET` | `/healthz` | Health check probe server & database pool | Publik |
| `GET` | `/api/v1/auth/login` | Inisiasi Login SSO Sumut (Redirect OAuth) | Publik |
| `GET` | `/api/v1/auth/callback` | Callback penukaran kode SSO & terbit JWT | Publik |
| `GET` | `/api/v1/auth/me` | Ambil profil pengguna & peran aktif | Bearer Auth |
| `GET` | `/api/v1/categories` | Daftar kategori aplikasi | Publik |
| `GET` | `/api/v1/apps` | Katalog aplikasi (tersaring sesuai role user) | Opsional Auth |
| `GET` | `/api/v1/apps/search?q=` | Pencarian cepat aplikasi (Trigram Similarity) | Opsional Auth |
| `GET` | `/api/v1/apps/:slug` | Detail lengkap aplikasi & panduan | Opsional Auth |
| `POST` | `/api/v1/apps/:id/click` | Rekam metrik kunjungan aplikasi | Opsional Auth |
| `POST` | `/api/v1/apps/:id/favorite` | Toggle pin aplikasi favorit | Bearer Auth |
| `GET` | `/api/v1/users/me/favorites` | Daftar aplikasi favorit pengguna | Bearer Auth |
| `GET` | `/api/v1/users/me/recents` | Riwayat aplikasi yang terakhir dibuka | Bearer Auth |
| `GET` | `/api/v1/services/status` | Rekapitulasi SLA uptime & latensi 30 hari | Publik |
| `GET` | `/api/v1/services/realtime-status` | **Streaming SSE Realtime Status Layanan** | Publik |
| `POST` | `/api/v1/ai/ask` | Rekomendasi aplikasi via asisten AI cerdas | Opsional Auth |
| `POST` | `/api/v1/feedbacks` | Kirim feedback / laporan bug aplikasi | Opsional Auth |
| `GET` | `/api/v1/openapi/apps` | Metadata katalog standar Satu Data Indonesia | Publik |
| `POST` | `/api/v1/admin/apps` | Tambah aplikasi baru ke katalog | Admin |
| `PUT` | `/api/v1/admin/apps/:id` | Perbarui data aplikasi | Admin |
| `DELETE` | `/api/v1/admin/apps/:id` | Nonaktifkan (Soft Delete) aplikasi | Admin |
| `POST` | `/api/v1/admin/apps/:id/logo` | **Upload & Kompresi Logo Aplikasi** | Admin |
| `POST` | `/api/v1/admin/apps/:id/probe` | Trigger manual health probe pada aplikasi | Admin |
| `PUT` | `/api/v1/admin/apps/reorder` | Simpan urutan posisi drag-and-drop aplikasi | Admin |
| `GET` | `/api/v1/admin/analytics/summary` | Dashboard ringkasan metrik analitik | Admin/Pimpinan |
| `GET` | `/api/v1/admin/reports/catalog/export` | Ekspor CSV Katalog Aplikasi | Admin/Pimpinan |
| `GET` | `/api/v1/admin/reports/analytics/export` | Ekspor CSV Analitik Penggunaan | Admin/Pimpinan |
| `GET` | `/api/v1/admin/audit-logs` | Inspeksi jejak log aktivitas admin | Admin |

---

## 🧪 Menjalankan Pengujian (Testing)

Proyek ini memiliki cakupan pengujian unit dan integrasi yang komprehensif:

```bash
# Menjalankan seluruh test suite
go test -v ./...

# Menjalankan test dengan laporan code coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 👥 Tim Pengembang & Kontak

* **Tim Pengembang**: Tim Inovasi & Pengolahan Data BPS Provinsi Sumatera Utara
* **Website Resmi**: [https://sumut.bps.go.id](https://sumut.bps.go.id)
* **Email Dinas**: `bps1200@bps.go.id`

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah lisensi [MIT License](LICENSE).
