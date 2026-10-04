import { createFileRoute, redirect } from '@tanstack/react-router'
import { Announcements } from '@/features/announcements'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/_authenticated/announcements/')({
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('announcements:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: Announcements,
})
