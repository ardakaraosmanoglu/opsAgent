import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function AuditLogPage() {
  const [tasks, setTasks] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getTasks().then(data => setTasks(data.tasks || [])).catch(console.error).finally(() => setLoading(false))
  }, [])

  return (
    <div className="container">
      <h1 className="page-title">Audit Log</h1>
      {loading ? <div>Loading...</div> : tasks.length === 0 ? <div className="card">No audit events</div> : (
        <table style={{width:'100%', borderCollapse:'collapse'}}>
          <thead>
            <tr style={{textAlign:'left', color:'#94a3b8', fontSize:12}}>
              <th style={{padding:'8px 12px'}}>Timestamp</th>
              <th style={{padding:'8px 12px'}}>Event</th>
              <th style={{padding:'8px 12px'}}>Actor</th>
              <th style={{padding:'8px 12px'}}>Message</th>
            </tr>
          </thead>
          <tbody>
            {tasks.map(task => (
              <tr key={task.id} style={{borderTop:'1px solid #334155'}}>
                <td style={{padding:'12px'}}>{new Date(task.created_at || Date.now()).toLocaleString()}</td>
                <td style={{padding:'12px'}}><span className="badge info">{task.status}</span></td>
                <td style={{padding:'12px'}}>{task.actor || 'system'}</td>
                <td style={{padding:'12px'}}>{task.description || task.type || 'Task'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
