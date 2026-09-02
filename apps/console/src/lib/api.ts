import type { Alert, AppData, DeploymentCorrelation, DependencyEdge, Incident, Overview, SLO, Trace } from '../types'

const baseUrl = import.meta.env.VITE_API_URL ?? ''

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${baseUrl}${path}`, { ...init, headers: { Accept: 'application/json', 'Content-Type': 'application/json', ...init?.headers } })
  if (!response.ok) {
    const problem = await response.json().catch(() => null) as { detail?: string } | null
    throw new Error(problem?.detail ?? `TraceForge API returned ${response.status}`)
  }
  return response.json() as Promise<T>
}

export async function loadAppData(signal?: AbortSignal): Promise<AppData> {
  const options = { signal }
  const [overview, traces, slos, topology, alerts, incidents, deployments] = await Promise.all([
    request<Overview>('/api/v1/overview', options), request<Trace[]>('/api/v1/traces?limit=40', options),
    request<SLO[]>('/api/v1/slos', options), request<DependencyEdge[]>('/api/v1/topology', options),
    request<Alert[]>('/api/v1/alerts', options), request<Incident[]>('/api/v1/incidents', options),
    request<DeploymentCorrelation[]>('/api/v1/deployments/correlations', options),
  ])
  return { overview, traces, slos, topology, alerts, incidents, deployments }
}

export const acknowledgeAlert = (id: string, owner: string) => request<Alert>(`/api/v1/alerts/${id}/acknowledge`, { method: 'POST', body: JSON.stringify({ owner }) })
export const transitionIncident = (id: string, status: Incident['status'], actor: string) => request<Incident>(`/api/v1/incidents/${id}/transition`, { method: 'POST', body: JSON.stringify({ status, actor }) })
