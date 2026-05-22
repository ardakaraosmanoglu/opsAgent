interface Props { label: string; value: string | number; unit?: string; color?: string }
export function MetricCard({ label, value, unit = '', color = '#e2e8f0' }: Props) {
  return (
    <div className="metric-card">
      <div className="metric-label">{label}</div>
      <div className="metric-value" style={{ color }}>{value}{unit}</div>
    </div>
  )
}
