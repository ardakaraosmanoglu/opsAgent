import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { MetricCard } from '../components/MetricCard'
import { api } from '../lib/api'

export default function DashboardPage() {
  const navigate = useNavigate()
  const [summary, setSummary] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getDashboardSummary().then(setSummary).catch(console.error).finally(() => setLoading(false))
  }, [])

  const quickActions = [
    { label: 'Analyze disk', prompt: 'Analyze disk usage and suggest cleanup' },
    { label: 'Check CPU', prompt: 'Check CPU usage and report' },
    { label: 'Show ports', prompt: 'Show open network ports' },
    { label: 'System info', prompt: 'Show system information' },
  ]

  function handleQuickAction(prompt: string) {
    navigate('/assistant', { state: { prompt } })
  }

  if (loading) return <div className="container">Loading...</div>

  return (
    <div className="container">
      <h1 className="page-title">Dashboard</h1>

      <div className="grid" style={{marginBottom:32}}>
        <MetricCard label="CPU" value={summary?.cpu_usage ?? '-'} unit="%" color="#60a5fa" />
        <MetricCard label="Memory" value={summary?.memory_usage ?? '-'} unit="%" color="#a78bfa" />
        <MetricCard label="Disk" value={summary?.disk_usage ?? '-'} unit="%" color="#f97316" />
        <MetricCard label="Load" value={summary?.load ?? '-'} color="#22d3ee" />
        <MetricCard label="Uptime" value={summary?.uptime ?? '-'} color="#86efac" />
        <MetricCard label="Open Ports" value={summary?.open_ports ?? '-'} color="#fcd34d" />
        <MetricCard label="Active Alerts" value={summary?.open_alerts ?? 0} color="#f87171" />
      </div>

      <div className="section">
        <h2 className="section-title">Quick Actions</h2>
        <div style={{display:'flex', gap:12, flexWrap:'wrap'}}>
          {quickActions.map(a => (
            <button key={a.label} className="secondary" onClick={() => handleQuickAction(a.prompt)}>{a.label}</button>
          ))}
        </div>
      </div>

      <div className="section">
        <h2 className="section-title">Recent Tasks</h2>
        <div className="card" style={{color:'#94a3b8'}}>No recent tasks</div>
      </div>
    </div>
  )
}
