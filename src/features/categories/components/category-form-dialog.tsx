import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2 } from 'lucide-react'
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
      <DialogContent className='max-w-md'>
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Edit Kategori Aplikasi' : 'Tambah Kategori Baru'}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Perbarui nama kelompok, deskripsi, atau urutan kategori aplikasi.'
              : 'Tambahkan kelompok kategori baru untuk mengelompokkan aplikasi dalam katalog.'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4 py-2'>
            <FormField
              control={form.control}
              name='nama'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Nama Kategori *</FormLabel>
                  <FormControl>
                    <Input placeholder='Contoh: Layanan Statistik, Manajemen Internal' {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className='grid grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='ikon'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Ikon Lucide</FormLabel>
                    <FormControl>
                      <Input placeholder='bar-chart-2' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='urutan'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nomor Urutan</FormLabel>
                    <FormControl>
                      <Input
                        type='number'
                        min={1}
                        value={field.value}
                        onChange={(e) => field.onChange(Number(e.target.value))}
                      />
                    </FormControl>
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
                  <FormLabel>Deskripsi Singkat</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder='Penjelasan kelompok layanan aplikasi ini...'
                      rows={3}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className='gap-2 pt-2'>
              <Button
                type='button'
                variant='outline'
                onClick={() => onOpenChange(false)}
                disabled={form.formState.isSubmitting}
              >
                Batal
              </Button>
              <Button type='submit' disabled={form.formState.isSubmitting}>
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
