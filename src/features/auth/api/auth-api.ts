import { apiClient } from '@/lib/api-client'
import { LoginResponse, StandardApiResponse, User } from '@/types/auth'

export interface ManualLoginCredentials {
  username: string
  password: string
}

export const authApi = {
  // Manual Login with Username/Email/NIP and Password
  manualLogin: async (credentials: ManualLoginCredentials): Promise<LoginResponse> => {
    const response = await apiClient.post<StandardApiResponse<LoginResponse>>(
      '/auth/login',
      credentials
    )
    return response.data.data
  },

  // Get SSO BPS Sumut Login URL
  getSSOLoginURL: async (): Promise<string> => {
    const response = await apiClient.get<StandardApiResponse<{ login_url: string; state: string }>>(
      '/auth/login?redirect=false'
    )
    return response.data.data.login_url
  },

  // Handle SSO OAuth Callback exchange
  handleSSOCallback: async (code: string, state?: string): Promise<LoginResponse> => {
    const params = new URLSearchParams({
      code,
      format: 'json',
    })
    if (state) {
      params.append('state', state)
    }

    const response = await apiClient.get<StandardApiResponse<LoginResponse>>(
      `/auth/callback?${params.toString()}`
    )
    return response.data.data
  },

  // Fetch current user profile
  getMe: async (): Promise<User> => {
    const response = await apiClient.get<StandardApiResponse<User>>('/auth/me')
    return response.data.data
  },

  // Logout from server session
  logout: async (): Promise<void> => {
    try {
      await apiClient.post('/auth/logout')
    } catch {
      // Ignore network error on logout
    }
  },
}
