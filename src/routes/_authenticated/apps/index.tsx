import { z } from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import { Apps } from '@/features/apps'

const appsSearchSchema = z.object({
  search: z.string().optional().catch(''),
  category: z.string().optional().catch(''),
})

export const Route = createFileRoute('/_authenticated/apps/')({
  validateSearch: appsSearchSchema,
  component: Apps,
})
