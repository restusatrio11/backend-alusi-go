import { z } from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import { Callback } from '@/features/auth/callback'

const searchSchema = z.object({
  code: z.string().optional(),
  state: z.string().optional(),
})

export const Route = createFileRoute('/(auth)/callback')({
  component: Callback,
  validateSearch: searchSchema,
})
