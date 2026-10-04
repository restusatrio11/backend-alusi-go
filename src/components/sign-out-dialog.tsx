import { useNavigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/auth-store'
import { authApi } from '@/features/auth/api/auth-api'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { toast } from 'sonner'

interface SignOutDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SignOutDialog({ open, onOpenChange }: SignOutDialogProps) {
  const navigate = useNavigate()
  const { auth } = useAuthStore()

  const handleSignOut = async () => {
    try {
      await authApi.logout()
    } catch {
      // ignore
    } finally {
      auth.reset()
      toast.info('Anda telah keluar dari sesi.')
      navigate({
        to: '/sign-in',
        replace: true,
      })
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      title='Keluar Sesi'
      desc='Apakah Anda yakin ingin keluar? Anda perlu login kembali untuk mengakses panel admin ALUSI.'
      confirmText='Keluar'
      cancelBtnText='Batal'
      destructive
      handleConfirm={handleSignOut}
      className='sm:max-w-sm'
    />
  )
}
