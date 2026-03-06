<template>
  <div class="admin-layout">
    <!-- Sidebar -->
    <aside class="admin-sidebar">
      <div class="admin-brand">
        <span class="admin-brand-icon">⚙</span>
        <div>
          <div class="admin-brand-title">Admin Panel</div>
          <div class="admin-brand-sub">Wanderlust</div>
        </div>
      </div>

      <nav class="admin-nav">
        <button class="admin-nav-item" :class="{ active: activeTab === 'destinations' }"
          @click="activeTab = 'destinations'">
          <span>🌍</span> Destinations
          <span class="admin-nav-count">{{ destinations.length }}</span>
        </button>
        <button class="admin-nav-item" :class="{ active: activeTab === 'reviews' }" @click="activeTab = 'reviews'">
          <span>💬</span> Reviews
          <span class="admin-nav-count">{{ reviews.length }}</span>
        </button>
      </nav>

      <div class="admin-sidebar-footer">
        <RouterLink to="/" class="btn btn-ghost btn-sm" style="width:100%;justify-content:center;">
          ← Back to Site
        </RouterLink>
      </div>
    </aside>

    <!-- Main -->
    <main class="admin-main">
      <!-- DESTINATIONS TAB -->
      <div v-if="activeTab === 'destinations'">
        <div class="admin-page-header">
          <div>
            <p class="section-label">✦ MANAGE CONTENT</p>
            <h1>Destinations</h1>
          </div>
          <button class="btn btn-primary btn-arrow" @click="openCreate">
            Add Destination
            <span class="btn-arrow-icon">↗</span>
          </button>
        </div>

        <div v-if="destLoading" class="spinner"></div>
        <div v-else-if="destError" class="alert alert-error">⚠ {{ destError }}</div>

        <div v-else class="admin-table-card">
          <table class="data-table">
            <thead>
              <tr>
                <th>Destination</th>
                <th>Country</th>
                <th>Rating</th>
                <th>Reviews</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="destinations.length === 0">
                <td colspan="5" style="text-align:center;color:var(--text-muted);padding:2rem;">No destinations yet.
                </td>
              </tr>
              <tr v-for="d in destinations" :key="d.id">
                <td>
                  <div class="dest-table-name">
                    <img :src="d.imageUrl" :alt="d.name" class="dest-table-thumb" @error="onImgError" />
                    <span style="color:var(--text-dark);font-weight:500;">{{ d.name }}</span>
                  </div>
                </td>
                <td><span class="badge badge-green">{{ d.country }}</span></td>
                <td>
                  <span class="stars" style="font-size:0.85rem;">{{ '★'.repeat(Math.round(d.averageRating || 0))
                  }}</span>
                  <span style="color:var(--text-muted);font-size:0.8rem;margin-left:4px;">{{ d.averageRating?.toFixed(1)
                    || '—' }}</span>
                </td>
                <td><span class="badge badge-orange">{{ d.reviewCount || 0 }}</span></td>
                <td>
                  <div class="flex gap-1">
                    <button class="btn btn-ghost btn-sm" @click="openEdit(d)">✏ Edit</button>
                    <button class="btn btn-danger btn-sm" @click="deleteDest(d.id)">🗑</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- REVIEWS TAB -->
      <div v-if="activeTab === 'reviews'">
        <div class="admin-page-header">
          <div>
            <p class="section-label">✦ MODERATION</p>
            <h1>All Reviews</h1>
          </div>
        </div>

        <div v-if="reviewLoading" class="spinner"></div>
        <div v-else-if="reviewError" class="alert alert-error">⚠ {{ reviewError }}</div>

        <div v-else class="admin-table-card">
          <table class="data-table">
            <thead>
              <tr>
                <th>Reviewer</th>
                <th>Destination</th>
                <th>Rating</th>
                <th>Review</th>
                <th>Date</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="reviews.length === 0">
                <td colspan="6" style="text-align:center;color:var(--text-muted);padding:2rem;">No reviews yet.</td>
              </tr>
              <tr v-for="r in reviews" :key="r.id">
                <td>
                  <div class="flex align-center gap-1">
                    <div class="reviewer-avatar" style="width:28px;height:28px;font-size:0.72rem;">{{
                      r.userName?.[0]?.toUpperCase() }}</div>
                    <span style="color:var(--text-dark);">{{ r.userName }}</span>
                  </div>
                </td>
                <td><span class="badge badge-purple">{{ r.destinationName || '—' }}</span></td>
                <td><span class="stars" style="font-size:0.85rem;">{{ '★'.repeat(r.rating) }}</span></td>
                <td style="max-width:220px;">
                  <span style="color:var(--text-body);font-size:0.85rem;">{{ r.text?.slice(0, 80) }}{{ r.text?.length >
                    80 ? '…' : '' }}</span>
                </td>
                <td style="color:var(--text-muted);font-size:0.8rem;">{{ formatDate(r.createdAt) }}</td>
                <td>
                  <button class="btn btn-danger btn-sm" @click="deleteReview(r.id)">🗑 Remove</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>

    <!-- ===== DESTINATION MODAL ===== -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h2>{{ editingId ? '✏ Edit Destination' : '🌍 Add Destination' }}</h2>
        <div v-if="modalError" class="alert alert-error">⚠ {{ modalError }}</div>
        <form @submit.prevent="saveDest">
          <div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem;">
            <div class="form-group">
              <label>Name</label>
              <input v-model="form.name" required placeholder="Paris" />
            </div>
            <div class="form-group">
              <label>Country</label>
              <input v-model="form.country" required placeholder="France" />
            </div>
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea v-model="form.description" required placeholder="Describe this destination..."
              rows="3"></textarea>
          </div>
          <div class="form-group">
            <label>Image URL</label>
            <input v-model="form.imageUrl" type="url" placeholder="https://..." />
          </div>
          <div v-if="form.imageUrl" class="img-preview-wrap">
            <img :src="form.imageUrl" alt="Preview" class="img-preview" @error="e => e.target.style.display = 'none'" />
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-ghost" @click="closeModal">Cancel</button>
            <button type="submit" class="btn btn-primary btn-arrow" :disabled="saving">
              {{ saving ? 'Saving...' : (editingId ? 'Save Changes' : 'Add Destination') }}
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

const activeTab = ref('destinations')
const destinations = ref([])
const reviews = ref([])
const destLoading = ref(false)
const reviewLoading = ref(false)
const destError = ref('')
const reviewError = ref('')
const showModal = ref(false)
const editingId = ref(null)
const saving = ref(false)
const modalError = ref('')
const form = reactive({ name: '', country: '', description: '', imageUrl: '' })

function formatDate(d) { return d ? new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }) : '—' }
function onImgError(e) { e.target.src = 'https://images.unsplash.com/photo-1488085061387-422e29b40080?w=80&q=60' }

function openCreate() {
  editingId.value = null
  form.name = ''; form.country = ''; form.description = ''; form.imageUrl = ''
  modalError.value = ''; showModal.value = true
}
function openEdit(d) {
  editingId.value = d.id
  form.name = d.name; form.country = d.country; form.description = d.description; form.imageUrl = d.imageUrl || ''
  modalError.value = ''; showModal.value = true
}
function closeModal() { showModal.value = false }

async function fetchDestinations() {
  destLoading.value = true; destError.value = ''
  try {
    const res = await api.get('/destinations', { params: { limit: 100 } })
    destinations.value = res.data.destinations
  } catch (e) { destError.value = e.response?.data?.error || 'Failed to load' }
  finally { destLoading.value = false }
}
async function fetchReviews() {
  reviewLoading.value = true; reviewError.value = ''
  try {
    const res = await api.get('/reviews')
    const data = res.data
    const rawReviews = data.reviews || data || []
    // Enrich reviews with destination name from our destinations list
    reviews.value = rawReviews.map(r => {
      const dest = destinations.value.find(d => d.id === r.destinationId)
      return { ...r, destinationName: dest?.name || r.destinationId }
    })
  } catch (e) { reviewError.value = e.response?.data?.error || 'Failed to load reviews' }
  finally { reviewLoading.value = false }
}
async function saveDest() {
  saving.value = true; modalError.value = ''
  try {
    if (editingId.value) await api.put(`/destinations/${editingId.value}`, form)
    else await api.post('/destinations', form)
    closeModal(); await fetchDestinations()
  } catch (e) { modalError.value = e.response?.data?.error || 'Failed to save' }
  finally { saving.value = false }
}
async function deleteDest(id) {
  if (!confirm('Delete this destination and all its reviews?')) return
  try { await api.delete(`/destinations/${id}`); await fetchDestinations() }
  catch (e) { destError.value = e.response?.data?.error || 'Failed to delete' }
}
async function deleteReview(id) {
  if (!confirm('Remove this review?')) return
  try { await api.delete(`/reviews/${id}`); await fetchReviews() }
  catch (e) { reviewError.value = e.response?.data?.error || 'Failed to delete' }
}

onMounted(async () => { await fetchDestinations(); fetchReviews() })
</script>

<style scoped>
.admin-layout {
  display: grid;
  grid-template-columns: 260px 1fr;
  min-height: calc(100vh - var(--nav-h));
}

/* Sidebar */
.admin-sidebar {
  background: #fff;
  border-right: 1px solid rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
  padding: 1.5rem;
  position: sticky;
  top: var(--nav-h);
  height: calc(100vh - var(--nav-h));
  overflow-y: auto;
}

.admin-brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.admin-brand-icon {
  width: 40px;
  height: 40px;
  background: var(--green-deep);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  color: #fff;
}

.admin-brand-title {
  font-weight: 700;
  color: var(--text-dark);
  font-size: 0.95rem;
}

.admin-brand-sub {
  font-size: 0.72rem;
  color: var(--text-muted);
}

.admin-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.admin-nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: var(--r-md);
  border: none;
  background: transparent;
  color: var(--text-body);
  font-size: 0.9rem;
  font-family: var(--font-body);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.admin-nav-item:hover {
  background: var(--bg-cream);
  color: var(--text-dark);
}

.admin-nav-item.active {
  background: rgba(45, 74, 62, 0.06);
  color: var(--green-deep);
  border: 1px solid rgba(45, 74, 62, 0.12);
}

.admin-nav-count {
  margin-left: auto;
  background: var(--bg-cream);
  color: var(--text-muted);
  font-size: 0.72rem;
  padding: 0.15rem 0.5rem;
  border-radius: var(--r-full);
}

.admin-sidebar-footer {
  margin-top: auto;
  padding-top: 1rem;
}

/* Main */
.admin-main {
  padding: 2rem;
  overflow-y: auto;
  background: var(--bg-cream);
}

.admin-page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 1.75rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.admin-table-card {
  background: #fff;
  border-radius: var(--r-lg);
  overflow: hidden;
  box-shadow: var(--shadow-card);
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.dest-table-name {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.dest-table-thumb {
  width: 44px;
  height: 44px;
  border-radius: var(--r-sm);
  object-fit: cover;
  flex-shrink: 0;
}

.img-preview-wrap {
  margin-bottom: 1rem;
  border-radius: var(--r-md);
  overflow: hidden;
}

.img-preview {
  width: 100%;
  height: 140px;
  object-fit: cover;
  display: block;
}

@media (max-width: 900px) {
  .admin-layout {
    grid-template-columns: 1fr;
  }

  .admin-sidebar {
    position: static;
    height: auto;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 1rem;
  }

  .admin-nav {
    flex-direction: row;
  }

  .admin-sidebar-footer {
    display: none;
  }
}
</style>
