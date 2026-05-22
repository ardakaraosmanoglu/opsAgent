import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

interface Command { cmd: string; risk: string }
interface Props {
  title: string
  summary: string
  commands: Command[]
  onApprove: () => void
  onReject: () => void
  loading?: boolean
}

const riskColors: Record<string, string> = {
  low: 'text-emerald-400',
  medium: 'text-amber-400',
  high: 'text-orange-400',
  critical: 'text-red-400',
}

export function ApprovalBox({ title, summary, commands, onApprove, onReject, loading }: Props) {
  return (
    <Card className="border-primary/50 mt-3">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm">{title}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm text-muted-foreground">{summary}</p>
        <div className="space-y-1.5">
          {commands.map((c, i) => (
            <pre key={i} className={cn(
              'rounded-md bg-black/40 border border-border px-3 py-2 text-xs font-mono overflow-x-auto',
              riskColors[c.risk] ?? 'text-muted-foreground'
            )}>
              {c.cmd}
            </pre>
          ))}
        </div>
        <div className="flex gap-2 pt-1">
          <Button size="sm" onClick={onApprove} disabled={loading}>Approve</Button>
          <Button size="sm" variant="destructive" onClick={onReject} disabled={loading}>Reject</Button>
        </div>
      </CardContent>
    </Card>
  )
}
