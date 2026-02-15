<script setup lang="ts">
import { ref, onMounted, inject } from 'vue'
import axios from 'axios'
import { ChefHat, Clock, CheckCircle2, Play } from 'lucide-vue-next'

const orders = ref<any[]>([])
const loading = ref(true)
const addToast = inject('addToast') as (msg: string, type?: any) => void

// Status constants matching backend
const STATUS_PAID = 'paid'
const STATUS_COOKING = 'cooking'
const STATUS_READY = 'ready'

const fetchKitchenOrders = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/v1/orders/all') 
    if (response.data.success) {
      orders.value = response.data.data.filter((o: any) => o.status === STATUS_PAID || o.status === STATUS_COOKING)
    }
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  } finally {
    loading.value = false
  }
}

const updateStatus = async (orderId: string, newStatus: string) => {
  try {
    // Note: We need this endpoint in the gateway!
    await axios.patch(`/api/v1/orders/${orderId}/status`, { status: newStatus })
    addToast('Order status updated', 'success')
    await fetchKitchenOrders()
  } catch (err: any) {
    addToast(err.response?.data?.error || 'Failed to update status', 'error')
  }
}

onMounted(fetchKitchenOrders)
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex items-center gap-4 mb-8">
      <div class="bg-primary p-3 rounded-2xl shadow-lg">
        <ChefHat class="w-8 h-8 text-primary-content" />
      </div>
      <div>
        <h1 class="text-3xl font-bold tracking-tight">Kitchen Dashboard</h1>
        <p class="text-base-content/60">Manage active cooking orders</p>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <div v-else-if="orders.length === 0" class="text-center py-20 bg-base-100 rounded-3xl border-2 border-dashed border-base-300">
      <Clock class="w-16 h-16 mx-auto mb-4 text-base-content/20" />
      <h2 class="text-xl font-semibold">No orders in queue</h2>
      <p class="text-base-content/60">Everything is cooked! Take a short break.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="order in orders" :key="order.order_id" class="card bg-base-100 shadow-xl border border-base-200">
        <div class="card-body">
          <div class="flex justify-between items-start mb-4">
            <div>
              <h2 class="card-title text-primary">{{ order.order_number }}</h2>
              <p class="text-xs opacity-60">{{ order.order_id.split('-')[0] }}</p>
            </div>
            <div class="badge" :class="order.status === STATUS_COOKING ? 'badge-warning' : 'badge-info'">
              {{ order.status === STATUS_COOKING ? 'Cooking' : 'Paid' }}
            </div>
          </div>

          <div class="divider my-0"></div>
          
          <div class="py-4 space-y-2">
            <div v-for="item in order.items" :key="item.product_id" class="flex justify-between items-center">
              <span class="font-medium"><span class="text-primary font-bold">x{{ item.quantity }}</span> {{ item.product_name }}</span>
            </div>
          </div>

          <div class="card-actions justify-end mt-4">
            <button 
              v-if="order.status === STATUS_PAID"
              @click="updateStatus(order.order_id, STATUS_COOKING)"
              class="btn btn-primary btn-block gap-2"
            >
              <Play class="w-4 h-4" /> Start Cooking
            </button>
            <button 
              v-if="order.status === STATUS_COOKING"
              @click="updateStatus(order.order_id, STATUS_READY)"
              class="btn btn-success btn-block gap-2"
            >
              <CheckCircle2 class="w-4 h-4" /> Mark as Ready
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
