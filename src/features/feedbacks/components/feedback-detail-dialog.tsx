import { useState, useEffect } from 'react'
import {
  MessageSquareWarning,
  CheckCircle2,
  Clock,
  AlertCircle,
  XCircle,
  ExternalLink,
  User,
  Mail,
  Phone,
  Layers,
  Loader2,
  FileText,
  ShieldAlert,
} from 'lucide-react'
import { toast } from 'sonner'
import { Feedback, FeedbackStatus, UpdateFeedbackStatusPayload } from '@/types/feedbacks'
import { feedbacksApi } from '../api/feedbacks-api'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

interface FeedbackDetailDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  feedback: Feedback | null
  onSuccess: () => void
}

export function FeedbackDetailDialog({
  open,
  onOpenChange,
  feedback,
  onSuccess,
}: FeedbackDetailDialogProps) {
  const [status, setStatus] = useState<FeedbackStatus>('pending')
  const [adminNotes, setAdminNotes] = useState('')
  const [isUpdating, setIsUpdating] = useState(false)

  useEffect(() => {
    if (feedback) {
      setStatus(feedback.status || 'pending')
      setAdminNotes(feedback.tanggapan_admin || '')
    }
  }, [feedback, open])

  const handleUpdate = async () => {
    if (!feedback) return
    setIsUpdating(true)

    try {
      const payload: UpdateFeedbackStatusPayload = {
        status,
        tanggapan_admin: adminNotes.trim() || undefined,
      }
      await feedbacksApi.updateFeedbackStatus(feedback.id, payload)
      toast.success('Status tiket laporan berhasil diperbarui!')
      onSuccess()
      onOpenChange(false)
    } catch (err: any) {
      toast.error('Gagal Memperbarui Tiket', {
        description: err?.message || 'Terjadi kesalahan.',
      })
    } finally {
      setIsUpdating(false)
    }
  }

  const getStatusBadge = (st: FeedbackStatus) => {
    switch (st) {
      case 'resolved':
        return (
          <Badge variant='outline' className='bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-xs px-2.5 py-1'>
            <CheckCircle2 className='mr-1.5 size-3.5' /> Selesai Ditindaklanjuti
          </Badge>
        )
      case 'in_progress':
        return (
          <Badge variant='outline' className='bg-blue-500/10 text-blue-600 border-blue-500/30 text-xs px-2.5 py-1'>
            <Clock className='mr-1.5 size-3.5' /> Sedang Diproses
          </Badge>
        )
      case 'closed':
        return (
          <Badge variant='outline' className='bg-muted text-muted-foreground text-xs px-2.5 py-1'>
            <XCircle className='mr-1.5 size-3.5' /> Ditutup / Ditolak
          </Badge>
        )
      case 'pending':
      default:
        return (
          <Badge variant='outline' className='bg-amber-500/10 text-amber-600 border-amber-500/30 text-xs px-2.5 py-1'>
            <AlertCircle className='mr-1.5 size-3.5' /> Menunggu Respon
          </Badge>
        )
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-3xl lg:max-w-4xl max-h-[90vh] overflow-y-auto p-6 md:p-8'>
        <DialogHeader className='space-y-1.5 pb-2 border-b'>
          <div className='flex items-center gap-2 text-primary'>
            <MessageSquareWarning className='size-5' />
            <DialogTitle className='text-xl font-bold tracking-tight'>
              Tiket Laporan & Masukan Pengguna #{feedback?.id}
            </DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            Diterima pada {feedback?.created_at ? new Date(feedback.created_at).toLocaleString('id-ID') : '-'}
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-5 py-3'>
          {/* Status and Category Summary Card */}
          <div className='flex flex-wrap items-center justify-between gap-3 p-4 bg-muted/30 rounded-xl border'>
            <div className='flex items-center gap-2.5'>
              {feedback && getStatusBadge(feedback.status)}
              <Badge variant='secondary' className='text-xs font-semibold uppercase px-2.5 py-1'>
                {feedback?.kategori || 'Masukan'}
              </Badge>
            </div>
            {feedback?.app && (
              <div className='flex items-center gap-2 text-xs font-semibold text-foreground bg-background px-3 py-1.5 rounded-lg border'>
                <Layers className='size-3.5 text-primary' />
                <span>Modul: {feedback.app.nama}</span>
              </div>
            )}
          </div>

          {/* Reporter Info Card */}
          <div className='grid grid-cols-1 sm:grid-cols-3 gap-3 p-4 bg-card rounded-xl border text-xs shadow-xs'>
            <div className='flex items-center gap-2.5'>
              <div className='p-2 rounded-lg bg-primary/10 text-primary'>
                <User className='size-4' />
              </div>
              <div className='min-w-0'>
                <p className='text-muted-foreground text-[11px]'>Nama Pelapor</p>
                <p className='font-semibold text-foreground truncate'>
                  {feedback?.nama_pelapor || 'Anonim'}
                </p>
              </div>
            </div>

            <div className='flex items-center gap-2.5'>
              <div className='p-2 rounded-lg bg-blue-500/10 text-blue-600'>
                <Mail className='size-4' />
              </div>
              <div className='min-w-0'>
                <p className='text-muted-foreground text-[11px]'>Email</p>
                <p className='font-semibold text-foreground truncate font-mono'>
                  {feedback?.email_pelapor || '-'}
                </p>
              </div>
            </div>

            <div className='flex items-center gap-2.5'>
              <div className='p-2 rounded-lg bg-emerald-500/10 text-emerald-600'>
                <Phone className='size-4' />
              </div>
              <div className='min-w-0'>
                <p className='text-muted-foreground text-[11px]'>Kontak / WhatsApp</p>
                <p className='font-semibold text-foreground truncate font-mono'>
                  {feedback?.kontak_pelapor || '-'}
                </p>
              </div>
            </div>
          </div>

          {/* Message Content */}
          <div className='space-y-2'>
            <span className='text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5'>
              <FileText className='size-3.5 text-primary' />
              Uraian Laporan Masalah / Saran
            </span>
            <div className='p-4 rounded-xl bg-muted/20 border space-y-2'>
              <h4 className='font-bold text-sm text-foreground'>{feedback?.judul}</h4>
              <p className='text-xs text-muted-foreground leading-relaxed whitespace-pre-line'>
                {feedback?.isi_laporan}
              </p>
            </div>
          </div>

          {/* Attachment Preview if any */}
          {feedback?.lampiran_url && (
            <div className='space-y-2'>
              <span className='text-xs font-semibold uppercase tracking-wider text-muted-foreground'>
                Lampiran Bukti Screenshot:
              </span>
              <div className='rounded-xl border p-3 bg-muted/20 space-y-2'>
                <img
                  src={feedback.lampiran_url}
                  alt='Lampiran'
                  className='max-h-72 rounded-lg object-contain bg-white p-1 border shadow-xs'
                />
                <a
                  href={feedback.lampiran_url}
                  target='_blank'
                  rel='noopener noreferrer'
                  className='inline-flex items-center gap-1.5 text-xs text-primary font-medium hover:underline'
                >
                  Buka Gambar Resolusi Penuh <ExternalLink className='size-3' />
                </a>
              </div>
            </div>
          )}

          {/* Admin Response & Action Form */}
          <div className='space-y-4 pt-3 border-t'>
            <span className='text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5'>
              <ShieldAlert className='size-3.5 text-primary' />
              Tindak Lanjut & Tanggapan Admin
            </span>

            <div className='grid grid-cols-1 sm:grid-cols-2 gap-4'>
              <div className='space-y-1.5'>
                <label className='text-xs font-medium text-foreground block'>
                  Status Progres Tiket
                </label>
                <Select
                  value={status}
                  onValueChange={(val) => setStatus(val as FeedbackStatus)}
                >
                  <SelectTrigger className='h-10 text-xs'>
                    <SelectValue placeholder='Pilih Status' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='pending'>🟡 Menunggu Respon (Pending)</SelectItem>
                    <SelectItem value='in_progress'>🔵 Sedang Ditangani (In Progress)</SelectItem>
                    <SelectItem value='resolved'>🟢 Selesai Ditindaklanjuti (Resolved)</SelectItem>
                    <SelectItem value='closed'>⚪ Ditutup / Ditolak (Closed)</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className='space-y-1.5'>
              <label className='text-xs font-medium text-foreground block'>
                Catatan Penyelesaian / Solusi Admin
              </label>
              <Textarea
                value={adminNotes}
                onChange={(e) => setAdminNotes(e.target.value)}
                placeholder='Tuliskan catatan perbaikan teknis, investigasi log, atau instruksi tindak lanjut...'
                rows={4}
                className='text-xs leading-relaxed resize-y'
              />
            </div>
          </div>
        </div>

        <DialogFooter className='gap-2.5 pt-4 border-t'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            disabled={isUpdating}
            className='h-10 px-4 text-xs'
          >
            Tutup
          </Button>
          <Button
            type='button'
            onClick={handleUpdate}
            disabled={isUpdating}
            className='h-10 px-5 text-xs font-semibold shadow-xs'
          >
            {isUpdating && <Loader2 className='mr-2 size-4 animate-spin' />}
            Simpan Status Tiket
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
