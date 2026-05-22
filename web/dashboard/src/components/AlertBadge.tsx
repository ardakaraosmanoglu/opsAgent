interface Props { severity: string; status?: string }
export function AlertBadge({ severity, status }: Props) {
  const cls = status ? `badge ${status}` : `badge ${severity}`
  return <span className={cls}>{status || severity}</span>
}
