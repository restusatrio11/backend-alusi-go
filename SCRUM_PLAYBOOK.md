# 📘 Scrum Master Playbook & Panduan Eksekusi Tim: Frontend Portal ALUSI BPS Sumut

> **Panduan Praktis Manajemen Proyek Tangkas (Agile/Scrum)**  
> Dibuat khusus untuk pengembang, Product Owner, dan Tim Inovasi BPS Provinsi Sumatera Utara dalam membangun **Frontend Portal Publik ALUSI**.

---

## 🎯 1. Peran & Tanggung Jawab (Scrum Roles)

| Peran | Siapa | Tugas Utama |
|---|---|---|
| **Product Owner (PO)** | Ketua Tim Inovasi / PJ Portal | Menentukan prioritas fitur (*What & Why*), menyetujui kriteria penerimaan (*Acceptance Criteria*), dan menerima hasil kerja tiap akhir sprint. |
| **Scrum Master (SM)** | Anda / Tech Lead | Memfasilitasi jalannya upacara Scrum (*How*), membersihkan hambatan (*blockers*), memastikan tim mematuhi *Definition of Done (DoD)*. |
| **Development Team** | Frontend Developer, UI/UX Designer, QA | Merancang arsitektur komponen, mengoding fitur, mengintegrasikan API backend, dan menguji aplikasi. |

---

## ⏱️ 2. Ritme & Upacara Scrum (Sprint Cadence - 2 Minggu per Sprint)

Setiap Sprint berdurasi **2 Minggu (10 Hari Kerja)** dengan jadwal upacara:

```
[ Hari 1: Senin ]   ➔ Sprint Planning (2 jam)
[ Tiap Hari: 08:30] ➔ Daily Standup (10-15 menit)
[ Hari 5: Jumat ]   ➔ Backlog Refinement / Grooming (1 jam)
[ Hari 10: Jumat ]  ➔ Sprint Review & Demo (1 jam)
[ Hari 10: Jumat ]  ➔ Sprint Retrospective (45 menit)
```

---

## 📋 3. Panduan Langkah Demi Langkah: "Apa yang Harus Saya Lakukan?"

### 🔹 Langkah 1: Persiapan Awal (Setup GitHub Project & Labels)
1. **Buka Repository GitHub**: Masuk ke repositori project.
2. **Setup Labels / Tags**: Buat label standar proyek (lihat daftar label di `SCRUM_BACKLOG_ISSUES.md`).
3. **Buat GitHub Project Board (Kanban)** dengan 5 Kolom:
   * 📥 **Product Backlog**: Kumpulan seluruh User Story & Ide masa depan.
   * 📋 **Todo (Sprint Backlog)**: Issue yang disepakati untuk dikerjakan di Sprint aktif.
   * 🏗️ **In Progress**: Issue yang sedang aktif dikerjakan (Maksimal 2 issue per developer / *WIP Limit*).
   * 🔍 **Review / PR**: Pull Request sedang ditinjau & di-linting.
   * 🧪 **Testing / QA**: Verifikasi kriteria penerimaan di staging/preview URL.
   * ✅ **Done**: Selesai memenuhi *Definition of Done (DoD)* dan di-merge.

---

### 🔹 Langkah 2: Hari ke-1 — Sprint Planning (Perencanaan Sprint)
**Tujuan**: Menentukan apa yang akan diselesaikan dalam 2 minggu ke depan.
1. **Pilih User Story dari Product Backlog** sesuai prioritas PO (misal: Sprint 1 Fokus Core Shell & Katalog).
2. **Breakdown Story menjadi Sub-tasks**:
   * Contoh: Story *"Katalog Aplikasi Grid"* dipecah menjadi:
     * *Task 1: Buat hook fetch `useApps` dengan TanStack Query.*
     * *Task 2: Desain komponen `AppCard` dengan Tailwind & Shadcn.*
     * *Task 3: Tambahkan filter kategori switch.*
3. **Estimasi Story Points** (Gunakan skala Fibonacci: 1, 2, 3, 5, 8, 13).
4. **Pindahkan issue terpilih ke kolom `Todo (Sprint Backlog)`**.
5. **Kunci Sprint Goal**: Satu kalimat tujuan utama sprint (contoh: *"Sprint 1: Pengguna dapat mencari dan melihat seluruh katalog aplikasi BPS Sumut secara instan"*).

---

### 🔹 Langkah 3: Hari ke-2 s/d Hari ke-9 — Eksekusi & Daily Standup
**Tujuan**: Sinkronisasi harian cepat (Maksimal 15 menit, berdiri/santai).
Setiap orang menjawab **3 Pertanyaan Kunci**:
1. *Apa yang sudah saya selesaikan kemarin?*
2. *Apa yang akan saya kerjakan hari ini?*
3. *Apakah ada kendala/hambatan (blocker) yang menghambat saya?*

> 💡 **Aturan Kerja Harian Developer**:
> 1. Pindahkan kartu dari `Todo` ke `In Progress` saat mulai ngoding.
> 2. Buat branch git baru dengan format: `feat/issue-ID-deskripsi` atau `fix/issue-ID-deskripsi`.
> 3. Buat Pull Request (PR) dan pindahkan kartu ke `Review / PR`.
> 4. Jalankan pengujian mandiri sebelum meminta review.

---

### 🔹 Langkah 4: Hari ke-10 — Sprint Review & Demo
**Tujuan**: Mendemonstrasikan produk yang berfungsi ke Product Owner & Stakeholder.
1. Buka aplikasi yang berjalan (*Working Software*), bukan sekadar slide presentasi.
2. Uji alur nyata:
   * Buka browser ➔ Cari aplikasi via `Cmd+K` ➔ Coba tanya AI Assistant ➔ Coba login SSO ➔ Klik aplikasi terproteksi.
3. Catat umpan balik dari PO untuk dimasukkan ke Product Backlog sprint berikutnya.

---

### 🔹 Langkah 5: Hari ke-10 — Sprint Retrospective
**Tujuan**: Evaluasi internal tim untuk peningkatan berkelanjutan (*Continuous Improvement*).
Diskusikan 3 hal:
1. 🟢 **What went well?** (Apa yang sudah berjalan bagus & perlu dipertahankan?)
2. 🔴 **What didn't go well?** (Apa yang kurang optimal / menghambat kecepatan tim?)
3. 🟡 **Action Items for next sprint** (1-2 tindakan konkret untuk perbaikan di sprint berikutnya).

---

## 🏆 4. Standar Mutu: Definition of Ready & Definition of Done

### ✅ Definition of Ready (DoR) — Sebelum Issue Ditarik ke Sprint:
* [ ] Memiliki deskripsi jelas dengan format: *Sebagai [Pengguna], saya ingin [Fitur] agar [Manfaat]*.
* [ ] Memiliki kriteria penerimaan (*Acceptance Criteria*) terukur.
* [ ] Endpoint backend API terkait sudah tersedia & teruji di Swagger/Postman.
* [ ] Story Point sudah diestimasi oleh tim.

### 🏁 Definition of Done (DoD) — Sebelum Issue Dinyatakan Selesai (Done):
* [ ] Kode lulus linting & formatting (ESLint, Prettier).
* [ ] Tidak ada error TypeScript / console log debug di browser.
* [ ] Responsive di Desktop (1920x1080), Tablet (iPad), dan Mobile (iPhone/Android).
* [ ] Lulus pengujian fungsional sesuai *Acceptance Criteria*.
* [ ] Pull Request sudah di-review dan di-merge ke branch `main` / `develop`.
* [ ] Berhasil di-deploy ke environment staging/production.

---

## 🏷️ 5. Cheat-Sheet Perintah Git Workflow Scrum

```bash
# 1. Update branch develop terbaru
git checkout develop
git pull origin develop

# 2. Buat branch fitur baru sesuai issue
git checkout -b feat/ALUSI-04-omni-search

# 3. Koding dan commit dengan referensi issue
git add .
git commit -m "feat(search): implement omnibox dialog with trigram debounce (#4)"

# 4. Push branch ke remote
git push origin feat/ALUSI-04-omni-search

# 5. Buka Pull Request (PR) di GitHub dan hubungkan ke Issue #4
```
