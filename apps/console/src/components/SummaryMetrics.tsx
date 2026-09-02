import { AlertTriangle, Boxes, Flame, Siren } from 'lucide-react'
import type { Overview } from '../types'

export function SummaryMetrics({ overview }: { overview: Overview }) {
  const items = [
    { label:'Observed services',value:overview.serviceCount,detail:'5 reporting in the window',icon:Boxes,tone:'teal' },
    { label:'Open alerts',value:overview.openAlerts,detail:'2 require ownership',icon:AlertTriangle,tone:'amber' },
    { label:'Active incidents',value:overview.activeIncidents,detail:'1 critical response',icon:Siren,tone:'red' },
    { label:'Budget risks',value:overview.errorBudgetRisks,detail:'Burn rate above 1x',icon:Flame,tone:'violet' },
  ]
  return <section className="metric-grid" aria-label="Reliability summary">{items.map(({label,value,detail,icon:Icon,tone}) => <article className="metric" key={label}><div className={`metric-icon ${tone}`}><Icon size={18}/></div><div><span>{label}</span><strong>{value}</strong><small>{detail}</small></div></article>)}</section>
}
