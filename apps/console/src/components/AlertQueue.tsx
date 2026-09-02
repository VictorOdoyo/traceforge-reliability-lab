import { Check, CircleAlert } from 'lucide-react'
import type { Alert } from '../types'
import { formatRelative } from '../lib/format'

export function AlertQueue({ alerts, onAcknowledge }: { alerts: Alert[]; onAcknowledge: (id:string)=>void }) {
  return <section className="panel" id="alerts"><div className="panel-heading"><div><h2>Alert queue</h2><p>Signals requiring ownership or action</p></div><span className="count-label">{alerts.filter(alert=>alert.state==='firing').length} unowned</span></div><div className="alert-list">{alerts.map(alert=><article className="alert-row" key={alert.id}><div className={`severity-marker ${alert.severity}`}><CircleAlert size={16}/></div><div className="alert-copy"><strong>{alert.title}</strong><span>{alert.service} · {formatRelative(alert.startedAt)}</span></div><span className={`badge ${alert.severity}`}>{alert.severity}</span>{alert.state==='firing'?<button className="secondary-button" onClick={()=>onAcknowledge(alert.id)}><Check size={14}/>Acknowledge</button>:<span className="owner">{alert.owner??alert.state}</span>}</article>)}</div></section>
}
