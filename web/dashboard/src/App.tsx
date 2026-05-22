import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useState, useEffect } from 'react'
import SetupPage from './pages/SetupPage'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import AlertsPage from './pages/AlertsPage'
import AssistantPage from './pages/AssistantPage'
import AuditLogPage from './pages/AuditLogPage'
import SettingsPage from './pages/SettingsPage'

function App() {
  const [loading, setLoading] = useState(true)
  const [setupRequired, setSetupRequired] = useState(true)
  const [authenticated, setAuthenticated] = useState(false)

  useEffect(() => {
    checkSetup()
  }, [])

  async function checkSetup() {
    try {
      const res = await fetch('/api/setup/required')
      const data = await res.json()
      setSetupRequired(data.required)
    } catch {
      setSetupRequired(true)
    } finally {
      setLoading(false)
    }
  }

  // Re-check setup status (called after setup completes)
  function refreshAuth() {
    checkSetup()
    setAuthenticated(false)
  }

  if (loading) return <div style={{padding:48, textAlign:'center'}}>Loading...</div>

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/setup" element={
          setupRequired ? <SetupPage onSetupDone={refreshAuth} /> : <Navigate to="/login" />
        } />
        <Route path="/login" element={authenticated ? <Navigate to="/" /> : <LoginPage onLogin={() => setAuthenticated(true)} />} />
        <Route path="/" element={
          !authenticated ? <Navigate to="/login" /> :
          setupRequired ? <Navigate to="/setup" /> :
          <DashboardPage />
        }>
          <Route index element={<Navigate to="/dashboard" />} />
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="alerts" element={<AlertsPage />} />
          <Route path="assistant" element={<AssistantPage />} />
          <Route path="audit" element={<AuditLogPage />} />
          <Route path="settings" element={<SettingsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
