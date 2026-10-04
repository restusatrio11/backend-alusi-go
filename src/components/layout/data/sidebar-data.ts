import {
  LayoutDashboard,
  Layers,
  FolderTree,
  Activity,
  Megaphone,
  MessageSquareWarning,
  LineChart,
  ShieldCheck,
  Settings,
  Palette,
  Monitor,
  Building2,
  KeyRound,
  Users2,
} from 'lucide-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  user: {
    name: 'Admin BPS Sumut',
    email: 'admin@bps.go.id',
    avatar: '',
  },
  teams: [
    {
      name: 'ALUSI BPS Sumut',
      logo: Building2,
      plan: 'Portal Admin Terpadu',
    },
  ],
  navGroups: [
    {
      title: 'Katalog & Layanan',
      items: [
        {
          title: 'Dashboard',
          url: '/',
          icon: LayoutDashboard,
        },
        {
          title: 'Katalog Aplikasi',
          url: '/apps',
          icon: Layers,
          permission: 'apps:view',
        },
        {
          title: 'Kategori Aplikasi',
          url: '/categories',
          icon: FolderTree,
          permission: 'categories:view',
        },
      ],
    },
    {
      title: 'Pemantauan & Interaksi',
      items: [
        {
          title: 'Monitoring Layanan',
          url: '/monitoring',
          icon: Activity,
          permission: 'monitoring:view',
        },
        {
          title: 'Pengumuman',
          url: '/announcements',
          icon: Megaphone,
          permission: 'announcements:view',
        },
        {
          title: 'Umpan Balik & Isu',
          url: '/feedbacks',
          icon: MessageSquareWarning,
          permission: 'feedbacks:view',
        },
      ],
    },
    {
      title: 'Analitik & Jejak Audit',
      items: [
        {
          title: 'Analitik Penggunaan',
          url: '/analytics',
          icon: LineChart,
          permission: 'analytics:view',
        },
        {
          title: 'Audit Trail',
          url: '/audit-logs',
          icon: ShieldCheck,
          permission: 'audit:view',
        },
      ],
    },
    {
      title: 'Manajemen Akses & RBAC',
      permissions: ['rbac:view', 'users:view'],
      items: [
        {
          title: 'Role & Hak Akses',
          url: '/rbac/roles',
          icon: KeyRound,
          permission: 'rbac:view',
        },
        {
          title: 'Pengguna & Peran',
          url: '/rbac/users',
          icon: Users2,
          permission: 'users:view',
        },
      ],
    },
    {
      title: 'Konfigurasi',
      items: [
        {
          title: 'Pengaturan',
          icon: Settings,
          items: [
            {
              title: 'Profil Akun',
              url: '/settings',
              icon: Settings,
            },
            {
              title: 'Tema & Tampilan',
              url: '/settings/appearance',
              icon: Palette,
            },
            {
              title: 'Display Layar',
              url: '/settings/display',
              icon: Monitor,
            },
          ],
        },
      ],
    },
  ],
}
