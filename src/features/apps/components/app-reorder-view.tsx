import { useState, useEffect } from 'react'
import { ArrowUp, ArrowDown, GripVertical, Check, RotateCcw, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { App } from '@/types/apps'
import { appsApi } from '../api/apps-api'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

interface AppReorderViewProps {
  initialApps: App[]
  onSuccess: () => void
}

export function AppReorderView({ initialApps, onSuccess }: AppReorderViewProps) {
  const [items, setItems] = useState<App[]>([])
  const [isSaving, setIsSaving] = useState(false)
  const [hasChanges, setHasChanges] = useState(false)

  useEffect(() => {
    setItems([...initialApps])
    setHasChanges(false)
  }, [initialApps])

  const moveItem = (fromIndex: number, toIndex: number) => {
    if (toIndex < 0 || toIndex >= items.length) return
    const updated = [...items]
    const [moved] = updated.splice(fromIndex, 1)
    updated.splice(toIndex, 0, moved)
    setItems(updated)
    setHasChanges(true)
  }

  const handleSave = async () => {
    setIsSaving(true)
    try {
      const appIds = items.map((app) => app.id)
      await appsApi.reorderApps(appIds)
      toast.success('Urutan katalog aplikasi berhasil disimpan!')
      setHasChanges(false)
      onSuccess()
    } catch (err: any) {
      toast.error('Gagal Menyimpan Urutan', {
        description: err?.message || 'Terjadi kesalahan.',
      })
    } finally {
      setIsSaving(false)
    }
  }

  const handleReset = () => {
    setItems([...initialApps])
    setHasChanges(false)
  }

  return (
    <div className='space-y-4 py-2'>
      <div className='flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 p-4 rounded-xl bg-card border shadow-xs'>
        <div>
          <h3 className='text-sm font-semibold'>Mode Pengurutan Tampilan Aplikasi</h3>
          <p className='text-xs text-muted-foreground'>
            Gunakan tombol panah Naik / Turun untuk menyesuaikan prioritas posisi aplikasi pada katalog portal ALUSI.
          </p>
        </div>
        <div className='flex items-center gap-2 self-end sm:self-auto'>
          <Button
            variant='outline'
            size='sm'
            onClick={handleReset}
            disabled={!hasChanges || isSaving}
          >
            <RotateCcw className='mr-1.5 size-3.5' />
            Reset
          </Button>
          <Button
            size='sm'
            onClick={handleSave}
            disabled={!hasChanges || isSaving}
            className='bg-primary text-primary-foreground'
          >
            {isSaving ? (
              <Loader2 className='mr-1.5 size-3.5 animate-spin' />
            ) : (
              <Check className='mr-1.5 size-3.5' />
            )}
            Simpan Urutan
          </Button>
        </div>
      </div>

      <div className='space-y-2 max-h-[65vh] overflow-y-auto pr-1'>
        {items.map((app, index) => (
          <Card
            key={app.id}
            className='transition-all duration-150 hover:border-primary/40 border shadow-xs'
          >
            <CardContent className='p-3 flex items-center justify-between gap-3'>
              <div className='flex items-center gap-3 min-w-0'>
                <div className='flex items-center text-muted-foreground cursor-grab active:cursor-grabbing px-1'>
                  <GripVertical className='size-4' />
                  <span className='font-mono text-xs font-bold w-6 text-center text-foreground'>
                    #{index + 1}
                  </span>
                </div>

                {app.ikon_url ? (
                  <img
                    src={app.ikon_url}
                    alt={app.nama}
                    className='size-9 rounded-lg object-contain bg-white p-1 border shadow-xs shrink-0'
                  />
                ) : (
                  <div className='size-9 rounded-lg bg-primary/10 text-primary flex items-center justify-center font-bold text-xs shrink-0'>
                    {app.nama.substring(0, 2).toUpperCase()}
                  </div>
                )}

                <div className='min-w-0'>
                  <div className='flex items-center gap-2'>
                    <h4 className='text-sm font-semibold truncate'>{app.nama}</h4>
                    {app.category && (
                      <Badge variant='outline' className='text-[10px] hidden sm:inline-flex'>
                        {app.category.nama}
                      </Badge>
                    )}
                  </div>
                  <p className='text-xs text-muted-foreground truncate'>{app.url}</p>
                </div>
              </div>

              <div className='flex items-center gap-1 shrink-0'>
                <Button
                  variant='ghost'
                  size='icon'
                  className='size-7'
                  disabled={index === 0}
                  onClick={() => moveItem(index, index - 1)}
                  title='Pindahkan ke atas'
                >
                  <ArrowUp className='size-4' />
                </Button>
                <Button
                  variant='ghost'
                  size='icon'
                  className='size-7'
                  disabled={index === items.length - 1}
                  onClick={() => moveItem(index, index + 1)}
                  title='Pindahkan ke bawah'
                >
                  <ArrowDown className='size-4' />
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
