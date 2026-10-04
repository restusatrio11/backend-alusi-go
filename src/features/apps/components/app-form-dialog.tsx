import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { App, Category, CreateAppPayload, UpdateAppPayload } from '@/types/apps'
import { appsApi } from '../api/apps-api'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'

const appFormSchema = z.object({
  nama: z.string().min(2, 'Nama aplikasi minimal 2 karakter.'),
  slug: z.string().optional(),
  category_id: z.number().min(1, 'Pilih kategori aplikasi.'),
  url: z.string().url('URL harus berupa alamat web yang valid (contoh: https://app.bps.go.id).'),
  target_pengguna: z.string(),
  pemilik: z.string().min(2, 'Unit pemilik / pengelola wajib diisi.'),
  kontak_admin: z.string().optional(),
  status_layanan: z.enum(['online', 'pemeliharaan', 'kendala', 'offline']),
  is_public: z.boolean(),
  aktif: z.boolean(),
  deskripsi: z.string().optional(),
})

type AppFormValues = z.infer<typeof appFormSchema>

interface AppFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  appToEdit?: App | null
  categories: Category[]
  onSuccess: () => void
}

export function AppFormDialog({
  open,
  onOpenChange,
  appToEdit,
  categories,
  onSuccess,
}: AppFormDialogProps) {
  const isEditing = !!appToEdit

  const form = useForm<AppFormValues>({
    resolver: zodResolver(appFormSchema),
    defaultValues: {
      nama: '',
      slug: '',
      category_id: categories[0]?.id || 0,
      url: '',
      target_pengguna: 'Semua',
      pemilik: 'BPS Provinsi Sumatera Utara',
      kontak_admin: '',
      status_layanan: 'online',
      is_public: true,
      aktif: true,
      deskripsi: '',
    },
  })

  useEffect(() => {
    if (appToEdit) {
      form.reset({
        nama: appToEdit.nama,
        slug: appToEdit.slug,
        category_id: appToEdit.category_id,
        url: appToEdit.url,
        target_pengguna: appToEdit.target_pengguna || 'Semua',
        pemilik: appToEdit.pemilik || 'BPS Provinsi Sumatera Utara',
        kontak_admin: appToEdit.kontak_admin || '',
        status_layanan: appToEdit.status_layanan || 'online',
        is_public: appToEdit.is_public ?? true,
        aktif: appToEdit.aktif ?? true,
        deskripsi: appToEdit.deskripsi || '',
      })
    } else {
      form.reset({
        nama: '',
        slug: '',
        category_id: categories[0]?.id || 0,
        url: '',
        target_pengguna: 'Semua',
        pemilik: 'BPS Provinsi Sumatera Utara',
        kontak_admin: '',
        status_layanan: 'online',
        is_public: true,
        aktif: true,
        deskripsi: '',
      })
    }
  }, [appToEdit, categories, form, open])

  const onSubmit = async (values: AppFormValues) => {
    try {
      if (isEditing && appToEdit) {
        const payload: UpdateAppPayload = {
          ...values,
          category_id: Number(values.category_id),
          deskripsi: values.deskripsi || undefined,
          kontak_admin: values.kontak_admin || undefined,
        }
        await appsApi.updateApp(appToEdit.id, payload)
        toast.success(`Aplikasi "${values.nama}" berhasil diperbarui!`)
      } else {
        const payload: CreateAppPayload = {
          ...values,
          category_id: Number(values.category_id),
          deskripsi: values.deskripsi || undefined,
          kontak_admin: values.kontak_admin || undefined,
        }
        await appsApi.createApp(payload)
        toast.success(`Aplikasi "${values.nama}" berhasil ditambahkan!`)
      }

      onSuccess()
      onOpenChange(false)
    } catch (err: any) {
      toast.error(
        isEditing ? 'Gagal memperbarui aplikasi' : 'Gagal menambahkan aplikasi',
        {
          description: err?.message || 'Terjadi kesalahan pada server backend.',
        }
      )
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-w-2xl max-h-[90vh] overflow-y-auto'>
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Edit Metadata Aplikasi' : 'Tambah Aplikasi Baru'}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Perbarui data katalog aplikasi, tautan URL, status layanan, dan target pengguna.'
              : 'Lengkapi formulir di bawah ini untuk mendaftarkan aplikasi baru ke dalam portal ALUSI.'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4 py-2'>
            <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='nama'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nama Aplikasi *</FormLabel>
                    <FormControl>
                      <Input placeholder='Contoh: SIMPEG BPS, Pojok Statistik' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='category_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Kategori Aplikasi *</FormLabel>
                    <Select
                      value={field.value ? String(field.value) : ''}
                      onValueChange={(val) => field.onChange(Number(val))}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder='Pilih Kategori' />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {categories.map((cat) => (
                          <SelectItem key={cat.id} value={String(cat.id)}>
                            {cat.nama}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name='url'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>URL Tautan Aplikasi *</FormLabel>
                  <FormControl>
                    <Input placeholder='https://layanan.bps.go.id' {...field} />
                  </FormControl>
                  <FormDescription>
                    Tautan langsung menuju aplikasi web atau dashboard terkait.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='target_pengguna'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Target Pengguna</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder='Pilih Sasaran' />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='Semua'>Semua (Publik & Pegawai)</SelectItem>
                        <SelectItem value='Internal BPS'>Internal Pegawai BPS</SelectItem>
                        <SelectItem value='Mitra'>Mitra Statistik / Lapangan</SelectItem>
                        <SelectItem value='Pimpinan'>Pimpinan & Manajerial</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='status_layanan'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Status Operasional Layanan</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder='Status' />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='online'>🟢 Operasional Normal (Online)</SelectItem>
                        <SelectItem value='pemeliharaan'>🟡 Pemeliharaan (Maintenance)</SelectItem>
                        <SelectItem value='kendala'>🟠 Gangguan / Kendala Teknis</SelectItem>
                        <SelectItem value='offline'>🔴 Tidak Aktif (Offline)</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='pemilik'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Unit Pengelola / Pemilik</FormLabel>
                    <FormControl>
                      <Input placeholder='Contoh: Tim TI / Bagian Umum' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='kontak_admin'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Kontak Admin / Helpdesk</FormLabel>
                    <FormControl>
                      <Input placeholder='Email / No WhatsApp PIC' {...field} />
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
                  <FormLabel>Deskripsi Aplikasi</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder='Jelaskan fungsi utama, fitur, dan cakupan penggunaan aplikasi...'
                      rows={3}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4 pt-2 border-t'>
              <FormField
                control={form.control}
                name='is_public'
                render={({ field }) => (
                  <FormItem className='flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs'>
                    <div className='space-y-0.5'>
                      <FormLabel className='text-sm font-medium'>
                        Dapat Diakses Publik
                      </FormLabel>
                      <FormDescription className='text-xs'>
                        Tampil di portal publik tanpa login SSO
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />

              {isEditing && (
                <FormField
                  control={form.control}
                  name='aktif'
                  render={({ field }) => (
                    <FormItem className='flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs'>
                      <div className='space-y-0.5'>
                        <FormLabel className='text-sm font-medium'>
                          Status Katalog Aktif
                        </FormLabel>
                        <FormDescription className='text-xs'>
                          Tampilkan aplikasi dalam katalog
                        </FormDescription>
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              )}
            </div>

            <DialogFooter className='gap-2 pt-4'>
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
                {isEditing ? 'Simpan Perubahan' : 'Tambah Aplikasi'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
