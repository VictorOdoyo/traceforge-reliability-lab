import { useCallback, useEffect, useState } from 'react'
import { loadAppData } from '../lib/api'
import type { AppData } from '../types'

export function useAppData() {
  const [data, setData] = useState<AppData | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [refreshing, setRefreshing] = useState(false)

  const refresh = useCallback(async () => {
    setRefreshing(true)
    try { setData(await loadAppData()); setError(null) }
    catch (caught) { setError(caught instanceof Error ? caught.message : 'Unable to load reliability data') }
    finally { setRefreshing(false) }
  }, [])

  useEffect(() => {
    const controller = new AbortController()
    loadAppData(controller.signal).then(setData).catch((caught: unknown) => {
      if (!controller.signal.aborted) setError(caught instanceof Error ? caught.message : 'Unable to load reliability data')
    })
    return () => controller.abort()
  }, [])

  return { data, error, refreshing, refresh, setData }
}
