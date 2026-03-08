import { api } from './client'

export interface User {
  id: number
  username: string
  nickname: string
  birth: string
  created_at: string
}

export function getSelf(token: string) {
  return api<User>('/users/', { token })
}

export function getUser(id: number, token: string) {
  return api<User>(`/users/${id}`, { token })
}

export function updateUser(token: string, data: { username?: string; nickname?: string; birth?: string }) {
  return api<void>('/users/', {
    method: 'PUT',
    token,
    body: JSON.stringify(data),
  })
}

export function updatePassword(token: string, password: string) {
  return api<void>('/users/password', {
    method: 'PATCH',
    token,
    body: JSON.stringify({ password }),
  })
}

export function deleteUser(token: string) {
  return api<void>('/users/', {
    method: 'DELETE',
    token,
  })
}
