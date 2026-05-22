import { useState, useRef, useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { ScrollArea } from '@/components/ui/scroll-area'
import { ApprovalBox } from '@/components/ApprovalBox'
import { Send } from 'lucide-react'
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
    <div className="flex flex-col h-[calc(100vh-8rem)] animate-fade-in">
      <div className="mb-4">
        <h1 className="text-xl font-semibold tracking-tight">Assistant</h1>
        <p className="text-xs text-muted-foreground mt-0.5">AI-powered system assistant</p>
      </div>

      <Card className="flex-1 flex flex-col min-h-0 border-border/60">
        <ScrollArea ref={scrollRef} className="flex-1">
          <div className="p-4 space-y-4">
            {messages.length === 0 && (
              <div className="text-center py-12 text-muted-foreground text-xs">
                Ask me anything about your system...
              </div>
            )}
            {messages.map((msg, i) => (
              <div key={i} className="space-y-1.5">
                <div className={`text-xs font-semibold ${msg.role === 'user' ? 'text-primary' : 'text-emerald-400'}`}>
                  {msg.role === 'user' ? 'You' : 'OpsAgent'}
                </div>
                <div className="text-xs pl-3 border-l border-border/50 leading-relaxed">{msg.content}</div>
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
              <div className="text-center py-8 text-muted-foreground text-xs animate-pulse">
                Thinking...
              </div>
            )}
            <div ref={bottomRef} />
          </div>
        </ScrollArea>

        <CardContent className="p-3 border-t border-border/60">
          <form onSubmit={handleSend} className="flex gap-2">
            <Textarea
              value={input}
              onChange={e => setInput(e.target.value)}
              placeholder="Ask me anything..."
              className="min-h-[36px] max-h-28 resize-none text-xs"
              onKeyDown={e => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault()
                  handleSend(e as any)
                }
              }}
            />
            <Button
              type="submit"
              size="icon"
              disabled={loading || !input.trim()}
              className="h-[36px] w-[36px] shrink-0"
            >
              <Send className="h-3.5 w-3.5" />
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
