import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'

interface Props { onLogin: () => void }
export default function LoginPage({ onLogin }: Props) {
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    if (!username || !password) {
      setError('Username and password are required')
      return
    }
    setLoading(true)
    try {
      const data = await api.login(username, password)
      localStorage.setItem('opsagent_token', data.token)
      onLogin()
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.message || 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container" style={{maxWidth:480, marginTop:80}}>
      <h1 className="page-title">Sign In</h1>
      <form onSubmit={handleSubmit} className="card">
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Username</label>
          <input type="text" value={username} onChange={e => setUsername(e.target.value)} placeholder="admin" />
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Password</label>
          <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="••••••••" />
        </div>
        {error && <div style={{color:'#f87171', marginBottom:16}}>{error}</div>}
        <button type="submit" className="primary" disabled={loading} style={{width:'100%'}}>
          {loading ? 'Signing in...' : 'Sign In'}
        </button>
      </form>
    </div>
  )
}
