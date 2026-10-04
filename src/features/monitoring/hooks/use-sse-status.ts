import { useEffect, useState, useRef, useCallback } from 'react'
import { toast } from 'sonner'
import { ServiceUptimeSummary, SSEStatusEvent } from '@/types/monitoring'
import { monitoringApi } from '../api/monitoring-api'

export function useSSEStatus() {
  const [summary, setSummary] = useState<ServiceUptimeSummary | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null)
  const eventSourceRef = useRef<EventSource | null>(null)

  const fetchInitialSummary = useCallback(async () => {
    try {
      const data = await monitoringApi.getServiceStatusSummary(30)
      setSummary(data)
      setLastUpdated(new Date())
    } catch {
      // Fallback
    }
  }, [])

  useEffect(() => {
    fetchInitialSummary()

    const sseURL = monitoringApi.getSSEStreamURL()
    const es = new EventSource(sseURL)
    eventSourceRef.current = es

    es.onopen = () => {
      setIsConnected(true)
    }

    es.addEventListener('snapshot', (event: MessageEvent) => {
      try {
        const payload: SSEStatusEvent = JSON.parse(event.data)
        if (payload.services) {
          setSummary(payload.services)
          setLastUpdated(new Date())
        }
      } catch {
        // parse error
      }
    })

    es.addEventListener('status_change', (event: MessageEvent) => {
      try {
        const payload: SSEStatusEvent = JSON.parse(event.data)
        setLastUpdated(new Date())

        if (payload.status === 'kendala' || payload.status === 'offline') {
          toast.warning(`Peringatan Status Layanan!`, {
            description: `Aplikasi #${payload.app_id} terdeteksi berstatus ${payload.status?.toUpperCase()}.`,
          })
        } else if (payload.status === 'online') {
          toast.success(`Layanan Pulih`, {
            description: `Aplikasi #${payload.app_id} kembali berstatus ONLINE.`,
          })
        }

        // Refresh snapshot
        fetchInitialSummary()
      } catch {
        // parse error
      }
    })

    es.onerror = () => {
      setIsConnected(false)
    }

    return () => {
      es.close()
      eventSourceRef.current = null
    }
  }, [fetchInitialSummary])

  return {
    summary,
    isConnected,
    lastUpdated,
    refresh: fetchInitialSummary,
  }
}
