<template>
  <div class="app-root">
    <AppNavbar />
    <main :class="['main-content', { 'main-content--full': isFullBleedPage }]">
      <RouterView v-slot="{ Component, route }">
        <Transition name="fade-slide" mode="out-in">
          <component :is="Component" :key="route.path" />
        </Transition>
      </RouterView>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppNavbar from './components/AppNavbar.vue'

const route = useRoute()
const fullBleedRoutes = ['/', '/login', '/register']
const isFullBleedPage = computed(() =>
  fullBleedRoutes.includes(route.path) || route.path.startsWith('/destinations/')
)
</script>

<style scoped>
.app-root {
  min-height: 100vh;
}

.main-content {
  padding-top: var(--nav-h);
}

.main-content--full {
  padding-top: 0;
}
</style>
