import { clearCookies } from '@/test-utils/cookies'
import { beforeEach, describe, expect, it, vi } from 'vitest'

async function importAuthStore() {
  const { useAuthStore } = await import('./auth-store')
  return useAuthStore
}

const sampleUser = {
  id: 1,
  nama: 'Admin BPS',
  email: 'admin@bps.go.id',
  username: 'admin',
  user_type: 'internal',
  status: 'active',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  roles: [{ id: 1, nama: 'admin', created_at: '2026-01-01T00:00:00Z' }],
}

describe('useAuthStore', () => {
  beforeEach(() => {
    clearCookies()
    if (typeof window !== 'undefined') {
      localStorage.clear()
    }
    vi.resetModules()
  })

  it('starts with an empty access token when nothing is persisted', async () => {
    const useAuthStore = await importAuthStore()

    expect(useAuthStore.getState().auth.accessToken).toBe('')
    expect(useAuthStore.getState().auth.user).toBeNull()
  })

  it('persists access token so a new store instance reads it back', async () => {
    const useAuthStore = await importAuthStore()
    useAuthStore.getState().auth.setAccessToken('session-token')

    vi.resetModules()
    const useAuthStoreAfterReload = await importAuthStore()

    expect(useAuthStoreAfterReload.getState().auth.accessToken).toBe(
      'session-token'
    )
  })

  it('clears persisted access token when resetAccessToken is used', async () => {
    const useAuthStore = await importAuthStore()
    useAuthStore.getState().auth.setAccessToken('to-clear')
    useAuthStore.getState().auth.resetAccessToken()

    vi.resetModules()
    const useAuthStoreAfterReload = await importAuthStore()

    expect(useAuthStoreAfterReload.getState().auth.accessToken).toBe('')
  })

  it('updates the signed-in user via setUser', async () => {
    const useAuthStore = await importAuthStore()

    useAuthStore.getState().auth.setUser({ ...sampleUser })

    expect(useAuthStore.getState().auth.user).toEqual(sampleUser)
  })

  it('reset clears user and access token and drops persistence', async () => {
    const useAuthStore = await importAuthStore()
    useAuthStore.getState().auth.setAccessToken('will-be-cleared')
    useAuthStore.getState().auth.setUser({ ...sampleUser })

    useAuthStore.getState().auth.reset()

    expect(useAuthStore.getState().auth.user).toBeNull()
    expect(useAuthStore.getState().auth.accessToken).toBe('')

    vi.resetModules()
    const useAuthStoreAfterReload = await importAuthStore()

    expect(useAuthStoreAfterReload.getState().auth.user).toBeNull()
    expect(useAuthStoreAfterReload.getState().auth.accessToken).toBe('')
  })

  it('correctly evaluates roles and permissions for admin and regular users', async () => {
    const useAuthStore = await importAuthStore()

    // 1. Admin user bypass
    useAuthStore.getState().auth.setUser({
      ...sampleUser,
      roles: [{ id: 1, nama: 'admin', created_at: '2026-01-01' }],
      permissions: [],
    })

    expect(useAuthStore.getState().auth.isAdmin()).toBe(true)
    expect(useAuthStore.getState().auth.hasRole('admin')).toBe(true)
    expect(useAuthStore.getState().auth.hasPermission('apps:create')).toBe(true)
    expect(useAuthStore.getState().auth.hasPermission('anything:custom')).toBe(true)

    // 2. Regular staff with specific permissions
    useAuthStore.getState().auth.setUser({
      ...sampleUser,
      roles: [{ id: 3, nama: 'pegawai', created_at: '2026-01-01' }],
      permissions: ['apps:view', 'announcements:view'],
    })

    expect(useAuthStore.getState().auth.isAdmin()).toBe(false)
    expect(useAuthStore.getState().auth.hasRole('pegawai')).toBe(true)
    expect(useAuthStore.getState().auth.hasRole('admin')).toBe(false)
    expect(useAuthStore.getState().auth.hasPermission('apps:view')).toBe(true)
    expect(useAuthStore.getState().auth.hasPermission('apps:create')).toBe(false)
    expect(useAuthStore.getState().auth.hasAnyPermission(['apps:create', 'apps:view'])).toBe(true)
    expect(useAuthStore.getState().auth.hasAllPermissions(['apps:view', 'announcements:view'])).toBe(true)
    expect(useAuthStore.getState().auth.hasAllPermissions(['apps:view', 'apps:create'])).toBe(false)
  })
})
