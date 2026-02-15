<script setup lang="ts">
import { RouterView } from 'vue-router'
import Navbar from './components/Navbar.vue'
import { ref, onMounted } from 'vue'
import { useAuthStore } from './store/auth'

const authStore = useAuthStore()

// Toast System
export interface Toast {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

const toasts = ref<Toast[]>([])
let toastId = 0

const addToast = (message: string, type: 'success' | 'error' | 'info' = 'success') => {
  const id = toastId++
  toasts.value.push({ id, message, type })
  setTimeout(() => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }, 3000)
}

// Provide toast function globally
import { provide } from 'vue'
provide('addToast', addToast)

onMounted(() => {
  if (authStore.isAuthenticated) {
    authStore.fetchMe()
  }
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-base-200">
    <Navbar />
    
    <main class="flex-1 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <RouterView />
    </main>

    <footer class="footer footer-center p-10 bg-base-100 text-base-content border-t border-base-300">
      <aside>
        <p class="font-bold text-primary">PIZZA GOOD Ltd.</p> 
        <p>Providing deliciousness since 2026</p>
      </aside>
    </footer>

    <!-- Toast Container -->
    <div class="toast toast-top toast-center z-[100]">
      <div v-for="toast in toasts" :key="toast.id" 
        class="alert shadow-lg border-none animate-bounce"
        :class="{
          'alert-success bg-green-500 text-white': toast.type === 'success',
          'alert-error bg-red-500 text-white': toast.type === 'error',
          'alert-info bg-blue-500 text-white': toast.type === 'info'
        }"
      >
        <span>{{ toast.message }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Scoped styles if needed */
</style>
