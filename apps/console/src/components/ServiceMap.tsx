import { Background, Controls, Handle, Position, ReactFlow, type Edge, type Node, type NodeProps } from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import type { DependencyEdge, ServiceReliability } from '../types'

type ServiceNodeData={label:string;health:number;errors:number}
function ServiceNode({data}:NodeProps<Node<ServiceNodeData>>){return <div className={`service-node ${data.health<85?'degraded':''}`}><Handle type="target" position={Position.Left}/><span>{data.label}</span><strong>{data.health}</strong><small>{data.errors} errors</small><Handle type="source" position={Position.Right}/></div>}
const nodeTypes={service:ServiceNode}

export function ServiceMap({ services, dependencies }: { services: ServiceReliability[]; dependencies: DependencyEdge[] }) {
  const positions=[[20,100],[230,30],[230,180],[450,30],[450,180]]
  const nodes:Node<ServiceNodeData>[] = services.map((service,index)=>({id:service.service,type:'service',position:{x:positions[index]?.[0]??index*170,y:positions[index]?.[1]??100},data:{label:service.service,health:service.healthScore,errors:service.errorCount}}))
  const edges:Edge[]=dependencies.map((edge,index)=>({id:`edge-${index}`,source:edge.source,target:edge.target,animated:edge.errors>0,label:`${edge.calls} calls`,style:{stroke:edge.errors>0?'#c63c3c':'#72a7a3',strokeWidth:edge.errors>0?2.5:1.5},labelStyle:{fontSize:10,fill:'#52616c'}}))
  return <section className="panel topology-panel" id="topology"><div className="panel-heading"><div><h2>Service topology</h2><p>Live dependency paths and error propagation</p></div><span className="live-label"><i/>live map</span></div><div className="topology-frame"><ReactFlow nodes={nodes} edges={edges} nodeTypes={nodeTypes} fitView minZoom={0.6} maxZoom={1.5} nodesDraggable={false} proOptions={{hideAttribution:true}}><Background color="#dbe3e8" gap={18}/><Controls showInteractive={false}/></ReactFlow></div></section>
}
