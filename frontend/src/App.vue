<script setup lang="ts">
import { RouterView } from 'vue-router'
import { accessToken } from './stores/auth'
import { clearAuth } from './stores/auth'

function logout() {
  clearAuth()
  window.location.href = '/login'
}
</script>

<template>
  <div id="app">
    <header v-if="accessToken" class="header">
      <span class="brand">SleepTracker</span>
      <button @click="logout" class="logout">Log out</button>
    </header>
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style>
* { box-sizing: border-box; }
body { margin: 0; font-family: system-ui, sans-serif; }
#app { min-height: 100vh; }
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1.5rem;
  background: #1a1a2e;
  color: white;
}
.brand { font-weight: bold; }
.logout { background: transparent; color: #fff; border: 1px solid #555; padding: 0.4rem 0.8rem; cursor: pointer; border-radius: 4px; }
.logout:hover { background: #333; }
main { padding: 0 1rem 2rem; }
</style>
