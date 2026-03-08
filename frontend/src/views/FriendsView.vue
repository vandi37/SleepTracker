<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { accessToken, user } from '../stores/auth'
import {
  getFriendships,
  requestFriend,
  acceptFriend,
  deleteFriendship,
  getFriendHistory,
  type Friend,
} from '../api/friends'

const friendships = ref<Friend[]>([])
const loading = ref(false)
const requestUserId = ref('')
const selectedFriend = ref<Friend | null>(null)
const friendSleeps = ref<import('../api/sleep').Sleep[]>([])
const friendPage = ref(0)

function load() {
  const t = accessToken.value
  if (!t) return
  loading.value = true
  getFriendships(t)
    .then((r) => { friendships.value = r.friendships })
    .catch(console.error)
    .finally(() => { loading.value = false })
}

async function sendRequest() {
  const t = accessToken.value
  const id = parseInt(requestUserId.value, 10)
  if (!t || isNaN(id)) return
  try {
    await requestFriend(t, id)
    requestUserId.value = ''
    load()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  }
}

async function accept(id: number) {
  const t = accessToken.value
  if (!t) return
  try {
    await acceptFriend(t, id)
    load()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  }
}

async function remove(id: number) {
  if (!confirm('Remove this friendship?')) return
  const t = accessToken.value
  if (!t) return
  try {
    await deleteFriendship(t, id)
    selectedFriend.value = null
    load()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  }
}

function otherUser(f: Friend) {
  const uid = user.value?.id ?? 0
  return f.user1.id === uid ? f.user2 : f.user1
}

function loadFriendHistory() {
  const t = accessToken.value
  const f = selectedFriend.value
  if (!t || !f) return
  getFriendHistory(t, f.friendship_id, friendPage.value)
    .then((r) => { friendSleeps.value = r.sleeps })
    .catch(console.error)
}

function selectFriend(f: Friend) {
  selectedFriend.value = f
  friendPage.value = 0
  if (f.is_accepted) loadFriendHistory()
}

function formatTime(m: number | null) {
  if (m == null) return '—'
  const h = Math.floor(m / 60)
  const mn = m % 60
  return `${h.toString().padStart(2, '0')}:${mn.toString().padStart(2, '0')}`
}

onMounted(load)
</script>

<template>
  <div class="friends-view">
    <nav class="nav">
      <router-link to="/dashboard">Dashboard</router-link>
      <router-link to="/sleep">Sleep</router-link>
      <router-link to="/friends">Friends</router-link>
      <router-link to="/profile">Profile</router-link>
    </nav>
    <h1>Friends</h1>
    <div class="request">
      <input v-model="requestUserId" type="number" placeholder="User ID" />
      <button @click="sendRequest">Send request</button>
    </div>
    <div v-if="loading">Loading…</div>
    <ul v-else class="list">
      <li v-for="f in friendships" :key="f.friendship_id" class="item">
        <span>{{ otherUser(f).nickname }} (@{{ otherUser(f).username }})</span>
        <span :class="{ accepted: f.is_accepted }">{{ f.is_accepted ? 'Accepted' : 'Pending' }}</span>
        <button v-if="!f.is_accepted && f.user2.id === user?.id" @click="accept(f.friendship_id)">Accept</button>
        <button @click="selectFriend(f)">View</button>
        <button @click="remove(f.friendship_id)">Remove</button>
      </li>
    </ul>
    <div v-if="selectedFriend" class="detail">
      <h2>{{ otherUser(selectedFriend).nickname }}'s sleep</h2>
      <template v-if="selectedFriend.is_accepted">
        <table v-if="friendSleeps.length">
          <thead><tr><th>Date</th><th>Sleep</th><th>Wake</th><th>Score</th></tr></thead>
          <tbody>
            <tr v-for="s in friendSleeps" :key="s.id">
              <td>{{ s.enter_date }}</td>
              <td>{{ formatTime(s.sleep_time) }}</td>
              <td>{{ formatTime(s.wake_time) }}</td>
              <td>{{ s.score }}</td>
            </tr>
          </tbody>
        </table>
        <p v-else>No sleep records.</p>
        <div class="pager">
          <button @click="friendPage--; loadFriendHistory()" :disabled="friendPage <= 0">Prev</button>
          <button @click="friendPage++; loadFriendHistory()">Next</button>
        </div>
      </template>
      <p v-else>Accept the request to view sleep data.</p>
    </div>
  </div>
</template>

<style scoped>
.friends-view { max-width: 600px; margin: 0 auto; padding: 1rem; }
.nav { display: flex; gap: 1rem; margin-bottom: 1.5rem; }
.nav a { color: #1a73e8; text-decoration: none; }
.nav a.router-link-active { font-weight: bold; }
.request { display: flex; gap: 0.5rem; margin-bottom: 1rem; }
.list { list-style: none; padding: 0; }
.item { display: flex; gap: 0.5rem; align-items: center; margin-bottom: 0.5rem; }
.accepted { color: green; }
.detail { margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid #ddd; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 0.5rem; border: 1px solid #ddd; }
.pager { margin-top: 0.5rem; }
button { padding: 0.3rem 0.6rem; cursor: pointer; margin-right: 0.3rem; }
</style>
