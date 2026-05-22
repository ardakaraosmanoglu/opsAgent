import { useState, useEffect } from 'react'
import { api } from '../lib/api'

export default function SettingsPage() {
  const [aiEnabled, setAiEnabled] = useState(false)
  const [provider, setProvider] = useState('')
  const [apiKey, setApiKey] = useState('')
  const [model, setModel] = useState('')
  const [saving, setSaving] = useState(false)
  const [msg, setMsg] = useState('')

  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [passwordMsg, setPasswordMsg] = useState('')
  const [changingPassword, setChangingPassword] = useState(false)

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

  async function handleChangePassword() {
    setChangingPassword(true)
    setPasswordMsg('')

    if (newPassword !== confirmPassword) {
      setPasswordMsg('Error: passwords do not match')
      setChangingPassword(false)
      return
    }

    if (newPassword.length < 8) {
      setPasswordMsg('Error: new password must be at least 8 characters')
      setChangingPassword(false)
      return
    }

    try {
      await api.updatePassword(currentPassword, newPassword)
      setPasswordMsg('Password changed successfully')
      setCurrentPassword('')
      setNewPassword('')
      setConfirmPassword('')
    } catch (err: any) {
      setPasswordMsg('Error: ' + err.message)
    } finally {
      setChangingPassword(false)
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

      <div className="card" style={{maxWidth:600, marginTop:24}}>
        <h2 className="section-title">Change Password</h2>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Current Password</label>
          <input type="password" value={currentPassword} onChange={e => setCurrentPassword(e.target.value)} placeholder="Enter current password" />
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>New Password</label>
          <input type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} placeholder="Enter new password" />
        </div>
        <div style={{marginBottom:16}}>
          <label style={{display:'block', marginBottom:4, fontSize:14}}>Confirm New Password</label>
          <input type="password" value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} placeholder="Confirm new password" />
        </div>
        {passwordMsg && <div style={{marginBottom:16, color: passwordMsg.startsWith('Error') ? '#f87171' : '#86efac'}}>{passwordMsg}</div>}
        <button className="primary" onClick={handleChangePassword} disabled={changingPassword}>{changingPassword ? 'Changing...' : 'Change Password'}</button>
      </div>
    </div>
  )
}
