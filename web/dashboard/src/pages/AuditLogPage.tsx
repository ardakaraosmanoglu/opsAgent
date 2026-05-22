import { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
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
    <div className="space-y-4">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Audit Log</h1>
        <p className="text-sm text-muted-foreground mt-1">{tasks.length} events recorded</p>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Task History</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {loading ? (
            <div className="py-12 text-center text-muted-foreground">Loading...</div>
          ) : tasks.length === 0 ? (
            <div className="py-12 text-center text-muted-foreground text-sm">No audit events</div>
          ) : (
            <ScrollArea className="h-[60vh]">
              <div className="divide-y divide-border">
                {tasks.map(task => (
                  <div key={task.id} className="flex items-center gap-4 px-4 py-3 hover:bg-muted/50 transition-colors">
                    <div className="text-xs text-muted-foreground w-36 shrink-0">
                      {new Date(task.created_at || Date.now()).toLocaleString()}
                    </div>
                    <div className="w-28 shrink-0">
                      <Badge variant="outline" className="text-xs capitalize">{task.status}</Badge>
                    </div>
                    <div className="w-24 shrink-0 text-xs text-muted-foreground">
                      {task.actor || 'system'}
                    </div>
                    <div className="text-sm truncate">
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
