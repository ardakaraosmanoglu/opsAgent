import { Card, CardContent } from '@/components/ui/card'

interface Props {
  label: string
  value: string | number
  unit?: string
  color?: string
}

export function MetricCard({ label, value, unit = '', color }: Props) {
  return (
    <Card className="hover:border-primary/30 transition-colors">
      <CardContent className="p-4">
        <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground mb-1">{label}</p>
        <p className="text-2xl font-bold tracking-tight" style={color ? { color } : undefined}>
          {value}{unit}
        </p>
      </CardContent>
    </Card>
  )
}
