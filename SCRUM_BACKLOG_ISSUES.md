# 🗂️ Product Backlog & Daftar Issue: Frontend Portal ALUSI BPS Sumut

Dokumen ini berisi seluruh daftar **Label / Tag**, **Epics**, dan **User Stories & Tasks** yang siap disalin/dibuat ke **GitHub Issues** atau **GitHub Projects**.

---

## 🏷️ 1. Standar Labels & Tags GitHub

Gunakan daftar tag warna standar berikut untuk kategorisasi issue di GitHub:

| Label | Warna Hex | Deskripsi |
|---|---|---|
| `epic` | `#3E4B9B` | Inisiatif fitur tingkat tinggi / modul besar |
| `sprint-1` | `#0E8A16` | Scope pengerjaan Sprint 1: Core Shell & Katalog |
| `sprint-2` | `#1D76DB` | Scope pengerjaan Sprint 2: AI & Personalisasi SSO |
| `sprint-3` | `#FBCA04` | Scope pengerjaan Sprint 3: Realtime SSE & Helpdesk |
| `sprint-4` | `#5319E7` | Scope pengerjaan Sprint 4: PWA, A11y & Polish |
| `type: feature` | `#A2EEEF` | Penambahan fitur baru antarmuka pengguna |
| `type: bug` | `#D73A4A` | Masalah atau galat pada tampilan / fungsi |
| `type: refactor`| `#D4C5F9` | Perapian kode tanpa mengubah fungsi |
| `type: docs` | `#0075CA` | Dokumentasi komponen & panduan |
| `layer: ui/ux` | `#BFD4F2` | Terkait slicing desain, animasi, styling Tailwind |
| `layer: api` | `#F9D0C4` | Terkait integrasi REST API & SSE backend |
| `layer: auth` | `#E99695` | Terkait SSO Keycloak, JWT Cookie & RBAC Guard |
| `priority: p1` | `#B60205` | Kritis / Blocker utama sistem |
| `priority: p2` | `#E4E669` | Prioritas menengah / fungsional penting |
| `priority: p3` | `#008672` | Peningkatan minor / estetika tambahan |

---

## 🚀 2. Breakdown Sprint & User Stories Backlog

---

### 🟢 SPRINT 1: Core Shell, Katalog Aplikasi & Pencarian
> **Sprint Goal**: Membangun fondasi arsitektur frontend, navigasi modern, dan katalog aplikasi dengan pencarian instan.

#### Issue #1: [Core] Inisialisasi Project React 19, Vite, Tailwind CSS v4 & Routing
* **Labels**: `sprint-1`, `type: feature`, `layer: ui/ux`, `priority: p1`
* **Story Points**: `3`
* **User Story**:
  > Sebagai developer, saya ingin struktur project frontend terkonfigurasi dengan React 19, Vite, Tailwind v4, dan React Router v7 agar pengembangan komponen modular berjalan stabil dan terstandarisasi.
* **Acceptance Criteria**:
  * [ ] Project Vite React 19 + TypeScript berhasil di-bootstrap.
  * [ ] Tailwind CSS v4 terpasang dengan palet warna resmi BPS/ALUSI.
  * [ ] React Router v7 terkonfigurasi untuk rute `/`, `/status`, `/guides`, `/login`.
  * [ ] TanStack Query Provider & Axios/Fetch API client terkonfigurasi dengan `baseURL` backend.

---

#### Issue #2: [Catalog] Integrasi API Kategori & Tabs Filter Aplikasi
* **Labels**: `sprint-1`, `type: feature`, `layer: api`, `layer: ui/ux`, `priority: p1`
* **Story Points**: `5`
* **User Story**:
  > Sebagai pengunjung, saya ingin melihat dan memfilter aplikasi berdasarkan kategori statistik (Sosial, Produksi, Nerwilis, IPDS, Umum) agar saya dapat menemukan aplikasi bidang terkait dengan mudah.
* **API Terkait**: `GET /api/v1/categories` & `GET /api/v1/apps?category_id={id}`
* **Acceptance Criteria**:
  * [ ] Menampilkan tab bar kategori horizontal dengan active pill state & badge jumlah aplikasi.
  * [ ] State kategori sinkron dengan URL parameter query (`?category=ipds`).
  * [ ] Tampilan loading skeleton saat transisi data kategori.

---

#### Issue #3: [Catalog] Komponen Kartu Aplikasi (AppLauncherCard)
* **Labels**: `sprint-1`, `type: feature`, `layer: ui/ux`, `priority: p1`
* **Story Points**: `5`
* **User Story**:
  > Sebagai pengguna, saya ingin melihat kartu aplikasi yang informatif (logo, nama, deskripsi singkat, kategori, badge publik/internal, tombol peluncur) agar mudah dikenali.
* **Acceptance Criteria**:
  * [ ] Render logo aplikasi teroptimasi dengan fallback default avatar.
  * [ ] Badge penanda `Publik` vs `Internal BPS`.
  * [ ] Tombol launch membuka URL aplikasi di tab baru (`target="_blank"`).
  * [ ] Hover elevation & smooth transition micro-interaction.

---

#### Issue #4: [Search] Omnibox Instant Search Palette (Cmd+K)
* **Labels**: `sprint-1`, `type: feature`, `layer: api`, `layer: ui/ux`, `priority: p1`
* **Story Points**: `8`
* **User Story**:
  > Sebagai pengguna, saya ingin membuka kotak pencarian sentral dengan menekan `Ctrl+K` / `Cmd+K` dan mengetikkan kata kunci agar mendapatkan hasil aplikasi secara real-time.
* **API Terkait**: `GET /api/v1/apps/search?q={query}`
* **Acceptance Criteria**:
  * [ ] Shortcut global `Cmd+K` (Mac) dan `Ctrl+K` (Windows/Linux) memicu modal dialog.
  * [ ] Debounce query input sebesar 200ms untuk efisiensi pemanggilan API backend.
  * [ ] Highlight teks kata kunci yang cocok pada hasil pencarian.
  * [ ] Navigasi keyboard penuh (`Arrow Up`, `Arrow Down`, `Enter` untuk buka, `Esc` untuk keluar).

---

### 🔵 SPRINT 2: AI Assistant Cerdas, Personalisasi & Autentikasi SSO
> **Sprint Goal**: Mengintegrasikan asisten AI Gemini, sistem single sign-on (SSO), serta fitur personalisasi aplikasi favorit dan riwayat pegawai.

#### Issue #5: [Auth] Integrasi Login SSO BPS Sumut (OAuth2 / OIDC Handshake)
* **Labels**: `sprint-2`, `type: feature`, `layer: auth`, `layer: api`, `priority: p1`
* **Story Points**: `8`
* **User Story**:
  > Sebagai pegawai BPS, saya ingin login menggunakan akun SSO BPS Sumut satu pintu agar data identitas (NIP, Nama, Satker, Role) terintegrasi otomatis di portal.
* **API Terkait**: `GET /api/v1/auth/login`, `GET /api/v1/auth/callback`, `GET /api/v1/users/me`
* **Acceptance Criteria**:
  * [ ] Tombol login di TopNav mengarahkan ke URL SSO BPS Sumut.
  * [ ] Halaman callback `/callback` menangkap auth code, memanggil API backend, dan menyimpan state sesi.
  * [ ] TopNav menampilkan profil user (Avatar, Nama, Badge Satker & Role) setelah login berhasil.
  * [ ] Tombol Logout menghapus cookie sesi dan mereset state autentikasi klien.

---

#### Issue #6: [RBAC] Auth Guard & Proteksi Peluncuran Modul Internal
* **Labels**: `sprint-2`, `type: feature`, `layer: auth`, `priority: p1`
* **Story Points**: `5`
* **User Story**:
  > Sebagai sistem, saya ingin memvalidasi status login dan kecocokan role pegawai sebelum membuka modul internal agar hak akses data statistik terjaga aman.
* **Acceptance Criteria**:
  * [ ] Klik pada aplikasi publik langsung membuka tautan.
  * [ ] Klik pada aplikasi internal saat belum login memicu modal peringatan SSO login.
  * [ ] Klik pada aplikasi internal saat role tidak cocok menampilkan modal `Akses Ditolak / Hubungi Admin Satker`.
  * [ ] Pengguna yang terotorisasi berhasil membuka aplikasi internal.

---

#### Issue #7: [AI] Spotlight Assistant Chat Drawer (Integrasi Gemini AI)
* **Labels**: `sprint-2`, `type: feature`, `layer: api`, `layer: ui/ux`, `priority: p8`
* **Story Points**: `8`
* **User Story**:
  > Sebagai pengguna, saya ingin mengajukan pertanyaan kebutuhan statistik/pekerjaan secara natural dan mendapatkan jawaban serta kartu rekomendasi aplikasi yang relevan.
* **API Terkait**: `POST /api/v1/ai/ask`
* **Acceptance Criteria**:
  * [ ] Spotlight bar AI di hero section & floating AI button di pojok bawah.
  * [ ] Rendering respons percakapan dengan format Markdown (bold, lists, link).
  * [ ] Rendering kartu aplikasi yang direkomendasikan langsung di dalam bubble AI chat.
  * [ ] Suggestion pills (contoh pertanyaan siap klik) untuk mempermudah pengguna baru.

---

#### Issue #8: [Personalization] Fitur Aplikasi Favorit & Telemetri Klik
* **Labels**: `sprint-2`, `type: feature`, `layer: api`, `priority: p2`
* **Story Points**: `5`
* **User Story**:
  > Sebagai pegawai login, saya ingin menandai aplikasi harian sebagai Favorit (bintang) dan melihat daftar 5 aplikasi terakhir dibuka agar pekerjaan lebih cepat.
* **API Terkait**: `POST /api/v1/apps/:id/favorite`, `POST /api/v1/apps/:id/click`, `GET /api/v1/users/me/recents`
* **Acceptance Criteria**:
  * [ ] Tombol bintang toggle favorit dengan *Optimistic UI Update* (berubah instan tanpa lag).
  * [ ] Widget horizontal `QuickFavoritesBar` di atas katalog menampilkan aplikasi ter-pin.
  * [ ] Setiap klik meluncurkan aplikasi otomatis memanggil API `click` secara asinkron (background).
  * [ ] Section `Riwayat Terakhir Diakses` memuat daftar 5 aplikasi terakhir.

---

### 🟡 SPRINT 3: Realtime Live Monitoring SSE, Broadcast & Helpdesk
> **Sprint Goal**: Menghadirkan monitoring kesehatan server realtime via Server-Sent Events, banner pengumuman dinamis, dan helpdesk umpan balik.

#### Issue #9: [SSE] Live Health Monitoring & Indikator Status Server Realtime
* **Labels**: `sprint-3`, `type: feature`, `layer: api`, `layer: ui/ux`, `priority: p1`
* **Story Points**: `8`
* **User Story**:
  > Sebagai pengunjung dan admin, saya ingin melihat status kesehatan modul aplikasi secara live tanpa perlu reload halaman agar mengetahui jika ada layanan yang maintenance/gangguan.
* **API Terkait**: `GET /api/v1/services/realtime-status` (EventSource SSE) & `GET /api/v1/services/status`
* **Acceptance Criteria**:
  * [ ] Koneksi persistent EventSource SSE terkelola dengan auto-reconnect logic jika jaringan terputus.
  * [ ] Titik status berdenyut (Pulsing Dot): Hijau (Normal), Kuning (Gangguan), Biru (Maintenance), Merah (Down).
  * [ ] Toast notification otomatis muncul di kanan atas saat status aplikasi kritis berubah.
  * [ ] Halaman SLA Publik `/status` menampilkan rekap 30 hari uptime.

---

#### Issue #10: [Broadcast] Header Banner Pengumuman & Dismissable State
* **Labels**: `sprint-3`, `type: feature`, `layer: api`, `priority: p2`
* **Story Points**: `3`
* **User Story**:
  > Sebagai pengguna, saya ingin melihat banner siaran pengumuman penting (jadwal survei nasional / rilis fitur / maintenance) di bagian atas portal.
* **API Terkait**: `GET /api/v1/announcements`
* **Acceptance Criteria**:
  * [ ] Banner sticky di header dengan tema warna dinamis (`info`, `warning`, `urgent`, `maintenance`).
  * [ ] Tombol silang `X` untuk menutup banner dengan preferensi tersimpan di `localStorage`.
  * [ ] Dukungan carousel jika terdapat lebih dari 1 pengumuman aktif.

---

#### Issue #11: [Helpdesk] Drawer Panduan Aplikasi & FAQ
* **Labels**: `sprint-3`, `type: feature`, `layer: api`, `layer: ui/ux`, `priority: p2`
* **Story Points**: `5`
* **User Story**:
  > Sebagai pengguna awam, saya ingin membaca petunjuk teknis langkah penggunaan aplikasi dan kontak PIC pengelola tanpa keluar dari portal.
* **API Terkait**: `GET /api/v1/apps/:slug/guides`
* **Acceptance Criteria**:
  * [ ] Tombol `Panduan` pada kartu membuka Slide-over Drawer.
  * [ ] Menampilkan informasi FAQ, panduan PDF download link, dan info PIC Satker pengembang.

---

#### Issue #12: [Feedback] Modal Formulir Pengajuan Tiket Kendala Layanan
* **Labels**: `sprint-3`, `type: feature`, `layer: api`, `layer: ui/ux`, `priority: p2`
* **Story Points**: `5`
* **User Story**:
  > Sebagai pengguna yang mengalami kendala teknis, saya ingin melaporkan bug atau saran pengembangan langsung ke tim IT BPS Sumut.
* **API Terkait**: `POST /api/v1/feedbacks`
* **Acceptance Criteria**:
  * [ ] Floating Helpdesk Button di pojok kanan bawah memicu modal dialog formulir.
  * [ ] Validasi input form: Kategori (Bug, Usulan Fitur, Data), Modul Terkait, Deskripsi, Kontak/Email.
  * [ ] Sukses submit menampilkan Nomor Tiket Pengaduan untuk tindak lanjut.

---

### 🟣 SPRINT 4: Aksesibilitas WCAG AA, PWA, Pengujian & Rilis
> **Sprint Goal**: Menjamin aksesibilitas standar tinggi, instalasi PWA di smartphone/desktop, pengujian end-to-end (E2E), dan deployment produksi.

#### Issue #13: [A11y] Audit Aksesibilitas WCAG 2.2 Level AA & Keyboard Navigation
* **Labels**: `sprint-4`, `type: refactor`, `layer: ui/ux`, `priority: p1`
* **Story Points**: `5`
* **User Story**:
  > Sebagai pengguna dengan kebutuhan khusus / pembaca layar, saya ingin seluruh tombol, modal, dan kartu memiliki ARIA label semantik dan rasio kontras warna yang nyaman.
* **Acceptance Criteria**:
  * [ ] Skor audit Google Lighthouse Accessibility 100/100.
  * [ ] Navigasi keyboard penuh (`Tab`, `Shift+Tab`, `Enter`, `Escape`, `Arrow keys`).
  * [ ] Focus trap di dalam modal Omnibox dan drawer AI.
  * [ ] Dukungan preferensi `prefers-reduced-motion` dan tema Dark/Light mode adaptif.

---

#### Issue #14: [PWA] Service Worker Caching & Installable Web App
* **Labels**: `sprint-4`, `type: feature`, `layer: ui/ux`, `priority: p2`
* **Story Points**: `5`
* **User Story**:
  > Sebagai pegawai mobile, saya ingin menginstall portal ALUSI ke layar utama smartphone/laptop dan dapat membuka katalog meski sinyal internet lemah/offline.
* **Acceptance Criteria**:
  * [ ] Web App Manifest (`manifest.json`) dengan ikon resolusi lengkap & theme color resmi.
  * [ ] Service Worker cache untuk asset statis (HTML, CSS, JS, logo aplikasi).
  * [ ] Banner prompt `Install ALUSI App` yang ramah di browser mobile.
  * [ ] Fallback offline state yang informatif jika koneksi internet terputus total.

---

#### Issue #15: [QA & Deploy] Pengujian E2E, UAT Internal & Production Release
* **Labels**: `sprint-4`, `type: docs`, `priority: p1`
* **Story Points**: `5`
* **User Story**:
  > Sebagai Product Owner, saya ingin memverifikasi seluruh fungsional portal bersama tim internal BPS sebelum dirilis resmi ke publik.
* **Acceptance Criteria**:
  * [ ] UAT (User Acceptance Testing) internal pegawai BPS lulus 100% skenario.
  * [ ] Bundle build produksi teroptimasi (< 150 KB gzipped).
  * [ ] Setup pipeline CI/CD deployment ke production server (Nginx / Vercel).
  * [ ] Dokumentasi penggunaan & serah terima sistem selesai.

---

## 📊 3. Ringkasan Estimasi Story Points & Kecepatan Sprint

| Sprint | Fokus Utama | Total Stories | Total Points | Estimasi Durasi |
|---|---|---|---|---|
| **Sprint 1** | Core Shell, Routing, Katalog, Filter, Omnibox Search | 4 Issues | 21 pts | 2 Minggu |
| **Sprint 2** | SSO BPS Auth, RBAC Guard, AI Gemini, Favorit, Recents | 4 Issues | 26 pts | 2 Minggu |
| **Sprint 3** | Live SSE Status, Banner Broadcast, Drawer Panduan, Feedback | 4 Issues | 21 pts | 2 Minggu |
| **Sprint 4** | Aksesibilitas WCAG AA, PWA Offline, UAT, Production Deploy | 3 Issues | 15 pts | 2 Minggu |
| **TOTAL** | **Full Portal MVP Release** | **15 Issues** | **83 pts** | **8 Minggu (4 Sprints)** |
