import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function SettingsPage() {
  const [aiEnabled, setAiEnabled] = useState(false)
  const [provider, setProvider] = useState('')
  const [apiKey, setApiKey] = useState('')
  const [model, setModel] = useState('')
  const [saving, setSaving] = useState(false)
  const [msg, setMsg] = useState('')

  useEffect(() => {
    api.getSettings().then(data => {
      const ai = data.ai || {}
      setAiEnabled(ai.enabled || false)
      setProvider(ai.provider || '')
      setApiKey(ai.api_key || '')
      setModel(ai.model || '')
    }).catch(console.error)
  }, [])

  async function handleSave() {
    setSaving(true)
    setMsg('')
    try {
      await api.updateAISettings({ enabled: aiEnabled, provider, api_key: apiKey, model })
      setMsg('Settings saved')
    } catch (err: any) {
      setMsg('Error: ' + err.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="container">
      <h1 className="page-title">Settings</h1>
      <div className="card" style={{maxWidth:600}}>
        <h2 className="section-title">AI Assistant</h2>
        <div style={{marginBottom:16}}>
          <label style={{display:'flex', alignItems:'center', gap:8, cursor:'pointer'}}>
            <input type="checkbox" checked={aiEnabled} onChange={e => setAiEnabled(e.target.checked)} style={{width:'auto'}} />
            Enable AI Assistant
          </label>
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Provider URL</label>
          <input type="text" value={provider} onChange={e => setProvider(e.target.value)} placeholder="https://api.openai.com/v1" />
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>API Key</label>
          <input type="password" value={apiKey} onChange={e => setApiKey(e.target.value)} placeholder="sk-..." />
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Model</label>
          <input type="text" value={model} onChange={e => setModel(e.target.value)} placeholder="gpt-4" />
        </div>
        {msg && <div style={{marginBottom:16, color: msg.startsWith('Error') ? '#f87171' : '#86efac'}}>{msg}</div>}
        <button className="primary" onClick={handleSave} disabled={saving}>{saving ? 'Saving...' : 'Save Settings'}</button>
      </div>
    </div>
  )
}
