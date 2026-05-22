import { useState, useEffect } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { AlertBadge } from '@/components/AlertBadge'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { api } from '@/lib/api'

type Filter = 'all' | 'open' | 'acknowledged' | 'resolved' | 'ignored'

export default function AlertsPage() {
  const [alerts, setAlerts] = useState<any[]>([])
  const [filter, setFilter] = useState<Filter>('all')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getAlerts()
      .then(data => setAlerts(data.alerts || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  const filtered = alerts.filter(a => filter === 'all' || a.status === filter)

  async function handleAction(id: number, action: 'acknowledge' | 'resolve' | 'ignore') {
    try {
      if (action === 'acknowledge') await api.acknowledgeAlert(id)
      else if (action === 'resolve') await api.resolveAlert(id)
      else await api.ignoreAlert(id)
      const data = await api.getAlerts()
      setAlerts(data.alerts || [])
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <div className="space-y-4 animate-fade-in">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Alerts</h1>
        <p className="text-xs text-muted-foreground mt-0.5">{alerts.length} total alerts</p>
      </div>

      <Tabs value={filter} onValueChange={(v) => setFilter(v as Filter)}>
        <TabsList className="h-8">
          {(['all', 'open', 'acknowledged', 'resolved', 'ignored'] as Filter[]).map(f => (
            <TabsTrigger key={f} value={f} className="text-xs capitalize h-7 px-2.5">
              {f}
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      {loading ? (
        <div className="py-12 text-center text-muted-foreground text-xs">Loading...</div>
      ) : filtered.length === 0 ? (
        <Card className="border-border/60">
          <CardContent className="py-12 text-center text-muted-foreground text-xs">
            No alerts
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-1.5">
          {filtered.map(alert => (
            <Card key={alert.id} className="border-border/60 p-3">
              <div className="flex items-center justify-between gap-3">
                <div className="flex items-center gap-2.5 min-w-0 flex-1">
                  <div className="shrink-0">
                    <AlertBadge severity={alert.severity} status={alert.status} />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs font-medium truncate">{alert.type}: {alert.message}</p>
                    <p className="text-xs text-muted-foreground mt-0.5">
                      {new Date(alert.created_at).toLocaleString()}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-1.5 shrink-0">
                  {alert.status === 'open' && (
                    <Button size="sm" variant="outline" className="h-7 text-xs" onClick={() => handleAction(alert.id, 'acknowledge')}>
                      Ack
                    </Button>
                  )}
                  {alert.status !== 'resolved' && (
                    <Button size="sm" className="h-7 text-xs" onClick={() => handleAction(alert.id, 'resolve')}>
                      Resolve
                    </Button>
                  )}
                  {alert.status !== 'ignored' && (
                    <Button size="sm" variant="destructive" className="h-7 text-xs" onClick={() => handleAction(alert.id, 'ignore')}>
                      Ignore
                    </Button>
                  )}
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
