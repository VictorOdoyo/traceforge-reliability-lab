export const durationMs = (nanoseconds: number) => nanoseconds / 1_000_000
export const formatDuration = (nanoseconds: number) => `${durationMs(nanoseconds).toFixed(0)} ms`
export const formatPercent = (value: number, digits = 1) => `${value.toFixed(digits)}%`
export const formatTime = (value: string) => new Intl.DateTimeFormat('en', { hour: '2-digit', minute: '2-digit' }).format(new Date(value))
export const formatRelative = (value: string) => {
  const minutes = Math.max(0, Math.round((Date.now() - new Date(value).getTime()) / 60_000))
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  return `${Math.round(minutes / 60)}h ago`
}
export const humanize = (value: string) => value.replace(/([a-z])([A-Z])/g, '$1 $2').replaceAll('-', ' ')
