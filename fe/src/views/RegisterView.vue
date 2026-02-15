<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UserPlus, Mail, Lock, User } from 'lucide-vue-next'
import axios from 'axios'

import { useAuthStore } from '../store/auth'

const router = useRouter()
const authStore = useAuthStore()
const name = ref('')
const email = ref('')
const password = ref('')
const role = ref('client')
const loading = ref(false)
const error = ref('')

const handleRegister = async () => {
  try {
    loading.value = true
    error.value = ''
    
    // 1. Register the user
    const regResp = await axios.post('/api/v1/auth/register', {
      name: name.value,
      email: email.value,
      password: password.value,
      role: role.value
    })
    
    if (regResp.data.success) {
      // 2. Automatically login after successful registration
      const loginSuccess = await authStore.login(email.value, password.value)
      
      if (loginSuccess) {
        // 3. Redirect based on role
        if (role.value === 'chef') {
          router.push('/kitchen')
        } else if (role.value === 'courier') {
          router.push('/logistics')
        } else {
          router.push('/')
        }
      } else {
        router.push('/login')
      }
    } else {
      error.value = regResp.data.error || 'Failed to create account'
    }
  } catch (err: any) {
    console.error('Registration failed:', err)
    error.value = 'Failed to create account. Please try again.'
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
          <div class="p-3 bg-secondary rounded-2xl text-secondary-content">
            <UserPlus class="w-8 h-8" />
          </div>
          <h1 class="text-3xl font-bold">Join Us</h1>
          <p class="text-base-content/60">Create your Pizza Good account</p>
        </div>

        <div v-if="error" class="alert alert-error mb-6">
          <span>{{ error }}</span>
        </div>

        <form @submit.prevent="handleRegister" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Full Name</span>
            </label>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-base-content/50">
                <User class="w-5 h-5" />
              </span>
              <input 
                v-model="name"
                type="text" 
                placeholder="John Doe" 
                class="input input-bordered w-full pl-10" 
                required 
              />
            </div>
          </div>

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

          <div class="form-control">
            <label class="label">
              <span class="label-text">I am a...</span>
            </label>
            <select v-model="role" class="select select-bordered w-full">
              <option value="client">Client (Hungry Customer)</option>
              <option value="chef">Chef (Master Pizza Maker)</option>
              <option value="courier">Courier (Fast Delivery)</option>
            </select>
          </div>

          <div class="card-actions mt-6">
            <button class="btn btn-secondary btn-block" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner"></span>
              Create Account
            </button>
          </div>
        </form>

        <div class="divider">OR</div>

        <p class="text-center text-sm">
          Already have an account? 
          <router-link to="/login" class="link link-secondary font-bold">Sign In</router-link>
        </p>
      </div>
    </div>
  </div>
</template>
