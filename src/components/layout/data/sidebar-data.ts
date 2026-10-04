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
        },
        {
          title: 'Kategori Aplikasi',
          url: '/categories',
          icon: FolderTree,
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
        },
        {
          title: 'Pengumuman',
          url: '/announcements',
          icon: Megaphone,
        },
        {
          title: 'Umpan Balik & Isu',
          url: '/feedbacks',
          icon: MessageSquareWarning,
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
        },
        {
          title: 'Audit Trail',
          url: '/audit-logs',
          icon: ShieldCheck,
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
