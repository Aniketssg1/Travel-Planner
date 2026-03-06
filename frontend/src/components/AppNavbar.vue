<template>
  <nav class="navbar" :class="navClass">
    <!-- Brand -->
    <RouterLink to="/" class="brand">
      <span class="brand-icon">✈</span>
      Wanderlust
    </RouterLink>

    <!-- Pill nav -->
    <div class="nav-pills">
      <RouterLink to="/">Explore</RouterLink>

      <template v-if="auth.isLoggedIn">
        <RouterLink to="/itineraries">My Trips</RouterLink>
        <RouterLink v-if="auth.isAdmin" to="/admin">Admin</RouterLink>
      </template>
      <template v-else>
        <RouterLink to="/login">Sign in</RouterLink>
      </template>
    </div>

    <!-- Right side -->
    <div class="flex align-center gap-1">
      <template v-if="auth.isLoggedIn">
        <div class="nav-user-badge">
          <div class="nav-avatar">{{ initials }}</div>
          <span>{{ auth.userName }}</span>
        </div>
        <button class="nav-logout" @click="handleLogout">Sign out</button>
      </template>
      <template v-else>
        <RouterLink to="/register" class="nav-cta">
          Get Started
          <span class="nav-cta-arrow">↗</span>
        </RouterLink>
      </template>
    </div>
  </nav>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const isScrolled = ref(false)

const initials = computed(() => {
  const name = auth.userName || ''
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
})

// Only use transparent style on the home page
const navClass = computed(() => {
  const isHomePage = route.path === '/'
  if (isHomePage && !isScrolled.value) return 'navbar--transparent'
  return 'navbar--solid'
})

function handleScroll() {
  isScrolled.value = window.scrollY > 40
}
function handleLogout() {
  auth.logout()
  router.push('/login')
}

onMounted(() => window.addEventListener('scroll', handleScroll))
onUnmounted(() => window.removeEventListener('scroll', handleScroll))
</script>
