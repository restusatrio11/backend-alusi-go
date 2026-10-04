import { useState } from 'react'
import { Download, FileSpreadsheet, Loader2 } from 'lucide-react'
import { reportsApi } from '@/features/reports/api/reports-api'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { toast } from 'sonner'

export function ExportReportButtons() {
  const [isExportingCatalog, setIsExportingCatalog] = useState(false)
  const [isExportingAnalytics, setIsExportingAnalytics] = useState(false)

  const handleExportCatalog = async () => {
    try {
      setIsExportingCatalog(true)
      await reportsApi.downloadCatalogCSV()
      toast.success('Ekspor CSV katalog aplikasi berhasil diunduh.')
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal mengekspor data katalog aplikasi.')
    } finally {
      setIsExportingCatalog(false)
    }
  }

  const handleExportAnalytics = async () => {
    try {
      setIsExportingAnalytics(true)
      await reportsApi.downloadAnalyticsCSV()
      toast.success('Ekspor CSV analitik portal berhasil diunduh.')
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal mengekspor data analitik portal.')
    } finally {
      setIsExportingAnalytics(false)
    }
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant='outline' size='sm' className='gap-1.5 h-8 text-xs font-medium'>
          {isExportingCatalog || isExportingAnalytics ? (
            <Loader2 className='size-3.5 animate-spin' />
          ) : (
            <Download className='size-3.5' />
          )}
          <span>Ekspor Laporan CSV</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end' className='w-56'>
        <DropdownMenuLabel className='text-xs'>Pilihan Format Data</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onClick={handleExportCatalog}
          disabled={isExportingCatalog}
          className='cursor-pointer text-xs'
        >
          <FileSpreadsheet className='size-4 mr-2 text-emerald-600' />
          <span>Katalog Aplikasi (.csv)</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={handleExportAnalytics}
          disabled={isExportingAnalytics}
          className='cursor-pointer text-xs'
        >
          <FileSpreadsheet className='size-4 mr-2 text-blue-600' />
          <span>Analitik & SLA (.csv)</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
