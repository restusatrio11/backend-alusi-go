import { useState } from 'react'
import { AlertTriangle, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { App } from '@/types/apps'
import { appsApi } from '../api/apps-api'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

interface AppDeleteDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  app: App | null
  onSuccess: () => void
}

export function AppDeleteDialog({
  open,
  onOpenChange,
  app,
  onSuccess,
}: AppDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState(false)

  const handleDelete = async () => {
    if (!app) return
    setIsDeleting(true)

    try {
      await appsApi.deleteApp(app.id)
      toast.success(`Aplikasi "${app.nama}" berhasil dinonaktifkan!`)
      onSuccess()
      onOpenChange(false)
    } catch (err: any) {
      toast.error('Gagal Menghapus Aplikasi', {
        description: err?.message || 'Terjadi kesalahan pada server backend.',
      })
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-w-md'>
        <DialogHeader>
          <div className='flex items-center gap-2 text-destructive mb-1'>
            <AlertTriangle className='size-5' />
            <DialogTitle>Nonaktifkan Aplikasi</DialogTitle>
          </div>
          <DialogDescription>
            Apakah Anda yakin ingin menonaktifkan aplikasi{' '}
            <span className='font-semibold text-foreground'>&quot;{app?.nama}&quot;</span>?
            Aplikasi tidak akan ditampilkan lagi pada katalog pengguna umum.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter className='gap-2 pt-2'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            disabled={isDeleting}
          >
            Batal
          </Button>
          <Button
            type='button'
            variant='destructive'
            onClick={handleDelete}
            disabled={isDeleting}
          >
            {isDeleting && <Loader2 className='mr-2 size-4 animate-spin' />}
            Ya, Nonaktifkan
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
