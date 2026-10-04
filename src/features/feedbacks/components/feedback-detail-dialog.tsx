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
          <Badge variant='outline' className='bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-xs'>
            <CheckCircle2 className='mr-1 size-3' /> Selesai Ditindaklanjuti
          </Badge>
        )
      case 'in_progress':
        return (
          <Badge variant='outline' className='bg-blue-500/10 text-blue-600 border-blue-500/30 text-xs'>
            <Clock className='mr-1 size-3' /> Sedang Diproses
          </Badge>
        )
      case 'closed':
        return (
          <Badge variant='outline' className='bg-muted text-muted-foreground text-xs'>
            <XCircle className='mr-1 size-3' /> Ditutup
          </Badge>
        )
      case 'pending':
      default:
        return (
          <Badge variant='outline' className='bg-amber-500/10 text-amber-600 border-amber-500/30 text-xs'>
            <AlertCircle className='mr-1 size-3' /> Menunggu Respon
          </Badge>
        )
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-w-2xl max-h-[90vh] overflow-y-auto'>
        <DialogHeader>
          <div className='flex items-center gap-2 text-primary mb-1'>
            <MessageSquareWarning className='size-5' />
            <DialogTitle>Detail Laporan / Masukan Pengguna</DialogTitle>
          </div>
          <DialogDescription>
            Tiket #{feedback?.id} - Diterima pada{' '}
            {feedback?.created_at ? new Date(feedback.created_at).toLocaleString('id-ID') : '-'}
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-4 py-2'>
          {/* Status and Category Badges */}
          <div className='flex flex-wrap items-center justify-between gap-2 p-3 bg-muted/40 rounded-xl border'>
            <div className='flex items-center gap-2'>
              {feedback && getStatusBadge(feedback.status)}
              <Badge variant='secondary' className='text-xs uppercase'>
                {feedback?.kategori || 'Masukan'}
              </Badge>
            </div>
            {feedback?.app && (
              <div className='flex items-center gap-1.5 text-xs font-semibold text-foreground'>
                <Layers className='size-3.5 text-muted-foreground' />
                <span>{feedback.app.nama}</span>
              </div>
            )}
          </div>

          {/* Reporter Info */}
          <div className='grid grid-cols-1 sm:grid-cols-3 gap-3 p-3 bg-card rounded-xl border text-xs'>
            <div className='flex items-center gap-2'>
              <User className='size-4 text-muted-foreground shrink-0' />
              <div>
                <p className='text-muted-foreground'>Nama Pelapor</p>
                <p className='font-semibold text-foreground'>{feedback?.nama_pelapor || 'Anonim'}</p>
              </div>
            </div>
            <div className='flex items-center gap-2'>
              <Mail className='size-4 text-muted-foreground shrink-0' />
              <div className='min-w-0'>
                <p className='text-muted-foreground'>Email</p>
                <p className='font-semibold text-foreground truncate'>{feedback?.email_pelapor || '-'}</p>
              </div>
            </div>
            <div className='flex items-center gap-2'>
              <Phone className='size-4 text-muted-foreground shrink-0' />
              <div>
                <p className='text-muted-foreground'>Kontak / WA</p>
                <p className='font-semibold text-foreground'>{feedback?.kontak_pelapor || '-'}</p>
              </div>
            </div>
          </div>

          {/* Title & Message Content */}
          <div className='space-y-1.5'>
            <h4 className='font-semibold text-sm text-foreground'>{feedback?.judul}</h4>
            <div className='p-3.5 rounded-xl bg-muted/30 border text-sm text-muted-foreground leading-relaxed whitespace-pre-line'>
              {feedback?.isi_laporan}
            </div>
          </div>

          {/* Attachment Preview if any */}
          {feedback?.lampiran_url && (
            <div className='space-y-1.5'>
              <span className='text-xs font-semibold text-muted-foreground'>Lampiran Gambar:</span>
              <div className='rounded-lg border p-2 bg-muted/20'>
                <img
                  src={feedback.lampiran_url}
                  alt='Lampiran'
                  className='max-h-60 rounded object-contain'
                />
                <a
                  href={feedback.lampiran_url}
                  target='_blank'
                  rel='noopener noreferrer'
                  className='inline-flex items-center gap-1 text-xs text-primary hover:underline mt-2'
                >
                  Buka Gambar Asli <ExternalLink className='size-3' />
                </a>
              </div>
            </div>
          )}

          {/* Admin Response & Status Form */}
          <div className='space-y-3 pt-3 border-t'>
            <h4 className='text-xs font-semibold uppercase tracking-wider text-muted-foreground'>
              Tindak Lanjut & Respon Admin
            </h4>

            <div className='grid grid-cols-1 sm:grid-cols-2 gap-3'>
              <div>
                <label className='text-xs font-medium mb-1 block'>Status Tiket</label>
                <Select
                  value={status}
                  onValueChange={(val) => setStatus(val as FeedbackStatus)}
                >
                  <SelectTrigger className='text-xs h-9'>
                    <SelectValue placeholder='Pilih Status' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='pending'>🟡 Pending / Baru</SelectItem>
                    <SelectItem value='in_progress'>🔵 Sedang Diproses</SelectItem>
                    <SelectItem value='resolved'>🟢 Selesai Ditindaklanjuti</SelectItem>
                    <SelectItem value='closed'>⚪ Ditutup / Ditolak</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div>
              <label className='text-xs font-medium mb-1 block'>
                Catatan Respon / Solusi Admin
              </label>
              <Textarea
                value={adminNotes}
                onChange={(e) => setAdminNotes(e.target.value)}
                placeholder='Tuliskan catatan perbaikan atau konfirmasi solusi untuk tiket ini...'
                rows={3}
              />
            </div>
          </div>
        </div>

        <DialogFooter className='gap-2 pt-2'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            disabled={isUpdating}
          >
            Tutup
          </Button>
          <Button type='button' onClick={handleUpdate} disabled={isUpdating}>
            {isUpdating && <Loader2 className='mr-2 size-4 animate-spin' />}
            Simpan Tindak Lanjut
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
