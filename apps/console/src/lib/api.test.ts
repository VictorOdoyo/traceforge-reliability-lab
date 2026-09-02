import { afterEach, describe, expect, it, vi } from 'vitest'
import { acknowledgeAlert, loadAppData, transitionIncident } from './api'

describe('reliability API client', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('aggregates every dashboard resource', async () => {
    const fetchMock = vi.fn((input: string | URL | Request) => {
      const path = new URL(String(input), 'http://localhost').pathname
      const body = path.endsWith('overview') ? { serviceCount: 5 } : []
      return Promise.resolve(new Response(JSON.stringify(body), { status: 200 }))
    })
    vi.stubGlobal('fetch', fetchMock)

    const result = await loadAppData()

    expect(fetchMock).toHaveBeenCalledTimes(7)
    expect(result.overview.serviceCount).toBe(5)
    expect(result.traces).toEqual([])
  })

  it('posts alert ownership', async () => {
    const fetchMock = vi.fn(() => Promise.resolve(new Response(JSON.stringify({ id: 'alert-1', state: 'acknowledged' }), { status: 200 })))
    vi.stubGlobal('fetch', fetchMock)

    await acknowledgeAlert('alert-1', 'Victor')

    expect(fetchMock).toHaveBeenCalledWith(expect.stringContaining('/alert-1/acknowledge'), expect.objectContaining({ method: 'POST', body: JSON.stringify({ owner: 'Victor' }) }))
  })

  it('posts incident transitions', async () => {
    const fetchMock = vi.fn(() => Promise.resolve(new Response(JSON.stringify({ id: 'incident-1', status: 'resolved' }), { status: 200 })))
    vi.stubGlobal('fetch', fetchMock)

    await transitionIncident('incident-1', 'resolved', 'Victor')

    expect(fetchMock).toHaveBeenCalledWith(expect.stringContaining('/incident-1/transition'), expect.objectContaining({ method: 'POST', body: JSON.stringify({ status: 'resolved', actor: 'Victor' }) }))
  })

  it('surfaces problem details', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(new Response(JSON.stringify({ detail: 'storage unavailable' }), { status: 503 }))))
    await expect(loadAppData()).rejects.toThrow('storage unavailable')
  })

  it('falls back to the response status for invalid problem bodies', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(new Response('unavailable', { status: 502 }))))
    await expect(loadAppData()).rejects.toThrow('TraceForge API returned 502')
  })
})
