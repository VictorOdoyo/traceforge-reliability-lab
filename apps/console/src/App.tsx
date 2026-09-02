import { useState } from 'react'
import './App.css'
import { acknowledgeAlert, transitionIncident } from './lib/api'
import { useAppData } from './hooks/useAppData'
import type { Incident } from './types'
import { Sidebar } from './components/Sidebar'
import { WorkspaceHeader } from './components/WorkspaceHeader'
import { SummaryMetrics } from './components/SummaryMetrics'
import { LatencyChart } from './components/LatencyChart'
import { ServiceHealth } from './components/ServiceHealth'
import { ErrorBudgets } from './components/ErrorBudgets'
import { TraceExplorer } from './components/TraceExplorer'
import { ServiceMap } from './components/ServiceMap'
import { AlertQueue } from './components/AlertQueue'
import { IncidentCommand } from './components/IncidentCommand'
import { DeploymentEvidence } from './components/DeploymentEvidence'

function App() {
  const { data, error, refreshing, refresh, setData } = useAppData()
  const [active, setActive] = useState('overview')
  const [window, setWindow] = useState('1h')
  const [notice, setNotice] = useState<string | null>(null)
  const navigate = (id: string) => { setActive(id); document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' }) }
  const ownAlert = async (id: string) => { try { const updated = await acknowledgeAlert(id, 'Victor'); setData(current => current ? { ...current, alerts: current.alerts.map(alert => alert.id === id ? updated : alert) } : current); setNotice('Alert assigned to Victor') } catch (caught) { setNotice(caught instanceof Error ? caught.message : 'Alert could not be acknowledged') } }
  const advanceIncident = async (id: string, status: Incident['status']) => { try { const updated = await transitionIncident(id, status, 'Victor'); setData(current => current ? { ...current, incidents: current.incidents.map(incident => incident.id === id ? updated : incident) } : current); setNotice(`Incident moved to ${status}`) } catch (caught) { setNotice(caught instanceof Error ? caught.message : 'Incident could not be updated') } }

  if (!data && !error) return <div className="loading-screen"><div className="pulse-logo">TF</div><strong>Loading reliability evidence</strong><span>Connecting traces, SLOs, and active response…</span></div>
  if (!data) return <div className="loading-screen error-screen"><strong>TraceForge could not load</strong><span>{error}</span><button className="primary-button" onClick={refresh}>Try again</button></div>

  return <div className="app-shell"><Sidebar active={active} onNavigate={navigate}/><main className="workspace" id="overview"><WorkspaceHeader window={window} onWindowChange={setWindow} refreshing={refreshing} onRefresh={refresh}/>{error && <div className="error-banner">Latest refresh failed: {error}</div>}<SummaryMetrics overview={data.overview}/><div className="dashboard-grid lead-grid"><LatencyChart traces={data.traces}/><ErrorBudgets slos={data.slos}/></div><ServiceHealth services={data.overview.services}/><div className="dashboard-grid investigation-grid"><ServiceMap services={data.overview.services} dependencies={data.topology}/><TraceExplorer traces={data.traces}/></div><div className="dashboard-grid operations-grid"><AlertQueue alerts={data.alerts} onAcknowledge={ownAlert}/><DeploymentEvidence deployments={data.deployments}/></div><IncidentCommand incidents={data.incidents} onTransition={advanceIncident}/><footer><span>TraceForge Reliability Lab</span><span>Synthetic production telemetry · No customer data</span></footer></main>{notice && <div className="toast" role="status"><span>{notice}</span><button aria-label="Dismiss notification" onClick={() => setNotice(null)}>×</button></div>}</div>
}

export default App
