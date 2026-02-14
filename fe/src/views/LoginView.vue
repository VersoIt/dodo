<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { LogIn, Mail, Lock } from 'lucide-vue-next'
import axios from 'axios'
import { useAuthStore } from '../store/auth'

const router = useRouter()
const authStore = useAuthStore()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await axios.post('/api/v1/auth/login', {
      email: email.value,
      password: password.value
    })
    
    // Backend returns { success: true, data: { token: "..." } }
    const result = response.data
    if (result.success && result.data?.token) {
      authStore.setToken(result.data.token)
      router.push('/')
    } else {
      error.value = result.error || 'Invalid response from server'
    }
  } catch (err: any) {
    console.error('Login failed:', err)
    error.value = 'Invalid email or password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex justify-center items-center py-12">
    <div class="card w-full max-w-md bg-base-100 shadow-2xl border border-base-200">
      <div class="card-body">
        <div class="flex flex-col items-center gap-2 mb-6">
          <div class="p-3 bg-primary rounded-2xl text-primary-content">
            <LogIn class="w-8 h-8" />
          </div>
          <h1 class="text-3xl font-bold">Welcome Back</h1>
          <p class="text-base-content/60">Log in to your Pizza Good account</p>
        </div>

        <div v-if="error" class="alert alert-error mb-6">
          <span>{{ error }}</span>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Email Address</span>
            </label>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-base-content/50">
                <Mail class="w-5 h-5" />
              </span>
              <input 
                v-model="email"
                type="email" 
                placeholder="you@example.com" 
                class="input input-bordered w-full pl-10" 
                required 
              />
            </div>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">Password</span>
              <a href="#" class="label-text-alt link link-primary">Forgot password?</a>
            </label>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-base-content/50">
                <Lock class="w-5 h-5" />
              </span>
              <input 
                v-model="password"
                type="password" 
                placeholder="••••••••" 
                class="input input-bordered w-full pl-10" 
                required 
              />
            </div>
          </div>

          <div class="card-actions mt-6">
            <button class="btn btn-primary btn-block" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner"></span>
              Sign In
            </button>
          </div>
        </form>

        <div class="divider">OR</div>

        <p class="text-center text-sm">
          Don't have an account? 
          <router-link to="/register" class="link link-primary font-bold">Create one now</router-link>
        </p>
      </div>
    </div>
  </div>
</template>
