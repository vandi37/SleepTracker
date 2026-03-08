<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { accessToken } from '../stores/auth'
import { getHistory, getScores } from '../api/sleep'
import type { Sleep, SleepScore } from '../api/sleep'

const sleeps = ref<Sleep[]>([])
const scores = ref<SleepScore[]>([])
const historyPage = ref(0)
const tablePage = ref(0)
const loading = ref(false)

function loadHistory() {
  const t = accessToken.value
  if (!t) return
  loading.value = true
  getHistory(t, historyPage.value)
    .then((r) => { sleeps.value = r.sleeps })
    .catch(console.error)
    .finally(() => { loading.value = false })
}

function loadScores() {
  const t = accessToken.value
  if (!t) return
  getScores(t, tablePage.value)
    .then((r) => { scores.value = r.scores })
    .catch(console.error)
}

function formatTime(mins: number | null) {
  if (mins == null) return '—'
  const h = Math.floor(mins / 60)
  const m = mins % 60
  return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}`
}

onMounted(() => {
  loadHistory()
  loadScores()
})
</script>

<template>
  <div class="dashboard">
    <nav class="nav">
      <router-link to="/dashboard">Dashboard</router-link>
      <router-link to="/sleep">Sleep</router-link>
      <router-link to="/friends">Friends</router-link>
      <router-link to="/profile">Profile</router-link>
    </nav>
    <h1>Dashboard</h1>
    <section>
      <h2>Recent sleep (week)</h2>
      <div v-if="loading">Loading…</div>
      <table v-else-if="sleeps.length">
        <thead>
          <tr><th>Date</th><th>Sleep</th><th>Wake</th><th>Score</th></tr>
        </thead>
        <tbody>
          <tr v-for="s in sleeps" :key="s.id">
            <td>{{ s.enter_date }}</td>
            <td>{{ formatTime(s.sleep_time) }}</td>
            <td>{{ formatTime(s.wake_time) }}</td>
            <td>{{ s.score }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else>No sleep records yet.</p>
      <div class="pager">
        <button @click="historyPage--; loadHistory()" :disabled="historyPage <= 0">Prev</button>
        <span>Page {{ historyPage }}</span>
        <button @click="historyPage++; loadHistory()">Next</button>
      </div>
    </section>
    <section>
      <h2>Scores (year)</h2>
      <table v-if="scores.length">
        <thead><tr><th>Date</th><th>Score</th></tr></thead>
        <tbody>
          <tr v-for="sc in scores" :key="sc.enter_date">
            <td>{{ sc.enter_date }}</td>
            <td>{{ sc.score }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else>No scores yet.</p>
      <div class="pager">
        <button @click="tablePage--; loadScores()" :disabled="tablePage <= 0">Prev</button>
        <span>Page {{ tablePage }}</span>
        <button @click="tablePage++; loadScores()">Next</button>
      </div>
    </section>
  </div>
</template>

<style scoped>
.dashboard { max-width: 800px; margin: 0 auto; padding: 1rem; }
.nav { display: flex; gap: 1rem; margin-bottom: 1.5rem; }
.nav a { color: #1a73e8; text-decoration: none; }
.nav a.router-link-active { font-weight: bold; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 0.5rem; border: 1px solid #ddd; text-align: left; }
.pager { margin-top: 0.5rem; display: flex; gap: 0.5rem; align-items: center; }
button { padding: 0.3rem 0.6rem; cursor: pointer; }
</style>
