<template>
  <div class="page">
    <!-- Header -->
    <div class="itin-header">
      <div>
        <p class="section-label">✦ YOUR ADVENTURES</p>
        <h1>My <em>Itineraries</em></h1>
        <p class="itin-sub">Plan, organise, and relive your journeys.</p>
      </div>
      <button class="btn btn-primary btn-arrow" @click="openCreate">
        New Itinerary
        <span class="btn-arrow-icon">↗</span>
      </button>
    </div>

    <div v-if="loading" class="spinner"></div>
    <div v-else-if="error" class="alert alert-error">⚠ {{ error }}</div>

    <!-- Empty state -->
    <div v-else-if="itineraries.length === 0" class="empty-state">
      <div class="empty-icon">🗺️</div>
      <h3>No itineraries yet</h3>
      <p>Create your first trip plan and start exploring!</p>
      <button class="btn btn-primary btn-arrow mt-3" @click="openCreate">
        Create My First Trip
        <span class="btn-arrow-icon">↗</span>
      </button>
    </div>

    <!-- ===== ITINERARY SECTIONS ===== -->
    <div v-else class="itin-list">
      <section v-for="(itin, idx) in itineraries" :key="itin.id" class="itin-section"
        :style="{ animationDelay: `${idx * 0.12}s` }">

        <!-- Section Header -->
        <div class="section-top">
          <div class="section-top-left">
            <span class="trip-badge">Trip #{{ idx + 1 }}</span>
            <h2 class="trip-name">{{ itin.name }}</h2>
            <div class="trip-dates">
              <span>{{ formatDate(itin.startDate) }}</span>
              <span class="trip-arrow">→</span>
              <span>{{ formatDate(itin.endDate) }}</span>
              <span class="trip-dur">{{ tripDuration(itin.startDate, itin.endDate) }}</span>
            </div>
          </div>
          <div class="section-top-right">
            <button class="btn btn-ghost btn-sm" @click="openEdit(itin)">✏ Edit</button>
            <button class="btn btn-danger btn-sm" @click="deleteItinerary(itin.id)">🗑</button>
          </div>
        </div>

        <!-- ===== JOURNEY MAP — Horizontal winding path ===== -->
        <div v-if="getDests(itin).length" class="journey-map" :key="`map-${itin.id}-${getDests(itin).length}`">
          <div class="map-hint">↔ Drag stops to reorder your route</div>

          <div class="map-wrap">
            <!-- SVG winding dashed path -->
            <svg class="map-svg" viewBox="0 0 1000 320" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M50,160 C150,50 250,50 350,160 C450,270 550,270 650,160 C750,50 850,50 950,160"
                class="map-path" />
              <!-- Animated airplane -->
              <text font-size="22" class="map-plane">
                <animateMotion dur="8s" repeatCount="indefinite"
                  path="M50,160 C150,50 250,50 350,160 C450,270 550,270 650,160 C750,50 850,50 950,160" />
                ✈
              </text>
            </svg>

            <!-- Destination stops positioned along the path -->
            <div class="map-stops">
              <div v-for="(dest, di) in getDests(itin)" :key="dest.id || di" class="map-stop"
                :style="stopPosition(di, getDests(itin).length)"
                :class="{ 'map-stop--dragging': dragState.itinId === itin.id && dragState.fromIdx === di }"
                draggable="true" @dragstart="onDragStart(itin, di, $event)" @dragover.prevent @drop="onDrop(itin, di)"
                @dragend="onDragEnd">
                <!-- Numbered marker -->
                <div class="stop-marker" :class="markerColor(di)">
                  <span>{{ di + 1 }}</span>
                </div>
                <!-- Info card below/above marker -->
                <div class="stop-info" :class="di % 2 === 0 ? 'stop-info--below' : 'stop-info--above'">
                  <strong>{{ dest.name }}</strong>
                  <span class="stop-country" v-if="dest.country">📍 {{ dest.country }}</span>
                  <span class="stop-rating" v-if="dest.averageRating">★ {{ dest.averageRating?.toFixed(1) }}</span>
                  <div class="stop-reorder">
                    <!-- <button v-if="di > 0" class="reorder-btn" @click.stop="moveStop(itin, di, di - 1)">◀</button>
                    <button v-if="di < getDests(itin).length - 1" class="reorder-btn"
                      @click.stop="moveStop(itin, di, di + 1)">▶</button> -->
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Save order -->
          <div v-if="changedItinIds.has(itin.id)" class="save-bar">
            <button class="btn btn-primary btn-sm" @click="saveOrder(itin)" :disabled="savingOrder === itin.id">
              {{ savingOrder === itin.id ? 'Saving...' : '💾 Save New Order' }}
            </button>
            <button class="btn btn-ghost btn-sm" @click="revertOrder(itin)">↩ Revert</button>
          </div>
        </div>

        <!-- No destinations -->
        <div v-else class="no-dest-msg">
          <span>No destinations added yet — </span>
          <button class="btn btn-ghost btn-sm" @click="openEdit(itin)">+ Add destinations</button>
        </div>
      </section>
    </div>

    <!-- ===== CREATE/EDIT MODAL ===== -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h2>{{ editingId ? '✏ Edit Itinerary' : '✈ New Itinerary' }}</h2>
        <div v-if="modalError" class="alert alert-error">⚠ {{ modalError }}</div>
        <form @submit.prevent="saveItinerary">
          <div class="form-group">
            <label>Trip Name</label>
            <input v-model="form.name" required placeholder="e.g. European Summer 2026" />
          </div>
          <div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem;">
            <div class="form-group">
              <label>Start Date</label>
              <input v-model="form.startDate" type="date" required />
            </div>
            <div class="form-group">
              <label>End Date</label>
              <input v-model="form.endDate" type="date" required />
            </div>
          </div>
          <div class="form-group">
            <label>Destinations (click in order)</label>
            <div class="dest-selector">
              <div v-for="d in allDestinations" :key="d.id" class="dest-option"
                :class="{ selected: form.destinationIds.includes(d.id) }" @click="toggleDest(d.id)">
                <span class="dest-option-num">
                  {{ form.destinationIds.includes(d.id) ? form.destinationIds.indexOf(d.id) + 1 : '+' }}
                </span>
                <span>{{ d.name }}</span>
                <span class="dest-option-country">{{ d.country }}</span>
              </div>
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-ghost" @click="closeModal">Cancel</button>
            <button type="submit" class="btn btn-primary btn-arrow" :disabled="saving">
              {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Create Trip') }}
              <span class="btn-arrow-icon">↗</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../api'

const itineraries = ref([])
const allDestinations = ref([])
const loading = ref(true)
const error = ref('')
const showModal = ref(false)
const editingId = ref(null)
const saving = ref(false)
const savingOrder = ref(null)
const modalError = ref('')
const form = reactive({ name: '', startDate: '', endDate: '', destinationIds: [] })

const changedItinIds = ref(new Set())
const originalOrders = ref({})
const dragState = reactive({ itinId: null, fromIdx: null })


const pathPoints = [
  { x: 5, y: 42 },
  { x: 15, y: 10 },
  { x: 25, y: 10 },
  { x: 35, y: 42 },
  { x: 45, y: 75 },
  { x: 55, y: 75 },
  { x: 65, y: 42 },
  { x: 75, y: 10 },
  { x: 85, y: 10 },
  { x: 95, y: 42 },
]

function stopPosition(idx, total) {
  // Distribute stops evenly across path points
  if (total <= 1) return { left: '50%', top: '42%' }
  const spacing = (pathPoints.length - 1) / (total - 1)
  const pi = Math.round(idx * spacing)
  const pt = pathPoints[Math.min(pi, pathPoints.length - 1)]
  return { left: `${pt.x}%`, top: `${pt.y}%` }
}

const colors = ['', 'marker--warm', 'marker--gold', 'marker--coral', 'marker--lime']
function markerColor(i) { return colors[i % colors.length] }

function getDests(itin) {
  const details = itin.destinationDetails || []
  const idOrder = itin.destinations || []
  if (!details.length || !idOrder.length) return details

  // Sort destinationDetails to match the order in the destinations ID array
  const orderMap = {}
  idOrder.forEach((id, i) => { orderMap[id] = i })
  const sorted = [...details].sort((a, b) => {
    const ai = orderMap[a.id] ?? 999
    const bi = orderMap[b.id] ?? 999
    return ai - bi
  })
  // Replace the array in-place so mutations (reorder) work on it
  itin.destinationDetails = sorted
  return sorted
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
function tripDuration(start, end) {
  if (!start || !end) return ''
  const days = Math.round((new Date(end) - new Date(start)) / 86400000)
  return days > 0 ? `${days}d` : ''
}


function toDateStr(d) {
  // Convert ISO timestamp or Date to YYYY-MM-DD string
  if (!d) return ''
  if (typeof d === 'string') return d.substring(0, 10)
  return new Date(d).toISOString().substring(0, 10)
}


function moveStop(itin, from, to) {
  // Save original BEFORE mutating
  if (!originalOrders.value[itin.id]) {
    originalOrders.value[itin.id] = getDests(itin).map(d => d.id)
  }
  const dests = getDests(itin)
  const item = dests.splice(from, 1)[0]
  dests.splice(to, 0, item)
  changedItinIds.value.add(itin.id)
  changedItinIds.value = new Set(changedItinIds.value)
}
function onDragStart(itin, idx, e) {
  dragState.itinId = itin.id; dragState.fromIdx = idx
  e.dataTransfer.effectAllowed = 'move'
}
function onDrop(itin, toIdx) {
  if (dragState.itinId !== itin.id || dragState.fromIdx === null) return
  if (dragState.fromIdx !== toIdx) moveStop(itin, dragState.fromIdx, toIdx)
  onDragEnd()
}
function onDragEnd() { dragState.itinId = null; dragState.fromIdx = null }

function revertOrder(itin) {
  const orig = originalOrders.value[itin.id]
  if (!orig) return
  getDests(itin).sort((a, b) => orig.indexOf(a.id) - orig.indexOf(b.id))
  changedItinIds.value.delete(itin.id)
  changedItinIds.value = new Set(changedItinIds.value)
  delete originalOrders.value[itin.id]
}
async function saveOrder(itin) {
  savingOrder.value = itin.id
  try {
    const orderedIds = getDests(itin).map(d => d.id)
    console.log('Saving order for', itin.id, ':', orderedIds)
    await api.put(`/itineraries/${itin.id}`, {
      name: itin.name,
      startDate: toDateStr(itin.startDate),
      endDate: toDateStr(itin.endDate),
      destinations: orderedIds
    })
    changedItinIds.value.delete(itin.id)
    changedItinIds.value = new Set(changedItinIds.value)
    delete originalOrders.value[itin.id]
    // Re-fetch to confirm persistence
    await fetchItineraries()
  } catch (e) {
    console.error('Save order error:', e.response?.data || e)
    error.value = e.response?.data?.error || 'Failed to save order'
  }
  finally { savingOrder.value = null }
}


function toggleDest(id) {
  const idx = form.destinationIds.indexOf(id)
  if (idx === -1) form.destinationIds.push(id)
  else form.destinationIds.splice(idx, 1)
}
function openCreate() {
  editingId.value = null
  form.name = ''; form.startDate = ''; form.endDate = ''; form.destinationIds = []
  modalError.value = ''; showModal.value = true
}
function openEdit(itin) {
  editingId.value = itin.id
  form.name = itin.name
  form.startDate = itin.startDate?.slice(0, 10) || ''
  form.endDate = itin.endDate?.slice(0, 10) || ''
  form.destinationIds = getDests(itin).map(d => d.id || d).filter(Boolean)
  modalError.value = ''; showModal.value = true
}
function closeModal() { showModal.value = false }

async function fetchItineraries() {
  loading.value = true; error.value = ''
  try {
    const [itinRes, destRes] = await Promise.all([
      api.get('/itineraries'),
      api.get('/destinations', { params: { limit: 100 } })
    ])
    const data = itinRes.data
    itineraries.value = data.itineraries || data || []
    allDestinations.value = destRes.data.destinations || destRes.data || []
  } catch (e) { error.value = e.response?.data?.error || 'Failed to load itineraries' }
  finally { loading.value = false }
}
async function saveItinerary() {
  saving.value = true; modalError.value = ''
  try {
    const payload = { name: form.name, startDate: form.startDate, endDate: form.endDate, destinations: form.destinationIds }
    if (editingId.value) await api.put(`/itineraries/${editingId.value}`, payload)
    else await api.post('/itineraries', payload)
    closeModal(); await fetchItineraries()
  } catch (e) { modalError.value = e.response?.data?.error || 'Failed to save' }
  finally { saving.value = false }
}
async function deleteItinerary(id) {
  if (!confirm('Delete this itinerary?')) return
  try { await api.delete(`/itineraries/${id}`); await fetchItineraries() }
  catch (e) { error.value = e.response?.data?.error || 'Failed to delete' }
}

onMounted(fetchItineraries)
</script>

<style scoped>
/* ─── HEADER ─── */
.itin-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 2.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.itin-header em {
  font-family: var(--font-display);
  color: var(--green-deep);
}

.itin-sub {
  color: var(--text-muted);
  margin-top: 0.4rem;
}

.itin-list {
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

/* ─── ITINERARY SECTION ─── */
.itin-section {
  background: #fff;
  border-radius: var(--r-xl);
  padding: 2rem 2.5rem;
  box-shadow: var(--shadow-card);
  border: 1px solid rgba(0, 0, 0, 0.04);
  animation: sectionIn 0.45s ease both;
}

@keyframes sectionIn {
  from {
    opacity: 0;
    transform: translateY(16px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.section-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.section-top-right {
  display: flex;
  gap: 0.5rem;
}

.trip-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: var(--bg-cream);
  border-radius: var(--r-full);
  font-size: 0.68rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}

.trip-name {
  font-size: 1.5rem;
  margin: 0 0 0.5rem;
  color: var(--text-dark);
}

.trip-dates {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.88rem;
  color: var(--text-body);
}

.trip-arrow {
  color: var(--orange-warm);
  font-weight: 700;
}

.trip-dur {
  background: var(--green-deep);
  color: #fff;
  padding: 0.2rem 0.6rem;
  border-radius: var(--r-full);
  font-size: 0.72rem;
  font-weight: 700;
  margin-left: 0.25rem;
}

/* ─── JOURNEY MAP ─── */
.journey-map {
  margin-top: 0.5rem;
}

.map-hint {
  font-size: 0.72rem;
  color: var(--text-muted);
  font-style: italic;
  text-align: right;
  margin-bottom: 0.25rem;
}

.map-wrap {
  position: relative;
  padding: 2rem 0;
}

.map-svg {
  width: 100%;
  height: auto;
  display: block;
}

.map-path {
  stroke: var(--green-deep);
  stroke-width: 2.5;
  stroke-dasharray: 12 8;
  fill: none;
  opacity: 0.35;
  stroke-linecap: round;
}

.map-plane {
  fill: var(--orange-warm);
}

/* Stop markers positioned absolutely over the SVG */
.map-stops {
  position: absolute;
  inset: 0;
}

.map-stop {
  position: absolute;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
  cursor: grab;
  transition: opacity 0.2s;
  z-index: 2;
}

.map-stop:active {
  cursor: grabbing;
}

.map-stop--dragging {
  opacity: 0.3;
}

/* Numbered circle marker */
.stop-marker {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: var(--green-deep);
  color: #fff;
  font-size: 0.82rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 14px rgba(45, 74, 62, 0.35), 0 0 0 4px rgba(45, 74, 62, 0.1);
  transition: transform 0.25s, box-shadow 0.25s;
  z-index: 2;
}

.map-stop:hover .stop-marker {
  transform: scale(1.15);
  box-shadow: 0 6px 20px rgba(45, 74, 62, 0.4), 0 0 0 5px rgba(45, 74, 62, 0.15);
}

.marker--warm {
  background: var(--orange-warm);
  box-shadow: 0 4px 14px rgba(232, 118, 75, 0.35), 0 0 0 4px rgba(232, 118, 75, 0.1);
}

.marker--gold {
  background: var(--gold);
  color: var(--text-dark);
  box-shadow: 0 4px 14px rgba(233, 196, 106, 0.35), 0 0 0 4px rgba(233, 196, 106, 0.1);
}

.marker--coral {
  background: var(--coral);
  box-shadow: 0 4px 14px rgba(231, 111, 81, 0.35), 0 0 0 4px rgba(231, 111, 81, 0.1);
}

.marker--lime {
  background: var(--green-accent);
  box-shadow: 0 4px 14px rgba(107, 186, 98, 0.35), 0 0 0 4px rgba(107, 186, 98, 0.1);
}

/* Info card that pops on hover or always shows */
.stop-info {
  background: #fff;
  border-radius: var(--r-md);
  padding: 0.5rem 0.75rem;
  box-shadow: var(--shadow-sm);
  border: 1px solid rgba(0, 0, 0, 0.06);
  text-align: center;
  min-width: 100px;
  max-width: 140px;
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  transition: box-shadow 0.2s;
}

.map-stop:hover .stop-info {
  box-shadow: var(--shadow-md);
}

.stop-info--above {
  order: -1;
}

.stop-info strong {
  font-size: 0.85rem;
  color: var(--text-dark);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stop-country {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.stop-rating {
  font-size: 0.7rem;
  color: var(--gold);
  font-weight: 600;
}

.stop-reorder {
  display: flex;
  gap: 0.25rem;
  justify-content: center;
  margin-top: 0.2rem;
}

.reorder-btn {
  background: var(--bg-cream);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 4px;
  font-size: 0.6rem;
  padding: 2px 8px;
  cursor: pointer;
  color: var(--text-muted);
  transition: all 0.15s;
  line-height: 1;
}

.reorder-btn:hover {
  background: var(--green-deep);
  color: #fff;
  border-color: var(--green-deep);
}

/* Save bar */
.save-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 1rem;
  padding: 0.75rem 1rem;
  background: rgba(45, 74, 62, 0.04);
  border-radius: var(--r-md);
  border: 1px dashed rgba(45, 74, 62, 0.2);
}

.no-dest-msg {
  padding: 1.5rem 0 0;
  color: var(--text-muted);
  font-size: 0.88rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* ─── MODAL DEST SELECTOR ─── */
.dest-selector {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 0.5rem;
  max-height: 220px;
  overflow-y: auto;
  padding: 0.5rem;
  background: var(--bg-cream);
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: var(--r-md);
}

.dest-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border-radius: var(--r-md);
  cursor: pointer;
  font-size: 0.85rem;
  color: var(--text-body);
  border: 1px solid transparent;
  transition: all 0.2s ease;
}

.dest-option:hover {
  background: #fff;
}

.dest-option.selected {
  background: rgba(45, 74, 62, 0.06);
  border-color: rgba(45, 74, 62, 0.2);
}

.dest-option-num {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #fff;
  border: 1.5px solid rgba(0, 0, 0, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  flex-shrink: 0;
  color: var(--green-deep);
  font-weight: 700;
}

.dest-option.selected .dest-option-num {
  background: var(--green-deep);
  border-color: var(--green-deep);
  color: #fff;
}

.dest-option-country {
  color: var(--text-muted);
  font-size: 0.72rem;
  margin-left: auto;
}

@media (max-width: 768px) {
  .itin-section {
    padding: 1.25rem;
  }

  .map-stops {
    position: static;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: center;
    margin-top: 1rem;
  }

  .map-stop {
    position: static !important;
    transform: none;
  }

  .stop-info--above {
    order: 0;
  }

  .map-svg {
    height: 100px;
  }
}
</style>
