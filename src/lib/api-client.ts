import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth-store'

const envUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const baseURL = envUrl.replace(/\/api\/v1\/?$/, '').replace(/\/+$/, '')

export const apiClient = axios.create({
  baseURL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request Interceptor: Attach Bearer JWT token if available & normalize URL path
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    if (config.url && !config.url.startsWith('/api/v1') && !config.url.startsWith('http')) {
      config.url = `/api/v1${config.url.startsWith('/') ? config.url : '/' + config.url}`
    }
    const token = useAuthStore.getState().auth.accessToken
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response Interceptor: Handle Unauthorized (401) and format errors
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error: AxiosError<{ message?: string; error?: { message?: string } }>) => {
    if (error.response?.status === 401) {
      // Clear invalid token
      useAuthStore.getState().auth.reset()

      // Only redirect if not already on the sign-in page or callback
      if (
        typeof window !== 'undefined' &&
        !window.location.pathname.includes('/sign-in') &&
        !window.location.pathname.includes('/callback')
      ) {
        const currentPath = encodeURIComponent(window.location.pathname + window.location.search)
        window.location.href = `/sign-in?redirect=${currentPath}`
      }
    }

    const errorMessage =
      error.response?.data?.message ||
      error.response?.data?.error?.message ||
      error.message ||
      'Terjadi kesalahan saat memproses permintaan.'

    return Promise.reject(new Error(errorMessage))
  }
)

export default apiClient
