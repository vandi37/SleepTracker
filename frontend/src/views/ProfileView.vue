<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { accessToken, user, setUser, clearAuth } from '../stores/auth'
import { getSelf, updateUser, updatePassword, deleteUser } from '../api/users'

const router = useRouter()
const nickname = ref('')
const birth = ref('')
const newPassword = ref('')
const loading = ref(false)
const profileLoading = ref(true)

onMounted(async () => {
  const t = accessToken.value
  if (!t) return
  try {
    const u = await getSelf(t)
    setUser(u)
    nickname.value = u.nickname
    birth.value = u.birth || ''
  } catch (e) {
    clearAuth()
    router.push('/login')
  } finally {
    profileLoading.value = false
  }
})

async function saveProfile() {
  const t = accessToken.value
  if (!t) return
  loading.value = true
  try {
    await updateUser(t, {
      username: user.value?.username,
      nickname: nickname.value,
      birth: birth.value || undefined,
    })
    const u = await getSelf(t)
    setUser(u)
    alert('Profile updated')
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  } finally {
    loading.value = false
  }
}

async function changePassword() {
  const t = accessToken.value
  if (!t || !newPassword.value) return
  loading.value = true
  try {
    await updatePassword(t, newPassword.value)
    newPassword.value = ''
    alert('Password updated')
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  } finally {
    loading.value = false
  }
}

async function deleteAccount() {
  if (!confirm('Delete your account? This cannot be undone.')) return
  const t = accessToken.value
  if (!t) return
  try {
    await deleteUser(t)
    clearAuth()
    router.push('/login')
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  }
}
</script>

<template>
  <div class="profile-view">
    <nav class="nav">
      <router-link to="/dashboard">Dashboard</router-link>
      <router-link to="/sleep">Sleep</router-link>
      <router-link to="/friends">Friends</router-link>
      <router-link to="/profile">Profile</router-link>
    </nav>
    <h1>Profile</h1>
    <div v-if="profileLoading">Loading…</div>
    <div v-else-if="user" class="form">
      <p><strong>Username:</strong> {{ user.username }}</p>
      <p><strong>ID:</strong> {{ user.id }}</p>
      <label>Nickname: <input v-model="nickname" /></label>
      <label>Birth: <input v-model="birth" type="date" /></label>
      <button @click="saveProfile" :disabled="loading">Save profile</button>
    </div>
    <section class="password">
      <h2>Change password</h2>
      <input v-model="newPassword" type="password" placeholder="New password" />
      <button @click="changePassword" :disabled="loading">Update</button>
    </section>
    <section class="danger">
      <button class="danger-btn" @click="deleteAccount">Delete account</button>
    </section>
  </div>
</template>

<style scoped>
.profile-view { max-width: 500px; margin: 0 auto; padding: 1rem; }
.nav { display: flex; gap: 1rem; margin-bottom: 1.5rem; }
.nav a { color: #1a73e8; text-decoration: none; }
.nav a.router-link-active { font-weight: bold; }
.form { display: flex; flex-direction: column; gap: 0.75rem; margin-bottom: 2rem; }
.form input { margin-left: 0.5rem; padding: 0.4rem; }
.password, .danger { margin-top: 1.5rem; }
.danger-btn { background: #c00; color: white; padding: 0.5rem 1rem; border: none; cursor: pointer; }
button { padding: 0.4rem 0.8rem; cursor: pointer; }
</style>
