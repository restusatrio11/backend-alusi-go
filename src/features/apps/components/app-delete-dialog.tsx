import { useState } from 'react'
import { AlertTriangle, Loader2, ShieldAlert } from 'lucide-react'
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
      <DialogContent className='sm:max-w-lg p-6 md:p-7'>
        <DialogHeader className='space-y-2 pb-2'>
          <div className='flex items-center gap-2.5 text-destructive'>
            <div className='p-2 rounded-full bg-destructive/10'>
              <AlertTriangle className='size-5 text-destructive' />
            </div>
            <DialogTitle className='text-lg font-bold'>Nonaktifkan Katalog Aplikasi</DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            Apakah Anda yakin ingin menonaktifkan aplikasi{' '}
            <strong className='text-foreground font-semibold'>&quot;{app?.nama}&quot;</strong>?
          </DialogDescription>
        </DialogHeader>

        <div className='p-3.5 rounded-xl border border-destructive/20 bg-destructive/5 text-xs text-muted-foreground space-y-1 my-1'>
          <p className='font-semibold text-destructive flex items-center gap-1.5'>
            <ShieldAlert className='size-3.5' /> Dampak Tindakan:
          </p>
          <p className='leading-relaxed'>
            Aplikasi tidak akan lagi tampil di halaman beranda atau pencarian publik portal ALUSI, namun riwayat log statistik klik sebelumnya tetap tersimpan secara aman.
          </p>
        </div>

        <DialogFooter className='gap-2.5 pt-3 border-t'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            disabled={isDeleting}
            className='h-10 px-4 text-xs'
          >
            Batal
          </Button>
          <Button
            type='button'
            variant='destructive'
            onClick={handleDelete}
            disabled={isDeleting}
            className='h-10 px-5 text-xs font-semibold shadow-xs'
          >
            {isDeleting && <Loader2 className='mr-2 size-4 animate-spin' />}
            Ya, Nonaktifkan Aplikasi
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
