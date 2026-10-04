import { useState } from 'react'
import { AlertTriangle, Loader2 } from 'lucide-react'
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
      <DialogContent className='max-w-md'>
        <DialogHeader>
          <div className='flex items-center gap-2 text-destructive mb-1'>
            <AlertTriangle className='size-5' />
            <DialogTitle>Hapus Pengumuman</DialogTitle>
          </div>
          <DialogDescription>
            Apakah Anda yakin ingin menghapus pengumuman{' '}
            <span className='font-semibold text-foreground'>&quot;{announcement?.judul}&quot;</span>?
            Banner ini tidak akan ditampilkan lagi pada portal pengguna.
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
            Ya, Hapus
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
