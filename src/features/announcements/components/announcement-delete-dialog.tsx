import { useState } from 'react'
import { AlertTriangle, Loader2, ShieldAlert } from 'lucide-react'
import { toast } from 'sonner'
import { Announcement } from '@/types/announcements'
import { announcementsApi } from '../api/announcements-api'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

interface AnnouncementDeleteDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  announcement: Announcement | null
  onSuccess: () => void
}

export function AnnouncementDeleteDialog({
  open,
  onOpenChange,
  announcement,
  onSuccess,
}: AnnouncementDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState(false)

  const handleDelete = async () => {
    if (!announcement) return
    setIsDeleting(true)

    try {
      await announcementsApi.deleteAnnouncement(announcement.id)
      toast.success(`Pengumuman "${announcement.judul}" berhasil dihapus!`)
      onSuccess()
      onOpenChange(false)
    } catch (err: any) {
      toast.error('Gagal Menghapus Pengumuman', {
        description: err?.message || 'Terjadi kesalahan.',
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
            <DialogTitle className='text-lg font-bold'>Hapus Siaran Pengumuman</DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            Apakah Anda yakin ingin menghapus siaran pengumuman{' '}
            <strong className='text-foreground font-semibold'>&quot;{announcement?.judul}&quot;</strong>?
          </DialogDescription>
        </DialogHeader>

        <div className='p-3.5 rounded-xl border border-destructive/20 bg-destructive/5 text-xs text-muted-foreground space-y-1 my-1'>
          <p className='font-semibold text-destructive flex items-center gap-1.5'>
            <ShieldAlert className='size-3.5' /> Dampak Penghapusan:
          </p>
          <p className='leading-relaxed'>
            Banner pengumuman ini akan langsung dihentikan dan tidak akan ditampilkan lagi pada portal pengguna publik maupun internal.
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
            Ya, Hapus Pengumuman
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
