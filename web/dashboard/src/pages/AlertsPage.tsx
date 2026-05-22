import { useState, useEffect } from 'react'
import { AlertBadge } from '../components/AlertBadge'
import { api } from '../lib/api'

type Filter = 'all' | 'open' | 'acknowledged' | 'resolved' | 'ignored'

export default function AlertsPage() {
  const [alerts, setAlerts] = useState<any[]>([])
  const [filter, setFilter] = useState<Filter>('all')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getAlerts().then(data => setAlerts(data.alerts || [])).catch(console.error).finally(() => setLoading(false))
  }, [])

  const filtered = alerts.filter(a => filter === 'all' || a.status === filter)

  async function handleAction(id: number, action: 'acknowledge' | 'resolve' | 'ignore') {
    try {
      if (action === 'acknowledge') await api.acknowledgeAlert(id)
      else if (action === 'resolve') await api.resolveAlert(id)
      else await api.ignoreAlert(id)
      const data = await api.getAlerts()
      setAlerts(data.alerts || [])
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <div className="container">
      <h1 className="page-title">Alerts</h1>
      <div style={{display:'flex', gap:8, marginBottom:24}}>
        {(['all','open','acknowledged','resolved','ignored'] as Filter[]).map(f => (
          <button key={f} className={filter === f ? 'primary' : 'secondary'} onClick={() => setFilter(f)}>{f}</button>
        ))}
      </div>
      {loading ? <div>Loading...</div> : filtered.length === 0 ? <div className="card">No alerts</div> : (
        <div style={{display:'flex', flexDirection:'column', gap:8}}>
          {filtered.map(alert => (
            <div key={alert.id} className="card" style={{display:'flex', justifyContent:'space-between', alignItems:'center'}}>
              <div>
                <AlertBadge severity={alert.severity} status={alert.status} />
                <span style={{marginLeft:12}}>{alert.type}: {alert.message}</span>
                <div style={{fontSize:12, color:'#64748b', marginTop:4}}>{new Date(alert.created_at).toLocaleString()}</div>
              </div>
              <div style={{display:'flex', gap:8}}>
                {alert.status === 'open' && <button className="secondary" onClick={() => handleAction(alert.id, 'acknowledge')}>Ack</button>}
                {alert.status !== 'resolved' && <button className="primary" onClick={() => handleAction(alert.id, 'resolve')}>Resolve</button>}
                {alert.status !== 'ignored' && <button className="danger" onClick={() => handleAction(alert.id, 'ignore')}>Ignore</button>}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
