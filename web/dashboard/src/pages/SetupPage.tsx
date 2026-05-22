import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'

interface Props {
  onSetupDone: () => void
}

export default function SetupPage({ onSetupDone }: Props) {
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    if (!username || !password) {
      setError('Username and password are required')
      return
    }
    if (password !== confirm) {
      setError('Passwords do not match')
      return
    }
    if (password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }
    setLoading(true)
    try {
      await api.createAdmin(username, password)
      onSetupDone()
      navigate('/login')
    } catch (err: any) {
      setError(err.message || 'Setup failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container" style={{maxWidth:480, marginTop:80}}>
      <h1 className="page-title">Initial Setup</h1>
      <p style={{color:'#94a3b8', marginBottom:24}}>Create your admin account to get started.</p>
      <form onSubmit={handleSubmit} className="card">
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Username</label>
          <input type="text" value={username} onChange={e => setUsername(e.target.value)} placeholder="admin" />
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Password</label>
          <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="••••••••" />
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Confirm Password</label>
          <input type="password" value={confirm} onChange={e => setConfirm(e.target.value)} placeholder="••••••••" />
        </div>
        {error && <div style={{color:'#f87171', marginBottom:16}}>{error}</div>}
        <button type="submit" className="primary" disabled={loading} style={{width:'100%'}}>
          {loading ? 'Creating...' : 'Create Admin Account'}
        </button>
      </form>
    </div>
  )
}
