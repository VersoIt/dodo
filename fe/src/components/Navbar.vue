<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ShoppingCart, User, Menu, X, Pizza, LogOut, Package, UserCircle } from 'lucide-vue-next'
import { useCartStore } from '../store/cart'
import { useAuthStore } from '../store/auth'

const router = useRouter()
const isMenuOpen = ref(false)
const cartStore = useCartStore()
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const navLinks = [
  { name: 'Menu', path: '/' },
  { name: 'About', path: '/about' },
]

const roleLinks = {
  chef: { name: 'Kitchen', path: '/kitchen' },
  courier: { name: 'Logistics', path: '/logistics' },
  admin: [
    { name: 'Kitchen', path: '/kitchen' },
    { name: 'Logistics', path: '/logistics' }
  ]
}

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}
</script>

<template>
  <div class="navbar bg-primary text-primary-content shadow-lg sticky top-0 z-50">
    <div class="container mx-auto px-4">
      <div class="flex-1">
        <router-link :to="authStore.user?.role === 'chef' ? '/kitchen' : (authStore.user?.role === 'courier' ? '/logistics' : '/')" class="btn btn-ghost normal-case text-xl flex items-center gap-2 p-0 hover:bg-transparent">
          <Pizza class="w-8 h-8" />
          <span class="font-bold tracking-tighter">PIZZA GOOD</span>
        </router-link>
      </div>
      
      <!-- Desktop Menu -->
      <div v-if="authStore.user?.role === 'client' || authStore.user?.role === 'manager' || !authStore.isAuthenticated" class="hidden md:flex items-center gap-2 px-4">
        <router-link v-for="link in navLinks" :key="link.path" :to="link.path" class="btn btn-ghost btn-sm rounded-lg">
          {{ link.name }}
        </router-link>
      </div>

      <div class="flex-none flex items-center gap-3">
        <!-- Dashboard links for Manager always visible, for others only if they are that role -->
        <div class="hidden md:flex items-center gap-2">
          <template v-if="authStore.user?.role === 'chef' || authStore.user?.role === 'manager'">
            <router-link to="/kitchen" class="btn btn-ghost btn-sm rounded-lg text-secondary font-bold">Kitchen</router-link>
          </template>
          <template v-if="authStore.user?.role === 'courier' || authStore.user?.role === 'manager'">
            <router-link to="/logistics" class="btn btn-ghost btn-sm rounded-lg text-accent font-bold">Logistics</router-link>
          </template>
        </div>

        <!-- User Name (Desktop) -->
        <span v-if="authStore.isAuthenticated" class="hidden md:block text-xs font-bold uppercase tracking-widest opacity-80">
          Hi, {{ authStore.user?.name || 'Friend' }}
        </span>

        <router-link v-if="authStore.isAuthenticated && authStore.user?.role === 'client'" to="/cart" class="btn btn-ghost btn-circle btn-sm">
          <div class="indicator">
            <ShoppingCart class="h-5 w-5" />
            <span v-if="cartStore.totalItems > 0" class="badge badge-sm indicator-item badge-secondary font-bold">{{ cartStore.totalItems }}</span>
          </div>
        </router-link>
        
        <div class="dropdown dropdown-end">
          <label tabindex="0" class="btn btn-ghost btn-circle btn-sm">
            <User v-if="!authStore.isAuthenticated" class="h-5 w-5" />
            <div v-else class="avatar placeholder">
              <div class="bg-primary-focus text-primary-content rounded-full w-8">
                <span class="text-xs font-bold uppercase">{{ authStore.user?.name?.[0] || 'U' }}</span>
              </div>
            </div>
          </label>
          <ul tabindex="0" class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow-2xl bg-base-100 text-base-content rounded-2xl w-56 border border-base-200">
            <template v-if="!authStore.isAuthenticated">
              <li class="menu-title px-4 py-2 font-bold text-primary">Guest</li>
              <li><router-link to="/login">Login</router-link></li>
              <li><router-link to="/register">Register</router-link></li>
            </template>
            <template v-else>
              <li class="menu-title px-4 py-2 font-bold text-primary truncate max-w-full">
                {{ authStore.user?.name || 'My Account' }}
              </li>
              <li><router-link to="/profile" class="py-3"><UserCircle class="w-4 h-4" /> Profile Info</router-link></li>
              <li v-if="authStore.user?.role === 'client'"><router-link to="/profile" class="py-3"><Package class="w-4 h-4" /> My Orders</router-link></li>
              <div class="divider my-0 opacity-50"></div>
              <li><a @click="handleLogout" class="text-error font-bold py-3"><LogOut class="w-4 h-4" /> Logout</a></li>
            </template>
          </ul>
        </div>

        <button class="btn btn-ghost btn-circle btn-sm md:hidden" @click="toggleMenu">
          <Menu v-if="!isMenuOpen" class="h-5 w-5" />
          <X v-else class="h-5 w-5" />
        </button>
      </div>
    </div>
  </div>

  <!-- Mobile Menu Drawer -->
  <div v-if="isMenuOpen" class="md:hidden bg-base-100 p-4 border-b border-base-200 shadow-xl">
    <ul class="menu w-full gap-2">
      <template v-if="authStore.user?.role === 'client' || authStore.user?.role === 'manager' || !authStore.isAuthenticated">
        <li v-for="link in navLinks" :key="link.path">
          <router-link :to="link.path" @click="isMenuOpen = false" class="py-3 font-semibold">{{ link.name }}</router-link>
        </li>
      </template>
      <li v-if="authStore.user?.role === 'chef' || authStore.user?.role === 'manager'">
        <router-link to="/kitchen" @click="isMenuOpen = false" class="py-3 font-bold text-secondary">Kitchen Dashboard</router-link>
      </li>
      <li v-if="authStore.user?.role === 'courier' || authStore.user?.role === 'manager'">
        <router-link to="/logistics" @click="isMenuOpen = false" class="py-3 font-bold text-accent">Logistics Dashboard</router-link>
      </li>
    </ul>
  </div>
</template>
