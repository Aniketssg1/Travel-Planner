<template>
  <div class="auth-page">
    <!-- Full-bleed photo background -->
    <img src="https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1600&q=80" alt="Travel"
      class="auth-bg-photo" />
    <div class="auth-bg-overlay"></div>

    <div class="auth-card-wrapper">
      <div class="auth-card">
        <!-- Left: branding panel -->
        <div class="auth-brand">
          <div>
            <p class="section-label-light">✦ WELCOME BACK</p>
            <h1 style="color:#fff;font-size:clamp(1.8rem,4vw,2.8rem);margin-bottom:0.75rem;">
              Your Next<br /><em>Adventure</em><br />Starts Here
            </h1>
            <p style="color:rgba(255,255,255,0.65);font-size:0.9rem;line-height:1.6;">
              Sign in to access your curated itineraries, reviews, and
              personalised destination recommendations.
            </p>
          </div>

          <!-- Floating destination chips -->
          <div class="auth-chips">
            <span class="auth-chip">🏔️ Himalayas</span>
            <span class="auth-chip">🏖️ Maldives</span>
            <span class="auth-chip">🗼 Paris</span>
          </div>
        </div>

        <!-- Right: form -->
        <div class="auth-form-panel">
          <h2 style="margin-bottom:0.5rem;">Sign In</h2>
          <p style="color:var(--text-muted);font-size:0.9rem;margin-bottom:1.75rem;">Enter your credentials to continue
          </p>

          <div v-if="error" class="alert alert-error">⚠ {{ error }}</div>

          <form @submit.prevent="handleLogin">
            <div class="form-group">
              <label>Email</label>
              <input v-model="email" type="email" required placeholder="you@example.com" />
            </div>
            <div class="form-group">
              <label>Password</label>
              <input v-model="password" type="password" required placeholder="••••••••" />
            </div>
            <button class="btn btn-primary btn-arrow btn-lg" type="submit" :disabled="loading"
              style="width:100%;justify-content:center;">
              {{ loading ? 'Signing in...' : 'Sign In' }}
              <span class="btn-arrow-icon">↗</span>
            </button>
          </form>

          <p style="text-align:center;margin-top:1.25rem;font-size:0.88rem;color:var(--text-muted);">
            Don't have an account?
            <RouterLink to="/register" style="color:var(--green-deep);font-weight:600;text-decoration:none;">Sign up
            </RouterLink>
          </p>

          <!-- <div class="auth-hint">
            <span style="font-size:0.75rem;">🔑</span>
            Admin: admin@travel.com / Admin@123
          </div> -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../api'

const router = useRouter()
const auth = useAuthStore()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  loading.value = true; error.value = ''
  try {
    const res = await api.post('/auth/login', { email: email.value, password: password.value })
    auth.setAuth(res.data.token, res.data.user)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Login failed'
  } finally { loading.value = false }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  padding: 2rem;
}

.auth-bg-photo {
  position: fixed;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
}

.auth-bg-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 1;
}

.auth-card-wrapper {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 880px;
}

.auth-card {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border-radius: var(--r-xl);
  overflow: hidden;
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.4);
}

.auth-brand {
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  padding: 2.5rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.auth-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1.5rem;
}

.auth-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  padding: 0.35rem 0.85rem;
  border-radius: var(--r-full);
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.78rem;
  font-weight: 500;
}

.auth-form-panel {
  background: #fff;
  padding: 2.5rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.auth-hint {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 1.25rem;
  padding: 0.6rem 0.85rem;
  background: var(--bg-cream);
  border-radius: var(--r-md);
  font-size: 0.75rem;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .auth-card {
    grid-template-columns: 1fr;
  }

  .auth-brand {
    display: none;
  }
}
</style>
