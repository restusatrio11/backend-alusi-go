import { describe, it, expect, beforeEach } from 'vitest'
import { render } from 'vitest-browser-react'
import { PermissionGuard } from './permission-guard'
import { useAuthStore } from '@/stores/auth-store'

describe('PermissionGuard component', () => {
  beforeEach(() => {
    useAuthStore.getState().auth.reset()
    if (typeof window !== 'undefined') {
      localStorage.clear()
    }
  })

  it('renders children when user has admin role', async () => {
    useAuthStore.getState().auth.setUser({
      id: 1,
      nama: 'Superadmin',
      email: 'admin@bps.go.id',
      user_type: 'internal',
      status: 'active',
      created_at: '2026-01-01',
      updated_at: '2026-01-01',
      roles: [{ id: 1, nama: 'admin', created_at: '2026-01-01' }],
      permissions: [],
    })

    const { getByRole } = await render(
      <PermissionGuard permission='apps:create'>
        <button id='action-btn'>Tambah Aplikasi</button>
      </PermissionGuard>
    )

    await expect.element(getByRole('button', { name: 'Tambah Aplikasi' })).toBeInTheDocument()
  })

  it('renders fallback or nothing when user lacks permission', async () => {
    useAuthStore.getState().auth.setUser({
      id: 2,
      nama: 'Pegawai Test',
      email: 'pegawai@bps.go.id',
      user_type: 'internal',
      status: 'active',
      created_at: '2026-01-01',
      updated_at: '2026-01-01',
      roles: [{ id: 3, nama: 'pegawai', created_at: '2026-01-01' }],
      permissions: ['apps:view'],
    })

    const { getByText } = await render(
      <PermissionGuard
        permission='apps:delete'
        fallback={<span id='no-access'>Akses Ditolak</span>}
      >
        <button id='action-btn'>Hapus Data</button>
      </PermissionGuard>
    )

    await expect.element(getByText('Akses Ditolak')).toBeInTheDocument()
  })

  it('renders children when user matches required role', async () => {
    useAuthStore.getState().auth.setUser({
      id: 3,
      nama: 'Pimpinan BPS',
      email: 'pimpinan@bps.go.id',
      user_type: 'internal',
      status: 'active',
      created_at: '2026-01-01',
      updated_at: '2026-01-01',
      roles: [{ id: 2, nama: 'pimpinan', created_at: '2026-01-01' }],
      permissions: [],
    })

    const { getByText } = await render(
      <PermissionGuard role='pimpinan'>
        <div id='pimpinan-panel'>Laporan Eksekutif</div>
      </PermissionGuard>
    )

    await expect.element(getByText('Laporan Eksekutif')).toBeInTheDocument()
  })
})
