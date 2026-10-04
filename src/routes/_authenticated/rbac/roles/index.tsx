import { createFileRoute, redirect } from '@tanstack/react-router'
import { RolesPermissionMatrix } from '@/features/rbac/roles'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/_authenticated/rbac/roles/')({
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('rbac:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: RolesPermissionMatrix,
})
