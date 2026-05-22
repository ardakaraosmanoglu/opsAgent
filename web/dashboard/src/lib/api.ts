const BASE = '/api'

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('opsagent_token')
  const headers: Record<string,string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string,string> || {}),
  }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${BASE}${path}`, { ...options, headers })
  if (res.status === 401) {
    localStorage.removeItem('opsagent_token')
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
    request<{token:string;user:{username:string;role:string}}>('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) }),
  logout: () => request('/auth/logout', { method: 'POST' }),
  getMe: () => request<{username:string;role:string}>('/auth/me'),

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
