<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="spinner" style="margin-top:6rem;"></div>

    <!-- Error -->
    <div v-else-if="error" class="page">
      <div class="alert alert-error">⚠ {{ error }}</div>
      <RouterLink to="/" class="btn btn-ghost mt-2">← Back to Explore</RouterLink>
    </div>

    <template v-else-if="destination">
      <!-- ===== HERO — Full-bleed photo ===== -->
      <section class="detail-hero" :style="heroStyle">
        <div class="detail-hero-overlay"></div>
        <div class="detail-hero-inner container">
          <RouterLink to="/" class="detail-back-btn btn btn-glass btn-sm">← All Destinations</RouterLink>
          <div class="detail-hero-text">
            <p class="section-label-light">📍 {{ destination.country }}</p>
            <h1 class="hero-display" style="font-size:clamp(2.2rem,6vw,4.5rem);">
              {{ destination.name }}
            </h1>
            <div class="detail-hero-meta">
              <span class="stars" style="font-size:1.1rem;">{{ starStr(destination.averageRating) }}</span>
              <span style="color:var(--gold);font-weight:700;font-size:1.1rem;">{{ destination.averageRating?.toFixed(1)
                || '—' }}</span>
              <span style="color:rgba(255,255,255,0.6);font-size:0.85rem;">{{ destination.reviewCount || 0 }}
                reviews</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Torn paper divider -->
      <div class="torn-divider"></div>

      <!-- ===== CONTENT ===== -->
      <div class="detail-body container">
        <div class="detail-grid">
          <!-- Left: about + reviews -->
          <div class="detail-main">
            <!-- About card -->
            <div class="detail-about-card">
              <p class="section-label">✦ ABOUT THIS DESTINATION</p>
              <p class="detail-desc">{{ destination.description }}</p>
            </div>

            <!-- Reviews section -->
            <div class="detail-reviews">
              <div class="flex align-center justify-between mb-2">
                <h2>Traveller <em>Reviews</em></h2>
                <span class="badge badge-orange">{{ reviews.length }} reviews</span>
              </div>

              <!-- Write review form -->
              <div v-if="auth.isLoggedIn && !hasReviewed" class="review-form-card">
                <p class="section-label">✍ SHARE YOUR EXPERIENCE</p>
                <div v-if="reviewError" class="alert alert-error">⚠ {{ reviewError }}</div>
                <div v-if="reviewSuccess" class="alert alert-success">✓ {{ reviewSuccess }}</div>
                <form @submit.prevent="submitReview">
                  <div class="form-group">
                    <label>Your Rating</label>
                    <div class="stars-input">
                      <span v-for="n in 5" :key="n" :class="{ active: n <= reviewForm.rating }"
                        @click="reviewForm.rating = n" @mouseover="hoverRating = n" @mouseleave="hoverRating = 0"
                        :style="{ color: n <= (hoverRating || reviewForm.rating) ? 'var(--gold)' : '#d4d0c8' }">★</span>
                    </div>
                  </div>
                  <div class="form-group">
                    <label>Your Review</label>
                    <textarea v-model="reviewForm.text" required
                      placeholder="Tell others what made this place special..." rows="3"></textarea>
                  </div>
                  <button class="btn btn-primary btn-arrow" type="submit"
                    :disabled="reviewLoading || !reviewForm.rating">
                    {{ reviewLoading ? 'Posting...' : 'Post Review' }}
                    <span class="btn-arrow-icon">↗</span>
                  </button>
                </form>
              </div>

              <div v-else-if="!auth.isLoggedIn" class="review-prompt">
                <span style="font-size:1.5rem;">🔐</span>
                <div>
                  <strong style="color:var(--text-dark);">Want to share your experience?</strong>
                  <p style="font-size:0.85rem;color:var(--text-muted);margin-top:0.15rem;">
                    <RouterLink to="/login" style="color:var(--green-deep);font-weight:600;text-decoration:none;">Sign
                      in</RouterLink> to write a review.
                  </p>
                </div>
              </div>

              <div v-else-if="hasReviewed" class="alert alert-info">✓ You've already reviewed this destination.</div>

              <!-- Review list -->
              <div v-if="reviews.length === 0" class="empty-state" style="padding:2rem 0;">
                <div class="empty-icon">💬</div>
                <p>No reviews yet. Be the first!</p>
              </div>
              <div v-else class="review-list">
                <div v-for="r in reviews" :key="r.id" class="review-item">
                  <div class="flex align-center justify-between">
                    <div class="flex align-center gap-1">
                      <div class="reviewer-avatar">{{ r.userName?.[0]?.toUpperCase() || '?' }}</div>
                      <strong style="color:var(--text-dark);">{{ r.userName }}</strong>
                    </div>
                    <span class="stars" style="font-size:0.9rem;">{{ '★'.repeat(r.rating) }}{{ '☆'.repeat(5 - r.rating)
                    }}</span>
                  </div>
                  <p style="margin-top:0.6rem;">{{ r.text }}</p>
                  <div class="review-meta">{{ formatDate(r.createdAt) }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Right: sidebar -->
          <aside class="detail-sidebar">
            <div class="sidebar-card">
              <p class="section-label">✦ QUICK FACTS</p>
              <div class="fact-item">
                <span class="fact-icon">🌍</span>
                <div>
                  <div class="fact-label">Country</div>
                  <div class="fact-value">{{ destination.country }}</div>
                </div>
              </div>
              <div class="fact-item">
                <span class="fact-icon">⭐</span>
                <div>
                  <div class="fact-label">Average Rating</div>
                  <div class="fact-value">{{ destination.averageRating?.toFixed(1) || 'Not rated' }}</div>
                </div>
              </div>
              <div class="fact-item">
                <span class="fact-icon">💬</span>
                <div>
                  <div class="fact-label">Total Reviews</div>
                  <div class="fact-value">{{ destination.reviewCount || 0 }}</div>
                </div>
              </div>
            </div>

            <div v-if="auth.isLoggedIn" class="sidebar-card">
              <p class="section-label">✦ PLAN YOUR TRIP</p>
              <p style="color:var(--text-muted);font-size:0.85rem;margin-bottom:1rem;">Add this to an itinerary.</p>
              <RouterLink to="/itineraries" class="btn btn-primary" style="width:100%;justify-content:center;">
                🗺 My Itineraries
              </RouterLink>
            </div>
          </aside>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../api'

const route = useRoute()
const auth = useAuthStore()
const destination = ref(null)
const reviews = ref([])
const loading = ref(true)
const error = ref('')
const reviewError = ref('')
const reviewSuccess = ref('')
const reviewLoading = ref(false)
const hoverRating = ref(0)
const reviewForm = reactive({ rating: 0, text: '' })

const heroStyle = computed(() => ({
  backgroundImage: destination.value?.imageUrl
    ? `url(${destination.value.imageUrl})`
    : `url(https://images.unsplash.com/photo-1488085061387-422e29b40080?w=1200&q=80)`
}))
const hasReviewed = computed(() =>
  auth.isLoggedIn && reviews.value.some(r => r.userId === auth.user?.id)
)

function starStr(r) { const n = Math.round(r || 0); return '★'.repeat(n) + '☆'.repeat(5 - n) }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' }) : '' }

async function fetchDestination() {
  loading.value = true; error.value = ''
  try {
    const res = await api.get(`/destinations/${route.params.id}`)
    destination.value = res.data.destination
    reviews.value = res.data.reviews || []
  } catch (e) {
    error.value = e.response?.data?.error || 'Failed to load destination'
  } finally { loading.value = false }
}

async function submitReview() {
  if (!reviewForm.rating) return
  reviewLoading.value = true; reviewError.value = ''; reviewSuccess.value = ''
  try {
    await api.post(`/destinations/${route.params.id}/reviews`, { rating: reviewForm.rating, text: reviewForm.text })
    reviewSuccess.value = 'Review posted! Thank you.'
    reviewForm.rating = 0; reviewForm.text = ''
    await fetchDestination()
  } catch (e) { reviewError.value = e.response?.data?.error || 'Failed to post review' }
  finally { reviewLoading.value = false }
}

onMounted(fetchDestination)
</script>

<style scoped>
/* Hero */
.detail-hero {
  min-height: 55vh;
  background-size: cover;
  background-position: center;
  position: relative;
  display: flex;
  align-items: flex-end;
}

.detail-hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7) 0%, rgba(0, 0, 0, 0.2) 50%, rgba(0, 0, 0, 0.1) 100%);
}

.detail-hero-inner {
  position: relative;
  z-index: 2;
  padding: calc(var(--nav-h) + 2rem) 1.5rem 3rem;
}

.detail-back-btn {
  margin-bottom: 1.5rem;
}

.detail-hero-meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 0.75rem;
}

/* Body */
.detail-body {
  padding: 2rem 1.5rem 5rem;
}

.detail-grid {
  margin-top: 2rem;
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 2.5rem;
  align-items: start;
}

.detail-about-card {
  background: #fff;
  border-radius: var(--r-lg);
  padding: 2rem;
  box-shadow: var(--shadow-card);
  margin-bottom: 2.5rem;
}

.detail-desc {
  color: var(--text-body);
  line-height: 1.8;
  margin-top: 0.75rem;
}

/* Reviews */
.review-form-card {
  background: #fff;
  border-radius: var(--r-lg);
  padding: 1.75rem;
  box-shadow: var(--shadow-card);
  margin-bottom: 1.5rem;
}

.review-prompt {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  background: #fff;
  border-radius: var(--r-lg);
  box-shadow: var(--shadow-sm);
  margin-bottom: 1.5rem;
}

/* Sidebar */
.sidebar-card {
  background: #fff;
  border-radius: var(--r-lg);
  padding: 1.5rem;
  box-shadow: var(--shadow-card);
  margin-bottom: 1.25rem;
}

.fact-item {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  padding: 0.75rem 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.fact-item:last-child {
  border-bottom: none;
}

.fact-icon {
  font-size: 1.2rem;
}

.fact-label {
  font-size: 0.72rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.fact-value {
  font-size: 0.95rem;
  color: var(--text-dark);
  font-weight: 600;
  margin-top: 0.1rem;
}

@media (max-width: 900px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }

  .detail-sidebar {
    order: -1;
  }
}
</style>
