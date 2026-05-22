import { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Settings2, Key, Server, Download } from 'lucide-react'
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
      setServiceStatus({ output: 'Failed: ' + err.message })
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
    <div className="space-y-6 max-w-2xl animate-fade-in">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Settings</h1>
        <p className="text-xs text-muted-foreground mt-0.5">Configure your OpsAgent installation</p>
      </div>

      {/* AI Assistant */}
      <Card className="border-border/60">
        <CardHeader className="pb-3">
          <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wide flex items-center gap-2">
            <Settings2 className="h-3.5 w-3.5" />
            AI Assistant
          </CardTitle>
          <CardDescription className="text-xs">Configure AI integration</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <Label htmlFor="ai-enabled" className="text-xs cursor-pointer">Enable AI Assistant</Label>
            <Switch id="ai-enabled" checked={aiEnabled} onCheckedChange={setAiEnabled} />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="provider" className="text-xs">Provider URL</Label>
            <Input id="provider" type="text" value={provider} onChange={e => setProvider(e.target.value)} placeholder="https://api.openai.com/v1" className="h-8 text-xs" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="apiKey" className="text-xs">API Key</Label>
            <Input id="apiKey" type="password" value={apiKey} onChange={e => setApiKey(e.target.value)} placeholder="sk-..." className="h-8 text-xs" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="model" className="text-xs">Model</Label>
            <Input id="model" type="text" value={model} onChange={e => setModel(e.target.value)} placeholder="gpt-4" className="h-8 text-xs" />
          </div>
          {msg && <Alert variant={msg.startsWith('Error') ? 'destructive' : 'default'} className="py-2"><AlertDescription className="text-xs">{msg}</AlertDescription></Alert>}
          <Button size="sm" onClick={handleSave} disabled={saving} className="h-8 text-xs">
            {saving ? 'Saving...' : 'Save Settings'}
          </Button>
        </CardContent>
      </Card>

      <Separator />

      {/* Change Password */}
      <Card className="border-border/60">
        <CardHeader className="pb-3">
          <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wide flex items-center gap-2">
            <Key className="h-3.5 w-3.5" />
            Change Password
          </CardTitle>
          <CardDescription className="text-xs">Update your admin password</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-1.5">
            <Label htmlFor="currentPassword" className="text-xs">Current Password</Label>
            <Input id="currentPassword" type="password" value={currentPassword} onChange={e => setCurrentPassword(e.target.value)} className="h-8 text-xs" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="newPassword" className="text-xs">New Password</Label>
            <Input id="newPassword" type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} className="h-8 text-xs" />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="confirmPassword" className="text-xs">Confirm New Password</Label>
            <Input id="confirmPassword" type="password" value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} className="h-8 text-xs" />
          </div>
          {passwordMsg && <Alert variant={passwordMsg.startsWith('Error') ? 'destructive' : 'default'} className="py-2"><AlertDescription className="text-xs">{passwordMsg}</AlertDescription></Alert>}
          <Button size="sm" onClick={handleChangePassword} disabled={changingPassword} className="h-8 text-xs">
            {changingPassword ? 'Changing...' : 'Change Password'}
          </Button>
        </CardContent>
      </Card>

      <Separator />

      {/* Service Status */}
      <Card className="border-border/60">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wide flex items-center gap-2">
                <Server className="h-3.5 w-3.5" />
                Service Status
              </CardTitle>
              <CardDescription className="text-xs mt-0.5">systemd service health</CardDescription>
            </div>
            <Button size="sm" variant="outline" onClick={refreshServiceStatus} disabled={serviceLoading} className="h-7 text-xs">
              {serviceLoading ? 'Loading...' : 'Refresh'}
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <ScrollArea className="h-36 rounded-md bg-black/40 border">
            <pre className="text-xs p-3 text-muted-foreground whitespace-pre-wrap font-mono leading-relaxed">
              {serviceStatus?.output || 'No output'}
            </pre>
          </ScrollArea>
        </CardContent>
      </Card>

      <Separator />

      {/* Agent Update */}
      <Card className="border-border/60">
        <CardHeader className="pb-3">
          <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wide flex items-center gap-2">
            <Download className="h-3.5 w-3.5" />
            Agent Update
          </CardTitle>
          <CardDescription className="text-xs">Check for and apply updates</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex gap-6 text-xs">
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
          {updateMsg && <Alert variant={updateMsg.includes('failed') || updateMsg.includes('Error') ? 'destructive' : 'default'} className="py-2"><AlertDescription className="text-xs">{updateMsg}</AlertDescription></Alert>}
          <div className="flex gap-2">
            <Button size="sm" variant="outline" onClick={handleCheckUpdate} disabled={updateChecking} className="h-8 text-xs">
              {updateChecking ? 'Checking...' : 'Check for Updates'}
            </Button>
            {updateInfo?.needs_update && (
              <Button size="sm" onClick={handleUpdateAgent} disabled={updating} className="h-8 text-xs">
                {updating ? 'Updating...' : 'Update Now'}
              </Button>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
