import type { ServiceReliability } from '../types'
import { formatPercent } from '../lib/format'

export function ServiceHealth({ services }: { services: ServiceReliability[] }) {
  return <section className="panel" id="services"><div className="panel-heading"><div><h2>Service health</h2><p>Latency, errors, and composite health by service</p></div><span className="count-label">{services.length} services</span></div><div className="table-wrap"><table><thead><tr><th>Service</th><th>Health</th><th>Error rate</th><th>Average</th><th>p95 latency</th><th>Signals</th></tr></thead><tbody>{services.map(service => <tr key={service.service}><td><strong>{service.service}</strong><span className="cell-subtitle">production</span></td><td><span className={`health-score ${service.healthScore < 85 ? 'risk' : ''}`}>{service.healthScore}</span></td><td>{formatPercent(service.errorRate)}</td><td>{service.averageLatencyMs.toFixed(0)} ms</td><td>{service.p95LatencyMs.toFixed(0)} ms</td><td>{service.traceCount} spans</td></tr>)}</tbody></table></div></section>
}
