<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../api/auth'
import { setTokens } from '../stores/auth'
import { getSelf } from '../api/users'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')

async function submit() {
  error.value = ''
  try {
    const res = await login(username.value, password.value)
    if (!res.access || !res.refresh) throw new Error('No tokens')
    const user = await getSelf(res.access)
    setTokens(res.access, res.refresh, user)
    router.push('/dashboard')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Login failed'
  }
}
</script>

<template>
  <div class="auth-page">
    <h1>SleepTracker</h1>
    <form @submit.prevent="submit" class="auth-form">
      <input v-model="username" placeholder="Username" required />
      <input v-model="password" type="password" placeholder="Password" required />
      <p v-if="error" class="error">{{ error }}</p>
      <button type="submit">Log in</button>
    </form>
    <p>Don't have an account? <router-link to="/register">Register</router-link></p>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  width: 280px;
}
.auth-form input {
  padding: 0.5rem 0.75rem;
  border: 1px solid #333;
  border-radius: 6px;
}
.auth-form button {
  padding: 0.6rem;
  background: #1a73e8;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.auth-form button:hover {
  background: #1557b0;
}
.error { color: #c00; font-size: 0.9rem; }
a { color: #1a73e8; }
</style>
