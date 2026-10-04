import { create } from 'zustand'
import { getCookie, setCookie, removeCookie } from '@/lib/cookies'
import { User } from '@/types/auth'

const ACCESS_TOKEN_KEY = 'alusi_access_token'
const USER_DATA_KEY = 'alusi_user_data'

interface AuthState {
  auth: {
    user: User | null
    setUser: (user: User | null) => void
    accessToken: string
    setAccessToken: (accessToken: string) => void
    setAuth: (token: string, user: User) => void
    resetAccessToken: () => void
    reset: () => void
    isAuthenticated: () => boolean
    isAdmin: () => boolean
    isPimpinan: () => boolean
    hasRole: (roleName: string) => boolean
    hasAnyRole: (roleNames: string[]) => boolean
    hasPermission: (permissionCode: string) => boolean
    hasAnyPermission: (permissionCodes: string[]) => boolean
    hasAllPermissions: (permissionCodes: string[]) => boolean
  }
}

export const useAuthStore = create<AuthState>()((set, get) => {
  // Initialize token and user from cookie / storage
  const cookieToken = getCookie(ACCESS_TOKEN_KEY)
  let initialToken = ''
  if (cookieToken) {
    try {
      initialToken = JSON.parse(cookieToken)
    } catch {
      initialToken = cookieToken
    }
  }

  let initialUser: User | null = null
  if (typeof window !== 'undefined') {
    const storedUser = localStorage.getItem(USER_DATA_KEY)
    if (storedUser) {
      try {
        initialUser = JSON.parse(storedUser)
      } catch {
        initialUser = null
      }
    }
  }

  return {
    auth: {
      user: initialUser,
      accessToken: initialToken,
      setUser: (user) => {
        if (user && typeof window !== 'undefined') {
          localStorage.setItem(USER_DATA_KEY, JSON.stringify(user))
        } else if (typeof window !== 'undefined') {
          localStorage.removeItem(USER_DATA_KEY)
        }
        set((state) => ({ ...state, auth: { ...state.auth, user } }))
      },
      setAccessToken: (accessToken) => {
        setCookie(ACCESS_TOKEN_KEY, JSON.stringify(accessToken))
        set((state) => ({ ...state, auth: { ...state.auth, accessToken } }))
      },
      setAuth: (token, user) => {
        setCookie(ACCESS_TOKEN_KEY, JSON.stringify(token))
        if (typeof window !== 'undefined') {
          localStorage.setItem(USER_DATA_KEY, JSON.stringify(user))
        }
        set((state) => ({
          ...state,
          auth: { ...state.auth, accessToken: token, user },
        }))
      },
      resetAccessToken: () => {
        removeCookie(ACCESS_TOKEN_KEY)
        set((state) => ({ ...state, auth: { ...state.auth, accessToken: '' } }))
      },
      reset: () => {
        removeCookie(ACCESS_TOKEN_KEY)
        if (typeof window !== 'undefined') {
          localStorage.removeItem(USER_DATA_KEY)
        }
        set((state) => ({
          ...state,
          auth: { ...state.auth, user: null, accessToken: '' },
        }))
      },
      isAuthenticated: () => {
        return !!get().auth.accessToken
      },
      isAdmin: () => {
        const user = get().auth.user
        if (!user || !user.roles) return false
        return user.roles.some((r) => r.nama === 'admin' || r.nama === 'superadmin')
      },
      isPimpinan: () => {
        const user = get().auth.user
        if (!user || !user.roles) return false
        return user.roles.some((r) => r.nama === 'pimpinan' || r.nama === 'admin')
      },
      hasRole: (roleName: string) => {
        const user = get().auth.user
        if (!user || !user.roles) return false
        return user.roles.some((r) => r.nama.toLowerCase() === roleName.toLowerCase())
      },
      hasAnyRole: (roleNames: string[]) => {
        const user = get().auth.user
        if (!user || !user.roles) return false
        return user.roles.some((r) =>
          roleNames.map((n) => n.toLowerCase()).includes(r.nama.toLowerCase())
        )
      },
      hasPermission: (permissionCode: string) => {
        const user = get().auth.user
        if (!user) return false
        // Admin always has full access
        if (user.roles?.some((r) => r.nama === 'admin' || r.nama === 'superadmin')) {
          return true
        }
        if (!user.permissions) return false
        return user.permissions.includes(permissionCode)
      },
      hasAnyPermission: (permissionCodes: string[]) => {
        const user = get().auth.user
        if (!user) return false
        if (user.roles?.some((r) => r.nama === 'admin' || r.nama === 'superadmin')) {
          return true
        }
        if (!user.permissions || user.permissions.length === 0) return false
        return permissionCodes.some((code) => user.permissions?.includes(code))
      },
      hasAllPermissions: (permissionCodes: string[]) => {
        const user = get().auth.user
        if (!user) return false
        if (user.roles?.some((r) => r.nama === 'admin' || r.nama === 'superadmin')) {
          return true
        }
        if (!user.permissions) return false
        return permissionCodes.every((code) => user.permissions?.includes(code))
      },
    },
  }
})
