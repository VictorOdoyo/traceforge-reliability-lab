import { afterEach, describe, expect, it, vi } from 'vitest'
import { durationMs, formatDuration, formatPercent, formatRelative, formatTime, humanize } from './format'

describe('telemetry formatting', () => {
  afterEach(() => vi.useRealTimers())

  it('converts Go durations from nanoseconds', () => {
    expect(durationMs(425_000_000)).toBe(425)
    expect(formatDuration(425_000_000)).toBe('425 ms')
  })

  it('formats percentages and identifiers', () => {
    expect(formatPercent(99.951, 2)).toBe('99.95%')
    expect(humanize('awaiting-Review')).toBe('awaiting Review')
  })

  it('formats absolute and relative times', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-09-02T08:30:00Z'))

    expect(formatTime('2026-09-02T08:10:00Z')).toMatch(/\d{2}:\d{2}/)
    expect(formatRelative('2026-09-02T08:30:00Z')).toBe('just now')
    expect(formatRelative('2026-09-02T08:10:00Z')).toBe('20m ago')
    expect(formatRelative('2026-09-02T05:30:00Z')).toBe('3h ago')
  })
})
