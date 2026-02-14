<script setup lang="ts">
import { useAuthStore } from '../store/auth'
import { User, Mail, Shield, Calendar } from 'lucide-vue-next'

const authStore = useAuthStore()
</script>

<template>
  <div class="max-w-2xl mx-auto py-8">
    <div class="card bg-base-100 shadow-2xl border border-base-200">
      <div class="card-body">
        <div class="flex flex-col items-center gap-4 mb-8">
          <div class="avatar placeholder">
            <div class="bg-neutral text-neutral-content rounded-full w-24">
              <span class="text-3xl uppercase">{{ authStore.user?.email?.[0] || 'U' }}</span>
            </div>
          </div>
          <h1 class="text-3xl font-bold">Your Profile</h1>
        </div>

        <div class="space-y-6">
          <div class="flex items-center gap-4 p-4 bg-base-200 rounded-xl">
            <Mail class="w-6 h-6 text-primary" />
            <div>
              <p class="text-xs font-bold uppercase text-base-content/50">Email Address</p>
              <p class="font-semibold">{{ authStore.user?.email || 'Not logged in' }}</p>
            </div>
          </div>

          <div class="flex items-center gap-4 p-4 bg-base-200 rounded-xl">
            <Shield class="w-6 h-6 text-secondary" />
            <div>
              <p class="text-xs font-bold uppercase text-base-content/50">Role</p>
              <p class="font-semibold">Customer</p>
            </div>
          </div>

          <div class="flex items-center gap-4 p-4 bg-base-200 rounded-xl">
            <Calendar class="w-6 h-6 text-accent" />
            <div>
              <p class="text-xs font-bold uppercase text-base-content/50">Account Created</p>
              <p class="font-semibold">February 2026</p>
            </div>
          </div>
        </div>

        <div class="card-actions mt-8 justify-end">
          <button class="btn btn-outline btn-error" @click="authStore.logout(); $router.push('/login')">Logout</button>
          <button class="btn btn-primary">Update Profile</button>
        </div>
      </div>
    </div>
  </div>
</template>
