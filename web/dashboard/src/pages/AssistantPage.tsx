import { useState, useRef, useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import { ApprovalBox } from '../components/ApprovalBox'
import { api } from '../lib/api'

interface Message { role: 'user' | 'assistant'; content: string; type?: string; tid?: string; summary?: string; commands?: any[] }

export default function AssistantPage() {
  const location = useLocation()
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState(location.state?.prompt || '')
  const [loading, setLoading] = useState(false)
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  async function handleSend(e: React.FormEvent) {
    e.preventDefault()
    if (!input.trim() || loading) return
    const userMsg = input
    setInput('')
    setMessages(m => [...m, { role: 'user', content: userMsg }])
    setLoading(true)
    try {
      const data = await api.sendMessage(userMsg)
      setMessages(m => [...m, {
        role: 'assistant',
        content: data.summary || data.type,
        type: data.type,
        tid: data.tid,
        summary: data.summary,
        commands: [],
      }])
    } catch (err: any) {
      setMessages(m => [...m, { role: 'assistant', content: err.message || 'Error' }])
    } finally {
      setLoading(false)
    }
  }

  async function handleApprove(tid: string) {
    try {
      await api.approveTask(parseInt(tid))
      setMessages(m => m.map(msg => msg.tid === tid ? { ...msg, content: 'Approved!' } : msg))
    } catch (err) {
      console.error(err)
    }
  }

  async function handleReject(tid: string) {
    try {
      await api.rejectTask(parseInt(tid))
      setMessages(m => m.map(msg => msg.tid === tid ? { ...msg, content: 'Rejected' } : msg))
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <div className="container" style={{display:'flex', flexDirection:'column', height:'calc(100vh - 48px)'}}>
      <h1 className="page-title">Assistant</h1>
      <div style={{flex:1, overflow:'auto', marginBottom:16}}>
        {messages.map((msg, i) => (
          <div key={i} style={{marginBottom:16}}>
            <div style={{fontWeight:600, color: msg.role === 'user' ? '#60a5fa' : '#86efac'}}>
              {msg.role === 'user' ? 'You' : 'OpsAgent'}
            </div>
            <div style={{marginTop:4, whiteSpace:'pre-wrap'}}>{msg.content}</div>
            {msg.type === 'plan' && msg.tid && (
              <ApprovalBox
                title="Proposed Plan"
                summary={msg.summary || ''}
                commands={msg.commands || []}
                onApprove={() => handleApprove(msg.tid!)}
                onReject={() => handleReject(msg.tid!)}
              />
            )}
          </div>
        ))}
        {loading && <div style={{color:'#64748b'}}>Thinking...</div>}
        <div ref={bottomRef} />
      </div>
      <form onSubmit={handleSend} style={{display:'flex', gap:8}}>
        <input value={input} onChange={e => setInput(e.target.value)} placeholder="Ask me anything..." style={{flex:1}} />
        <button type="submit" className="primary" disabled={loading || !input.trim()}>Send</button>
      </form>
    </div>
  )
}
