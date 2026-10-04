import { createFileRoute } from '@tanstack/react-router'
import { Feedbacks } from '@/features/feedbacks'

export const Route = createFileRoute('/_authenticated/feedbacks/')({
  component: Feedbacks,
})
