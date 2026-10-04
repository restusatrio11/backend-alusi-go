import { useState } from 'react'
import { AlertTriangle, Loader2, ShieldAlert } from 'lucide-react'
import { toast } from 'sonner'
import { Category } from '@/types/apps'
import { categoriesApi } from '../api/categories-api'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

interface CategoryDeleteDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  category: Category | null
  onSuccess: () => void
}

export function CategoryDeleteDialog({
  open,
  onOpenChange,
  category,
  onSuccess,
}: CategoryDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState(false)

  const handleDelete = async () => {
    if (!category) return
    setIsDeleting(true)

    try {
      await categoriesApi.deleteCategory(category.id)
      toast.success(`Kategori "${category.nama}" berhasil dihapus!`)
      onSuccess()
      onOpenChange(false)
    } catch (err: any) {
      toast.error('Gagal Menghapus Kategori', {
        description: err?.message || 'Kategori ini mungkin masih memiliki aplikasi terhubung.',
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
            <DialogTitle className='text-lg font-bold'>Hapus Kategori Aplikasi</DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            Apakah Anda yakin ingin menghapus kategori{' '}
            <strong className='text-foreground font-semibold'>&quot;{category?.nama}&quot;</strong>?
          </DialogDescription>
        </DialogHeader>

        <div className='p-3.5 rounded-xl border border-destructive/20 bg-destructive/5 text-xs text-muted-foreground space-y-1 my-1'>
          <p className='font-semibold text-destructive flex items-center gap-1.5'>
            <ShieldAlert className='size-3.5' /> Perhatian:
          </p>
          <p className='leading-relaxed'>
            Pastikan tidak ada aplikasi yang masih terhubung ke kategori ini sebelum melakukan penghapusan untuk menjaga integritas data katalog.
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
            Ya, Hapus Kategori
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
