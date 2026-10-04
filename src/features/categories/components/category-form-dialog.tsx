import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { FolderTree, Sparkles, Hash, FileText, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Category, CategoryPayload } from '@/types/apps'
import { categoriesApi } from '../api/categories-api'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'

const categoryFormSchema = z.object({
  nama: z.string().min(2, 'Nama kategori minimal 2 karakter.'),
  slug: z.string().optional(),
  deskripsi: z.string().optional(),
  ikon: z.string().optional(),
  urutan: z.number().min(1, 'Urutan minimal 1.'),
})

type CategoryFormValues = z.infer<typeof categoryFormSchema>

interface CategoryFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  categoryToEdit?: Category | null
  onSuccess: () => void
}

export function CategoryFormDialog({
  open,
  onOpenChange,
  categoryToEdit,
  onSuccess,
}: CategoryFormDialogProps) {
  const isEditing = !!categoryToEdit

  const form = useForm<CategoryFormValues>({
    resolver: zodResolver(categoryFormSchema),
    defaultValues: {
      nama: '',
      slug: '',
      deskripsi: '',
      ikon: '',
      urutan: 1,
    },
  })

  useEffect(() => {
    if (categoryToEdit) {
      form.reset({
        nama: categoryToEdit.nama,
        slug: categoryToEdit.slug,
        deskripsi: categoryToEdit.deskripsi || '',
        ikon: categoryToEdit.ikon || '',
        urutan: categoryToEdit.urutan || 1,
      })
    } else {
      form.reset({
        nama: '',
        slug: '',
        deskripsi: '',
        ikon: '',
        urutan: 1,
      })
    }
  }, [categoryToEdit, form, open])

  const onSubmit = async (values: CategoryFormValues) => {
    try {
      const payload: CategoryPayload = {
        ...values,
        deskripsi: values.deskripsi || undefined,
        ikon: values.ikon || undefined,
      }

      if (isEditing && categoryToEdit) {
        await categoriesApi.updateCategory(categoryToEdit.id, payload)
        toast.success(`Kategori "${values.nama}" berhasil diperbarui!`)
      } else {
        await categoriesApi.createCategory(payload)
        toast.success(`Kategori "${values.nama}" berhasil ditambahkan!`)
      }

      onSuccess()
      onOpenChange(false)
    } catch (err: any) {
      toast.error('Gagal Menyimpan Kategori', {
        description: err?.message || 'Terjadi kesalahan pada server backend.',
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-xl lg:max-w-2xl max-h-[90vh] overflow-y-auto p-6 md:p-8'>
        <DialogHeader className='space-y-1.5 pb-2 border-b'>
          <div className='flex items-center gap-2 text-primary'>
            <FolderTree className='size-5' />
            <DialogTitle className='text-xl font-bold tracking-tight'>
              {isEditing ? 'Edit Kategori Aplikasi' : 'Tambah Kategori Aplikasi Baru'}
            </DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            {isEditing
              ? 'Perbarui nama kelompok, deskripsi tugas pokok fungsi, atau urutan indeks kategori.'
              : 'Tambahkan kelompok kategori baru untuk menata dan mengklasifikasikan seluruh aplikasi dalam katalog portal.'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-5 pt-3'>
            <FormField
              control={form.control}
              name='nama'
              render={({ field }) => (
                <FormItem>
                  <FormLabel className='text-xs font-medium'>Nama Kategori *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder='Contoh: Layanan Statistik, Manajemen Internal, Tata Usaha'
                      className='h-10 text-xs'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className='grid grid-cols-1 sm:grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='ikon'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                      <Sparkles className='size-3.5 text-muted-foreground' />
                      Ikon Identifier (Lucide Icon)
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder='bar-chart-2, database, layers'
                        className='h-10 text-xs font-mono'
                        {...field}
                      />
                    </FormControl>
                    <FormDescription className='text-[11px]'>
                      Nama identifier ikon dari pustaka Lucide
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='urutan'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                      <Hash className='size-3.5 text-muted-foreground' />
                      Nomor Urutan Tampilan
                    </FormLabel>
                    <FormControl>
                      <Input
                        type='number'
                        min={1}
                        className='h-10 text-xs font-mono'
                        value={field.value}
                        onChange={(e) => field.onChange(Number(e.target.value))}
                      />
                    </FormControl>
                    <FormDescription className='text-[11px]'>
                      Urutan prioritas penataan di menu navigasi
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name='deskripsi'
              render={({ field }) => (
                <FormItem>
                  <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                    <FileText className='size-3.5 text-muted-foreground' />
                    Deskripsi Singkat Kategori
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder='Jelaskan kelompok fungsi modul atau sistem yang termasuk ke dalam kategori ini...'
                      rows={4}
                      className='text-xs leading-relaxed resize-y'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className='gap-2.5 pt-4 border-t'>
              <Button
                type='button'
                variant='outline'
                onClick={() => onOpenChange(false)}
                disabled={form.formState.isSubmitting}
                className='h-10 px-4 text-xs'
              >
                Batal
              </Button>
              <Button
                type='submit'
                disabled={form.formState.isSubmitting}
                className='h-10 px-5 text-xs font-semibold shadow-xs'
              >
                {form.formState.isSubmitting && (
                  <Loader2 className='mr-2 size-4 animate-spin' />
                )}
                {isEditing ? 'Simpan Perubahan' : 'Tambah Kategori'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
