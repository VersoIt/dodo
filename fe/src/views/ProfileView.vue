<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../store/auth'
import { Mail, Shield, Calendar, User, Package } from 'lucide-vue-next'
import axios from 'axios'

const authStore = useAuthStore()
const activeTab = ref('profile')
const orders = ref<any[]>([])
const loadingOrders = ref(false)

const fetchOrders = async () => {
  try {
    loadingOrders.value = true
    const response = await axios.get('/api/v1/orders/my')
    if (response.data.success) {
      orders.value = response.data.data || []
    }
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  } finally {
    loadingOrders.value = false
  }
}

onMounted(() => {
  authStore.fetchMe()
})
</script>

<template>
  <div class="max-w-4xl mx-auto py-8 px-4">
    <div class="flex flex-col md:flex-row gap-8">
      <!-- Sidebar -->
      <div class="w-full md:w-64 space-y-2">
        <button 
          @click="activeTab = 'profile'"
          class="btn btn-ghost w-full justify-start gap-3"
          :class="{ 'btn-active bg-primary/10 text-primary': activeTab === 'profile' }"
        >
          <User class="w-5 h-5" />
          Profile Info
        </button>
        <button 
          @click="activeTab = 'orders'; fetchOrders()"
          class="btn btn-ghost w-full justify-start gap-3"
          :class="{ 'btn-active bg-primary/10 text-primary': activeTab === 'orders' }"
        >
          <Package class="w-5 h-5" />
          My Orders
        </button>
        <div class="divider"></div>
        <button 
          @click="authStore.logout(); $router.push('/login')"
          class="btn btn-ghost w-full justify-start gap-3 text-error hover:bg-error/10"
        >
          Logout
        </button>
      </div>

      <!-- Main Content -->
      <div class="flex-1">
        <div v-if="activeTab === 'profile'" class="card bg-base-100 shadow-xl border border-base-200">
          <div class="card-body">
            <div class="flex items-center gap-6 mb-8">
              <div class="avatar placeholder">
                <div class="bg-primary text-primary-content rounded-2xl w-20">
                  <span class="text-3xl font-bold uppercase">{{ authStore.user?.name?.[0] || authStore.user?.email?.[0] || 'U' }}</span>
                </div>
              </div>
              <div>
                <h1 class="text-3xl font-bold">{{ authStore.user?.name || 'Pizza Lover' }}</h1>
                <p class="text-base-content/60">{{ authStore.user?.email }}</p>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="p-4 bg-base-200 rounded-2xl space-y-1">
                <p class="text-xs font-bold uppercase text-base-content/40">Full Name</p>
                <p class="font-semibold">{{ authStore.user?.name || 'Not set' }}</p>
              </div>
              <div class="p-4 bg-base-200 rounded-2xl space-y-1">
                <p class="text-xs font-bold uppercase text-base-content/40">Email</p>
                <p class="font-semibold">{{ authStore.user?.email }}</p>
              </div>
              <div class="p-4 bg-base-200 rounded-2xl space-y-1">
                <p class="text-xs font-bold uppercase text-base-content/40">Role</p>
                <div class="badge badge-secondary">{{ authStore.user?.role || 'client' }}</div>
              </div>
              <div class="p-4 bg-base-200 rounded-2xl space-y-1">
                <p class="text-xs font-bold uppercase text-base-content/40">Member Since</p>
                <p class="font-semibold">Feb 2026</p>
              </div>
            </div>

            <div class="card-actions mt-8 justify-end">
              <button class="btn btn-primary">Edit Profile</button>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'orders'" class="space-y-4">
          <h2 class="text-2xl font-bold mb-4">Order History</h2>
          
          <div v-if="loadingOrders" class="flex justify-center py-12">
            <span class="loading loading-spinner loading-lg text-primary"></span>
          </div>

          <div v-else-if="orders.length === 0" class="card bg-base-100 border border-dashed border-base-300 py-12">
            <div class="flex flex-col items-center text-center">
              <Package class="w-16 h-16 text-base-content/20 mb-4" />
              <h3 class="text-xl font-bold">No orders yet</h3>
              <p class="text-base-content/60 mb-6">Hungry? Treat yourself to a delicious pizza!</p>
              <router-link to="/" class="btn btn-primary">Browse Menu</router-link>
            </div>
          </div>

          <div v-else v-for="order in orders" :key="order.order_id" class="card bg-base-100 shadow-md border border-base-200 hover:border-primary/30 transition-colors">
            <div class="card-body p-6 flex flex-row justify-between items-center">
              <div>
                <div class="flex items-center gap-3 mb-1">
                  <span class="font-mono font-bold text-lg">#{{ order.order_number }}</span>
                  <div class="badge badge-outline uppercase text-[10px] font-bold" 
                    :class="{ 
                      'badge-success': order.status === 'paid',
                      'badge-warning': order.status === 'pending'
                    }">
                    {{ order.status }}
                  </div>
                </div>
                <p class="text-sm text-base-content/60">Feb 14, 2026</p>
              </div>
              <div class="text-right">
                <p class="text-xl font-black text-primary">${{ order.final_price?.toFixed(2) }}</p>
                <router-link :to="'/order/' + order.order_id" class="btn btn-ghost btn-xs text-primary">Details</router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
