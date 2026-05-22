import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
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
      <div className="header" style={{marginBottom:24, marginTop:-24, marginLeft:-24, marginRight:-24}}>
        <div style={{fontSize:18, fontWeight:600}}>OpsAgent</div>
        <nav>
          <Link to="/dashboard" style={{color:'#60a5fa'}}>Dashboard</Link>
          <Link to="/alerts" style={{marginLeft:24}}>Alerts</Link>
          <Link to="/assistant" style={{marginLeft:24}}>Assistant</Link>
          <Link to="/audit" style={{marginLeft:24}}>Audit</Link>
          <Link to="/settings" style={{marginLeft:24}}>Settings</Link>
        </nav>
      </div>

      <h1 className="page-title">Dashboard</h1>

      <div className="grid" style={{marginBottom:32}}>
        <MetricCard label="CPU" value={summary?.cpu ?? '-'} unit="%" color="#60a5fa" />
        <MetricCard label="Memory" value={summary?.memory ?? '-'} unit="%" color="#a78bfa" />
        <MetricCard label="Disk" value={summary?.disk ?? '-'} unit="%" color="#f97316" />
        <MetricCard label="Load" value={summary?.load ?? '-'} color="#22d3ee" />
        <MetricCard label="Uptime" value={summary?.uptime ?? '-'} color="#86efac" />
        <MetricCard label="Open Ports" value={summary?.open_ports ?? '-'} color="#fcd34d" />
        <MetricCard label="Active Alerts" value={summary?.active_alerts ?? 0} color="#f87171" />
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
