const API_BASE = '/api'

export interface JsonError {
  code: number
  message: string
  context?: Record<string, unknown>
}

export async function api<T>(
  path: string,
  opts: RequestInit & { token?: string } = {}
): Promise<T> {
  const { token, ...init } = opts
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...init.headers as Record<string, string>,
  }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${API_BASE}${path}`, { ...init, headers })

  if (!res.ok) {
    const data = await res.json().catch(() => null) as JsonError | null
    throw new Error(data?.message || `HTTP ${res.status}`)
  }
  if (res.status === 204) return undefined as T
  const data = await res.json().catch(() => null)
  return data as T
}
