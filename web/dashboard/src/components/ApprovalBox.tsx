import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

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
    <Card className="border-primary/30 bg-primary/5 mt-2">
      <CardHeader className="pb-2 pt-3 px-4">
        <CardTitle className="text-xs font-medium">{title}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-3 space-y-2">
        <p className="text-xs text-muted-foreground">{summary}</p>
        <div className="space-y-1">
          {commands.map((c, i) => (
            <pre key={i} className={`rounded-md bg-black/40 border border-border/50 px-3 py-2 text-xs font-mono overflow-x-auto ${riskColors[c.risk] ?? 'text-muted-foreground'}`}>
              {c.cmd}
            </pre>
          ))}
        </div>
        <div className="flex gap-2 pt-1">
          <Button size="sm" className="h-7 text-xs" onClick={onApprove} disabled={loading}>
            Approve
          </Button>
          <Button size="sm" variant="destructive" className="h-7 text-xs" onClick={onReject} disabled={loading}>
            Reject
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
