import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { MetricCard } from '@/components/MetricCard'
import { HardDrive, Cpu, Network, Info } from 'lucide-react'
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
    { label: 'Analyze disk', prompt: 'Analyze disk usage and suggest cleanup', icon: HardDrive },
    { label: 'Check CPU', prompt: 'Check CPU usage and report', icon: Cpu },
    { label: 'Show ports', prompt: 'Show open network ports', icon: Network },
    { label: 'System info', prompt: 'Show system information', icon: Info },
  ]

  function handleQuickAction(prompt: string) {
    navigate('/assistant', { state: { prompt } })
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-sm text-muted-foreground">Loading...</p>
      </div>
    )
  }

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-xs text-muted-foreground mt-0.5">System overview</p>
      </div>

      {/* Metrics Grid */}
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-7 gap-2">
        <MetricCard label="CPU" value={summary?.cpu_usage ?? '-'} unit="%" color="#60a5fa" />
        <MetricCard label="Memory" value={summary?.memory_usage ?? '-'} unit="%" color="#a78bfa" />
        <MetricCard label="Disk" value={summary?.disk_usage ?? '-'} unit="%" color="#f97316" />
        <MetricCard label="Load" value={summary?.load ?? '-'} color="#22d3ee" />
        <MetricCard label="Uptime" value={summary?.uptime ?? '-'} color="#86efac" />
        <MetricCard label="Ports" value={summary?.open_ports ?? '-'} color="#fcd34d" />
        <MetricCard label="Alerts" value={summary?.open_alerts ?? 0} color="#f87171" />
      </div>

      {/* Quick Actions */}
      <Card className="border-border/60">
        <CardHeader className="pb-2">
          <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Quick Actions</CardTitle>
        </CardHeader>
        <CardContent className="pt-0">
          <div className="flex flex-wrap gap-2">
            {quickActions.map(({ label, prompt, icon: Icon }) => (
              <Button
                key={label}
                variant="outline"
                size="sm"
                className="h-7 text-xs gap-1.5"
                onClick={() => handleQuickAction(prompt)}
              >
                <Icon className="h-3 w-3" />
                {label}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Recent Tasks */}
      <Card className="border-border/60">
        <CardHeader className="pb-2">
          <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Recent Tasks</CardTitle>
        </CardHeader>
        <CardContent className="pt-0">
          <p className="text-xs text-muted-foreground">No recent tasks</p>
        </CardContent>
      </Card>
    </div>
  )
}
