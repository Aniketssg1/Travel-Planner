import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const Home = () => import('../views/Home.vue')
const Login = () => import('../views/Login.vue')
const Register = () => import('../views/Register.vue')
const DestinationDetail = () => import('../views/DestinationDetail.vue')
const MyItineraries = () => import('../views/MyItineraries.vue')
const AdminDashboard = () => import('../views/AdminDashboard.vue')
const NotFound = () => import('../views/NotFound.vue')

const routes = [
    { path: '/', name: 'Home', component: Home },
    { path: '/login', name: 'Login', component: Login },
    { path: '/register', name: 'Register', component: Register },
    { path: '/destinations/:id', name: 'DestinationDetail', component: DestinationDetail },
    {
        path: '/itineraries',
        name: 'MyItineraries',
        component: MyItineraries,
        meta: { requiresAuth: true },
    },
    {
        path: '/admin',
        name: 'AdminDashboard',
        component: AdminDashboard,
        meta: { requiresAuth: true, requiresAdmin: true },
    },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, from, next) => {
    const auth = useAuthStore()

    if (to.meta.requiresAuth && !auth.isLoggedIn) {
        return next({ name: 'Login', query: { redirect: to.fullPath } })
    }
    if (to.meta.requiresAdmin && !auth.isAdmin) {
        return next({ name: 'Home' })
    }
    if ((to.name === 'Login' || to.name === 'Register') && auth.isLoggedIn) {
        return next({ name: 'Home' })
    }
    next()
})

export default router
