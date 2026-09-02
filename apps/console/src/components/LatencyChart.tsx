import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import type { Trace } from '../types'
import { durationMs, formatTime } from '../lib/format'

export function LatencyChart({ traces }: { traces: Trace[] }) {
  const data = [...traces].reverse().map(trace => ({ time: formatTime(trace.startedAt), latency: Math.round(durationMs(trace.duration)), error: trace.hasError ? 1 : 0 }))
  return <section className="panel chart-panel"><div className="panel-heading"><div><h2>End-to-end latency</h2><p>Trace duration across the selected window</p></div><span className="legend"><i/>p95 threshold 650 ms</span></div><div className="chart-frame"><ResponsiveContainer width="100%" height="100%"><LineChart data={data} margin={{left:-18,right:14,top:10,bottom:0}}><CartesianGrid stroke="#e7edf1" vertical={false}/><XAxis dataKey="time" tick={{fontSize:11,fill:'#65717e'}} tickLine={false} axisLine={false}/><YAxis tick={{fontSize:11,fill:'#65717e'}} tickLine={false} axisLine={false} unit="ms"/><Tooltip contentStyle={{borderRadius:6,border:'1px solid #dbe3e8',fontSize:12}}/><Line type="monotone" dataKey="latency" stroke="#087f78" strokeWidth={2.5} dot={(props) => { const { cx, cy, payload } = props as { cx:number;cy:number;payload:{error:number} }; return <circle cx={cx} cy={cy} r={payload.error ? 4 : 2.5} fill={payload.error ? '#c63c3c' : '#087f78'} /> }} /></LineChart></ResponsiveContainer></div></section>
}
