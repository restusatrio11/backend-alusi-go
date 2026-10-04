-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- 1. Satuan Kerja (Satker) BPS
CREATE TABLE IF NOT EXISTS satker (
    id SERIAL PRIMARY KEY,
    kode VARCHAR(10) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Roles
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(50) NOT NULL UNIQUE,
    deskripsi VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Users (Sinkronisasi JIT dari SSO BPS)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    sso_sub VARCHAR(255) NOT NULL UNIQUE,
    nip VARCHAR(50) UNIQUE,
    nama VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    satker_id INT REFERENCES satker(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. User Roles Mapping
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- 5. Categories
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    deskripsi TEXT,
    urutan INT NOT NULL DEFAULT 0,
    ikon VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 6. Applications (Katalog Aplikasi)
CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    nama VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    url TEXT NOT NULL,
    deskripsi TEXT,
    ikon_url TEXT,
    target_pengguna VARCHAR(100) NOT NULL DEFAULT 'Semua',
    pemilik VARCHAR(255) NOT NULL DEFAULT 'BPS RI',
    kontak_admin VARCHAR(255),
    urutan INT NOT NULL DEFAULT 0,
    status_layanan VARCHAR(20) NOT NULL DEFAULT 'online',
    is_public BOOLEAN NOT NULL DEFAULT true,
    aktif BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. App Role Access Restrictions
CREATE TABLE IF NOT EXISTS app_access (
    app_id INT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (app_id, role_id)
);

-- 8. User Favorites
CREATE TABLE IF NOT EXISTS favorites (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    app_id INT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, app_id)
);

-- 9. Click Logs (Partitioned / Volume Analytics)
CREATE TABLE IF NOT EXISTS click_logs (
    id BIGSERIAL,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    app_id INT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    satker_id INT REFERENCES satker(id) ON DELETE SET NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, clicked_at)
) PARTITION BY RANGE (clicked_at);

-- Create default initial partitions for click_logs
CREATE TABLE IF NOT EXISTS click_logs_2026_q4 PARTITION OF click_logs
    FOR VALUES FROM ('2026-10-01') TO ('2027-01-01');
CREATE TABLE IF NOT EXISTS click_logs_default PARTITION OF click_logs
    DEFAULT;

-- 10. Status Checks (Automated Health Monitoring Logs)
CREATE TABLE IF NOT EXISTS status_checks (
    id BIGSERIAL,
    app_id INT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    http_status_code INT,
    latency_ms INT NOT NULL,
    error_message TEXT,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, checked_at)
) PARTITION BY RANGE (checked_at);

-- Create default initial partitions for status_checks
CREATE TABLE IF NOT EXISTS status_checks_2026_q4 PARTITION OF status_checks
    FOR VALUES FROM ('2026-10-01') TO ('2027-01-01');
CREATE TABLE IF NOT EXISTS status_checks_default PARTITION OF status_checks
    DEFAULT;

-- 11. Application Guides & FAQs
CREATE TABLE IF NOT EXISTS guides (
    id SERIAL PRIMARY KEY,
    app_id INT NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    judul VARCHAR(255) NOT NULL,
    konten TEXT NOT NULL,
    urutan INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 12. Announcements
CREATE TABLE IF NOT EXISTS announcements (
    id SERIAL PRIMARY KEY,
    judul VARCHAR(255) NOT NULL,
    konten TEXT NOT NULL,
    tipe VARCHAR(30) NOT NULL DEFAULT 'info',
    app_id INT REFERENCES apps(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    starts_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 13. Feedbacks & Issue Reports
CREATE TABLE IF NOT EXISTS feedbacks (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    app_id INT REFERENCES apps(id) ON DELETE SET NULL,
    kategori VARCHAR(50) NOT NULL DEFAULT 'kendala',
    pesan TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    catatan_admin TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 14. Admin Audit Logs (Append-Only)
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    entity VARCHAR(50) NOT NULL,
    entity_id VARCHAR(50),
    payload_before JSONB,
    payload_after JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indices for High-Performance Queries & Search
CREATE INDEX IF NOT EXISTS idx_apps_category_aktif ON apps (category_id, aktif);
CREATE INDEX IF NOT EXISTS idx_apps_urutan ON apps (urutan ASC);
CREATE INDEX IF NOT EXISTS idx_apps_trgm_search ON apps USING gin (nama gin_trgm_ops, deskripsi gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_users_sso_sub ON users (sso_sub);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_satker ON users (satker_id);

CREATE INDEX IF NOT EXISTS idx_click_logs_app_time ON click_logs (app_id, clicked_at DESC);
CREATE INDEX IF NOT EXISTS idx_click_logs_user ON click_logs (user_id);

CREATE INDEX IF NOT EXISTS idx_status_checks_app_time ON status_checks (app_id, checked_at DESC);
CREATE INDEX IF NOT EXISTS idx_announcements_active_time ON announcements (is_active, starts_at, ends_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs (entity, entity_id);

-- Initial Master Data (Roles & Default Satker Sumut)
INSERT INTO roles (id, nama, deskripsi) VALUES
    (1, 'admin', 'Administrator Portal'),
    (2, 'pimpinan', 'Pimpinan / Eselon BPS'),
    (3, 'pegawai', 'Pegawai BPS Provinsi dan Satker Kabupaten/Kota'),
    (4, 'petugas_lapangan', 'Petugas Lapangan Sensus dan Survei'),
    (5, 'mitra', 'Mitra Statistik Non-Pegawai')
ON CONFLICT (id) DO NOTHING;

INSERT INTO satker (id, kode, nama) VALUES
    (1, '1200', 'BPS Provinsi Sumatera Utara'),
    (2, '1201', 'BPS Kabupaten Nias'),
    (3, '1202', 'BPS Kabupaten Mandailing Natal'),
    (4, '1203', 'BPS Kabupaten Tapanuli Selatan'),
    (5, '1204', 'BPS Kabupaten Tapanuli Tengah'),
    (6, '1205', 'BPS Kabupaten Tapanuli Utara'),
    (7, '1206', 'BPS Kabupaten Toba'),
    (8, '1207', 'BPS Kabupaten Labuhanbatu'),
    (9, '1208', 'BPS Kabupaten Asahan'),
    (10, '1209', 'BPS Kabupaten Simalungun'),
    (11, '1210', 'BPS Kabupaten Dairi'),
    (12, '1211', 'BPS Kabupaten Karo'),
    (13, '1212', 'BPS Kabupaten Deli Serdang'),
    (14, '1213', 'BPS Kabupaten Langkat'),
    (15, '1214', 'BPS Kabupaten Nias Selatan'),
    (16, '1215', 'BPS Kabupaten Humbang Hasundutan'),
    (17, '1216', 'BPS Kabupaten Pakpak Bharat'),
    (18, '1217', 'BPS Kabupaten Samosir'),
    (19, '1218', 'BPS Kabupaten Serdang Bedagai'),
    (20, '1219', 'BPS Kabupaten Batu Bara'),
    (21, '1220', 'BPS Kabupaten Padang Lawas Utara'),
    (22, '1221', 'BPS Kabupaten Padang Lawas'),
    (23, '1222', 'BPS Kabupaten Labuhanbatu Selatan'),
    (24, '1223', 'BPS Kabupaten Labuhanbatu Utara'),
    (25, '1224', 'BPS Kabupaten Nias Utara'),
    (26, '1225', 'BPS Kabupaten Nias Barat'),
    (27, '1271', 'BPS Kota Sibolga'),
    (28, '1272', 'BPS Kota Tanjungbalai'),
    (29, '1273', 'BPS Kota Pematangsiantar'),
    (30, '1274', 'BPS Kota Tebing Tinggi'),
    (31, '1275', 'BPS Kota Medan'),
    (32, '1276', 'BPS Kota Binjai'),
    (33, '1277', 'BPS Kota Padangsidimpuan'),
    (34, '1278', 'BPS Kota Gunungsitoli')
ON CONFLICT (id) DO NOTHING;

INSERT INTO categories (id, nama, slug, deskripsi, urutan, ikon) VALUES
    (1, 'Survei & Sensus', 'survei-sensus', 'Aplikasi pengumpulan data dan survei statistik', 1, 'clipboard-document-list'),
    (2, 'Diseminasi & Publikasi', 'diseminasi-publikasi', 'Portal data, visualisasi, dan rilis statistik', 2, 'chart-bar'),
    (3, 'Kepegawaian & SDM', 'kepegawaian-sdm', 'Presensi, kinerja, cuti, dan administrasi pegawai', 3, 'user-group'),
    (4, 'Keuangan & Pengadaan', 'keuangan-pengadaan', 'Anggaran, SPK, pembayaran, dan pertanggungjawaban', 4, 'banknotes'),
    (5, 'Tata Kelola & Internal', 'tata-kelola-internal', 'Persuratan dinas, arsip, dan kolaborasi tim', 5, 'folder')
ON CONFLICT (id) DO NOTHING;
