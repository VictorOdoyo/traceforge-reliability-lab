import { Activity, BellRing, Boxes, Gauge, GitCommitHorizontal, Network, Siren, Waves } from 'lucide-react'

const items = [
  { id: 'overview', label: 'Overview', icon: Gauge }, { id: 'services', label: 'Services', icon: Boxes },
  { id: 'traces', label: 'Traces', icon: Activity }, { id: 'topology', label: 'Topology', icon: Network },
  { id: 'slos', label: 'SLOs', icon: Waves }, { id: 'alerts', label: 'Alerts', icon: BellRing },
  { id: 'incidents', label: 'Incidents', icon: Siren }, { id: 'deployments', label: 'Deployments', icon: GitCommitHorizontal },
]

export function Sidebar({ active, onNavigate }: { active: string; onNavigate: (id: string) => void }) {
  return <aside className="sidebar">
    <div className="brand"><div className="brand-mark"><Activity size={20} /></div><div><strong>TraceForge</strong><span>Reliability Lab</span></div></div>
    <nav aria-label="Reliability sections">{items.map(({ id, label, icon: Icon }) => <button key={id} className={active === id ? 'nav-item active' : 'nav-item'} onClick={() => onNavigate(id)}><Icon size={17} /><span>{label}</span></button>)}</nav>
    <div className="environment"><span className="status-dot" />Production telemetry<span>synthetic dataset</span></div>
  </aside>
}
