import { createFileRoute, redirect } from '@tanstack/react-router'
import { Categories } from '@/features/categories'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/_authenticated/categories/')({
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('categories:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: Categories,
})
