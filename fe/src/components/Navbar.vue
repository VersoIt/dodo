<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ShoppingCart, User, Menu, X, Pizza } from 'lucide-vue-next'
import { useCartStore } from '../store/cart'

const router = useRouter()
const isMenuOpen = ref(false)
const cartStore = useCartStore()

const navLinks = [
  { name: 'Menu', path: '/' },
  { name: 'About', path: '/about' },
]

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}
</script>

<template>
  <div class="navbar bg-primary text-primary-content shadow-lg sticky top-0 z-50">
    <div class="container mx-auto">
      <div class="flex-1">
        <router-link to="/" class="btn btn-ghost normal-case text-xl flex items-center gap-2">
          <Pizza class="w-8 h-8" />
          <span class="font-bold tracking-tight">PIZZA GOOD</span>
        </router-link>
      </div>
      
      <!-- Desktop Menu -->
      <div class="hidden md:flex items-center gap-4 px-4">
        <router-link v-for="link in navLinks" :key="link.path" :to="link.path" class="btn btn-ghost">
          {{ link.name }}
        </router-link>
      </div>

      <div class="flex-none flex items-center gap-2">
        <router-link to="/cart" class="btn btn-ghost btn-circle">
          <div class="indicator">
            <ShoppingCart class="h-5 w-5" />
            <span v-if="cartStore.totalItems > 0" class="badge badge-sm indicator-item badge-secondary">{{ cartStore.totalItems }}</span>
          </div>
        </router-link>
        
        <div class="dropdown dropdown-end">
          <label tabindex="0" class="btn btn-ghost btn-circle avatar">
            <User class="h-5 w-5" />
          </label>
          <ul tabindex="0" class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 text-base-content rounded-box w-52">
            <li><router-link to="/login">Login</router-link></li>
            <li><router-link to="/register">Register</router-link></li>
            <li><a>Profile</a></li>
            <li><a>My Orders</a></li>
            <div class="divider my-0"></div>
            <li><a>Logout</a></li>
          </ul>
        </div>

        <button class="btn btn-ghost btn-circle md:hidden" @click="toggleMenu">
          <Menu v-if="!isMenuOpen" class="h-5 w-5" />
          <X v-else class="h-5 w-5" />
        </button>
      </div>
    </div>
  </div>

  <!-- Mobile Menu Drawer -->
  <div v-if="isMenuOpen" class="md:hidden bg-base-200 p-4 border-b border-base-300">
    <ul class="menu w-full gap-2">
      <li v-for="link in navLinks" :key="link.path">
        <router-link :to="link.path" @click="isMenuOpen = false">{{ link.name }}</router-link>
      </li>
    </ul>
  </div>
</template>
