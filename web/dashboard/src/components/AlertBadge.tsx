import { Badge } from '@/components/ui/badge'

interface Props { severity: string; status?: string }

export function AlertBadge({ severity, status }: Props) {
  let variant: 'default' | 'secondary' | 'destructive' | 'outline' | 'success' = 'secondary'
  let text = severity

  if (status) {
    text = status
    if (status === 'open') variant = 'default'
    else if (status === 'acknowledged') variant = 'secondary'
    else if (status === 'resolved') variant = 'success'
    else if (status === 'ignored') variant = 'outline'
  } else {
    if (severity === 'critical' || severity === 'error') variant = 'destructive'
    else if (severity === 'warning') variant = 'secondary'
    else if (severity === 'info') variant = 'default'
    else if (severity === 'success') variant = 'success'
  }

  return <Badge variant={variant} className="text-xs capitalize">{text}</Badge>
}
