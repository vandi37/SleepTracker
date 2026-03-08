import { ref, computed } from 'vue'
import type { User } from '../api/users'
import { getSelf } from '../api/users'

const ACCESS_KEY = 'st_access'
const REFRESH_KEY = 'st_refresh'
const USER_KEY = 'st_user'

function loadStorage() {
  return {
    access: localStorage.getItem(ACCESS_KEY),
    refresh: localStorage.getItem(REFRESH_KEY),
    user: (() => {
      try {
        const u = localStorage.getItem(USER_KEY)
        return u ? JSON.parse(u) as User : null
      } catch {
        return null
      }
    })(),
  }
}

const accessToken = ref<string | null>(loadStorage().access)
const refreshToken = ref<string | null>(loadStorage().refresh)
const user = ref<User | null>(loadStorage().user)

export function setTokens(access: string, refresh: string, u: User | null) {
  accessToken.value = access
  refreshToken.value = refresh
  user.value = u
  if (access) localStorage.setItem(ACCESS_KEY, access)
  else localStorage.removeItem(ACCESS_KEY)
  if (refresh) localStorage.setItem(REFRESH_KEY, refresh)
  else localStorage.removeItem(REFRESH_KEY)
  if (u) localStorage.setItem(USER_KEY, JSON.stringify(u))
  else localStorage.removeItem(USER_KEY)
}

export function clearAuth() {
  setTokens('', '', null)
}

export function setUser(u: User) {
  user.value = u
  localStorage.setItem(USER_KEY, JSON.stringify(u))
}

export async function ensureUser() {
  const token = accessToken.value
  if (!token || user.value) return user.value
  try {
    const u = await getSelf(token)
    setUser(u)
    return u
  } catch {
    clearAuth()
    return null
  }
}

export const isLoggedIn = computed(() => !!accessToken.value)
export { accessToken, refreshToken, user }
