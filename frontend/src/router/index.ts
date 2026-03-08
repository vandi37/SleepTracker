import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn, ensureUser } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', component: () => import('../views/LoginView.vue'), meta: { guest: true } },
    { path: '/register', component: () => import('../views/RegisterView.vue'), meta: { guest: true } },
    { path: '/dashboard', component: () => import('../views/DashboardView.vue'), meta: { auth: true } },
    { path: '/sleep', component: () => import('../views/SleepView.vue'), meta: { auth: true } },
    { path: '/friends', component: () => import('../views/FriendsView.vue'), meta: { auth: true } },
    { path: '/profile', component: () => import('../views/ProfileView.vue'), meta: { auth: true } },
  ],
})

router.beforeEach(async (to) => {
  if (to.meta.auth && !isLoggedIn.value) return { path: '/login' }
  if (to.meta.guest && isLoggedIn.value) return { path: '/dashboard' }
  if (to.meta.auth) await ensureUser()
  return true
})

export default router
