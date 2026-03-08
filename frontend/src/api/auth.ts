import { api } from './client'

export interface UserWithToken {
  id: number
  access?: string
  expires?: string
  refresh?: string
}

export interface UserReq {
  username: string
  password: string
  nickname?: string
  birth?: string
}

export function register(req: UserReq) {
  return api<UserWithToken>('/register', {
    method: 'POST',
    body: JSON.stringify(req),
  })
}

export function login(username: string, password: string) {
  return api<UserWithToken>('/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
}

export function refresh(token: string) {
  return api<UserWithToken>('/refresh', {
    method: 'POST',
    body: JSON.stringify({ token }),
  })
}
