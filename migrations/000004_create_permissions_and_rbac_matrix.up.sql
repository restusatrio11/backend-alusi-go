-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    kode VARCHAR(100) NOT NULL UNIQUE,
    nama VARCHAR(150) NOT NULL,
    modul VARCHAR(100) NOT NULL,
    deskripsi VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

-- Indexes for fast RBAC checks
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions (permission_id);

-- Seed Granular Permissions across all portal modules
INSERT INTO permissions (kode, nama, modul, deskripsi) VALUES
    -- Dashboard & Executive Analytics
    ('dashboard:view', 'Lihat Executive Dashboard', 'Dashboard', 'Mengakses grafik ringkasan metrik eksekutif dan KPI rollup'),
    ('analytics:view', 'Lihat Analitik Penggunaan', 'Analitik', 'Mengakses statistik tren klik harian dan peringkat pemanfaatan'),
    ('reports:export', 'Ekspor Laporan CSV', 'Laporan', 'Mengunduh laporan inventaris katalog dan rekapitulasi analitik'),

    -- Katalog Aplikasi & Panduan
    ('apps:view', 'Lihat Katalog Aplikasi', 'Katalog Aplikasi', 'Melihat daftar seluruh katalog aplikasi admin'),
    ('apps:create', 'Tambah Aplikasi Baru', 'Katalog Aplikasi', 'Mendaftarkan metadata aplikasi baru ke dalam katalog portal'),
    ('apps:update', 'Edit Metadata Aplikasi', 'Katalog Aplikasi', 'Mengubah data aplikasi, tautan URL, status layanan, dan pemilik'),
    ('apps:delete', 'Nonaktifkan / Hapus Aplikasi', 'Katalog Aplikasi', 'Menonaktifkan aplikasi dari katalog publik'),
    ('apps:reorder', 'Ubah Urutan Prioritas', 'Katalog Aplikasi', 'Mengatur posisi dan urutan prioritas peluncuran aplikasi'),
    ('apps:upload_logo', 'Upload & Kompresi Logo', 'Katalog Aplikasi', 'Mengunggah dan mengoptimalkan gambar ikon logo aplikasi'),
    ('guides:manage', 'Kelola Panduan & FAQ', 'Panduan Aplikasi', 'Menambah, mengubah, dan menghapus dokumen petunjuk dan FAQ'),

    -- Kategori Aplikasi
    ('categories:view', 'Lihat Kategori Aplikasi', 'Kategori', 'Melihat daftar kelompok kategori aplikasi'),
    ('categories:create', 'Tambah Kategori', 'Kategori', 'Menambahkan kelompok kategori aplikasi baru'),
    ('categories:update', 'Edit Kategori', 'Kategori', 'Mengubah nama, ikon, atau urutan kelompok kategori'),
    ('categories:delete', 'Hapus Kategori', 'Kategori', 'Menghapus kelompok kategori aplikasi'),

    -- Pemantauan Layanan (Health Monitoring)
    ('monitoring:view', 'Lihat Monitoring Realtime', 'Monitoring Layanan', 'Mengakses dashboard streaming status layanan dan SLA realtime'),
    ('monitoring:probe', 'Trigger Probe Ping Manual', 'Monitoring Layanan', 'Menjalankan pengecekan kesehatan server backend secara manual'),

    -- Siaran Pengumuman
    ('announcements:view', 'Lihat Daftar Pengumuman', 'Pengumuman', 'Melihat daftar seluruh siaran banner pengumuman'),
    ('announcements:create', 'Buat Pengumuman Baru', 'Pengumuman', 'Menerbitkan siaran pengumuman maintenance atau rilis fitur'),
    ('announcements:update', 'Edit Pengumuman', 'Pengumuman', 'Mengubah judul, pesan, warna, atau jadwal siaran pengumuman'),
    ('announcements:delete', 'Hapus Pengumuman', 'Pengumuman', 'Menghentikan dan menghapus banner siaran pengumuman'),

    -- Umpan Balik & Isu
    ('feedbacks:view', 'Lihat Masukan & Laporan Isu', 'Umpan Balik', 'Melihat tiket keluhan, bug, dan saran dari pengguna'),
    ('feedbacks:respond', 'Tindak Lanjut & Respon Tiket', 'Umpan Balik', 'Mengubah status tiket dan memberikan respon solusi'),

    -- Jejak Audit Forensik
    ('audit:view', 'Lihat Jejak Audit Forensik', 'Audit Trail', 'Memeriksa log aktivitas keamanan dan manipulasi data sistem'),

    -- Manajemen RBAC & Hak Akses
    ('rbac:view', 'Lihat Manajemen Hak Akses', 'Manajemen RBAC', 'Mengakses menu pengelolaan role, permission, dan pengguna'),
    ('rbac:manage_roles', 'Kelola Matriks Permission Role', 'Manajemen RBAC', 'Mengubah centang hak akses perizinan untuk setiap role'),
    ('rbac:assign_users', 'Tetapkan Role Pengguna', 'Manajemen RBAC', 'Menetapkan dan mengubah penugasan role pada akun pegawai')
ON CONFLICT (kode) DO UPDATE
SET nama = EXCLUDED.nama,
    modul = EXCLUDED.modul,
    deskripsi = EXCLUDED.deskripsi;

-- Seed default Role Permissions mapping:
-- 1. Role 'admin' (ID: 1) -> Mendapatkan SEMUA permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, p.id FROM permissions p
ON CONFLICT DO NOTHING;

-- 2. Role 'pimpinan' (ID: 2) -> Dashboard, Analytics, Reports, Monitoring, Announcements, Feedbacks, Audit view
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, p.id FROM permissions p
WHERE p.kode IN (
    'dashboard:view',
    'analytics:view',
    'reports:export',
    'apps:view',
    'categories:view',
    'monitoring:view',
    'announcements:view',
    'feedbacks:view',
    'audit:view'
)
ON CONFLICT DO NOTHING;

-- 3. Role 'pegawai' (ID: 3) -> Dashboard view, Apps view, Monitoring view
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, p.id FROM permissions p
WHERE p.kode IN (
    'apps:view',
    'categories:view',
    'monitoring:view'
)
ON CONFLICT DO NOTHING;
