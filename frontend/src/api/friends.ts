import { api } from './client'
import type { User } from './users'

export interface Friend {
  friendship_id: number
  user1: User
  user2: User
  is_accepted: boolean
  created_at: string
  updated_at: string
}

export interface FriendshipsResp {
  friendships: Friend[]
}

export function requestFriend(token: string, userId: number) {
  return api<{ id: number }>('/friends/request', {
    method: 'POST',
    token,
    body: JSON.stringify({ id: userId }),
  })
}

export function acceptFriend(token: string, friendshipId: number) {
  return api<void>('/friends/accept', {
    method: 'POST',
    token,
    body: JSON.stringify({ id: friendshipId }),
  })
}

export function getFriendships(token: string, limit = 20, offset = 0) {
  return api<FriendshipsResp>(`/friends/?limit=${limit}&offset=${offset}`, { token })
}

export function deleteFriendship(token: string, friendshipId: number) {
  return api<void>(`/friends/${friendshipId}`, {
    method: 'DELETE',
    token,
  })
}

export function getFriendHistory(token: string, friendshipId: number, page: number) {
  return api<{ sleeps: import('./sleep').Sleep[] }>(
    `/friends/${friendshipId}/history/${page}`,
    { token }
  )
}

export function getFriendScores(token: string, friendshipId: number, page: number) {
  return api<{ scores: import('./sleep').SleepScore[] }>(
    `/friends/${friendshipId}/table/${page}`,
    { token }
  )
}
