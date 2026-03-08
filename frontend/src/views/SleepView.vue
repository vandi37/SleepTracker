<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { accessToken } from '../stores/auth'
import {
  getHistory,
  enterSleep,
  updateSleep,
  deleteSleep,
  type Sleep,
} from '../api/sleep'

const sleeps = ref<Sleep[]>([])
const page = ref(0)
const loading = ref(false)
const editing = ref<Sleep | null>(null)
const form = ref({
  enter_date: '',
  sleep_time_str: '',
  wake_time_str: '',
  score: 80,
})

function formatTime(mins: number | null) {
  if (mins == null) return ''
  const h = Math.floor(mins / 60)
  const m = mins % 60
  return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}`
}

function parseTime(s: string): number | null {
  if (!s) return null
  const parts = s.split(':').map(Number)
  const h = parts[0], m = parts[1]
  if (h == null || m == null || isNaN(h) || isNaN(m)) return null
  return h * 60 + m
}

function load() {
  const t = accessToken.value
  if (!t) return
  loading.value = true
  getHistory(t, page.value)
    .then((r) => { sleeps.value = r.sleeps })
    .catch(console.error)
    .finally(() => { loading.value = false })
}

async function submit() {
  const t = accessToken.value
  if (!t) return
  try {
    const data = {
      enter_date: form.value.enter_date,
      sleep_time: parseTime(form.value.sleep_time_str) ?? undefined,
      wake_time: parseTime(form.value.wake_time_str) ?? undefined,
      score: form.value.score,
    }
    if (editing.value) {
      await updateSleep(t, editing.value.id, data)
    } else {
      await enterSleep(t, data)
    }
    editing.value = null
    form.value = { enter_date: '', sleep_time_str: '', wake_time_str: '', score: 80 }
    load()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  }
}

function edit(s: Sleep) {
  editing.value = s
  form.value = {
    enter_date: s.enter_date,
    sleep_time_str: s.sleep_time != null ? formatTime(s.sleep_time) : '',
    wake_time_str: s.wake_time != null ? formatTime(s.wake_time) : '',
    score: s.score,
  }
}

function cancelEdit() {
  editing.value = null
  form.value = { enter_date: '', sleep_time_str: '', wake_time_str: '', score: 80 }
}

async function remove(s: Sleep) {
  if (!confirm('Delete this record?')) return
  const t = accessToken.value
  if (!t) return
  try {
    await deleteSleep(t, s.id)
    load()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed')
  }
}

onMounted(load)
</script>

<template>
  <div class="sleep-view">
    <nav class="nav">
      <router-link to="/dashboard">Dashboard</router-link>
      <router-link to="/sleep">Sleep</router-link>
      <router-link to="/friends">Friends</router-link>
      <router-link to="/profile">Profile</router-link>
    </nav>
    <h1>Sleep records</h1>
    <form @submit.prevent="submit" class="form">
      <input v-model="form.enter_date" type="date" required placeholder="Date" />
      <input v-model="form.sleep_time_str" type="time" placeholder="Sleep" />
      <input v-model="form.wake_time_str" type="time" placeholder="Wake" />
      <input v-model.number="form.score" type="number" min="0" max="100" required />
      <div class="row">
        <button type="submit">{{ editing ? 'Update' : 'Add' }}</button>
        <button v-if="editing" type="button" @click="cancelEdit">Cancel</button>
      </div>
    </form>
    <p class="hint">Time format: HH:MM. Score 0–100. Sleep/wake optional.</p>
    <div v-if="loading">Loading…</div>
    <table v-else-if="sleeps.length">
      <thead><tr><th>Date</th><th>Sleep</th><th>Wake</th><th>Score</th><th></th></tr></thead>
      <tbody>
        <tr v-for="s in sleeps" :key="s.id">
          <td>{{ s.enter_date }}</td>
          <td>{{ formatTime(s.sleep_time) || '—' }}</td>
          <td>{{ formatTime(s.wake_time) || '—' }}</td>
          <td>{{ s.score }}</td>
          <td>
            <button @click="edit(s)">Edit</button>
            <button @click="remove(s)">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
    <p v-else>No records.</p>
    <div class="pager">
      <button @click="page--; load()" :disabled="page <= 0">Prev</button>
      <span>Page {{ page }}</span>
      <button @click="page++; load()">Next</button>
    </div>
  </div>
</template>

<style scoped>
.sleep-view { max-width: 700px; margin: 0 auto; padding: 1rem; }
.nav { display: flex; gap: 1rem; margin-bottom: 1.5rem; }
.nav a { color: #1a73e8; text-decoration: none; }
.nav a.router-link-active { font-weight: bold; }
.form { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-bottom: 0.5rem; }
.form input { padding: 0.4rem; }
.row { display: flex; gap: 0.5rem; }
.hint { font-size: 0.85rem; color: #666; margin-bottom: 1rem; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 0.5rem; border: 1px solid #ddd; }
.pager { margin-top: 0.5rem; display: flex; gap: 0.5rem; align-items: center; }
button { padding: 0.3rem 0.6rem; cursor: pointer; margin-right: 0.3rem; }
</style>
