import { Card } from '@/components/ui/card'

interface Props {
  label: string
  value: string | number
  unit?: string
  color?: string
}

export function MetricCard({ label, value, unit = '', color }: Props) {
  return (
    <Card className="border-border/60 p-3 hover:border-primary/30 transition-colors duration-150">
      <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground mb-1">{label}</p>
      <p className="text-xl font-bold tracking-tight leading-none" style={color ? { color } : undefined}>
        {value}{unit}
      </p>
    </Card>
  )
}
