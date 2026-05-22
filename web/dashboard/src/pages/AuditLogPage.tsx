import { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Activity } from 'lucide-react'
import { api } from '@/lib/api'

export default function AuditLogPage() {
  const [tasks, setTasks] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getTasks()
      .then(data => setTasks(data.tasks || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="space-y-4 animate-fade-in">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Audit Log</h1>
        <p className="text-xs text-muted-foreground mt-0.5">{tasks.length} events recorded</p>
      </div>

      <Card className="border-border/60">
        <CardHeader className="pb-3">
          <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wide flex items-center gap-2">
            <Activity className="h-3.5 w-3.5" />
            Task History
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {loading ? (
            <div className="py-12 text-center text-muted-foreground text-xs">Loading...</div>
          ) : tasks.length === 0 ? (
            <div className="py-12 text-center text-muted-foreground text-xs">No audit events</div>
          ) : (
            <ScrollArea className="h-[calc(100vh-16rem)]">
              <div className="divide-y divide-border/50">
                {tasks.map(task => (
                  <div key={task.id} className="flex items-center gap-3 px-4 py-2.5 hover:bg-muted/30 transition-colors">
                    <div className="text-xs text-muted-foreground w-36 shrink-0 font-mono">
                      {new Date(task.created_at || Date.now()).toLocaleString()}
                    </div>
                    <div className="w-20 shrink-0">
                      <Badge variant="outline" className="text-xs capitalize font-normal">
                        {task.status}
                      </Badge>
                    </div>
                    <div className="w-20 shrink-0 text-xs text-muted-foreground">
                      {task.actor || 'system'}
                    </div>
                    <div className="text-xs truncate text-foreground/80">
                      {task.description || task.type || 'Task'}
                    </div>
                  </div>
                ))}
              </div>
            </ScrollArea>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
