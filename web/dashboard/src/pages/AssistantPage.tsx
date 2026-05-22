import { useState, useRef, useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { ScrollArea } from '@/components/ui/scroll-area'
import { ApprovalBox } from '@/components/ApprovalBox'
import { api } from '@/lib/api'

interface Message {
  role: 'user' | 'assistant'
  content: string
  type?: string
  tid?: string
  summary?: string
  commands?: any[]
}

export default function AssistantPage() {
  const location = useLocation()
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState(location.state?.prompt || '')
  const [loading, setLoading] = useState(false)
  const bottomRef = useRef<HTMLDivElement>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

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
    <div className="flex flex-col h-[calc(100vh-8rem)]">
      <div className="mb-4">
        <h1 className="text-2xl font-bold tracking-tight">Assistant</h1>
        <p className="text-sm text-muted-foreground mt-1">AI-powered system assistant</p>
      </div>

      <Card className="flex-1 flex flex-col min-h-0">
        <ScrollArea className="flex-1 p-4" ref={scrollRef}>
          <div className="space-y-4">
            {messages.length === 0 && (
              <div className="text-center py-12 text-muted-foreground text-sm">
                Ask me anything about your system...
              </div>
            )}
            {messages.map((msg, i) => (
              <div key={i} className="space-y-2">
                <div className={`text-xs font-medium ${msg.role === 'user' ? 'text-primary' : 'text-emerald-400'}`}>
                  {msg.role === 'user' ? 'You' : 'OpsAgent'}
                </div>
                <div className="text-sm pl-3 border-l-2 border-border">{msg.content}</div>
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
            {loading && (
              <div className="text-center py-8 text-muted-foreground text-sm animate-pulse">
                Thinking...
              </div>
            )}
            <div ref={bottomRef} />
          </div>
        </ScrollArea>

        <CardContent className="p-4 border-t">
          <form onSubmit={handleSend} className="flex gap-2">
            <Textarea
              value={input}
              onChange={e => setInput(e.target.value)}
              placeholder="Ask me anything..."
              className="min-h-9 max-h-32 resize-none"
              onKeyDown={e => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault()
                  handleSend(e as any)
                }
              }}
            />
            <Button type="submit" size="sm" disabled={loading || !input.trim()} className="shrink-0">
              Send
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
