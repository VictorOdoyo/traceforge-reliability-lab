import type { SLO } from '../types'
import { formatPercent } from '../lib/format'

const attainment = (slo:SLO) => slo.totalEvents ? slo.goodEvents / slo.totalEvents * 100 : 100
const budget = (slo:SLO) => { const allowed=slo.totalEvents*(1-slo.target/100); const remaining=(allowed-(slo.totalEvents-slo.goodEvents))/allowed*100; return Math.max(0,remaining) }
const burn = (slo:SLO) => ((slo.totalEvents-slo.goodEvents)/slo.totalEvents)/(1-slo.target/100)

export function ErrorBudgets({ slos }: { slos: SLO[] }) {
  return <section className="panel" id="slos"><div className="panel-heading"><div><h2>Error budgets</h2><p>30-day objective health and consumption</p></div></div><div className="budget-list">{slos.map(slo => { const remaining=budget(slo); const rate=burn(slo); return <article className="budget-row" key={slo.id}><div className="budget-title"><strong>{slo.name}</strong><span>{slo.service} · {formatPercent(slo.target,2)} target</span></div><div className="budget-measure"><div><span>Attainment</span><strong>{formatPercent(attainment(slo),3)}</strong></div><div><span>Burn rate</span><strong className={rate>1?'danger':''}>{rate.toFixed(2)}x</strong></div><div className="budget-bar-wrap"><span>{formatPercent(remaining,0)} remaining</span><div className="budget-bar"><i className={remaining<25?'danger':''} style={{width:`${remaining}%`}}/></div></div></div></article> })}</div></section>
}
