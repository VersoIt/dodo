import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/cart',
      name: 'cart',
      component: () => import('../views/CartView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue')
    },
    {
      path: '/order/:id',
      name: 'order',
      component: () => import('../views/OrderView.vue')
    },
    {
      path: '/order/:id/success',
      name: 'order-success',
      component: () => import('../views/OrderSuccessView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue')
    },
    {
      path: '/kitchen',
      name: 'kitchen',
      component: () => import('../views/KitchenView.vue'),
      meta: { requiresAuth: true, role: 'chef' }
    },
    {
      path: '/logistics',
      name: 'logistics',
      component: () => import('../views/LogisticsView.vue'),
      meta: { requiresAuth: true, role: 'courier' }
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/manager',
      name: 'manager',
      component: () => import('../views/ManagerView.vue'),
      meta: { requiresAuth: true, role: 'manager' }
    }
  ]
})

router.beforeEach((to, from, next) => {
  const userStr = localStorage.getItem('user')
  const user = userStr && userStr !== 'null' ? JSON.parse(userStr) : {}
  const isAuthenticated = !!localStorage.getItem('token')

  // 1. Auth check
  if (to.meta.requiresAuth && !isAuthenticated) {
    return next({ name: 'login' })
  }

  // 2. Role-based access control
  if (to.meta.role && user.role !== to.meta.role && user.role !== 'manager') {
    return next({ name: 'home' })
  }

  // 3. Smart redirect for staff from Home page
  if (to.name === 'home' && isAuthenticated) {
    if (user.role === 'chef') {
      return next({ name: 'kitchen' })
    }
    if (user.role === 'courier') {
      return next({ name: 'logistics' })
    }
  }

  next()
})

export default router
