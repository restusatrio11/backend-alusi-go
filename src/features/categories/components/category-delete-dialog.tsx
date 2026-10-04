import { useState } from 'react'
import { AlertTriangle, Loader2 } from 'lucide-react'
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
      <DialogContent className='max-w-md'>
        <DialogHeader>
          <div className='flex items-center gap-2 text-destructive mb-1'>
            <AlertTriangle className='size-5' />
            <DialogTitle>Hapus Kategori</DialogTitle>
          </div>
          <DialogDescription>
            Apakah Anda yakin ingin menghapus kategori{' '}
            <span className='font-semibold text-foreground'>&quot;{category?.nama}&quot;</span>?
            Pastikan tidak ada aplikasi yang masih menggunakan kategori ini sebelum dihapus.
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
