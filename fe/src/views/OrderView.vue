<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import { CheckCircle2, Clock, Truck, PackageCheck, CreditCard, MapPin, ChefHat } from 'lucide-vue-next'

const route = useRoute()
const orderId = route.params.id as string
const order = ref<any>(null)
const loading = ref(true)
let pollInterval: any = null

const fetchOrder = async () => {
  try {
    const response = await axios.get(`/api/v1/orders/${orderId}`)
    if (response.data.success) {
      order.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to fetch order:', error)
  } finally {
    loading.value = false
  }
}

// Map string statuses to index for the progress bar
const statusMap: Record<string, number> = {
  'created': 0,
  'paid': 1,
  'cooking': 2,
  'ready': 3,
  'delivering': 4,
  'completed': 5,
  'canceled': -1
}

const getStatusIndex = computed(() => {
  if (!order.value) return 0
  return statusMap[order.value.status] ?? 0
})

const steps = [
  { label: 'Received', icon: PackageCheck },
  { label: 'Paid', icon: CreditCard },
  { label: 'Cooking', icon: ChefHat },
  { label: 'Ready', icon: CheckCircle2 },
  { label: 'On the Way', icon: Truck },
  { label: 'Enjoy!', icon: CheckCircle2 }
]

const handlePayment = async () => {
  try {
    const resp = await axios.post(`/api/v1/orders/${orderId}/pay`)
    if (resp.data.success) {
      fetchOrder()
    }
  } catch (error) {
    console.error('Payment failed')
  }
}

onMounted(() => {
  fetchOrder()
  pollInterval = setInterval(fetchOrder, 5000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<template>
  <div class="max-w-4xl mx-auto py-8 px-4 pb-20">
    <div v-if="loading" class="flex justify-center py-20">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>

    <div v-else-if="order" class="space-y-8 animate-in fade-in duration-500">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <div class="flex items-center gap-3">
            <h1 class="text-3xl font-black tracking-tight uppercase">Order Details</h1>
            <div class="badge badge-primary font-bold uppercase text-[10px] tracking-widest px-3 py-3">#{{ order.order_number }}</div>
          </div>
          <p class="text-base-content/60 mt-1">Status tracked in real-time</p>
        </div>
        <div class="badge badge-lg py-4 px-6 gap-2 font-bold uppercase text-xs" 
          :class="order.status === 'completed' ? 'badge-success' : 'badge-warning'">
          {{ order.status }}
        </div>
      </div>

      <!-- Tracking Steps -->
      <div class="card bg-base-100 shadow-2xl border border-base-200 overflow-hidden">
        <div class="card-body p-8">
          <h2 class="card-title mb-8 text-xl font-bold flex items-center gap-2">
            <Clock class="w-5 h-5 text-primary" /> Track Order
          </h2>
          <ul class="steps steps-vertical md:steps-horizontal w-full">
            <li 
              v-for="(step, index) in steps" 
              :key="index"
              class="step"
              :class="{ 'step-primary': index <= getStatusIndex }"
            >
              <div class="flex flex-col items-center gap-1">
                <component :is="step.icon" class="w-5 h-5 mb-1" :class="index <= getStatusIndex ? 'text-primary' : 'opacity-20'" />
                <span class="text-[10px] font-black uppercase tracking-tighter">{{ step.label }}</span>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Order Summary -->
        <div class="card bg-base-100 shadow-xl border border-base-200 h-fit">
          <div class="card-body p-8">
            <h2 class="card-title mb-6 font-bold">Items Summary</h2>
            <div class="divide-y divide-base-200">
              <div v-for="item in order.items" :key="item.product_id" class="py-4 flex justify-between items-center">
                <div>
                  <p class="font-bold"><span class="text-primary font-black">x{{ item.quantity }}</span> {{ item.product_name }}</p>
                </div>
              </div>
              <div class="pt-6 mt-4 space-y-2">
                <div class="flex justify-between text-2xl font-black pt-4">
                  <span>Total</span>
                  <span class="text-primary">${{ order.final_price ? order.final_price.toFixed(2) : '0.00' }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Delivery Info -->
        <div class="space-y-6">
          <div class="card bg-base-100 shadow-xl border border-base-200">
            <div class="card-body p-8">
              <h2 class="card-title mb-4 font-bold flex items-center gap-2">
                <MapPin class="w-5 h-5 text-error" /> Delivery Address
              </h2>
              <div class="p-4 bg-base-200/50 rounded-2xl border border-base-300">
                <p class="font-bold text-lg">{{ order.address?.street }}</p>
                <p class="text-base-content/60">{{ order.address?.city }}</p>
              </div>
            </div>
          </div>

          <div v-if="order.status === 'created'" class="card bg-primary text-primary-content shadow-xl shadow-primary/20 overflow-hidden">
            <div class="card-body p-8">
              <h2 class="card-title text-2xl font-black uppercase tracking-tighter">Action Required</h2>
              <p class="opacity-80">This order is not paid yet. Pay now to start cooking!</p>
              <div class="card-actions mt-4">
                <button @click="handlePayment" class="btn btn-neutral btn-block btn-lg rounded-xl shadow-xl">Pay Now</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.step::after {
  @apply bg-base-300 opacity-20;
}
.step-primary::after {
  @apply bg-primary opacity-100;
}
</style>
