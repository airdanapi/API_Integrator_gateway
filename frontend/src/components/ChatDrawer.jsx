import { useCallback, useEffect, useState } from 'react'
import {
  fetchChatConversations,
  fetchChatHistory,
  markChatRead,
  sendChatMessage,
} from '../services/chat'

const POLL_INTERVAL_MS = 2_000

// ─── ChatDrawer ─────────────────────────────────────────────────────────────

function ChatDrawer({ apiClient, currentUser }) {
  const [open, setOpen] = useState(false)
  const [conversations, setConversations] = useState([])
  const [totalUnread, setTotalUnread] = useState(0)
  const [activeConv, setActiveConv] = useState(null)
  const [messages, setMessages] = useState([])
  const [draft, setDraft] = useState('')
  const [error, setError] = useState(null)
  const [sending, setSending] = useState(false)

  // ── poll conversations ──────────────────────────────────────────────────────
  const loadConversations = useCallback(async () => {
    try {
      const result = await fetchChatConversations(apiClient)
      setConversations(result.conversations ?? [])
      setTotalUnread(result.total_unread ?? 0)
      setError(null)
    } catch {
      setError('Gagal memuat chat.')
    }
  }, [apiClient])

  useEffect(() => {
    const initialLoad = setTimeout(() => {
      loadConversations()
    }, 0)
    const timer = setInterval(() => loadConversations(), POLL_INTERVAL_MS)
    return () => {
      clearTimeout(initialLoad)
      clearInterval(timer)
    }
  }, [loadConversations])

  // ── open conversation ───────────────────────────────────────────────────────
  async function handleSelectConversation(conv) {
    setActiveConv(conv)
    try {
      const history = await fetchChatHistory(apiClient, conv.conversation_id)
      setMessages(history.messages ?? [])
    } catch {
      setMessages([])
    }
    try {
      await markChatRead(apiClient, conv.conversation_id)
    } catch {
      // silent
    }
  }

  // ── send message ────────────────────────────────────────────────────────────
  async function handleSend(e) {
    e.preventDefault()
    if (!draft.trim() || !activeConv || sending) return
    setSending(true)
    try {
      const result = await sendChatMessage(apiClient, {
        to_username: activeConv.target_username,
        to_app_name: activeConv.target_app_name,
        message: draft.trim(),
      })
      setMessages((prev) => [...prev, result.message])
      setDraft('')
      setTotalUnread(result.total_unread ?? 0)
    } catch {
      // keep draft
    } finally {
      setSending(false)
    }
  }

  // ── render ──────────────────────────────────────────────────────────────────
  return (
    <div className="relative">
      {/* Trigger button */}
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="relative rounded-xl border border-slate-300 px-3 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-50"
        aria-label="Chat"
      >
        Chat
        {totalUnread > 0 && (
          <span className="absolute -right-1.5 -top-1.5 flex h-5 min-w-5 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-bold text-white">
            {totalUnread > 99 ? '99+' : totalUnread}
          </span>
        )}
      </button>

      {/* Drawer panel */}
      {open && (
        <div className="absolute right-0 top-full z-50 mt-2 flex h-[28rem] w-80 flex-col rounded-2xl border border-slate-200 bg-white shadow-xl">
          {/* Header */}
          <div className="flex items-center justify-between border-b border-slate-100 px-4 py-3">
            <h3 className="text-sm font-bold text-slate-800">
              {activeConv ? activeConv.target_app_name : 'Percakapan'}
            </h3>
            <button
              type="button"
              onClick={() => { setOpen(false); setActiveConv(null); setMessages([]) }}
              className="text-xs font-semibold text-slate-500 hover:text-slate-700"
            >
              Tutup
            </button>
          </div>

          {/* Error */}
          {error && (
            <div role="alert" className="border-b border-red-100 bg-red-50 px-4 py-2 text-xs text-red-700">
              {error}
            </div>
          )}

          {/* Content */}
          <div className="flex-1 overflow-y-auto">
            {!activeConv && (
              /* Conversation list */
              conversations.length === 0 && !error ? (
                <p className="px-4 py-8 text-center text-sm text-slate-400">
                  Belum ada percakapan.
                </p>
              ) : (
                <ul className="divide-y divide-slate-50">
                  {conversations.map((conv) => (
                    <li key={conv.conversation_id}>
                      <button
                        type="button"
                        onClick={() => handleSelectConversation(conv)}
                        className="flex w-full items-center justify-between px-4 py-3 text-left transition hover:bg-slate-50"
                      >
                        <div>
                          <p className="text-sm font-semibold text-slate-800">
                            {conv.target_app_name}
                          </p>
                          {conv.latest_message && (
                            <p className="mt-0.5 truncate text-xs text-slate-500">
                              {conv.latest_message.message}
                            </p>
                          )}
                        </div>
                        {conv.unread_count > 0 && (
                          <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-bold text-white">
                            {conv.unread_count}
                          </span>
                        )}
                      </button>
                    </li>
                  ))}
                </ul>
              )
            )}

            {activeConv && (
              /* Message history */
              <div className="flex flex-col gap-2 px-4 py-3">
                {messages.map((msg) => {
                  const isMine = msg.from_user === currentUser?.username
                  return (
                    <div
                      key={msg.id}
                      className={`max-w-[80%] rounded-xl px-3 py-2 text-sm ${
                        isMine
                          ? 'ml-auto bg-blue-500 text-white'
                          : 'bg-slate-100 text-slate-800'
                      }`}
                    >
                      <p>{msg.message}</p>
                      <p className={`mt-0.5 text-[10px] ${isMine ? 'text-blue-100' : 'text-slate-400'}`}>
                        {new Date(msg.timestamp).toLocaleTimeString('id-ID', {
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </p>
                    </div>
                  )
                })}
              </div>
            )}
          </div>

          {/* Send form */}
          {activeConv && (
            <form
              onSubmit={handleSend}
              className="flex items-center gap-2 border-t border-slate-100 px-4 py-3"
            >
              <label htmlFor="chat-draft" className="sr-only">Pesan</label>
              <input
                id="chat-draft"
                type="text"
                aria-label="Pesan"
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
                placeholder="Tulis pesan..."
                className="flex-1 rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-blue-400 focus:ring-1 focus:ring-blue-400"
              />
              <button
                type="submit"
                disabled={sending || !draft.trim()}
                className="rounded-lg bg-blue-500 px-4 py-2 text-sm font-bold text-white transition hover:bg-blue-600 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Kirim
              </button>
            </form>
          )}
        </div>
      )}
    </div>
  )
}

export default ChatDrawer
