import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('token') || null)
    const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

    const isLoggedIn = computed(() => !!token.value)
    const isAdmin = computed(() => user.value?.role === 'admin')
    const userName = computed(() => user.value?.name || '')

    function setAuth(newToken, newUser) {
        token.value = newToken
        user.value = newUser
        localStorage.setItem('token', newToken)
        localStorage.setItem('user', JSON.stringify(newUser))
    }

    function logout() {
        token.value = null
        user.value = null
        localStorage.removeItem('token')
        localStorage.removeItem('user')
    }

    return { token, user, isLoggedIn, isAdmin, userName, setAuth, logout }
})
