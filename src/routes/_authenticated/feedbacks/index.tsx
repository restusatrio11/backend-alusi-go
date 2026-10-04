import { createFileRoute, redirect } from '@tanstack/react-router'
import { Feedbacks } from '@/features/feedbacks'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/_authenticated/feedbacks/')({
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('feedbacks:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: Feedbacks,
})
