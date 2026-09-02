import { GitCommitHorizontal, TrendingUp } from 'lucide-react'
import type { DeploymentCorrelation } from '../types'
import { formatPercent, formatRelative } from '../lib/format'

export function DeploymentEvidence({ deployments }: { deployments: DeploymentCorrelation[] }) {
  return <section className="panel" id="deployments"><div className="panel-heading"><div><h2>Deployment evidence</h2><p>Changes inside the active correlation window</p></div></div><div className="deployment-list">{deployments.map(({deployment,errorRate,errorTraces,totalTraces})=><article className="deployment" key={deployment.id}><div className="deployment-icon"><GitCommitHorizontal size={18}/></div><div><strong>{deployment.service} <span>{deployment.version}</span></strong><p>{deployment.changeSummary}</p><small>{deployment.commitSha} · {formatRelative(deployment.deployedAt)} · {deployment.deployedBy}</small></div><div className={errorRate>10?'correlation risk':'correlation'}><TrendingUp size={15}/><strong>{formatPercent(errorRate)}</strong><span>{errorTraces}/{totalTraces} error traces</span></div></article>)}</div></section>
}
