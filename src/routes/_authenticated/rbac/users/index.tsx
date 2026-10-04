import { createFileRoute, redirect } from '@tanstack/react-router'
import { UserRoleManagement } from '@/features/rbac/users'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/_authenticated/rbac/users/')({
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('users:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: UserRoleManagement,
})
