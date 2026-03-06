<template>
  <div class="auth-page">
    <img src="https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?w=1600&q=80" alt="Travel"
      class="auth-bg-photo" />
    <div class="auth-bg-overlay"></div>

    <div class="auth-card-wrapper">
      <div class="auth-card">
        <!-- Left: branding -->
        <div class="auth-brand">
          <div>
            <p class="section-label-light">✦ JOIN THE COMMUNITY</p>
            <h1 style="color:#fff;font-size:clamp(1.8rem,4vw,2.8rem);margin-bottom:0.75rem;">
              Start Your<br /><em>Wanderlust</em><br />Journey
            </h1>
            <p style="color:rgba(255,255,255,0.65);font-size:0.9rem;line-height:1.6;">
              Create an account to save destinations, build itineraries,
              and share your travel experiences.
            </p>
          </div>
          <div class="auth-perks">
            <div class="auth-perk"><span>🌍</span> Explore curated destinations</div>
            <div class="auth-perk"><span>🗺️</span> Build custom itineraries</div>
            <div class="auth-perk"><span>⭐</span> Rate & review places</div>
          </div>
        </div>

        <!-- Right: form -->
        <div class="auth-form-panel">
          <h2 style="margin-bottom:0.5rem;">Create Account</h2>
          <p style="color:var(--text-muted);font-size:0.9rem;margin-bottom:1.75rem;">Fill in your details to get started
          </p>

          <div v-if="error" class="alert alert-error">⚠ {{ error }}</div>
          <div v-if="success" class="alert alert-success">✓ {{ success }}</div>

          <form @submit.prevent="handleRegister">
            <div class="form-group">
              <label>Full Name</label>
              <input v-model="name" required placeholder="Your name" />
            </div>
            <div class="form-group">
              <label>Email</label>
              <input v-model="email" type="email" required placeholder="you@example.com" />
            </div>
            <div class="form-group">
              <label>Password</label>
              <input v-model="password" type="password" required placeholder="Min 6 characters" />
            </div>
            <button class="btn btn-warm btn-arrow btn-lg" type="submit" :disabled="loading"
              style="width:100%;justify-content:center;">
              {{ loading ? 'Creating...' : 'Create Account' }}
              <span class="btn-arrow-icon">↗</span>
            </button>
          </form>

          <p style="text-align:center;margin-top:1.25rem;font-size:0.88rem;color:var(--text-muted);">
            Already have an account?
            <RouterLink to="/login" style="color:var(--orange-warm);font-weight:600;text-decoration:none;">Sign in
            </RouterLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleRegister() {
  loading.value = true; error.value = ''; success.value = ''
  try {
    await api.post('/auth/register', { name: name.value, email: email.value, password: password.value })
    success.value = 'Account created! Redirecting...'
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    error.value = e.response?.data?.error || 'Registration failed'
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

.auth-perks {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
  margin-top: 1.5rem;
}

.auth-perk {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.85rem;
}

.auth-form-panel {
  background: #fff;
  padding: 2.5rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
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
