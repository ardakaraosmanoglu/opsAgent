import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { MetricCard } from '@/components/MetricCard'
import { api } from '@/lib/api'

export default function DashboardPage() {
  const navigate = useNavigate()
  const [summary, setSummary] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getDashboardSummary()
      .then(setSummary)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  const quickActions = [
    { label: 'Analyze disk', prompt: 'Analyze disk usage and suggest cleanup', icon: '💾' },
    { label: 'Check CPU', prompt: 'Check CPU usage and report', icon: '⚡' },
    { label: 'Show ports', prompt: 'Show open network ports', icon: '🔌' },
    { label: 'System info', prompt: 'Show system information', icon: 'ℹ️' },
  ]

  function handleQuickAction(prompt: string) {
    navigate('/assistant', { state: { prompt } })
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-sm text-muted-foreground mt-1">System overview and quick actions</p>
      </div>

      {/* Metrics Grid */}
      <div className="grid gap-3 grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6">
        <MetricCard label="CPU" value={summary?.cpu_usage ?? '-'} unit="%" color="#60a5fa" />
        <MetricCard label="Memory" value={summary?.memory_usage ?? '-'} unit="%" color="#a78bfa" />
        <MetricCard label="Disk" value={summary?.disk_usage ?? '-'} unit="%" color="#f97316" />
        <MetricCard label="Load" value={summary?.load ?? '-'} color="#22d3ee" />
        <MetricCard label="Uptime" value={summary?.uptime ?? '-'} color="#86efac" />
        <MetricCard label="Open Ports" value={summary?.open_ports ?? '-'} color="#fcd34d" />
        <MetricCard label="Active Alerts" value={summary?.open_alerts ?? 0} color="#f87171" />
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Quick Actions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {quickActions.map(a => (
              <Button
                key={a.label}
                variant="outline"
                size="sm"
                onClick={() => handleQuickAction(a.prompt)}
              >
                <span className="mr-1">{a.icon}</span>
                {a.label}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Recent Tasks */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Recent Tasks</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No recent tasks</p>
        </CardContent>
      </Card>
    </div>
  )
}
