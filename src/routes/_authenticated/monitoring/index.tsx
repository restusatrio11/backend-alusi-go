import { createFileRoute, redirect } from '@tanstack/react-router'
import { Monitoring } from '@/features/monitoring'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/_authenticated/monitoring/')({
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('monitoring:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: Monitoring,
})
