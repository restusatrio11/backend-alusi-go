import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Info,
  AlertTriangle,
  AlertOctagon,
  CheckCircle2,
  ExternalLink,
  Loader2,
} from 'lucide-react'
import { toast } from 'sonner'
import { App } from '@/types/apps'
import {
  Announcement,
  AnnouncementType,
  CreateAnnouncementPayload,
  UpdateAnnouncementPayload,
} from '@/types/announcements'
import { announcementsApi } from '../api/announcements-api'
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

const announcementSchema = z.object({
  judul: z.string().min(2, 'Judul pengumuman minimal 2 karakter.'),
  pesan: z.string().min(5, 'Isi pesan pengumuman minimal 5 karakter.'),
  tipe: z.enum(['info', 'warning', 'danger', 'success']),
  app_id: z.string(),
  tautan_url: z.string().optional(),
  tautan_teks: z.string().optional(),
  mulai_pada: z.string().optional(),
  berakhir_pada: z.string().optional(),
  aktif: z.boolean(),
})

type AnnouncementFormValues = z.infer<typeof announcementSchema>

interface AnnouncementFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  announcementToEdit?: Announcement | null
  apps: App[]
  onSuccess: () => void
}

export function AnnouncementFormDialog({
  open,
  onOpenChange,
  announcementToEdit,
  apps,
  onSuccess,
}: AnnouncementFormDialogProps) {
  const isEditing = !!announcementToEdit

  const form = useForm<AnnouncementFormValues>({
    resolver: zodResolver(announcementSchema),
    defaultValues: {
      judul: '',
      pesan: '',
      tipe: 'info',
      app_id: 'all',
      tautan_url: '',
      tautan_teks: '',
      mulai_pada: '',
      berakhir_pada: '',
      aktif: true,
    },
  })

  useEffect(() => {
    if (announcementToEdit) {
      form.reset({
        judul: announcementToEdit.judul,
        pesan: announcementToEdit.pesan,
        tipe: announcementToEdit.tipe || 'info',
        app_id: announcementToEdit.app_id ? String(announcementToEdit.app_id) : 'all',
        tautan_url: announcementToEdit.tautan_url || '',
        tautan_teks: announcementToEdit.tautan_teks || '',
        mulai_pada: announcementToEdit.mulai_pada ? announcementToEdit.mulai_pada.slice(0, 16) : '',
        berakhir_pada: announcementToEdit.berakhir_pada ? announcementToEdit.berakhir_pada.slice(0, 16) : '',
        aktif: announcementToEdit.aktif ?? true,
      })
    } else {
      form.reset({
        judul: '',
        pesan: '',
        tipe: 'info',
        app_id: 'all',
        tautan_url: '',
        tautan_teks: '',
        mulai_pada: '',
        berakhir_pada: '',
        aktif: true,
      })
    }
  }, [announcementToEdit, form, open])

  const watchedValues = form.watch()

  const onSubmit = async (values: AnnouncementFormValues) => {
    try {
      const appIdNumber = values.app_id === 'all' ? null : Number(values.app_id)

      if (isEditing && announcementToEdit) {
        const payload: UpdateAnnouncementPayload = {
          ...values,
          app_id: appIdNumber,
          tautan_url: values.tautan_url || undefined,
          tautan_teks: values.tautan_teks || undefined,
          mulai_pada: values.mulai_pada ? new Date(values.mulai_pada).toISOString() : undefined,
          berakhir_pada: values.berakhir_pada ? new Date(values.berakhir_pada).toISOString() : undefined,
        }
        await announcementsApi.updateAnnouncement(announcementToEdit.id, payload)
        toast.success(`Pengumuman "${values.judul}" berhasil diperbarui!`)
      } else {
        const payload: CreateAnnouncementPayload = {
          ...values,
          app_id: appIdNumber,
          tautan_url: values.tautan_url || undefined,
          tautan_teks: values.tautan_teks || undefined,
          mulai_pada: values.mulai_pada ? new Date(values.mulai_pada).toISOString() : undefined,
          berakhir_pada: values.berakhir_pada ? new Date(values.berakhir_pada).toISOString() : undefined,
        }
        await announcementsApi.createAnnouncement(payload)
        toast.success(`Pengumuman "${values.judul}" berhasil diterbitkan!`)
      }

      onSuccess()
      onOpenChange(false)
    } catch (err: any) {
      toast.error('Gagal Menyimpan Pengumuman', {
        description: err?.message || 'Terjadi kesalahan.',
      })
    }
  }

  const getPreviewBannerStyle = (tipe: AnnouncementType) => {
    switch (tipe) {
      case 'danger':
        return 'bg-destructive/10 border-destructive/30 text-destructive'
      case 'warning':
        return 'bg-amber-500/10 border-amber-500/30 text-amber-800 dark:text-amber-400'
      case 'success':
        return 'bg-emerald-500/10 border-emerald-500/30 text-emerald-800 dark:text-emerald-400'
      case 'info':
      default:
        return 'bg-blue-500/10 border-blue-500/30 text-blue-800 dark:text-blue-400'
    }
  }

  const getPreviewIcon = (tipe: AnnouncementType) => {
    switch (tipe) {
      case 'danger':
        return <AlertOctagon className='size-5 text-destructive shrink-0 mt-0.5' />
      case 'warning':
        return <AlertTriangle className='size-5 text-amber-500 shrink-0 mt-0.5' />
      case 'success':
        return <CheckCircle2 className='size-5 text-emerald-500 shrink-0 mt-0.5' />
      case 'info':
      default:
        return <Info className='size-5 text-blue-500 shrink-0 mt-0.5' />
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-w-2xl max-h-[90vh] overflow-y-auto'>
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Edit Pengumuman Broadcast' : 'Buat Pengumuman Baru'}
          </DialogTitle>
          <DialogDescription>
            Siarkan banner pengumuman pemeliharaan, rilis fitur baru, atau informasi penting kepada pengguna portal.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4 py-2'>
            {/* Live Banner Preview */}
            <div className='space-y-1.5'>
              <span className='text-xs font-semibold text-muted-foreground uppercase tracking-wider'>
                Live Preview Tampilan Banner
              </span>
              <div
                className={`p-4 rounded-xl border flex items-start gap-3 transition-colors ${getPreviewBannerStyle(
                  watchedValues.tipe
                )}`}
              >
                {getPreviewIcon(watchedValues.tipe)}
                <div className='flex-1 min-w-0 space-y-1'>
                  <h4 className='font-semibold text-sm leading-tight'>
                    {watchedValues.judul || 'Judul Pengumuman / Pemeliharaan'}
                  </h4>
                  <p className='text-xs opacity-90 leading-relaxed'>
                    {watchedValues.pesan ||
                      'Pesan pengumuman atau rincian jadwal maintenance akan ditampilkan di sini kepada pengguna.'}
                  </p>
                  {watchedValues.tautan_url && (
                    <div className='pt-1'>
                      <span className='inline-flex items-center text-xs font-semibold underline gap-1'>
                        {watchedValues.tautan_teks || 'Lihat Selengkapnya'}
                        <ExternalLink className='size-3' />
                      </span>
                    </div>
                  )}
                </div>
              </div>
            </div>

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4 pt-2'>
              <FormField
                control={form.control}
                name='judul'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Judul Pengumuman *</FormLabel>
                    <FormControl>
                      <Input placeholder='Pemeliharaan Server SIMPEG' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='tipe'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tipe & Warna Banner</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder='Pilih Tipe' />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='info'>ℹ️ Informasi (Biru)</SelectItem>
                        <SelectItem value='warning'>⚠️ Peringatan / Maintenance (Kuning)</SelectItem>
                        <SelectItem value='danger'>🚨 Gangguan Kritis (Merah)</SelectItem>
                        <SelectItem value='success'>✅ Pengumuman Sukses (Hijau)</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name='pesan'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Isi Pesan Pengumuman *</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder='Tuliskan informasi lengkap mengenai pengumuman atau jadwal pemeliharaan...'
                      rows={3}
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
                name='app_id'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Target Cakupan Aplikasi</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder='Semua Aplikasi' />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='all'>🌐 Semua Aplikasi (Global Portal)</SelectItem>
                        {apps.map((app) => (
                          <SelectItem key={app.id} value={String(app.id)}>
                            📱 {app.nama}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='aktif'
                render={({ field }) => (
                  <FormItem className='flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs'>
                    <div className='space-y-0.5'>
                      <FormLabel className='text-sm font-medium'>
                        Status Siaran Aktif
                      </FormLabel>
                      <FormDescription className='text-xs'>
                        Tayangkan banner saat ini
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
            </div>

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='tautan_url'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tautan URL Aksi (Opsional)</FormLabel>
                    <FormControl>
                      <Input placeholder='https://layanan.bps.go.id/info' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='tautan_teks'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Label Teks Tombol Aksi</FormLabel>
                    <FormControl>
                      <Input placeholder='Contoh: Baca Selengkapnya' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='mulai_pada'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Jadwal Mulai Tayang (Opsional)</FormLabel>
                    <FormControl>
                      <Input type='datetime-local' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='berakhir_pada'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Jadwal Berakhir Tayang (Opsional)</FormLabel>
                    <FormControl>
                      <Input type='datetime-local' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

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
                {isEditing ? 'Simpan Perubahan' : 'Terbitkan Pengumuman'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
