import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Layers,
  Globe,
  Tag,
  Users,
  Activity,
  Building2,
  Phone,
  FileText,
  Eye,
  CheckCircle2,
  Loader2,
} from 'lucide-react'
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
      <DialogContent className='sm:max-w-3xl lg:max-w-4xl max-h-[90vh] overflow-y-auto p-6 md:p-8'>
        <DialogHeader className='space-y-1.5 pb-2 border-b'>
          <div className='flex items-center gap-2 text-primary'>
            <Layers className='size-5' />
            <DialogTitle className='text-xl font-bold tracking-tight'>
              {isEditing ? 'Edit Metadata Aplikasi' : 'Tambah Aplikasi Baru'}
            </DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            {isEditing
              ? 'Perbarui informasi katalog aplikasi, tautan URL, status layanan, dan target pengguna secara lengkap.'
              : 'Lengkapi formulir di bawah ini untuk mendaftarkan modul atau sistem aplikasi baru ke dalam katalog ALUSI.'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-6 pt-3'>
            {/* Section 1: Identitas & Klasifikasi */}
            <div className='space-y-4'>
              <div className='flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground'>
                <Tag className='size-3.5 text-primary' />
                <span>Identitas & Klasifikasi Modul</span>
              </div>

              <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
                <FormField
                  control={form.control}
                  name='nama'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className='text-xs font-medium'>Nama Aplikasi *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder='Contoh: SIMPEG BPS, Pojok Statistik, ARON'
                          className='h-10 text-xs'
                          {...field}
                        />
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
                      <FormLabel className='text-xs font-medium'>Kategori Aplikasi *</FormLabel>
                      <Select
                        value={field.value ? String(field.value) : ''}
                        onValueChange={(val) => field.onChange(Number(val))}
                      >
                        <FormControl>
                          <SelectTrigger className='h-10 text-xs'>
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
                    <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                      <Globe className='size-3.5 text-muted-foreground' />
                      URL Tautan Aplikasi Web *
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder='https://layanan.bps.go.id atau https://sumut.bps.go.id/app'
                        className='h-10 text-xs font-mono'
                        {...field}
                      />
                    </FormControl>
                    <FormDescription className='text-[11px]'>
                      Tautan absolut yang dibuka ketika pengguna mengklik peluncuran modul ini.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Section 2: Sasaran & Operasional */}
            <div className='space-y-4 pt-2 border-t'>
              <div className='flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground'>
                <Activity className='size-3.5 text-primary' />
                <span>Sasaran Pengguna & Status Layanan</span>
              </div>

              <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
                <FormField
                  control={form.control}
                  name='target_pengguna'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                        <Users className='size-3.5 text-muted-foreground' />
                        Target Sasaran Pengguna
                      </FormLabel>
                      <Select value={field.value} onValueChange={field.onChange}>
                        <FormControl>
                          <SelectTrigger className='h-10 text-xs'>
                            <SelectValue placeholder='Pilih Sasaran' />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value='Semua'>🌐 Semua (Publik & Pegawai)</SelectItem>
                          <SelectItem value='Internal BPS'>🔒 Internal Pegawai BPS</SelectItem>
                          <SelectItem value='Mitra'>📋 Mitra Statistik / Lapangan</SelectItem>
                          <SelectItem value='Pimpinan'>👔 Pimpinan & Manajerial</SelectItem>
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
                      <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                        <Activity className='size-3.5 text-muted-foreground' />
                        Status Operasional Layanan
                      </FormLabel>
                      <Select value={field.value} onValueChange={field.onChange}>
                        <FormControl>
                          <SelectTrigger className='h-10 text-xs'>
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
                      <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                        <Building2 className='size-3.5 text-muted-foreground' />
                        Unit Pengelola / Pemilik Modul *
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder='Contoh: Tim TI / Bagian Umum / Nerwilis'
                          className='h-10 text-xs'
                          {...field}
                        />
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
                      <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                        <Phone className='size-3.5 text-muted-foreground' />
                        Kontak Admin / Helpdesk
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder='Email atau No. WhatsApp PIC layanan'
                          className='h-10 text-xs'
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>

            {/* Section 3: Deskripsi & Visibilitas */}
            <div className='space-y-4 pt-2 border-t'>
              <FormField
                control={form.control}
                name='deskripsi'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className='text-xs font-medium flex items-center gap-1.5'>
                      <FileText className='size-3.5 text-muted-foreground' />
                      Deskripsi & Ringkasan Fitur
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder='Jelaskan fungsi utama, fitur unggulan, dan cakupan pemanfaatan aplikasi ini...'
                        rows={4}
                        className='text-xs leading-relaxed resize-y'
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
                <FormField
                  control={form.control}
                  name='is_public'
                  render={({ field }) => (
                    <FormItem className='flex flex-row items-center justify-between rounded-xl border bg-muted/30 p-4 shadow-xs'>
                      <div className='space-y-0.5 pr-2'>
                        <FormLabel className='text-xs font-semibold flex items-center gap-1.5'>
                          <Eye className='size-3.5 text-primary' />
                          Visibilitas Publik
                        </FormLabel>
                        <FormDescription className='text-[11px] leading-snug'>
                          Dapat dilihat dan diakses pada portal publik tanpa login SSO
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

                {isEditing ? (
                  <FormField
                    control={form.control}
                    name='aktif'
                    render={({ field }) => (
                      <FormItem className='flex flex-row items-center justify-between rounded-xl border bg-muted/30 p-4 shadow-xs'>
                        <div className='space-y-0.5 pr-2'>
                          <FormLabel className='text-xs font-semibold flex items-center gap-1.5'>
                            <CheckCircle2 className='size-3.5 text-emerald-600' />
                            Status Katalog Aktif
                          </FormLabel>
                          <FormDescription className='text-[11px] leading-snug'>
                            Menampilkan aplikasi dalam daftar pencarian katalog
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
                ) : (
                  <div className='rounded-xl border border-dashed p-4 flex items-center text-xs text-muted-foreground bg-muted/10'>
                    Aplikasi baru secara default akan berstatus aktif di katalog.
                  </div>
                )}
              </div>
            </div>

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
                {isEditing ? 'Simpan Perubahan' : 'Tambah Aplikasi'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
