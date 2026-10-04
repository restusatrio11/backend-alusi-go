import { z } from 'zod'
import { createFileRoute, redirect } from '@tanstack/react-router'
import { Apps } from '@/features/apps'
import { useAuthStore } from '@/stores/auth-store'

const appsSearchSchema = z.object({
  search: z.string().optional().catch(''),
  category: z.string().optional().catch(''),
})

export const Route = createFileRoute('/_authenticated/apps/')({
  validateSearch: appsSearchSchema,
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('apps:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: Apps,
})
