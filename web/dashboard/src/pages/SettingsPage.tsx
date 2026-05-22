import { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { ScrollArea } from '@/components/ui/scroll-area'
import { api } from '@/lib/api'

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

  const [serviceStatus, setServiceStatus] = useState<any>(null)
  const [serviceLoading, setServiceLoading] = useState(true)

  const [updateInfo, setUpdateInfo] = useState<any>(null)
  const [updateChecking, setUpdateChecking] = useState(false)
  const [updating, setUpdating] = useState(false)
  const [updateMsg, setUpdateMsg] = useState('')

  useEffect(() => {
    api.getSettings().then(data => {
      const ai = data.ai || {}
      setAiEnabled(ai.enabled || false)
      setProvider(ai.provider || '')
      setApiKey(ai.api_key || '')
      setModel(ai.model || '')
    }).catch(console.error)

    api.getServiceStatus().then(data => {
      setServiceStatus(data)
    }).catch(err => {
      setServiceStatus({ output: 'Failed to load service status: ' + err.message })
    }).finally(() => setServiceLoading(false))

    api.checkForUpdate().then(data => {
      setUpdateInfo(data)
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

  async function handleCheckUpdate() {
    setUpdateChecking(true)
    setUpdateMsg('')
    try {
      const data = await api.checkForUpdate()
      setUpdateInfo(data)
      setUpdateMsg(data.needs_update ? 'Update available: ' + data.latest_version : 'You are on the latest version.')
    } catch (err: any) {
      setUpdateMsg('Update check failed: ' + err.message)
    } finally {
      setUpdateChecking(false)
    }
  }

  async function handleUpdateAgent() {
    setUpdating(true)
    setUpdateMsg('')
    try {
      await api.updateAgent()
      setUpdateMsg('Agent updated successfully. Service restarted.')
    } catch (err: any) {
      setUpdateMsg('Update failed: ' + err.message)
    } finally {
      setUpdating(false)
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

  function refreshServiceStatus() {
    setServiceLoading(true)
    api.getServiceStatus().then(d => { setServiceStatus(d); setServiceLoading(false) }).catch(e => { setServiceStatus({ output: 'Error: ' + e.message }); setServiceLoading(false) })
  }

  return (
    <div className="space-y-6 max-w-2xl">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
        <p className="text-sm text-muted-foreground mt-1">Configure your OpsAgent installation</p>
      </div>

      {/* AI Assistant */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">AI Assistant</CardTitle>
          <CardDescription>Configure the AI assistant integration</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <Label htmlFor="ai-enabled" className="cursor-pointer">Enable AI Assistant</Label>
            <Switch id="ai-enabled" checked={aiEnabled} onCheckedChange={setAiEnabled} />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="provider">Provider URL</Label>
            <Input id="provider" type="text" value={provider} onChange={e => setProvider(e.target.value)} placeholder="https://api.openai.com/v1" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="apiKey">API Key</Label>
            <Input id="apiKey" type="password" value={apiKey} onChange={e => setApiKey(e.target.value)} placeholder="sk-..." />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="model">Model</Label>
            <Input id="model" type="text" value={model} onChange={e => setModel(e.target.value)} placeholder="gpt-4" />
          </div>
          {msg && <Alert variant={msg.startsWith('Error') ? 'destructive' : 'default'}><AlertDescription>{msg}</AlertDescription></Alert>}
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Saving...' : 'Save Settings'}</Button>
        </CardContent>
      </Card>

      <Separator />

      {/* Change Password */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Change Password</CardTitle>
          <CardDescription>Update your admin account password</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="currentPassword">Current Password</Label>
            <Input id="currentPassword" type="password" value={currentPassword} onChange={e => setCurrentPassword(e.target.value)} placeholder="Enter current password" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="newPassword">New Password</Label>
            <Input id="newPassword" type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} placeholder="Enter new password" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="confirmPassword">Confirm New Password</Label>
            <Input id="confirmPassword" type="password" value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} placeholder="Confirm new password" />
          </div>
          {passwordMsg && <Alert variant={passwordMsg.startsWith('Error') ? 'destructive' : 'default'}><AlertDescription>{passwordMsg}</AlertDescription></Alert>}
          <Button onClick={handleChangePassword} disabled={changingPassword}>{changingPassword ? 'Changing...' : 'Change Password'}</Button>
        </CardContent>
      </Card>

      <Separator />

      {/* Service Status */}
      <Card>
        <CardHeader className="flex-row items-center justify-between space-y-0">
          <div>
            <CardTitle className="text-base">Service Status</CardTitle>
            <CardDescription>systemd service health</CardDescription>
          </div>
          <Button size="sm" variant="outline" onClick={refreshServiceStatus} disabled={serviceLoading}>
            {serviceLoading ? 'Loading...' : 'Refresh'}
          </Button>
        </CardHeader>
        <CardContent>
          <ScrollArea className="h-48 rounded-md bg-black/40 border">
            <pre className="text-xs p-3 text-muted-foreground whitespace-pre-wrap font-mono">
              {serviceStatus?.output || 'No output'}
            </pre>
          </ScrollArea>
        </CardContent>
      </Card>

      <Separator />

      {/* Agent Update */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Agent Update</CardTitle>
          <CardDescription>Check for and apply agent updates</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-6 text-sm">
            <div>
              <span className="text-muted-foreground">Current: </span>
              <span className="font-medium">0.1.0</span>
            </div>
            <div>
              <span className="text-muted-foreground">Latest: </span>
              <span className="font-medium">{updateInfo?.latest_version || 'Unknown'}</span>
            </div>
            {updateInfo?.needs_update && (
              <div className="text-amber-400 font-medium">Update available!</div>
            )}
          </div>
          {updateMsg && <Alert variant={updateMsg.includes('failed') || updateMsg.includes('Error') ? 'destructive' : 'default'}><AlertDescription>{updateMsg}</AlertDescription></Alert>}
          <div className="flex gap-2">
            <Button size="sm" variant="outline" onClick={handleCheckUpdate} disabled={updateChecking}>
              {updateChecking ? 'Checking...' : 'Check for Updates'}
            </Button>
            {updateInfo?.needs_update && (
              <Button size="sm" onClick={handleUpdateAgent} disabled={updating}>
                {updating ? 'Updating...' : 'Update Now'}
              </Button>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
