const BASE = '/api'

const TOKEN_KEY = 'opsagent_access_token'
const REFRESH_KEY = 'opsagent_refresh_token'

function getAccessToken() {
  return localStorage.getItem(TOKEN_KEY)
}

function getRefreshToken() {
  return localStorage.getItem(REFRESH_KEY)
}

function setTokens(access: string, refresh: string) {
  localStorage.setItem(TOKEN_KEY, access)
  localStorage.setItem(REFRESH_KEY, refresh)
}

function clearTokens() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(REFRESH_KEY)
}

async function refreshAccessToken(): Promise<boolean> {
  const refresh = getRefreshToken()
  if (!refresh) return false

  try {
    const res = await fetch(`${BASE}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refresh }),
    })
    if (!res.ok) return false

    const data = await res.json()
    if (data.access_token) {
      setTokens(data.access_token, data.refresh_token || refresh)
      return true
    }
    return false
  } catch {
    return false
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getAccessToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  }
  if (token) headers['Authorization'] = `Bearer ${token}`

  let res = await fetch(`${BASE}${path}`, { ...options, headers })

  // If 401, try to refresh token
  if (res.status === 401) {
    const refreshed = await refreshAccessToken()
    if (refreshed) {
      // Retry with new token
      const newToken = getAccessToken()
      headers['Authorization'] = `Bearer ${newToken}`
      res = await fetch(`${BASE}${path}`, { ...options, headers })
    }
  }

  if (res.status === 401) {
    clearTokens()
    window.location.href = '/login'
    throw new Error('unauthorized')
  }
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Request failed: ${res.status}`)
  }
  return res.json()
}

export const api = {
  // Auth
  getSetupRequired: () => request<{required:boolean}>('/setup/required'),
  createAdmin: (username: string, password: string) =>
    request<{success:boolean}>('/setup/admin', { method: 'POST', body: JSON.stringify({ username, password }) }),

  login: (username: string, password: string) =>
    request<{access_token:string;refresh_token:string;user:{id:string;username:string;email:string;name:string}}>(
      '/auth/login',
      { method: 'POST', body: JSON.stringify({ username, password }) }
    ).then(data => {
      if (data.access_token) {
        setTokens(data.access_token, data.refresh_token || '')
      }
      return data
    }),

  logout: () => {
    clearTokens()
    return request('/auth/logout', { method: 'POST' })
  },

  getMe: () => request<{id:string;username:string;email:string;name:string}>('/auth/me'),

  // Dashboard
  getDashboardSummary: () => request<any>('/dashboard/summary'),
  getSystemInfo: () => request<any>('/system/info'),

  // Metrics
  getLatestMetrics: () => request<any>('/metrics/latest'),
  getMetricsHistory: (limit = 100) => request<any>(`/metrics/history?limit=${limit}`),

  // Alerts
  getAlerts: (limit = 50, offset = 0) => request<any>(`/alerts?limit=${limit}&offset=${offset}`),
  getAlert: (id: number) => request<any>(`/alerts/${id}`),
  acknowledgeAlert: (id: number) => request('/alerts/'+id+'?acknowledge=1', { method: 'POST' }),
  resolveAlert: (id: number) => request('/alerts/'+id+'?resolve=1', { method: 'POST' }),
  ignoreAlert: (id: number) => request('/alerts/'+id+'?ignore=1', { method: 'POST' }),

  // Assistant
  sendMessage: (message: string) =>
    request<{type:string;tid:string;summary:string;requires_approval:boolean}>('/assistant/message', { method: 'POST', body: JSON.stringify({ message }) }),

  // Tasks
  getTasks: (limit = 50, offset = 0) => request<any>(`/tasks?limit=${limit}&offset=${offset}`),
  getTask: (id: number) => request<any>(`/tasks/${id}`),
  approveTask: (id: number) => request(`/tasks/${id}/approve`, { method: 'POST' }),
  rejectTask: (id: number) => request(`/tasks/${id}/reject`, { method: 'POST' }),
  runTask: (id: number) => request<any>(`/tasks/${id}/run`, { method: 'POST' }),
  getTaskExecutions: (id: number) => request<any>(`/tasks/${id}/executions`),

  // Settings
  getSettings: () => request<any>('/settings'),
  updateAISettings: (data: {enabled:boolean;provider:string;api_key:string;model:string}) =>
    request('/settings/ai', { method: 'POST', body: JSON.stringify(data) }),
  updatePassword: (currentPassword: string, newPassword: string) =>
    request('/settings/password', { method: 'POST', body: JSON.stringify({ current_password: currentPassword, new_password: newPassword }) }),

  // System
  getServiceStatus: () => request<any>('/system/service-status'),
  checkForUpdate: () => request<any>('/system/update-check'),
  updateAgent: () => request<any>('/system/update', { method: 'POST' }),
}
