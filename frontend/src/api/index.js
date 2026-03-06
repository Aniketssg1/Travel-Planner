import axios from 'axios'

const api = axios.create({
    baseURL: '/api',
    headers: { 'Content-Type': 'application/json' },
    timeout: 30000, // 30 seconds
})

// ── Request Interceptor ────────────────────────────────────
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// ── Response Interceptor with Retry ────────────────────────
const MAX_RETRIES = 3
const RETRY_DELAY = 1000 // ms (will use exponential backoff)

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const config = error.config

        // Auth error – redirect to login
        if (error.response?.status === 401) {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            window.location.href = '/login'
            return Promise.reject(error)
        }

        // Retry on network errors or 429/5xx (idempotent methods only)
        const isRetryable =
            (!error.response || error.response.status === 429 || error.response.status >= 500) &&
            ['get', 'head', 'options'].includes((config.method || '').toLowerCase())

        if (isRetryable && (!config._retryCount || config._retryCount < MAX_RETRIES)) {
            config._retryCount = (config._retryCount || 0) + 1
            const delay = RETRY_DELAY * Math.pow(2, config._retryCount - 1)
            await new Promise((resolve) => setTimeout(resolve, delay))
            return api(config)
        }

        return Promise.reject(error)
    }
)

export default api
