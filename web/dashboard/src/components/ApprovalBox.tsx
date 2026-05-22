interface Command { cmd: string; risk: string }
interface Props {
  title: string
  summary: string
  commands: Command[]
  onApprove: () => void
  onReject: () => void
  loading?: boolean
}
export function ApprovalBox({ title, summary, commands, onApprove, onReject, loading }: Props) {
  return (
    <div className="approval-box">
      <h3>{title}</h3>
      <p style={{color:'#94a3b8', marginTop:8}}>{summary}</p>
      <div className="approval-commands">
        {commands.map((c, i) => (
          <div key={i} className={`cmd-block risk-${c.risk}`}>{c.cmd}</div>
        ))}
      </div>
      <div className="approval-actions">
        <button className="primary" onClick={onApprove} disabled={loading}>Approve</button>
        <button className="danger" onClick={onReject} disabled={loading}>Reject</button>
      </div>
    </div>
  )
}
