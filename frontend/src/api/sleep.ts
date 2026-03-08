import { api } from './client'

export interface Sleep {
  id: number
  user_id: number
  sleep_time: number | null
  wake_time: number | null
  score: number
  enter_date: string
  created_at: string
}

export interface SleepScore {
  score: number
  enter_date: string
}

export interface SleepsResp {
  sleeps: Sleep[]
}

export interface ScoresResp {
  scores: SleepScore[]
}

export function enterSleep(
  token: string,
  data: {
    sleep_time?: number | null
    wake_time?: number | null
    score: number
    enter_date: string
  }
) {
  return api<{ id: number }>('/sleep/', {
    method: 'POST',
    token,
    body: JSON.stringify(data),
  })
}

export function updateSleep(
  token: string,
  id: number,
  data: {
    sleep_time?: number | null
    wake_time?: number | null
    score?: number
    enter_date?: string
  }
) {
  return api<void>(`/sleep/${id}`, {
    method: 'PUT',
    token,
    body: JSON.stringify(data),
  })
}

export function deleteSleep(token: string, id: number) {
  return api<void>(`/sleep/${id}`, {
    method: 'DELETE',
    token,
  })
}

export function getHistory(token: string, page: number) {
  return api<SleepsResp>(`/history/${page}`, { token })
}

export function getScores(token: string, page: number) {
  return api<ScoresResp>(`/table/${page}`, { token })
}
