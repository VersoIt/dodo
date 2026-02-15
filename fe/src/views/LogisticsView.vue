<script setup lang="ts">
import { ref, onMounted, inject } from 'vue'
import axios from 'axios'
import { Truck, MapPin, CheckCircle2, Package } from 'lucide-vue-next'

const orders = ref<any[]>([])
const loading = ref(true)
const addToast = inject('addToast') as (msg: string, type?: any) => void

// Status constants
const STATUS_READY = 'ready'
const STATUS_DELIVERING = 'delivering'
const STATUS_COMPLETED = 'completed'

const fetchLogisticsOrders = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/v1/orders/all') // Fetch all active orders
    if (response.data.success) {
      orders.value = response.data.data.filter((o: any) => o.status === STATUS_READY || o.status === STATUS_DELIVERING)
    }
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  } finally {
    loading.value = false
  }
}

const updateStatus = async (orderId: string, newStatus: string) => {
  try {
    await axios.patch(`/api/v1/orders/${orderId}/status`, { status: newStatus })
    addToast('Delivery status updated', 'success')
    await fetchLogisticsOrders()
  } catch (err: any) {
    addToast(err.response?.data?.error || 'Failed to update status', 'error')
  }
}

onMounted(fetchLogisticsOrders)
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex items-center gap-4 mb-8">
      <div class="bg-secondary p-3 rounded-2xl shadow-lg">
        <Truck class="w-8 h-8 text-secondary-content" />
      </div>
      <div>
        <h1 class="text-3xl font-bold tracking-tight">Logistics Dashboard</h1>
        <p class="text-base-content/60">Manage deliveries and routes</p>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <span class="loading loading-spinner loading-lg text-secondary"></span>
    </div>

    <div v-else-if="orders.length === 0" class="text-center py-20 bg-base-100 rounded-3xl border-2 border-dashed border-base-300">
      <Package class="w-16 h-16 mx-auto mb-4 text-base-content/20" />
      <h2 class="text-xl font-semibold">No deliveries available</h2>
      <p class="text-base-content/60">All pizzas have been delivered. Good job!</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div v-for="order in orders" :key="order.order_id" class="card bg-base-100 shadow-xl border border-base-200 overflow-hidden">
        <div class="bg-secondary/10 px-6 py-3 flex justify-between items-center border-b border-secondary/20">
           <span class="font-bold text-secondary">{{ order.order_number }}</span>
           <div class="badge badge-secondary" v-if="order.status === STATUS_READY">Ready for pickup</div>
           <div class="badge badge-info" v-else>On the way</div>
        </div>
        
        <div class="card-body">
          <div class="flex items-start gap-3 mb-4">
            <MapPin class="w-5 h-5 text-error mt-1 flex-shrink-0" />
            <div>
              <p class="font-bold">Delivery Address:</p>
              <p class="text-base-content/70">{{ order.address.street }}, {{ order.address.city }}</p>
              <p v-if="order.address.house" class="text-sm">House: {{ order.address.house }}, Apt: {{ order.address.apartment }}</p>
            </div>
          </div>

          <div class="divider my-0"></div>

          <div class="py-4">
            <p class="text-sm font-semibold mb-2">Order Contents:</p>
            <div class="flex flex-wrap gap-2">
              <span v-for="item in order.items" :key="item.product_id" class="badge badge-outline">
                {{ item.quantity }}x {{ item.product_name }}
              </span>
            </div>
          </div>

          <div class="card-actions justify-end mt-4">
            <button 
              v-if="order.status === STATUS_READY"
              @click="updateStatus(order.order_id, STATUS_DELIVERING)"
              class="btn btn-secondary btn-block gap-2"
            >
              <Truck class="w-4 h-4" /> Pick Up Order
            </button>
            <button 
              v-if="order.status === STATUS_DELIVERING"
              @click="updateStatus(order.order_id, STATUS_COMPLETED)"
              class="btn btn-success btn-block gap-2"
            >
              <CheckCircle2 class="w-4 h-4" /> Delivered
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
