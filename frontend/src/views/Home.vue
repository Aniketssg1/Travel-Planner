<template>
  <div>
    <!-- ===== HERO — Full-bleed photo (BromoRise style) ===== -->
    <section class="hero">
      <img src="https://images.unsplash.com/photo-1506929562872-bb421503ef21?w=1600&q=80"
        alt="Breathtaking sunrise landscape" class="hero-photo" />
      <div class="hero-overlay"></div>

      <!-- Content over hero photo -->
      <div class="hero-inner container">
        <div class="hero-text">
          <p class="section-label-light">✦ EXPLORE · DISCOVER · TRAVEL</p>
          <h1 class="hero-display">
            <span class="wave-line" style="--wave-delay: 0s;">Unforgettable</span><br />
            <span class="wave-line accent" style="--wave-delay: 0.4s;">Destinations</span><br />
            <span class="wave-line thin" style="--wave-delay: 0.8s;">Await You</span>
          </h1>
          <p class="hero-sub">
            Travel through stunning landscapes, vibrant cultures, and timeless beauty
            with expertly curated destinations.
          </p>
          <div class="hero-actions">
            <button class="btn btn-light btn-arrow btn-lg" @click="scrollToDestinations">
              Explore Now
              <span class="btn-arrow-icon">↗</span>
            </button>
            <RouterLink v-if="!auth.isLoggedIn" to="/register" class="btn btn-glass btn-lg">
              Join Wanderlust
            </RouterLink>
          </div>
        </div>

        <!-- Floating cards + social proof (BromoRise) -->
        <div class="hero-bottom">
          <!-- Social proof -->
          <!-- <div class="hero-social-proof">
            <div class="avatar-stack">
              <div class="avatar-stack-item" style="background:#2d4a3e;">A</div>
              <div class="avatar-stack-item" style="background:#e8764b;">M</div>
              <div class="avatar-stack-item" style="background:#6bba62;">K</div>
              <div class="avatar-stack-item" style="background:#e9c46a;color:#1a1a1a;">R</div>
            </div>
            <div class="avatar-stack-count">{{ stats.destinations }}+ Destinations</div>
          </div> -->

          <!-- Floating preview cards -->
          <!-- <div class="hero-preview-cards">
            <div v-for="(img, i) in previewImages" :key="i" class="hero-float-card"
              :style="{ animationDelay: `${i * 0.15}s` }">
              <img :src="img.src" :alt="img.label" />
              <span class="hero-float-label">{{ img.label }}</span>
            </div>
          </div> -->
        </div>
      </div>
    </section>

    <!-- Wave divider: smooth multi-layer transition from hero to cream -->
    <div class="wave-divider">
      <svg viewBox="0 0 1440 150" preserveAspectRatio="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M0,80 C180,150 360,0 540,80 C720,160 900,20 1080,80 C1260,140 1380,40 1440,80 L1440,150 L0,150 Z"
          fill="rgba(245,240,232,0.35)" />
        <path d="M0,100 C200,140 400,40 600,90 C800,140 1000,50 1200,100 C1320,130 1400,70 1440,100 L1440,150 L0,150 Z"
          fill="rgba(245,240,232,0.6)" />
        <path d="M0,120 C160,135 320,95 480,115 C640,135 800,90 960,115 C1120,140 1280,100 1440,120 L1440,150 L0,150 Z"
          fill="#f5f0e8" />
      </svg>
    </div>

    <!-- ===== DESTINATIONS SECTION — Cream background ===== -->
    <section class="dest-section" ref="destSection">
      <div class="container">
        <div class="section-header">
          <p class="section-label">✦ HANDPICKED DESTINATIONS</p>
          <h2>Where Will Your<br /><em>Next Journey</em> Take You?</h2>
        </div>

        <!-- Filter Bar -->
        <div class="filter-bar">
          <div class="form-group">
            <label>🔍 Search</label>
            <input v-model="filters.name" @input="debouncedFetch" placeholder="Paris, Tokyo, Bali..." />
          </div>
          <div class="form-group">
            <label>🌍 Country</label>
            <input v-model="filters.country" @input="debouncedFetch" placeholder="France, Japan..." />
          </div>
          <div class="form-group">
            <label>⭐ Min Rating</label>
            <select v-model="filters.minRating" @change="fetchDestinations">
              <option value="">Any</option>
              <option value="1">1+ stars</option>
              <option value="2">2+ stars</option>
              <option value="3">3+ stars</option>
              <option value="4">4+ stars</option>
            </select>
          </div>
          <div style="display:flex;align-items:flex-end;">
            <button class="btn btn-ghost btn-sm" @click="clearFilters">Clear</button>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="spinner"></div>

        <!-- Error -->
        <div v-else-if="error" class="alert alert-error">⚠ {{ error }}</div>

        <!-- Grid -->
        <template v-else>
          <div v-if="destinations.length === 0" class="empty-state">
            <div class="empty-icon">🗺️</div>
            <h3>No destinations found</h3>
            <p>Try adjusting your filters to discover more places.</p>
          </div>

          <div v-else class="dest-grid">
            <RouterLink v-for="dest in destinations" :key="dest.id" :to="`/destinations/${dest.id}`" class="dest-card">
              <div class="dest-card-img-wrap">
                <img class="dest-card-img" :src="dest.imageUrl || fallbackImg(dest.name)" :alt="dest.name"
                  @error="onImgError" />
                <div class="dest-card-rating-badge" v-if="dest.averageRating">
                  <span class="stars" style="font-size:0.7rem;">★</span>
                  {{ dest.averageRating.toFixed(1) }}
                </div>
              </div>
              <div class="dest-card-body">
                <h3>{{ dest.name }}</h3>
                <div class="country">{{ dest.country }}</div>
                <div class="flex align-center justify-between mt-1">
                  <div class="flex align-center gap-1">
                    <span class="stars" style="font-size:0.8rem;">{{ starStr(dest.averageRating) }}</span>
                    <span style="font-size:0.78rem;color:var(--text-muted);">
                      {{ dest.reviewCount ? `(${dest.reviewCount})` : '' }}
                    </span>
                  </div>
                  <span class="dest-arrow">→</span>
                </div>
              </div>
            </RouterLink>
          </div>

          <!-- Pagination -->
          <div v-if="pagination.totalPages > 1" class="pagination">
            <button :disabled="pagination.page === 1" @click="goToPage(pagination.page - 1)">← Prev</button>
            <button v-for="p in pagination.totalPages" :key="p" :class="{ active: p === pagination.page }"
              @click="goToPage(p)">{{ p }}</button>
            <button :disabled="pagination.page === pagination.totalPages" @click="goToPage(pagination.page + 1)">Next
              →</button>
          </div>
        </template>
      </div>
    </section>

    <!-- ===== DASHED TRAVEL ITINERARY PATH ===== -->
    <section class="itinerary-path-section">
      <div class="container">
        <div class="section-header" style="text-align:center;">
          <p class="section-label">{{ journeyLabel }}</p>
          <h2 v-if="auth.isLoggedIn && userJourneyStops.length">
            My <em>Itinerary</em> Route
          </h2>
          <h2 v-else>Follow the <em>Path</em></h2>
          <!-- Itinerary selector for logged-in users -->
          <div v-if="auth.isLoggedIn && userItineraries.length > 1" class="itin-selector">
            <label>Viewing:</label>
            <select v-model="selectedItinIdx" class="itin-select">
              <option :value="-1">All Trips Combined</option>
              <option v-for="(it, idx) in userItineraries" :key="it.id" :value="idx">
                {{ it.name }}
              </option>
            </select>
          </div>
        </div>

        <div class="journey-map" :key="journeyKey">
          <!-- SVG winding dashed path with stops -->
          <svg class="journey-svg" viewBox="0 0 1000 320" fill="none" xmlns="http://www.w3.org/2000/svg">
            <!-- Dashed path -->
            <path d="M50,160 C150,50 250,50 350,160 C450,270 550,270 650,160 C750,50 850,50 950,160"
              class="journey-line" />
            <!-- Animated airplane along path -->
            <text font-size="24" class="journey-plane">
              <animateMotion dur="8s" repeatCount="indefinite"
                path="M50,160 C150,50 250,50 350,160 C450,270 550,270 650,160 C750,50 850,50 950,160" />
              ✈
            </text>
          </svg>

          <!-- Dynamic destination stops -->
          <div class="journey-stops">
            <div v-for="(stop, i) in journeyStops" :key="i" class="journey-stop" :style="stop.style">
              <div class="stop-marker" :class="stop.markerClass">
                <span class="stop-num">{{ i + 1 }}</span>
              </div>
              <div class="stop-label">
                <strong>{{ stop.name }}</strong>
                <span>{{ stop.country }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Link to itineraries -->
        <div v-if="auth.isLoggedIn" style="text-align:center;margin-top:1.5rem;">
          <RouterLink to="/itineraries" class="btn btn-ghost">
            {{ userJourneyStops.length ? '✏ Manage My Itineraries' : '✈ Create Your First Itinerary' }}
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- ===== BOTTOM CTA ===== -->
    <section class="cta-section" v-if="!auth.isLoggedIn">
      <div class="container">
        <div class="cta-card">
          <img src="https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?w=1200&q=80" alt="Travel"
            class="cta-bg" />
          <div class="cta-overlay"></div>
          <div class="cta-content">
            <p class="section-label-light">✦ START YOUR JOURNEY</p>
            <h2 class="text-white" style="font-size:clamp(1.5rem,4vw,2.8rem);">
              Experience the Magic of<br /><em>Travel</em> Today
            </h2>
            <p style="color:rgba(255,255,255,0.75);max-width:400px;margin:1rem auto 0;">
              Discover breathtaking views, unforgettable moments, and
              the world's most iconic destinations.
            </p>
            <RouterLink to="/register" class="btn btn-light btn-arrow btn-lg mt-3">
              Reserve Your Tour Now
              <span class="btn-arrow-icon">↗</span>
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== FOOTER ===== -->
    <footer class="site-footer">
      <div class="container footer-inner">
        <div class="footer-brand">
          <div style="display:flex;align-items:center;gap:0.5rem;margin-bottom:0.75rem;">
            <span style="font-size:1.2rem;">✈</span>
            <strong style="font-family:var(--font-display);font-size:1.2rem;color:var(--text-dark);">Wanderlust</strong>
          </div>
          <p style="font-size:0.85rem;max-width:260px;">
            We provide curated travel experiences with expert guides,
            safe transportation, and unforgettable moments.
          </p>
        </div>
        <div class="footer-links">
          <h4>Quick Links</h4>
          <RouterLink to="/">Home</RouterLink>
          <RouterLink to="/login">Login</RouterLink>
          <RouterLink to="/register">Register</RouterLink>
        </div>
        <div class="footer-links">
          <h4>Contact</h4>
          <span>aniketghage@gmail.com</span>
          <span>+91 8830436673</span>
          <span>India</span>
        </div>
      </div>
      <div class="footer-bottom container">
        <span>© 2026 Wanderlust. All rights reserved.</span>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import api from '../api'

const auth = useAuthStore()
const destSection = ref(null)
const destinations = ref([])
const loading = ref(false)
const error = ref('')
const pagination = reactive({ page: 1, totalPages: 1, total: 0 })
const filters = reactive({ name: '', country: '', minRating: '' })
const stats = reactive({ destinations: 0 })


const userItineraries = ref([])
const selectedItinIdx = ref(-1) // -1 = all combined

const markerColors = ['', 'stop-marker--warm', 'stop-marker--gold', 'stop-marker--coral', 'stop-marker--lime']

// Pre-computed stop positions along the winding path
const stopPositions = [
  { left: '5%', top: '42%' },
  { left: '18%', top: '12%' },
  { left: '35%', top: '42%' },
  { left: '50%', top: '75%' },
  { left: '65%', top: '42%' },
  { left: '78%', top: '12%' },
  { left: '92%', top: '42%' },
]

const defaultStops = [
  { name: 'Rome', country: 'Italy 🇮🇹' },
  { name: 'Tokyo', country: 'Japan 🇯🇵' },
  { name: 'Bali', country: 'Indonesia 🇮🇩' },
  { name: 'Paris', country: 'France 🇫🇷' },
  { name: 'Cape Town', country: 'South Africa 🇿🇦' },
]

// Get user's destinations from itineraries
const userJourneyStops = computed(() => {
  if (!userItineraries.value.length) return []

  let dests = []
  if (selectedItinIdx.value === -1) {
    // All combined — unique destinations across all itineraries
    const seen = new Set()
    for (const it of userItineraries.value) {
      for (const d of (it.destinationDetails || it.destinations || [])) {
        const key = d.id || d.name
        if (!seen.has(key)) {
          seen.add(key)
          dests.push({ name: d.name || d, country: d.country || '' })
        }
      }
    }
  } else {
    const it = userItineraries.value[selectedItinIdx.value]
    const itDests = it?.destinationDetails || it?.destinations || []
    dests = itDests.map(d => ({
      name: d.name || d,
      country: d.country || ''
    }))
  }
  return dests.slice(0, 7) // max 7 stops on path
})

// Final stops to display (user's or default)
const journeyStops = computed(() => {
  const rawStops = userJourneyStops.value.length > 0 ? userJourneyStops.value : defaultStops
  return rawStops.map((s, i) => ({
    ...s,
    style: stopPositions[i % stopPositions.length],
    markerClass: markerColors[i % markerColors.length],
  }))
})

// Key to force re-render SVG when itinerary changes
const journeyKey = computed(() => `journey-${selectedItinIdx.value}-${userJourneyStops.value.length}`)

const journeyLabel = computed(() => {
  return auth.isLoggedIn && userJourneyStops.value.length ? '✦ MY JOURNEY' : '✦ EXPLORE THE WORLD'
})

async function fetchUserItineraries() {
  if (!auth.isLoggedIn) return
  try {
    const res = await api.get('/itineraries')
    const data = res.data
    userItineraries.value = data.itineraries || data || []
  } catch (e) {
    // Silently fail — just show default stops
    userItineraries.value = []
  }
}

const previewImages = [
  { src: 'https://images.unsplash.com/photo-1552832230-c0197dd311b5?w=400&q=80', label: 'Rome' },
  { src: 'https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=400&q=80', label: 'Tokyo' },
  { src: 'https://images.unsplash.com/photo-1537996194471-e657df975ab4?w=400&q=80', label: 'Bali' },
]

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { pagination.page = 1; fetchDestinations() }, 400)
}

async function fetchDestinations() {
  loading.value = true; error.value = ''
  try {
    const params = { page: pagination.page, limit: 9 }
    if (filters.name) params.name = filters.name
    if (filters.country) params.country = filters.country
    if (filters.minRating) params.minRating = filters.minRating
    const res = await api.get('/destinations', { params })
    destinations.value = res.data.destinations
    pagination.page = res.data.pagination.page
    pagination.totalPages = res.data.pagination.totalPages
    pagination.total = res.data.pagination.total
    stats.destinations = res.data.pagination.total
  } catch (e) {
    error.value = e.response?.data?.error || 'Failed to load destinations'
  } finally { loading.value = false }
}

function goToPage(p) { pagination.page = p; fetchDestinations() }
function clearFilters() {
  filters.name = ''; filters.country = ''; filters.minRating = ''
  pagination.page = 1; fetchDestinations()
}
function starStr(r) {
  const n = Math.round(r || 0)
  return '★'.repeat(n) + '☆'.repeat(5 - n)
}
function fallbackImg(name) {
  return `https://source.unsplash.com/600x400/?${encodeURIComponent(name + ' travel')}`
}
function onImgError(e) {
  e.target.src = 'https://images.unsplash.com/photo-1488085061387-422e29b40080?w=600&q=80'
}
function scrollToDestinations() {
  destSection.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

onMounted(() => {
  fetchDestinations()
  fetchUserItineraries()
})
</script>

<style scoped>
/* ====== HERO ====== */
.hero {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  overflow: hidden;
}

.hero-photo {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
}

/* Gentle wave-float animation for hero text */
.wave-line {
  display: inline-block;
  animation: waveFloat 4s ease-in-out infinite;
  animation-delay: var(--wave-delay, 0s);
}

@keyframes waveFloat {

  0%,
  100% {
    transform: translateY(0);
  }

  25% {
    transform: translateY(-6px);
  }

  50% {
    transform: translateY(0);
  }

  75% {
    transform: translateY(4px);
  }
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom,
      rgba(0, 0, 0, 0.15) 0%,
      rgba(0, 0, 0, 0.3) 50%,
      rgba(0, 0, 0, 0.55) 100%);
  z-index: 1;
}

.hero-inner {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 100vh;
  padding-top: calc(var(--nav-h) + 4rem);
  padding-bottom: 3rem;
}

.hero-text {
  max-width: 650px;
}

.hero-sub {
  color: rgba(255, 255, 255, 0.75);
  font-size: 1.05rem;
  line-height: 1.7;
  margin: 1.25rem 0 2rem;
  max-width: 480px;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

/* Bottom area: social proof + floating cards */
.hero-bottom {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 2rem;
  flex-wrap: wrap;
}

.hero-social-proof {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: var(--r-full);
  padding: 0.5rem 1.25rem 0.5rem 0.5rem;
}

.hero-social-proof .avatar-stack-count {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.82rem;
}

/* Floating preview cards */
.hero-preview-cards {
  display: flex;
  gap: 0.75rem;
}

.hero-float-card {
  width: 140px;
  border-radius: var(--r-md);
  overflow: hidden;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
  position: relative;
  animation: floatUp 0.6s ease both;
  transition: transform 0.3s ease;
}

.hero-float-card:hover {
  transform: translateY(-6px);
}

.hero-float-card img {
  width: 100%;
  height: 180px;
  object-fit: cover;
  display: block;
}

.hero-float-label {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 0.5rem;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
  color: #fff;
  font-size: 0.78rem;
  font-weight: 600;
  text-align: center;
}

/* Wave divider */
.wave-divider {
  position: relative;
  z-index: 3;
  margin-top: -100px;
  margin-bottom: -1px;
  line-height: 0;
}

.wave-divider svg {
  display: block;
  width: 100%;
  height: 150px;
}

@keyframes floatUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ====== JOURNEY MAP (Dashed itinerary path) ====== */
.itinerary-path-section {
  padding: 3rem 0 4rem;
  background: var(--bg-cream);
  overflow: hidden;
}

.journey-map {
  position: relative;
  margin-top: 1.5rem;
}

.journey-svg {
  width: 100%;
  height: auto;
  display: block;
}

.journey-line {
  stroke: var(--green-deep);
  stroke-width: 2.5;
  stroke-dasharray: 12 8;
  fill: none;
  opacity: 0.45;
  stroke-linecap: round;
  animation: drawPath 3s ease forwards;
}

@keyframes drawPath {
  from {
    stroke-dashoffset: 2000;
  }

  to {
    stroke-dashoffset: 0;
  }
}

.journey-plane {
  fill: var(--orange-warm);
}

/* Stops container */
.journey-stops {
  position: absolute;
  inset: 0;
}

.journey-stop {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.4rem;
  transform: translateX(-50%);
  animation: stopPop 0.5s ease both;
}

.journey-stop:nth-child(1) {
  animation-delay: 0.3s;
}

.journey-stop:nth-child(2) {
  animation-delay: 0.6s;
}

.journey-stop:nth-child(3) {
  animation-delay: 0.9s;
}

.journey-stop:nth-child(4) {
  animation-delay: 1.2s;
}

.journey-stop:nth-child(5) {
  animation-delay: 1.5s;
}

@keyframes stopPop {
  from {
    opacity: 0;
    transform: translateX(-50%) scale(0.5);
  }

  to {
    opacity: 1;
    transform: translateX(-50%) scale(1);
  }
}

.stop-marker {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--green-deep);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(45, 74, 62, 0.35), 0 0 0 4px rgba(45, 74, 62, 0.1);
  transition: transform 0.3s ease;
  cursor: default;
}

.stop-marker:hover {
  transform: scale(1.15);
}

.stop-marker--warm {
  background: var(--orange-warm);
  box-shadow: 0 4px 16px rgba(232, 118, 75, 0.35), 0 0 0 4px rgba(232, 118, 75, 0.1);
}

.stop-marker--gold {
  background: var(--gold);
  color: var(--text-dark);
  box-shadow: 0 4px 16px rgba(233, 196, 106, 0.35), 0 0 0 4px rgba(233, 196, 106, 0.1);
}

.stop-marker--coral {
  background: var(--coral);
  box-shadow: 0 4px 16px rgba(231, 111, 81, 0.35), 0 0 0 4px rgba(231, 111, 81, 0.1);
}

.stop-marker--lime {
  background: var(--green-accent);
  box-shadow: 0 4px 16px rgba(107, 186, 98, 0.35), 0 0 0 4px rgba(107, 186, 98, 0.1);
}

.stop-num {
  font-size: 0.85rem;
  font-weight: 800;
  font-family: var(--font-body);
}

.stop-label {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
}

.stop-label strong {
  font-size: 0.88rem;
  color: var(--text-dark);
}

.stop-label span {
  font-size: 0.72rem;
  color: var(--text-muted);
}

/* Itinerary selector dropdown */
.itin-selector {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.itin-selector label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.itin-select {
  padding: 0.4rem 1rem 0.4rem 0.75rem;
  border: 1.5px solid rgba(0, 0, 0, 0.1);
  border-radius: var(--r-full);
  font-size: 0.85rem;
  font-family: var(--font-body);
  color: var(--text-dark);
  background: #fff;
  cursor: pointer;
  transition: border-color 0.2s;
}

.itin-select:focus {
  outline: none;
  border-color: var(--green-deep);
  box-shadow: 0 0 0 3px rgba(45, 74, 62, 0.1);
}

@media (max-width: 768px) {
  .journey-stops {
    display: none;
  }

  .journey-svg {
    height: 120px;
  }

  .itinerary-path-section {
    padding: 2rem 0;
  }
}

/* ====== DESTINATIONS SECTION ====== */
.dest-section {
  padding: 4rem 0;
}

.section-header {
  margin-bottom: 2rem;
}

.section-header em {
  font-family: var(--font-display);
  color: var(--green-deep);
}

/* ====== CTA SECTION ====== */
.cta-section {
  padding: 2rem 0 5rem;
}

.cta-card {
  position: relative;
  border-radius: var(--r-xl);
  overflow: hidden;
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cta-bg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cta-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
}

.cta-content {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 3rem;
}

/* ====== FOOTER ====== */
.site-footer {
  background: var(--bg-white);
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  padding: 3rem 0 0;
  position: relative;
  z-index: 10;
}

.footer-inner {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 2rem;
  padding-bottom: 2rem;
}

.footer-links {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.footer-links h4 {
  font-family: var(--font-body);
  font-size: 0.82rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-dark);
  margin-bottom: 0.5rem;
}

.footer-links a,
.footer-links span {
  font-size: 0.88rem;
  color: var(--text-muted);
  text-decoration: none;
  transition: color 0.2s;
  cursor: pointer;
  position: relative;
  z-index: 1;
}

.footer-links a:hover {
  color: var(--green-deep);
  text-decoration: underline;
}

.footer-bottom {
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  padding: 1.25rem 0;
  font-size: 0.78rem;
  color: var(--text-muted);
}

@media (max-width: 900px) {
  .hero-preview-cards {
    display: none;
  }

  .hero-inner {
    padding-top: calc(var(--nav-h) + 2rem);
  }

  .footer-inner {
    grid-template-columns: 1fr;
  }
}
</style>
